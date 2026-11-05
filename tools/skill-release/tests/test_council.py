"""Exercise the real council orchestration with a deterministic fake transport."""

import copy
import json
from pathlib import Path
import sys
import tempfile
import threading
import unittest

ROOT = Path(__file__).resolve().parents[3]
COUNCIL = ROOT / "plugins/ai-orchestration"
sys.path.insert(0, str(COUNCIL / "skills/llm-council/scripts"))
from llm_council.client import LiteLLMClient, TransportResponse
from llm_council.orchestrator import (
    CouncilOrchestrator,
    RequestBudget,
    compute_plan_fingerprint,
)
from llm_council.errors import (
    BudgetBlocked,
    Cancelled,
    InsufficientMembers,
    PreflightBlocked,
    ValidationError,
)
from llm_council.persistence import persist_result

FIXTURES = COUNCIL / "evals/fixtures/llm-council"


def fixture(name):
    return json.loads((FIXTURES / (name + ".json")).read_text())


BASE = fixture("architecture-consensus")


class StableOrder:
    def shuffle(self, items):
        pass


class Transport:
    def __init__(self, mode="normal"):
        self.mode = mode
        self.posts = []
        self.lock = threading.Lock()
        self.attempts = {}

    def __call__(self, method, url, headers, body, *timeouts):
        if method == "GET":
            return TransportResponse(
                200,
                {},
                json.dumps(
                    {"data": [{"id": x} for x in BASE["gateway"]["models"]]}
                ).encode(),
            )
        payload = json.loads(body)
        system = payload["messages"][0]["content"]
        user = json.loads(payload["messages"][1]["content"])
        model = payload["model"]
        stage = (
            "repair"
            if "invalid_response" in user
            else "chairman"
            if model == "fixture-chair"
            else "critique"
            if "critic" in system.lower()
            else "independent"
        )
        with self.lock:
            self.posts.append(payload)
            key = (stage, model)
            self.attempts[key] = self.attempts.get(key, 0) + 1
        if (
            self.mode == "partial"
            and model == "fixture-model-b"
            and stage == "independent"
        ):
            return TransportResponse(503, {"retry-after": "0"}, b"{}")
        if (
            self.mode == "quorum"
            and model in ("fixture-model-a", "fixture-model-b")
            and stage == "independent"
        ):
            return TransportResponse(503, {"retry-after": "0"}, b"{}")
        if self.mode == "chairman" and stage == "chairman":
            raise TimeoutError("fixture ambiguous timeout")
        if stage == "repair":
            content = fixture("malformed-model-output")["events"][1]["content"]
        elif stage == "chairman":
            content = copy.deepcopy(BASE["gateway"]["events"][6]["content"])
        elif stage == "critique":
            content = copy.deepcopy(BASE["gateway"]["events"][3]["content"])
            labels = [c["candidate_label"] for c in user["candidates"]]
            content["evaluations"] = [
                e for e in content["evaluations"] if e["candidate_label"] in labels
            ]
            content["ranking"] = [
                label for label in content["ranking"] if label in labels
            ]
        else:
            content = copy.deepcopy(
                BASE["gateway"]["events"][
                    ["fixture-model-a", "fixture-model-b", "fixture-model-c"].index(
                        model
                    )
                ]["content"]
            )
        if (
            self.mode == "malformed"
            and stage == "independent"
            and model == "fixture-model-a"
        ):
            content = fixture("malformed-model-output")["events"][0]["raw_content"]
        if (
            self.mode == "injection"
            and stage == "independent"
            and model == "fixture-model-a"
        ):
            content = fixture("prompt-injection-output")["candidate"]
            content = {k: v for k, v in content.items() if k != "candidate_label"}
        if self.mode == "minority":
            minor = fixture("minority-security-finding")
            if stage == "chairman":
                content = minor["chairman"]
            elif stage == "critique":
                content["ranking"] = minor["peer_rankings"][
                    ["fixture-model-a", "fixture-model-b", "fixture-model-c"].index(
                        model
                    )
                ]
                for e in content["evaluations"]:
                    e["scores"] = {k: 4 for k in e["scores"]}
        text = content if isinstance(content, str) else json.dumps(content)
        return TransportResponse(
            200,
            {},
            json.dumps(
                {
                    "choices": [{"message": {"content": text}}],
                    "usage": {"prompt_tokens": 10, "completion_tokens": 10},
                }
            ).encode(),
        )


def council(mode="normal", event=None):
    transport = Transport(mode)
    client = LiteLLMClient(
        "https://gateway.example",
        "fixture-key",
        transport=transport,
        sleep=lambda _: None,
    )
    return CouncilOrchestrator(
        client, random_source=StableOrder(), cancellation_event=event
    ), transport


