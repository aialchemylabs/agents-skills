"""Runtime configuration, URL policy, and immutable safety ceilings."""

from __future__ import annotations

import json
import os
import re
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, Mapping
from urllib.parse import urlsplit, urlunsplit

from .errors import AuthenticationBlocked, ValidationError

HARD_CAPS: dict[str, int] = {
    "members": 5,
    "max_concurrency": 5,
    "connect_timeout_seconds": 60,
    "request_timeout_seconds": 600,
    "max_transient_retries": 2,
    "max_inference_requests": 25,
    "max_input_characters": 500_000,
    "max_output_tokens": 50_000,
    "max_output_tokens_per_call": 8_192,
}

DEFAULTS: dict[str, Any] = {
    "execution": {
        "max_concurrency": 3,
        "request_timeout_seconds": 120,
        "connect_timeout_seconds": 10,
        "max_transient_retries": 1,
        "temperature": 0.2,
        "stage_output_tokens": {
            "independent": 1500,
            "critique": 1200,
            "revision": 1200,
            "chairman": 1800,
        },
    },
    "budget": {
        "max_inference_requests": 10,
        "max_input_characters": 120_000,
        "max_output_tokens": 12_000,
        "max_usd": None,
        "require_known_cost": False,
    },
    "policy": {
        "minimum_successful_members": 2,
        "require_distinct_families": True,
        "allow_perspective_only_fallback": False,
        "run_cross_critique": True,
        "run_addy_revision": "auto",
        "expose_model_identities_after_synthesis": True,
        "persist": "none",
        "output_directory": None,
    },
}

_FORBIDDEN_CONFIG_KEYS = {
    "api_key",
    "apikey",
    "authorization",
    "credential",
    "credentials",
    "master_key",
    "password",
    "secret",
    "token",
}


@dataclass(frozen=True)
class RuntimeConfig:
    base_url: str
    api_key: str | None
    policy: dict[str, Any] = field(default_factory=dict)
    config_path: Path | None = None


_POLICY_DEFAULT_SECTIONS = {"execution", "budget", "policy"}
_POLICY_MAXIMUMS = {
    "max_concurrency",
    "connect_timeout_seconds",
    "request_timeout_seconds",
    "max_transient_retries",
    "max_inference_requests",
    "max_input_characters",
    "max_output_tokens",
    "max_output_tokens_per_call",
    "max_usd",
}
_SECRET_CONFIG_PATTERNS = (
    re.compile(r"-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----", re.I),
    re.compile(r"\b(?:sk|ghp|gho|github_pat|xox[baprs])[-_][A-Za-z0-9_-]{12,}\b"),
    re.compile(r"\bAuthorization\s*:\s*Bearer\s+\S+", re.I),
)


def normalize_base_url(value: str, *, allow_localhost_http: bool = True) -> str:
    """Normalize a LiteLLM URL and reject credentials or unsafe cleartext hosts."""
    if not isinstance(value, str) or not value.strip():
        raise ValidationError("LLM_COUNCIL_BASE_URL must be a non-empty URL")
    parsed = urlsplit(value.strip())
    if parsed.username or parsed.password:
        raise ValidationError("Gateway URL must not contain credentials")
    if parsed.query or parsed.fragment:
        raise ValidationError("Gateway URL must not contain a query or fragment")
    if parsed.scheme not in {"https", "http"} or not parsed.hostname:
        raise ValidationError("Gateway URL must use http or https and include a host")
    local_hosts = {"localhost", "127.0.0.1", "::1"}
    if parsed.scheme != "https" and not (
        allow_localhost_http and parsed.hostname.lower() in local_hosts
    ):
        raise ValidationError("Gateway URL must use HTTPS except for localhost development")
    path = parsed.path.rstrip("/")
    return urlunsplit((parsed.scheme.lower(), parsed.netloc.lower(), path, "", ""))


def _reject_secret_keys(value: Any, location: str = "config") -> None:
    if isinstance(value, Mapping):
        for key, child in value.items():
            normalized = str(key).lower().replace("-", "_")
            if normalized in _FORBIDDEN_CONFIG_KEYS:
                raise ValidationError(f"{location} contains forbidden secret field '{key}'")
            _reject_secret_keys(child, f"{location}.{key}")
    elif isinstance(value, list):
        for index, child in enumerate(value):
            _reject_secret_keys(child, f"{location}[{index}]")
    elif isinstance(value, str) and any(
        pattern.search(value) for pattern in _SECRET_CONFIG_PATTERNS
    ):
        raise ValidationError(f"{location} contains a possible credential")


