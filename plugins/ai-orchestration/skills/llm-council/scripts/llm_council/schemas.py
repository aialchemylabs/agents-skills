"""Strict versioned validators for manifests and model/result objects."""

from __future__ import annotations

import copy
import json
import re
from pathlib import Path
from typing import Any, Callable, Iterable, Mapping

from .config import DEFAULTS, HARD_CAPS
from .errors import ValidationError

MODES = {"decision", "architecture", "code-review", "red-team", "research", "benchmark", "auto"}
RUN_STATUSES = {"complete", "partial", "failed", "cancelled", "blocked"}
STAGE_STATUSES = {"pending", "running", "complete", "partial", "failed", "skipped", "cancelled"}
COST_STATUSES = {"reported", "estimated", "unknown"}
PERSISTENCE_MODES = {"none", "summary", "full"}
CRITIQUE_CRITERIA = (
    "evidence_quality",
    "correctness",
    "feasibility",
    "risk_awareness",
    "instruction_fit",
)
DEFAULT_RUBRIC_WEIGHTS = {
    "evidence_quality": 0.30,
    "correctness": 0.25,
    "feasibility": 0.15,
    "risk_awareness": 0.15,
    "instruction_fit": 0.15,
}

_SECRET_VALUE_PATTERNS = (
    re.compile(r"-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----", re.I),
    re.compile(r"\b(?:sk|ghp|gho|github_pat|xox[baprs])[-_][A-Za-z0-9_-]{12,}\b"),
    re.compile(r"\bAuthorization\s*:\s*Bearer\s+\S+", re.I),
    re.compile(
        r"\b(?:api[_ -]?key|access[_ -]?token|client[_ -]?secret|password)\s*[:=]\s*"
        r"['\"]?[A-Za-z0-9_./+=-]{12,}",
        re.I,
    ),
)
_FORBIDDEN_MANIFEST_KEYS = {
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


def _error(path: str, message: str) -> ValidationError:
    return ValidationError(f"{path}: {message}")


def _object(value: Any, path: str) -> dict[str, Any]:
    if not isinstance(value, dict):
        raise _error(path, "must be an object")
    return value


def _array(value: Any, path: str) -> list[Any]:
    if not isinstance(value, list):
        raise _error(path, "must be an array")
    return value


def _allowed(obj: Mapping[str, Any], allowed: set[str], path: str) -> None:
    unknown = set(obj) - allowed
    if unknown:
        raise _error(path, "unknown field(s): " + ", ".join(sorted(unknown)))


def _string(value: Any, path: str, *, allow_empty: bool = False) -> str:
    if not isinstance(value, str) or (not allow_empty and not value.strip()):
        raise _error(path, "must be a non-empty string")
    return value


def _number(
    value: Any,
    path: str,
    *,
    minimum: float | None = None,
    maximum: float | None = None,
) -> float:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise _error(path, "must be a number")
    if minimum is not None and value < minimum:
        raise _error(path, f"must be at least {minimum}")
    if maximum is not None and value > maximum:
        raise _error(path, f"must be at most {maximum}")
    return float(value)


def _integer(
    value: Any,
    path: str,
    *,
    minimum: int | None = None,
    maximum: int | None = None,
) -> int:
    if isinstance(value, bool) or not isinstance(value, int):
        raise _error(path, "must be an integer")
    if minimum is not None and value < minimum:
        raise _error(path, f"must be at least {minimum}")
    if maximum is not None and value > maximum:
        raise _error(path, f"must be at most {maximum}")
    return value


def _string_list(value: Any, path: str, *, allow_empty: bool = True) -> list[str]:
    items = _array(value, path)
    if not allow_empty and not items:
        raise _error(path, "must not be empty")
    for index, item in enumerate(items):
        _string(item, f"{path}[{index}]")
    return items


def _deep_merge(default: Mapping[str, Any], supplied: Mapping[str, Any]) -> dict[str, Any]:
    result = copy.deepcopy(dict(default))
    for key, value in supplied.items():
        if isinstance(value, dict) and isinstance(result.get(key), dict):
            result[key] = _deep_merge(result[key], value)
        else:
            result[key] = copy.deepcopy(value)
    return result


def _scan_for_secrets(value: Any, path: str = "$") -> None:
    if isinstance(value, Mapping):
        for key, child in value.items():
            normalized = str(key).lower().replace("-", "_")
            if normalized in _FORBIDDEN_MANIFEST_KEYS:
                raise _error(path, f"forbidden credential field '{key}'")
            _scan_for_secrets(child, f"{path}.{key}")
    elif isinstance(value, list):
        for index, child in enumerate(value):
            _scan_for_secrets(child, f"{path}[{index}]")
    elif isinstance(value, str):
        if any(pattern.search(value) for pattern in _SECRET_VALUE_PATTERNS):
            raise _error(path, "possible credential detected; sanitize before egress")


def _validate_output_directory(value: Any, path: str) -> str:
    text = _string(value, path)
    resolved = Path(text).expanduser().resolve()
    home = Path.home().resolve()
    cwd = Path.cwd().resolve()
    unsafe = {Path("/"), home, cwd}
    if resolved in unsafe or len(resolved.parts) < 3:
        raise _error(path, "must not be a filesystem root, home directory, or repository root")
    return str(resolved)


def _planned_stages(manifest: Mapping[str, Any]) -> dict[str, bool]:
    policy = manifest["policy"]
    mode = manifest["mode"]
    critique = bool(policy["run_cross_critique"])
    addy_setting = policy["run_addy_revision"]
    has_proposer = any(member.get("role") == "proposer" for member in manifest["roster"]["members"])
    revision = addy_setting is True or (
        addy_setting == "auto" and (mode == "code-review" or has_proposer)
    )
    if mode == "benchmark" and not critique:
        revision = False
    chairman = bool(manifest["roster"].get("chairman")) and not (
        mode == "benchmark" and not critique
    )
    return {
        "independent": True,
        "cross_critique": critique,
        "revision": revision,
        "chairman": chairman,
    }


def planned_stages(manifest: Mapping[str, Any]) -> dict[str, bool]:
    return _planned_stages(manifest)


def planned_base_requests(manifest: Mapping[str, Any]) -> int:
    stages = _planned_stages(manifest)
    count = len(manifest["roster"]["members"])
    if stages["cross_critique"]:
        count += len(manifest["roster"]["members"])
    if stages["revision"]:
        count += 1
    if stages["chairman"]:
        count += 1
    return count


def planned_output_tokens(manifest: Mapping[str, Any]) -> int:
    stages = _planned_stages(manifest)
    caps = manifest["execution"]["stage_output_tokens"]
    members = len(manifest["roster"]["members"])
    total = members * caps["independent"]
    if stages["cross_critique"]:
        total += members * caps["critique"]
    if stages["revision"]:
        total += caps["revision"]
    if stages["chairman"]:
        total += caps["chairman"]
    return total


def validate_manifest(manifest: Any, command: str = "validate") -> dict[str, Any]:
    """Return a normalized manifest or raise a precise, non-secret error."""
    if command not in {"validate", "preflight", "run"}:
        raise ValidationError("command must be validate, preflight, or run")
    obj = copy.deepcopy(_object(manifest, "$"))
    _allowed(
        obj,
        {
            "schema_version",
            "request_id",
            "mode",
            "question",
            "decision_criteria",
            "context",
            "roster",
            "execution",
            "budget",
            "approval",
            "policy",
            "data_egress",
        },
        "$",
    )
    _scan_for_secrets(obj)
    if isinstance(obj.get("schema_version"), bool) or obj.get("schema_version") != 1:
        raise _error("$.schema_version", "must equal 1")
    _string(obj.get("request_id"), "$.request_id")
    mode = _string(obj.get("mode"), "$.mode")
    if mode not in MODES:
        raise _error("$.mode", f"must be one of {', '.join(sorted(MODES))}")
    _string(obj.get("question"), "$.question")
    _string_list(
        obj.get("decision_criteria"),
        "$.decision_criteria",
        allow_empty=False,
    )

    context = _object(obj.get("context", {}), "$.context")
    _allowed(context, {"instructions", "evidence", "artifact_refs"}, "$.context")
    context.setdefault("instructions", "Read-only analysis. Do not invent unavailable evidence.")
    _string(context["instructions"], "$.context.instructions", allow_empty=True)
    evidence = _array(context.setdefault("evidence", []), "$.context.evidence")
    evidence_ids: set[str] = set()
    for index, item in enumerate(evidence):
        item = _object(item, f"$.context.evidence[{index}]")
        _allowed(item, {"id", "source", "content"}, f"$.context.evidence[{index}]")
        evidence_id = _string(item.get("id"), f"$.context.evidence[{index}].id")
        if evidence_id in evidence_ids:
            raise _error(f"$.context.evidence[{index}].id", "must be unique")
        evidence_ids.add(evidence_id)
        _string(item.get("source"), f"$.context.evidence[{index}].source")
        _string(item.get("content"), f"$.context.evidence[{index}].content", allow_empty=True)
    _string_list(context.setdefault("artifact_refs", []), "$.context.artifact_refs")
    obj["context"] = context

    roster = _object(obj.get("roster"), "$.roster")
    _allowed(roster, {"members", "chairman", "fallback_chairman"}, "$.roster")
    members = _array(roster.get("members"), "$.roster.members")
    if not 2 <= len(members) <= HARD_CAPS["members"]:
        raise _error("$.roster.members", "must contain two to five members")
    member_ids: set[str] = set()
    models: set[str] = set()
    families: list[str] = []
    for index, member in enumerate(members):
        path = f"$.roster.members[{index}]"
        member = _object(member, path)
        _allowed(member, {"member_id", "model", "family", "angle", "role"}, path)
        member_id = _string(member.get("member_id"), f"{path}.member_id")
        if member_id in member_ids:
            raise _error(f"{path}.member_id", "must be unique")
        member_ids.add(member_id)
        model = _string(member.get("model"), f"{path}.model")
        family = _string(member.get("family", "unverified"), f"{path}.family").lower()
        member["family"] = family
        families.append(family)
        _string(member.get("angle"), f"{path}.angle")
        role = member.get("role", "member")
        if role not in {"member", "proposer", "reviewer"}:
            raise _error(f"{path}.role", "must be member, proposer, or reviewer")
        member["role"] = role
        if model in models:
            raise _error(f"{path}.model", "must identify a distinct model; repeated roles are not a multi-model council")
        models.add(model)
    if sum(member.get("role") == "proposer" for member in members) > 1:
        raise _error("$.roster.members", "at most one member may be the proposer")
    roster.setdefault("chairman", None)
    chairman = roster["chairman"]
    if chairman is not None:
        chairman = _object(chairman, "$.roster.chairman")
        _allowed(chairman, {"model", "family"}, "$.roster.chairman")
        _string(chairman.get("model"), "$.roster.chairman.model")
        chairman["family"] = _string(
            chairman.get("family", "unverified"), "$.roster.chairman.family"
        ).lower()
    fallback = roster.get("fallback_chairman")
    if fallback is not None:
        if chairman is None:
            raise _error("$.roster.fallback_chairman", "requires a primary chairman")
        fallback = _object(fallback, "$.roster.fallback_chairman")
        _allowed(fallback, {"model", "family"}, "$.roster.fallback_chairman")
        _string(fallback.get("model"), "$.roster.fallback_chairman.model")
        fallback["family"] = _string(
            fallback.get("family", "unverified"), "$.roster.fallback_chairman.family"
        ).lower()
    roster.setdefault("fallback_chairman", None)
    obj["roster"] = roster

    supplied_execution = _object(obj.get("execution", {}), "$.execution")
    execution = _deep_merge(DEFAULTS["execution"], supplied_execution)
    if "max_concurrency" not in supplied_execution:
        execution["max_concurrency"] = min(
            DEFAULTS["execution"]["max_concurrency"], len(members)
        )
    _allowed(
        execution,
        {
            "round_number",
            "max_concurrency",
            "request_timeout_seconds",
            "connect_timeout_seconds",
            "max_transient_retries",
            "temperature",
            "stage_output_tokens",
        },
        "$.execution",
    )
    _integer(execution["round_number"], "$.execution.round_number", minimum=1, maximum=3)
    _integer(
        execution["max_concurrency"],
        "$.execution.max_concurrency",
        minimum=1,
        maximum=min(HARD_CAPS["max_concurrency"], len(members)),
    )
    _integer(
        execution["request_timeout_seconds"],
        "$.execution.request_timeout_seconds",
        minimum=1,
        maximum=HARD_CAPS["request_timeout_seconds"],
    )
    _integer(
        execution["connect_timeout_seconds"],
        "$.execution.connect_timeout_seconds",
        minimum=1,
        maximum=HARD_CAPS["connect_timeout_seconds"],
    )
    if execution["connect_timeout_seconds"] > execution["request_timeout_seconds"]:
        raise _error("$.execution.connect_timeout_seconds", "must not exceed request timeout")
    _integer(
        execution["max_transient_retries"],
        "$.execution.max_transient_retries",
        minimum=0,
        maximum=HARD_CAPS["max_transient_retries"],
    )
    _number(execution["temperature"], "$.execution.temperature", minimum=0, maximum=2)
    token_caps = _object(execution["stage_output_tokens"], "$.execution.stage_output_tokens")
    _allowed(token_caps, {"independent", "critique", "revision", "chairman"}, "$.execution.stage_output_tokens")
    for stage in ("independent", "critique", "revision", "chairman"):
        _integer(
            token_caps.get(stage),
            f"$.execution.stage_output_tokens.{stage}",
            minimum=1,
            maximum=HARD_CAPS["max_output_tokens_per_call"],
        )
    obj["execution"] = execution

    budget = _deep_merge(DEFAULTS["budget"], _object(obj.get("budget", {}), "$.budget"))
    _allowed(
        budget,
        {
            "max_inference_requests",
            "max_input_characters",
            "max_output_tokens",
            "max_usd",
            "require_known_cost",
        },
        "$.budget",
    )
    _integer(
        budget["max_inference_requests"],
        "$.budget.max_inference_requests",
        minimum=1,
        maximum=HARD_CAPS["max_inference_requests"],
    )
    _integer(
        budget["max_input_characters"],
        "$.budget.max_input_characters",
        minimum=1,
        maximum=HARD_CAPS["max_input_characters"],
    )
    _integer(
        budget["max_output_tokens"],
        "$.budget.max_output_tokens",
        minimum=1,
        maximum=HARD_CAPS["max_output_tokens"],
    )
    if budget["max_usd"] is not None:
        _number(budget["max_usd"], "$.budget.max_usd", minimum=0)
    if not isinstance(budget["require_known_cost"], bool):
        raise _error("$.budget.require_known_cost", "must be boolean")
    obj["budget"] = budget

    policy = _deep_merge(DEFAULTS["policy"], _object(obj.get("policy", {}), "$.policy"))
    _allowed(
        policy,
        {
            "minimum_successful_members",
            "require_distinct_families",
            "allow_perspective_only_fallback",
            "run_cross_critique",
            "score_candidates",
            "anonymize_candidates",
            "run_addy_revision",
            "expose_model_identities_after_synthesis",
            "persist",
            "output_directory",
            "rubric_weights",
        },
        "$.policy",
    )
    _integer(
        policy["minimum_successful_members"],
        "$.policy.minimum_successful_members",
        minimum=2,
        maximum=len(members),
    )
    for key in (
        "require_distinct_families",
        "allow_perspective_only_fallback",
        "run_cross_critique",
        "score_candidates",
        "anonymize_candidates",
        "expose_model_identities_after_synthesis",
    ):
        if not isinstance(policy[key], bool):
            raise _error(f"$.policy.{key}", "must be boolean")
    if mode != "benchmark" and not policy["run_cross_critique"]:
        raise _error(
            "$.policy.run_cross_critique",
            "may be false only in benchmark mode",
        )
    if not (
        isinstance(policy["run_addy_revision"], bool)
        or policy["run_addy_revision"] == "auto"
    ):
        raise _error("$.policy.run_addy_revision", "must be true, false, or 'auto'")
    if policy["persist"] not in PERSISTENCE_MODES:
        raise _error("$.policy.persist", "must be none, summary, or full")
    if policy["persist"] != "none":
        policy["output_directory"] = _validate_output_directory(
            policy.get("output_directory"), "$.policy.output_directory"
        )
    elif policy.get("output_directory") is not None:
        raise _error("$.policy.output_directory", "must be null when persistence is none")
    weights = policy.get("rubric_weights", DEFAULT_RUBRIC_WEIGHTS)
    weights = _object(weights, "$.policy.rubric_weights")
    _allowed(weights, set(CRITIQUE_CRITERIA), "$.policy.rubric_weights")
    for criterion in CRITIQUE_CRITERIA:
        _number(weights.get(criterion), f"$.policy.rubric_weights.{criterion}", minimum=0)
    if abs(sum(float(weights[item]) for item in CRITIQUE_CRITERIA) - 1.0) > 1e-9:
        raise _error("$.policy.rubric_weights", "weights must sum to 1.0")
    policy["rubric_weights"] = {key: float(weights[key]) for key in CRITIQUE_CRITERIA}
    obj["policy"] = policy

    if policy["require_distinct_families"]:
        verified = all(family not in {"", "unknown", "unverified"} for family in families)
        if not verified:
            raise _error("$.roster.members", "verified family metadata is required")
        if len(set(families)) < 2:
            raise _error("$.roster.members", "at least two distinct families are required")
    # Distinct accessible models may share a provider. Label diversity honestly;
    # provider diversity is a preference unless explicitly required by the caller.

    data_egress = _object(obj.get("data_egress", {}), "$.data_egress")
    _allowed(
        data_egress,
        {
            "classification",
            "contains_personal_data",
            "contains_proprietary_data",
            "contains_regulated_data",
            "authorization_confirmed",
        },
        "$.data_egress",
    )
    data_egress.setdefault("classification", "general")
    if data_egress["classification"] not in {
        "general",
        "public",
        "proprietary",
        "personal",
        "regulated",
    }:
        raise _error("$.data_egress.classification", "unsupported classification")
    for key in (
        "contains_personal_data",
        "contains_proprietary_data",
        "contains_regulated_data",
        "authorization_confirmed",
    ):
        data_egress.setdefault(key, False)
        if not isinstance(data_egress[key], bool):
            raise _error(f"$.data_egress.{key}", "must be boolean")
    if (
        data_egress["contains_personal_data"]
        or data_egress["contains_proprietary_data"]
        or data_egress["contains_regulated_data"]
    ) and not data_egress["authorization_confirmed"]:
        raise _error("$.data_egress", "sensitive data egress requires explicit authorization")
    obj["data_egress"] = data_egress

    input_characters = len(
        json.dumps(
            {"question": obj["question"], "criteria": obj["decision_criteria"], "context": context},
            ensure_ascii=False,
            separators=(",", ":"),
        )
    )
    if input_characters > budget["max_input_characters"]:
        raise _error("$.context", "question and evidence exceed max_input_characters")
    if planned_base_requests(obj) > budget["max_inference_requests"]:
        raise _error("$.budget.max_inference_requests", "is below the planned base request count")
    if planned_output_tokens(obj) > budget["max_output_tokens"]:
        raise _error("$.budget.max_output_tokens", "is below planned requested output tokens")

    approval = _object(obj.get("approval", {}), "$.approval")
    _allowed(approval, {"inference_approved", "plan_fingerprint", "authority"}, "$.approval")
    approval.setdefault("authority", "interactive")
    if approval["authority"] not in {"interactive", "autonomous"}:
        raise _error("$.approval.authority", "must be interactive or autonomous")
    approval.setdefault("inference_approved", False)
    approval.setdefault("plan_fingerprint", None)
    if not isinstance(approval["inference_approved"], bool):
        raise _error("$.approval.inference_approved", "must be boolean")
    if approval["plan_fingerprint"] is not None:
        fingerprint = _string(approval["plan_fingerprint"], "$.approval.plan_fingerprint")
        if not re.fullmatch(r"sha256:[0-9a-f]{64}", fingerprint):
            raise _error("$.approval.plan_fingerprint", "must be a sha256 fingerprint")
    if command == "run":
        if not approval["inference_approved"] and approval["authority"] != "autonomous":
            raise _error("$.approval.inference_approved", "must be true for run")
        if not approval["plan_fingerprint"] and approval["authority"] != "autonomous":
            raise _error("$.approval.plan_fingerprint", "is required for run")
    obj["approval"] = approval
    return obj


def _validate_string_array_field(obj: Mapping[str, Any], field: str, path: str) -> None:
    _string_list(obj.get(field), f"{path}.{field}")


def validate_member_response(
    value: Any, evidence_ids: Iterable[str] | None = None
) -> dict[str, Any]:
    obj = copy.deepcopy(_object(value, "member"))
    _allowed(
        obj,
        {"recommendation", "rationale", "evidence", "assumptions", "risks", "unknowns", "confidence"},
        "member",
    )
    _string(obj.get("recommendation"), "member.recommendation")
    for field in ("rationale", "assumptions", "unknowns"):
        _validate_string_array_field(obj, field, "member")
    evidence = _array(obj.get("evidence"), "member.evidence")
    allowed_evidence = set(evidence_ids) if evidence_ids is not None else None
    for index, item in enumerate(evidence):
        item = _object(item, f"member.evidence[{index}]")
        _allowed(item, {"evidence_id", "claim"}, f"member.evidence[{index}]")
        evidence_id = _string(
            item.get("evidence_id"), f"member.evidence[{index}].evidence_id"
        )
        if allowed_evidence is not None and evidence_id not in allowed_evidence:
            raise _error(
                f"member.evidence[{index}].evidence_id",
                "does not refer to supplied evidence",
            )
        _string(item.get("claim"), f"member.evidence[{index}].claim")
    risks = _array(obj.get("risks"), "member.risks")
    for index, item in enumerate(risks):
        item = _object(item, f"member.risks[{index}]")
        _allowed(item, {"risk", "severity", "mitigation"}, f"member.risks[{index}]")
        _string(item.get("risk"), f"member.risks[{index}].risk")
        if item.get("severity") not in {"low", "medium", "high", "critical"}:
            raise _error(f"member.risks[{index}].severity", "has an invalid severity")
        _string(item.get("mitigation"), f"member.risks[{index}].mitigation")
    _number(obj.get("confidence"), "member.confidence", minimum=0, maximum=1)
    return obj


def validate_critique(value: Any, labels: Iterable[str], *, partial: bool = False) -> dict[str, Any]:
    expected = list(labels)
    expected_set = set(expected)
    obj = copy.deepcopy(_object(value, "critique"))
    _allowed(obj, {"evaluations", "ranking", "revised_position", "confidence", "partial"}, "critique")
    evaluations = _array(obj.get("evaluations"), "critique.evaluations")
    seen: set[str] = set()
    for index, evaluation in enumerate(evaluations):
        path = f"critique.evaluations[{index}]"
        evaluation = _object(evaluation, path)
        _allowed(
            evaluation,
            {"candidate_label", "strengths", "weaknesses", "unsupported_claims", "scores"},
            path,
        )
        label = _string(evaluation.get("candidate_label"), f"{path}.candidate_label")
        if label not in expected_set or label in seen:
            raise _error(f"{path}.candidate_label", "must be a unique known candidate label")
        seen.add(label)
        for field in ("strengths", "weaknesses", "unsupported_claims"):
            _string_list(evaluation.get(field), f"{path}.{field}")
        if "scores" in evaluation:
            scores = _object(evaluation["scores"], f"{path}.scores")
            _allowed(scores, set(CRITIQUE_CRITERIA), f"{path}.scores")
            for criterion in CRITIQUE_CRITERIA:
                _integer(scores.get(criterion), f"{path}.scores.{criterion}", minimum=1, maximum=5)
    is_partial = bool(obj.get("partial", partial))
    if not is_partial and seen != expected_set:
        raise _error("critique.evaluations", "must evaluate every candidate")
    ranking = _string_list(obj.get("ranking", []), "critique.ranking")
    if len(set(ranking)) != len(ranking):
        raise _error("critique.ranking", "contains duplicate labels")
    if not set(ranking).issubset(expected_set):
        raise _error("critique.ranking", "contains an unknown label")
    if "ranking" in obj and not is_partial and set(ranking) != expected_set:
        raise _error("critique.ranking", "must rank every candidate")
    _string(obj.get("revised_position"), "critique.revised_position")
    _number(obj.get("confidence"), "critique.confidence", minimum=0, maximum=1)
    obj["partial"] = is_partial
    return obj


def validate_critique_response(
    value: Any,
    candidate_labels: Iterable[str],
    allow_partial: bool = False,
) -> dict[str, Any]:
    return validate_critique(value, candidate_labels, partial=allow_partial)


def validate_revision(value: Any) -> dict[str, Any]:
    obj = copy.deepcopy(_object(value, "revision"))
    _allowed(
        obj,
        {"accepted_feedback", "rejected_feedback", "revised_proposal", "remaining_risks", "confidence"},
        "revision",
    )
    for field, keys in (
        ("accepted_feedback", {"feedback", "response"}),
        ("rejected_feedback", {"feedback", "reason"}),
    ):
        rows = _array(obj.get(field), f"revision.{field}")
        for index, row in enumerate(rows):
            row = _object(row, f"revision.{field}[{index}]")
            _allowed(row, keys, f"revision.{field}[{index}]")
            for key in keys:
                _string(row.get(key), f"revision.{field}[{index}].{key}")
    _string(obj.get("revised_proposal"), "revision.revised_proposal")
    _string_list(obj.get("remaining_risks"), "revision.remaining_risks")
    _number(obj.get("confidence"), "revision.confidence", minimum=0, maximum=1)
    return obj


def validate_revision_response(value: Any) -> dict[str, Any]:
    return validate_revision(value)


def validate_chairman(value: Any) -> dict[str, Any]:
    obj = copy.deepcopy(_object(value, "chairman"))
    _allowed(
        obj,
        {
            "recommendation",
            "why",
            "consensus",
            "disagreements",
            "minority_report",
            "risks",
            "next_steps",
            "confidence",
            "material_unknowns",
        },
        "chairman",
    )
    _string(obj.get("recommendation"), "chairman.recommendation")
    for field in ("why", "consensus", "minority_report", "next_steps", "material_unknowns"):
        _string_list(obj.get(field), f"chairman.{field}")
    disagreements = _array(obj.get("disagreements"), "chairman.disagreements")
    for index, item in enumerate(disagreements):
        path = f"chairman.disagreements[{index}]"
        item = _object(item, path)
        _allowed(item, {"issue", "positions", "resolution"}, path)
        _string(item.get("issue"), f"{path}.issue")
        _string_list(item.get("positions"), f"{path}.positions")
        _string(item.get("resolution"), f"{path}.resolution")
    risks = _array(obj.get("risks"), "chairman.risks")
    for index, item in enumerate(risks):
        path = f"chairman.risks[{index}]"
        item = _object(item, path)
        _allowed(item, {"risk", "mitigation", "owner"}, path)
        _string(item.get("risk"), f"{path}.risk")
        _string(item.get("mitigation"), f"{path}.mitigation")
        _string(item.get("owner"), f"{path}.owner", allow_empty=True)
    _number(obj.get("confidence"), "chairman.confidence", minimum=0, maximum=1)
    return obj


def validate_chairman_response(value: Any) -> dict[str, Any]:
    return validate_chairman(value)


def validate_result(value: Any) -> dict[str, Any]:
    """Validate the stable top-level result contract used by render/persistence."""
    obj = copy.deepcopy(_object(value, "$"))
    required = {
        "schema_version",
        "request_id",
        "run_id",
        "status",
        "mode",
        "diversity",
        "stages",
        "candidates",
        "critiques",
        "ranking",
        "chairman",
        "benchmark",
        "usage",
        "failures",
        "persistence",
    }
    _allowed(obj, required | {"revision", "transitions", "identity_mapping", "execution_report"}, "$")
    if "execution_report" in obj:
        report = _object(obj["execution_report"], "$.execution_report")
        _allowed(report, {"backend", "roster", "round_number", "maximum_rounds", "stop_reason", "plan_fingerprint"}, "$.execution_report")
        if report.get("backend") != "litellm":
            raise _error("$.execution_report.backend", "must be litellm for gateway results")
        _integer(report.get("round_number"), "$.execution_report.round_number", minimum=1, maximum=3)
        if report.get("maximum_rounds") != 3:
            raise _error("$.execution_report.maximum_rounds", "must equal 3")
        _string(report.get("stop_reason"), "$.execution_report.stop_reason")
        _string(report.get("plan_fingerprint"), "$.execution_report.plan_fingerprint")
        for index, member in enumerate(_array(report.get("roster"), "$.execution_report.roster")):
            path = f"$.execution_report.roster[{index}]"
            _allowed(_object(member, path), {"member_id", "model", "family", "angle", "role"}, path)
            for field in ("member_id", "model", "family", "angle", "role"):
                _string(member.get(field), f"{path}.{field}")
    missing = required - set(obj)
    if missing:
        raise _error("$", "missing field(s): " + ", ".join(sorted(missing)))
    if isinstance(obj["schema_version"], bool) or obj["schema_version"] != 1:
        raise _error("$.schema_version", "must equal 1")
    _string(obj["request_id"], "$.request_id")
    _string(obj["run_id"], "$.run_id")
    if obj["status"] not in RUN_STATUSES:
        raise _error("$.status", "invalid run status")
    if obj["mode"] not in MODES:
        raise _error("$.mode", "invalid mode")
    diversity = _object(obj["diversity"], "$.diversity")
    _allowed(diversity, {"status", "families", "limitations"}, "$.diversity")
    if diversity.get("status") not in {"verified", "unverified", "perspective-only"}:
        raise _error("$.diversity.status", "invalid diversity status")
    _string_list(diversity.get("families"), "$.diversity.families")
    _string_list(diversity.get("limitations"), "$.diversity.limitations")
    stages = _object(obj["stages"], "$.stages")
    required_stages = {
        "preflight",
        "independent",
        "cross_critique",
        "revision",
        "chairman",
    }
    missing_stages = required_stages - set(stages)
    if missing_stages:
        raise _error(
            "$.stages", "missing stage(s): " + ", ".join(sorted(missing_stages))
        )
    for stage, stage_value in stages.items():
        stage_obj = _object(stage_value, f"$.stages.{stage}")
        _allowed(
            stage_obj,
            {"status", "reason", "call_ids"},
            f"$.stages.{stage}",
        )
        if stage_obj.get("status") not in STAGE_STATUSES:
            raise _error(f"$.stages.{stage}.status", "invalid stage status")
        if "reason" in stage_obj and stage_obj["reason"] is not None:
            _string(stage_obj["reason"], f"$.stages.{stage}.reason")
        if "call_ids" in stage_obj:
            _string_list(stage_obj["call_ids"], f"$.stages.{stage}.call_ids")

    candidates = _array(obj["candidates"], "$.candidates")
    labels: list[str] = []
    for index, candidate in enumerate(candidates):
        path = f"$.candidates[{index}]"
        candidate = _object(candidate, path)
        _allowed(
            candidate,
            {"candidate_label", "member_id", "model", "family", "response"},
            path,
        )
        label = _string(candidate.get("candidate_label"), f"{path}.candidate_label")
        if label in labels:
            raise _error(f"{path}.candidate_label", "must be unique")
        labels.append(label)
        for field in ("member_id", "model", "family"):
            if field in candidate and candidate[field] is not None:
                _string(candidate[field], f"{path}.{field}")
        validate_member_response(candidate.get("response"))

    critiques = _array(obj["critiques"], "$.critiques")
    for index, critique in enumerate(critiques):
        path = f"$.critiques[{index}]"
        critique = _object(critique, path)
        _allowed(
            critique,
            {
                "critic_member_id",
                "reviewer_member_id",
                "critic_candidate_label",
                "response",
            },
            path,
        )
        for field in ("critic_member_id", "reviewer_member_id", "critic_candidate_label"):
            if field in critique and critique[field] is not None:
                _string(critique[field], f"{path}.{field}")
        validate_critique_response(
            critique.get("response"),
            labels,
            allow_partial=bool(critique.get("response", {}).get("partial", False))
            if isinstance(critique.get("response"), dict)
            else False,
        )

    ranking = _object(obj["ranking"], "$.ranking")
    _allowed(
        ranking,
        {
            "advisory",
            "subjective",
            "critic_count",
            "expected_critic_count",
            "coverage",
            "weights",
            "candidates",
            "order",
            "limitations",
        },
        "$.ranking",
    )
    if "advisory" in ranking and not isinstance(ranking["advisory"], bool):
        raise _error("$.ranking.advisory", "must be boolean")
    if "subjective" in ranking and not isinstance(ranking["subjective"], bool):
        raise _error("$.ranking.subjective", "must be boolean")
    for field in ("critic_count", "expected_critic_count"):
        if field in ranking:
            _integer(ranking[field], f"$.ranking.{field}", minimum=0)
    if "coverage" in ranking:
        coverage = ranking["coverage"]
        if isinstance(coverage, dict):
            _allowed(
                coverage,
                {"valid_critics", "expected_critics"},
                "$.ranking.coverage",
            )
            for field in ("valid_critics", "expected_critics"):
                _integer(
                    coverage.get(field),
                    f"$.ranking.coverage.{field}",
                    minimum=0,
                )
        else:
            _number(coverage, "$.ranking.coverage", minimum=0, maximum=1)
    if "weights" in ranking:
        weights = _object(ranking["weights"], "$.ranking.weights")
        _allowed(weights, set(CRITIQUE_CRITERIA), "$.ranking.weights")
        for criterion in CRITIQUE_CRITERIA:
            _number(
                weights.get(criterion),
                f"$.ranking.weights.{criterion}",
                minimum=0,
            )
        if abs(sum(float(weights[item]) for item in CRITIQUE_CRITERIA) - 1.0) > 1e-9:
            raise _error("$.ranking.weights", "must sum to 1.0")
    ranking_candidates = _object(ranking.get("candidates", {}), "$.ranking.candidates")
    if not set(ranking_candidates).issubset(set(labels)):
        raise _error("$.ranking.candidates", "contains an unknown candidate label")
    for label, metrics in ranking_candidates.items():
        path = f"$.ranking.candidates.{label}"
        metrics = _object(metrics, path)
        _allowed(
            metrics,
            {
                "weighted_score",
                "rank_score",
                "evaluation_coverage",
                "ranking_coverage",
            },
            path,
        )
        for field in ("weighted_score", "rank_score"):
            if metrics.get(field) is not None:
                _number(metrics[field], f"{path}.{field}", minimum=0, maximum=1)
        for field in ("evaluation_coverage", "ranking_coverage"):
            if field in metrics:
                _integer(metrics[field], f"{path}.{field}", minimum=0)
    if "order" in ranking:
        order = _string_list(ranking["order"], "$.ranking.order")
        if len(order) != len(set(order)) or not set(order).issubset(set(labels)):
            raise _error("$.ranking.order", "must contain unique known labels")
    if "limitations" in ranking:
        _string_list(ranking["limitations"], "$.ranking.limitations")

    chairman = _object(obj["chairman"], "$.chairman")
    if chairman:
        _allowed(
            chairman,
            {
                "status",
                "model",
                "family",
                "used_fallback",
                "request_id",
                "response",
                "host_synthesis_required",
                "host_synthesis_label",
                "participated_as_member",
            },
            "$.chairman",
        )
        if chairman.get("status") not in {"complete", "failed", "skipped"}:
            raise _error("$.chairman.status", "must be complete, failed, or skipped")
        for field in ("model", "family", "request_id", "host_synthesis_label"):
            if field in chairman and chairman[field] is not None:
                _string(chairman[field], f"$.chairman.{field}")
        for field in (
            "used_fallback",
            "host_synthesis_required",
            "participated_as_member",
        ):
            if field in chairman and not isinstance(chairman[field], bool):
                raise _error(f"$.chairman.{field}", "must be boolean")
        if chairman["status"] == "complete":
            validate_chairman_response(chairman.get("response"))
        elif chairman.get("response") is not None:
            raise _error(
                "$.chairman.response", "must be null or omitted unless status is complete"
            )

    benchmark = _object(obj["benchmark"], "$.benchmark")
    _allowed(benchmark, {"calls", "caveats"}, "$.benchmark")
    calls = _array(benchmark.get("calls"), "$.benchmark.calls")
    for index, call in enumerate(calls):
        path = f"$.benchmark.calls[{index}]"
        call = _object(call, path)
        _allowed(
            call,
            {
                "model",
                "family",
                "member_id",
                "candidate_label",
                "stage",
                "status",
                "latency_ms",
                "input_tokens",
                "output_tokens",
                "total_tokens",
                "cost_usd",
                "cost_status",
                "retry_count",
                "request_id",
                "http_status",
                "quality",
            },
            path,
        )
        for field in ("model", "family", "member_id", "candidate_label", "request_id"):
            if field in call and call[field] is not None:
                _string(call[field], f"{path}.{field}")
        _string(call.get("stage"), f"{path}.stage")
        if call.get("status") not in {"success", "failed"}:
            raise _error(f"{path}.status", "must be success or failed")
        _integer(call.get("latency_ms"), f"{path}.latency_ms", minimum=0)
        for field in ("input_tokens", "output_tokens", "total_tokens"):
            if call.get(field) is not None:
                _integer(call[field], f"{path}.{field}", minimum=0)
        if call.get("cost_status") not in COST_STATUSES:
            raise _error(f"{path}.cost_status", "invalid cost status")
        if call["cost_status"] == "unknown":
            if call.get("cost_usd") is not None:
                raise _error(f"{path}.cost_usd", "must be null when cost is unknown")
        else:
            _number(call.get("cost_usd"), f"{path}.cost_usd", minimum=0)
        _integer(call.get("retry_count"), f"{path}.retry_count", minimum=0)
        if "http_status" in call:
            _integer(call["http_status"], f"{path}.http_status", minimum=100, maximum=599)
        quality = _object(call.get("quality"), f"{path}.quality")
        _allowed(
            quality,
            {"peer_score", "rank_score", "chairman_selected", "subjective"},
            f"{path}.quality",
        )
        for field in ("peer_score", "rank_score"):
            if quality.get(field) is not None:
                _number(quality[field], f"{path}.quality.{field}", minimum=0, maximum=1)
        for field in ("chairman_selected", "subjective"):
            if not isinstance(quality.get(field), bool):
                raise _error(f"{path}.quality.{field}", "must be boolean")
    if "caveats" in benchmark:
        _string_list(benchmark["caveats"], "$.benchmark.caveats")

    usage = _object(obj["usage"], "$.usage")
    _allowed(
        usage,
        {
            "inference_requests",
            "requested_output_tokens",
            "input_tokens",
            "output_tokens",
            "total_tokens",
            "cost_usd",
            "reported_cost_usd",
            "cost_status",
            "unknown_cost_calls",
        },
        "$.usage",
    )
    for field in (
        "inference_requests",
        "requested_output_tokens",
        "input_tokens",
        "output_tokens",
        "total_tokens",
        "unknown_cost_calls",
    ):
        if field in usage:
            _integer(usage[field], f"$.usage.{field}", minimum=0)
    if usage.get("cost_status") not in COST_STATUSES:
        raise _error("$.usage.cost_status", "invalid cost status")
    if usage["cost_status"] == "unknown":
        if usage.get("cost_usd") is not None:
            raise _error("$.usage.cost_usd", "must be null when cost is unknown")
    else:
        _number(usage.get("cost_usd"), "$.usage.cost_usd", minimum=0)
    if "reported_cost_usd" in usage:
        _number(usage["reported_cost_usd"], "$.usage.reported_cost_usd", minimum=0)

    failures = _array(obj["failures"], "$.failures")
    for index, failure in enumerate(failures):
        path = f"$.failures[{index}]"
        failure = _object(failure, path)
        _allowed(
            failure,
            {"stage", "subject", "model", "category", "message"},
            path,
        )
        for field in ("stage", "category", "message"):
            _string(failure.get(field), f"{path}.{field}")
        for field in ("subject", "model"):
            if field in failure and failure[field] is not None:
                _string(failure[field], f"{path}.{field}")

    persistence = obj["persistence"]
    persistence = _object(persistence, "$.persistence")
    _allowed(persistence, {"mode", "status", "paths"}, "$.persistence")
    if persistence.get("mode") not in PERSISTENCE_MODES:
        raise _error("$.persistence.mode", "invalid persistence mode")
    if "status" in persistence and persistence["status"] not in {
        "not_requested",
        "pending",
        "complete",
        "failed",
    }:
        raise _error("$.persistence.status", "invalid persistence status")
    _string_list(persistence.get("paths", []), "$.persistence.paths")

    revision = obj.get("revision")
    if revision is not None:
        revision = _object(revision, "$.revision")
        _allowed(revision, {"proposer_member_id", "response"}, "$.revision")
        if revision.get("proposer_member_id") is not None:
            _string(revision["proposer_member_id"], "$.revision.proposer_member_id")
        validate_revision_response(revision.get("response"))
    if "transitions" in obj:
        transitions = _array(obj["transitions"], "$.transitions")
        for index, transition in enumerate(transitions):
            path = f"$.transitions[{index}]"
            transition = _object(transition, path)
            _allowed(
                transition,
                {"stage", "status", "monotonic_ms", "reason", "call_ids"},
                path,
            )
            _string(transition.get("stage"), f"{path}.stage")
            _string(transition.get("status"), f"{path}.status")
            _integer(transition.get("monotonic_ms"), f"{path}.monotonic_ms", minimum=0)
            if transition.get("reason") is not None:
                _string(transition["reason"], f"{path}.reason")
            _string_list(transition.get("call_ids", []), f"{path}.call_ids")
    if "identity_mapping" in obj:
        mapping = _object(obj["identity_mapping"], "$.identity_mapping")
        for member_id, label in mapping.items():
            _string(member_id, "$.identity_mapping key")
            _string(label, f"$.identity_mapping.{member_id}")
            if label not in labels:
                raise _error(
                    f"$.identity_mapping.{member_id}",
                    "must refer to a known candidate label",
                )
    return obj


def parse_model_json(text: str) -> dict[str, Any]:
    """Parse one JSON object, allowing deterministic extraction from surrounding text."""
    if not isinstance(text, str) or not text.strip():
        raise ValidationError("Model returned empty structured content")
    try:
        value = json.loads(text)
        return _object(value, "model_output")
    except json.JSONDecodeError:
        pass
    decoder = json.JSONDecoder()
    objects: list[dict[str, Any]] = []
    for match in re.finditer(r"\{", text):
        try:
            value, _ = decoder.raw_decode(text[match.start() :])
        except json.JSONDecodeError:
            continue
        if isinstance(value, dict):
            objects.append(value)
    unique = {json.dumps(item, sort_keys=True, separators=(",", ":")) for item in objects}
    if len(unique) != 1:
        raise ValidationError("Model output did not contain exactly one JSON object")
    return json.loads(next(iter(unique)))
