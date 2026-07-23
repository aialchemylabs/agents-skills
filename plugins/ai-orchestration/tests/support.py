from __future__ import annotations

import copy
import json
import sys
import threading
from collections import deque
from pathlib import Path
from typing import Any


PLUGIN_ROOT = Path(__file__).resolve().parents[1]
SKILL_ROOT = PLUGIN_ROOT / "skills" / "llm-council"
SCRIPTS_ROOT = SKILL_ROOT / "scripts"
FIXTURE_ROOT = PLUGIN_ROOT / "evals" / "fixtures" / "llm-council"

if str(SCRIPTS_ROOT) not in sys.path:
    sys.path.insert(0, str(SCRIPTS_ROOT))


def load_fixture(name: str) -> dict[str, Any]:
    with (FIXTURE_ROOT / name).open(encoding="utf-8") as handle:
        return json.load(handle)


def valid_manifest() -> dict[str, Any]:
    return copy.deepcopy(load_fixture("architecture-consensus.json")["manifest"])


def member_response(recommendation: str = "Use the bounded hybrid.") -> dict[str, Any]:
    return {
        "recommendation": recommendation,
        "rationale": ["The recommendation follows the supplied evidence."],
        "evidence": [{"evidence_id": "E1", "claim": "E1 supports the decision."}],
        "assumptions": ["The supplied boundary is complete."],
        "risks": [
            {
                "risk": "The assumption may be incomplete.",
                "severity": "medium",
                "mitigation": "Validate it before implementation.",
            }
        ],
        "unknowns": ["Observed workload shape."],
        "confidence": 0.8,
    }


def critique_response(labels: list[str]) -> dict[str, Any]:
    return {
        "evaluations": [
            {
                "candidate_label": label,
                "strengths": ["Uses evidence."],
                "weaknesses": ["One material unknown remains."],
                "unsupported_claims": [],
                "scores": {
                    "evidence_quality": 4,
                    "correctness": 4,
                    "feasibility": 4,
                    "risk_awareness": 4,
                    "instruction_fit": 4,
                },
            }
            for label in labels
        ],
        "ranking": labels,
        "revised_position": "Prefer the highest-evidence candidate.",
        "confidence": 0.8,
    }


def revision_response() -> dict[str, Any]:
    return {
        "accepted_feedback": [
            {"feedback": "Add a rollback gate.", "response": "The proposal now has one."}
        ],
        "rejected_feedback": [
            {
                "feedback": "Replace the datastore.",
                "reason": "No supplied evidence requires replacement.",
            }
        ],
        "revised_proposal": "Migrate in one reversible, measured slice.",
        "remaining_risks": ["Rollback may be slow."],
        "confidence": 0.82,
    }


def chairman_response() -> dict[str, Any]:
    return {
        "recommendation": "Use the bounded hybrid.",
        "why": ["It best matches the supplied synchronous and asynchronous needs."],
        "consensus": ["Separate the two delivery requirements."],
        "disagreements": [
            {
                "issue": "How much infrastructure to add.",
                "positions": ["Candidate A prefers less.", "Candidate B specifies an outbox."],
                "resolution": "Keep the event path narrow.",
            }
        ],
        "minority_report": ["Do not generalize beyond the supplied boundary."],
        "risks": [
            {
                "risk": "Unknown workload shape.",
                "mitigation": "Measure before scaling.",
                "owner": "service owner",
            }
        ],
        "next_steps": ["Run a bounded prototype."],
        "confidence": 0.84,
        "material_unknowns": ["Peak notification volume."],
    }


