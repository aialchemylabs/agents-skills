from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path

from support import valid_manifest

from llm_council.config import (
    RuntimeConfig,
    apply_runtime_policy,
    load_runtime_config,
    normalize_base_url,
)
from llm_council.errors import AuthenticationBlocked, ValidationError
from llm_council.schemas import validate_manifest


class RuntimeConfigTests(unittest.TestCase):
    def test_url_policy_requires_https_except_explicit_local_development(self) -> None:
        self.assertEqual(
            normalize_base_url("https://GATEWAY.invalid/base/"),
            "https://gateway.invalid/base",
        )
        self.assertEqual(
            normalize_base_url("http://127.0.0.1:4000/"),
            "http://127.0.0.1:4000",
        )
        with self.assertRaises(ValidationError):
            normalize_base_url("http://gateway.invalid")
        with self.assertRaises(ValidationError):
            normalize_base_url("https://user:password@gateway.invalid")

    def test_missing_scoped_key_never_falls_back_to_master_key(self) -> None:
        with self.assertRaises(AuthenticationBlocked):
            load_runtime_config(
                {
                    "LLM_COUNCIL_BASE_URL": "https://gateway.invalid",
                    "LITELLM_MASTER_KEY": "unit-test-master-placeholder",
                },
                require_key=True,
            )

    def test_config_rejects_secret_fields_and_unsupported_pricing(self) -> None:
        cases = [
            {"schema_version": 1, "api_key": "unit-test-placeholder"},
            {"schema_version": 1, "policy": {"defaults": {"token": "placeholder"}}},
            {"schema_version": 1, "pricing": {"fixture-model": 0.1}},
        ]
        with tempfile.TemporaryDirectory(prefix="llm-council-config-") as temporary:
            path = Path(temporary) / "policy.json"
            for payload in cases:
                with self.subTest(payload=payload):
                    path.write_text(json.dumps(payload), encoding="utf-8")
                    with self.assertRaises(ValidationError):
                        load_runtime_config(
                            {
                                "LLM_COUNCIL_BASE_URL": "https://gateway.invalid",
                                "LLM_COUNCIL_API_KEY": "unit-test-placeholder",
                                "LLM_COUNCIL_CONFIG": str(path),
                            },
                            require_key=True,
                        )

    def test_deployment_defaults_fill_only_omitted_manifest_fields(self) -> None:
        manifest = valid_manifest()
        manifest["execution"].pop("max_concurrency")
        manifest["policy"].pop("expose_model_identities_after_synthesis")
        runtime = RuntimeConfig(
            base_url="https://gateway.invalid",
            api_key="unit-test-placeholder",
            policy={
                "defaults": {
                    "execution": {"max_concurrency": 2},
                    "policy": {"expose_model_identities_after_synthesis": False},
                }
            },
        )
        effective = apply_runtime_policy(manifest, runtime)
        normalized = validate_manifest(effective, command="preflight")
        self.assertEqual(normalized["execution"]["max_concurrency"], 2)
        self.assertFalse(
            normalized["policy"]["expose_model_identities_after_synthesis"]
        )

        explicit = valid_manifest()
        effective = apply_runtime_policy(explicit, runtime)
        normalized = validate_manifest(effective, command="preflight")
        self.assertEqual(normalized["execution"]["max_concurrency"], 3)
        self.assertTrue(
            normalized["policy"]["expose_model_identities_after_synthesis"]
        )

    def test_deployment_maximums_reject_instead_of_silently_capping(self) -> None:
        manifest = valid_manifest()
        runtime = RuntimeConfig(
            base_url="https://gateway.invalid",
            api_key="unit-test-placeholder",
            policy={
                "maximums": {
                    "max_concurrency": 2,
                    "max_inference_requests": 8,
                    "max_output_tokens_per_call": 700,
                }
            },
        )
        with self.assertRaises(ValidationError):
            apply_runtime_policy(manifest, runtime)
        self.assertEqual(manifest["execution"]["max_concurrency"], 3)
        self.assertEqual(manifest["budget"]["max_inference_requests"], 10)

    def test_deployment_max_usd_fills_null_and_rejects_larger_request(self) -> None:
        runtime = RuntimeConfig(
            base_url="https://gateway.invalid",
            api_key="unit-test-placeholder",
            policy={"maximums": {"max_usd": 1.25}},
        )
        manifest = valid_manifest()
        effective = apply_runtime_policy(manifest, runtime)
        self.assertEqual(effective["budget"]["max_usd"], 1.25)

        manifest["budget"]["max_usd"] = 2.0
        with self.assertRaises(ValidationError):
            apply_runtime_policy(manifest, runtime)


if __name__ == "__main__":
    unittest.main()