def _validate_policy_config(policy: dict[str, Any]) -> None:
    unknown = set(policy) - {"defaults", "maximums"}
    if unknown:
        raise ValidationError(
            "Unknown config.policy field(s): " + ", ".join(sorted(unknown))
        )
    defaults = policy.get("defaults", {})
    if not isinstance(defaults, dict):
        raise ValidationError("config.policy.defaults must be an object")
    unknown_defaults = set(defaults) - _POLICY_DEFAULT_SECTIONS
    if unknown_defaults:
        raise ValidationError(
            "Unknown config.policy.defaults section(s): "
            + ", ".join(sorted(unknown_defaults))
        )
    for section, value in defaults.items():
        if not isinstance(value, dict):
            raise ValidationError(f"config.policy.defaults.{section} must be an object")
    maximums = policy.get("maximums", {})
    if not isinstance(maximums, dict):
        raise ValidationError("config.policy.maximums must be an object")
    unknown_maximums = set(maximums) - _POLICY_MAXIMUMS
    if unknown_maximums:
        raise ValidationError(
            "Unknown config.policy.maximums field(s): "
            + ", ".join(sorted(unknown_maximums))
        )
    for key, value in maximums.items():
        if key == "max_usd":
            if (
                isinstance(value, bool)
                or not isinstance(value, (int, float))
                or value < 0
            ):
                raise ValidationError("config.policy.maximums.max_usd must be non-negative")
        elif key == "max_transient_retries":
            if isinstance(value, bool) or not isinstance(value, int) or value < 0:
                raise ValidationError(
                    "config.policy.maximums.max_transient_retries must be a non-negative integer"
                )
        elif isinstance(value, bool) or not isinstance(value, int) or value < 1:
            raise ValidationError(f"config.policy.maximums.{key} must be a positive integer")
    hard_names = {
        "max_concurrency": "max_concurrency",
        "connect_timeout_seconds": "connect_timeout_seconds",
        "request_timeout_seconds": "request_timeout_seconds",
        "max_transient_retries": "max_transient_retries",
        "max_inference_requests": "max_inference_requests",
        "max_input_characters": "max_input_characters",
        "max_output_tokens": "max_output_tokens",
        "max_output_tokens_per_call": "max_output_tokens_per_call",
    }
    for key, hard_name in hard_names.items():
        if key in maximums and maximums[key] > HARD_CAPS[hard_name]:
            raise ValidationError(
                f"config.policy.maximums.{key} may lower but not raise the hard cap"
            )


def _fill_missing(target: dict[str, Any], defaults: Mapping[str, Any]) -> dict[str, Any]:
    result = json.loads(json.dumps(target))
    for key, value in defaults.items():
        if key not in result:
            result[key] = json.loads(json.dumps(value))
        elif isinstance(result[key], dict) and isinstance(value, Mapping):
            result[key] = _fill_missing(result[key], value)
    return result


def apply_runtime_policy(
    manifest: Mapping[str, Any],
    runtime: RuntimeConfig,
) -> dict[str, Any]:
    """Apply non-secret deployment defaults and enforce stricter maximums."""
    effective = json.loads(json.dumps(manifest))
    defaults = runtime.policy.get("defaults", {})
    effective = _fill_missing(effective, defaults)
    # Import lazily to avoid config/schemas import cycles.
    from .schemas import validate_manifest

    normalized = validate_manifest(effective, command="preflight")
    maximums = runtime.policy.get("maximums", {})
    paths = {
        "max_concurrency": ("execution", "max_concurrency"),
        "connect_timeout_seconds": ("execution", "connect_timeout_seconds"),
        "request_timeout_seconds": ("execution", "request_timeout_seconds"),
        "max_transient_retries": ("execution", "max_transient_retries"),
        "max_inference_requests": ("budget", "max_inference_requests"),
        "max_input_characters": ("budget", "max_input_characters"),
        "max_output_tokens": ("budget", "max_output_tokens"),
    }
    for key, path in paths.items():
        if key in maximums and normalized[path[0]][path[1]] > maximums[key]:
            raise ValidationError(
                f"Manifest {path[0]}.{path[1]} exceeds the deployment policy maximum"
            )
    if "max_output_tokens_per_call" in maximums:
        for stage, value in normalized["execution"]["stage_output_tokens"].items():
            if value > maximums["max_output_tokens_per_call"]:
                raise ValidationError(
                    f"Manifest execution.stage_output_tokens.{stage} exceeds "
                    "the deployment policy maximum"
                )
    if "max_usd" in maximums:
        requested = normalized["budget"]["max_usd"]
        if requested is None:
            effective.setdefault("budget", {})["max_usd"] = maximums["max_usd"]
        elif requested > maximums["max_usd"]:
            raise ValidationError("Manifest budget.max_usd exceeds the deployment policy maximum")
    return effective


def load_runtime_config(
    env: Mapping[str, str] | None = None,
    *,
    require_key: bool = False,
) -> RuntimeConfig:
    """Load only skill-specific environment variables and a non-secret JSON policy."""
    source = os.environ if env is None else env
    base_url = normalize_base_url(source.get("LLM_COUNCIL_BASE_URL", ""))
    api_key = source.get("LLM_COUNCIL_API_KEY")
    if require_key and not api_key:
        raise AuthenticationBlocked("A scoped LLM_COUNCIL_API_KEY is required")
    config_path_raw = source.get("LLM_COUNCIL_CONFIG")
    policy: dict[str, Any] = {}
    config_path: Path | None = None
    if config_path_raw:
        config_path = Path(config_path_raw).expanduser().resolve()
        try:
            loaded = json.loads(config_path.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError) as exc:
            raise ValidationError("LLM_COUNCIL_CONFIG is not readable valid JSON") from exc
        if not isinstance(loaded, dict):
            raise ValidationError("LLM_COUNCIL_CONFIG must contain a JSON object")
        _reject_secret_keys(loaded)
        unknown = set(loaded) - {"schema_version", "policy"}
        if unknown:
            raise ValidationError(
                "Unknown configuration field(s): " + ", ".join(sorted(unknown))
            )
        if loaded.get("schema_version", 1) != 1:
            raise ValidationError("Unsupported configuration schema_version")
        if not isinstance(loaded.get("policy", {}), dict):
            raise ValidationError("config.policy must be an object")
        policy = dict(loaded.get("policy", {}))
        _validate_policy_config(policy)
    return RuntimeConfig(
        base_url=base_url,
        api_key=api_key,
        policy=policy,
        config_path=config_path,
    )