def valid_result(status: str = "complete") -> dict[str, Any]:
    chairman = (
        {
            "status": "complete",
            "model": "fixture-chair",
            "family": "family-d",
            "used_fallback": False,
            "participated_as_member": False,
            "request_id": "fixture-chair-request",
            "response": chairman_response(),
        }
        if status != "failed"
        else {}
    )
    return {
        "schema_version": 1,
        "request_id": "fixture-request",
        "run_id": "fixture-run",
        "status": status,
        "mode": "architecture",
        "diversity": {
            "status": "verified",
            "families": ["family-a", "family-b", "family-c"],
            "limitations": [],
        },
        "stages": {
            "preflight": {"status": "complete"},
            "independent": {"status": "complete"},
            "cross_critique": {"status": "complete"},
            "revision": {"status": "skipped"},
            "chairman": {"status": "complete" if chairman else "failed"},
        },
        "candidates": [
            {
                "candidate_label": "Candidate A",
                "member_id": "architect",
                "model": "fixture-model-a",
                "family": "family-a",
                "response": member_response(),
            }
        ],
        "critiques": [
            {
                "reviewer_member_id": "architect",
                "response": critique_response(["Candidate A"]),
            }
        ],
        "ranking": {
            "candidates": {
                "Candidate A": {"weighted_score": 0.75, "rank_score": 1.0}
            },
            "coverage": {"valid_critics": 1, "expected_critics": 1},
            "advisory": True,
        },
        "chairman": chairman,
        "benchmark": {
            "calls": [
                {
                    "model": "fixture|model-a",
                    "family": "family-a",
                    "stage": "independent",
                    "status": "success",
                    "latency_ms": 12,
                    "input_tokens": 10,
                    "output_tokens": 5,
                    "total_tokens": 15,
                    "cost_usd": None,
                    "cost_status": "unknown",
                    "retry_count": 0,
                    "quality": {
                        "peer_score": 0.75,
                        "rank_score": 1.0,
                        "chairman_selected": True,
                        "subjective": True,
                    },
                }
            ]
        },
        "usage": {
            "inference_requests": 7,
            "requested_output_tokens": 4000,
            "input_tokens": 70,
            "output_tokens": 35,
            "total_tokens": 105,
            "cost_usd": None,
            "cost_status": "unknown",
        },
        "failures": [],
        "persistence": {"mode": "none", "paths": []},
    }


def openai_response(
    content: dict[str, Any] | str,
    *,
    prompt_tokens: int = 10,
    completion_tokens: int = 5,
    cost: float | None = None,
) -> bytes:
    text = content if isinstance(content, str) else json.dumps(content)
    payload: dict[str, Any] = {
        "id": "fixture-request",
        "choices": [{"message": {"role": "assistant", "content": text}}],
        "usage": {
            "prompt_tokens": prompt_tokens,
            "completion_tokens": completion_tokens,
            "total_tokens": prompt_tokens + completion_tokens,
        },
    }
    if cost is not None:
        payload["response_cost"] = cost
    return json.dumps(payload).encode("utf-8")


def model_list_response(models: list[str]) -> bytes:
    return json.dumps({"data": [{"id": model} for model in models]}).encode("utf-8")


class RecordingTransport:
    """Thread-safe scripted transport implementing the runtime's injected contract."""

    def __init__(self, responses: list[Any]):
        self.responses = deque(responses)
        self.calls: list[dict[str, Any]] = []
        self.lock = threading.Lock()

    def __call__(
        self,
        method: str,
        url: str,
        headers: dict[str, str],
        body: bytes | None,
        timeout_seconds: float,
    ) -> Any:
        with self.lock:
            self.calls.append(
                {
                    "method": method,
                    "url": url,
                    "headers": dict(headers),
                    "body": body,
                    "timeout_seconds": timeout_seconds,
                }
            )
            if not self.responses:
                raise AssertionError(f"Unexpected transport call: {method} {url}")
            response = self.responses.popleft()
        if isinstance(response, BaseException):
            raise response
        if callable(response):
            return response(method, url, headers, body, timeout_seconds)
        return response


def transport_response(
    status: int,
    body: bytes | str | dict[str, Any],
    headers: dict[str, str] | None = None,
) -> Any:
    from llm_council.client import TransportResponse

    if isinstance(body, dict):
        body = json.dumps(body).encode("utf-8")
    elif isinstance(body, str):
        body = body.encode("utf-8")
    return TransportResponse(status=status, headers=headers or {}, body=body)