def approved(orchestrator):
    manifest = copy.deepcopy(BASE["manifest"])
    manifest["approval"]["inference_approved"] = True
    manifest["approval"]["plan_fingerprint"] = compute_plan_fingerprint(
        manifest, orchestrator.client.base_url
    )
    return manifest


class CouncilTests(unittest.TestCase):
    def test_consensus_and_request_contract(self):
        runner, transport = council()
        manifest = approved(runner)
        result = runner.run(manifest)
        self.assertEqual(result["status"], BASE["expected"]["status"])
        self.assertEqual(
            len(transport.posts), BASE["expected"]["base_inference_requests"]
        )
        self.assertIn(
            BASE["expected"]["recommendation_contains"], json.dumps(result["chairman"])
        )
        for request in transport.posts:
            self.assertNotIn("tools", request)
            self.assertIn("max_tokens", request)
            self.assertNotIn("fixture-key", json.dumps(request))
        self.assertEqual(result["usage"]["cost_status"], "unknown")
        self.assertEqual(persist_result(manifest, result)["paths"], [])

    def test_partial_provider_failure_and_retry_limit(self):
        runner, transport = council("partial")
        result = runner.run(approved(runner))
        expected = fixture("partial-provider-failure")["expected"]
        self.assertEqual(result["status"], expected["status"])
        self.assertEqual(len(result["candidates"]), expected["successful_members"])
        self.assertEqual(transport.attempts[("independent", "fixture-model-b")], 2)
        self.assertEqual(
            result["stages"]["chairman"]["status"], expected["chairman_status"]
        )

    def test_malformed_output_repair_is_bounded_and_minimal(self):
        runner, transport = council("malformed")
        result = runner.run(approved(runner))
        self.assertEqual(result["status"], "complete")
        repairs = [
            json.loads(r["messages"][1]["content"])
            for r in transport.posts
            if "invalid_response" in json.loads(r["messages"][1]["content"])
        ]
        self.assertEqual(len(repairs), 1)
        self.assertEqual(
            set(repairs[0]),
            set(fixture("malformed-model-output")["repair_request_must_include_only"]),
        )

    def test_chairman_timeout_returns_partial_without_retry(self):
        runner, transport = council("chairman")
        result = runner.run(approved(runner))
        expected = fixture("chairman-failure")["expected"]
        self.assertEqual(result["status"], expected["status"])
        self.assertEqual(result["stages"]["chairman"]["status"], "failed")
        self.assertEqual(transport.attempts[("chairman", "fixture-chair")], 1)

    def test_injected_candidate_is_only_data(self):
        runner, transport = council("injection")
        result = runner.run(approved(runner))
        self.assertEqual(result["status"], "complete")
        critic = next(
            r
            for r in transport.posts
            if "critic" in r["messages"][0]["content"].lower()
        )
        self.assertIn(
            fixture("prompt-injection-output")["candidate"]["recommendation"],
            critic["messages"][1]["content"],
        )
        self.assertNotIn("fixture-key", critic["messages"][1]["content"])
        self.assertNotIn("tools", critic)

    def test_chairman_preserves_minority_over_rank(self):
        runner, _ = council("minority")
        result = runner.run(approved(runner))
        self.assertIn("tenant predicate", json.dumps(result["chairman"]))
        self.assertIn("minority_report", result["chairman"]["response"])

    def test_approval_quorum_budget_and_cancellation(self):
        runner, transport = council()
        with self.assertRaises(ValidationError):
            runner.run(BASE["manifest"])
        manifest = approved(runner)
        manifest["question"] += " changed"
        with self.assertRaises(PreflightBlocked):
            runner.run(manifest)
        self.assertEqual(transport.posts, [])
        runner, _ = council("quorum")
        with self.assertRaises(InsufficientMembers):
            runner.run(approved(runner))
        event = threading.Event()
        event.set()
        runner, transport = council(event=event)
        with self.assertRaises(Cancelled):
            runner.run(approved(runner))
        self.assertEqual(transport.posts, [])
        budget = RequestBudget(
            max_requests=1,
            max_usd=None,
            cancellation_event=threading.Event(),
            max_output_tokens=100,
        )
        budget.before_attempt("test", "fake", 100)
        with self.assertRaises(BudgetBlocked):
            budget.before_attempt("test", "fake", 1)

    def test_summary_persistence_is_explicit_and_does_not_overwrite(self):
        runner, _ = council()
        manifest = approved(runner)
        result = runner.run(manifest)
        with tempfile.TemporaryDirectory() as directory:
            manifest["policy"]["persist"] = "summary"
            manifest["policy"]["output_directory"] = directory
            persisted = persist_result(manifest, result)
            self.assertEqual(len(persisted["paths"]), 1)
            content = Path(persisted["paths"][0]).read_text()
            self.assertNotIn("fixture-key", content)
            with self.assertRaises(Exception):
                persist_result(manifest, result)


if __name__ == "__main__":
    unittest.main()
