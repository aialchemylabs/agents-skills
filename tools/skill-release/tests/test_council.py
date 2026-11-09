"""Exercise the real council orchestration with a deterministic fake transport."""

import copy
import io
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
from llm_council.cli import main
from llm_council.render import render_result
from llm_council.schemas import validate_manifest, validate_result, planned_base_requests

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
        self.data = fixture("non-code-decision") if mode == "non-code" else BASE
        self.posts = []
        self.gets = 0
        self.lock = threading.Lock()
        self.attempts = {}

    def __call__(self, method, url, headers, body, *timeouts):
        if method == "GET":
            self.gets += 1
            if self.mode == "auth":
                return TransportResponse(401, {}, b'{"error":"unauthorized"}')
            models = self.data["gateway"]["models"]
            if self.mode == "unavailable":
                models = models[1:]
            return TransportResponse(
                200,
                {},
                json.dumps(
                    {"data": [{"id": x} for x in models]}
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
        if self.mode == "critics" and stage == "critique":
            raise TimeoutError("fixture critique timeout")
        if stage == "repair":
            content = fixture("malformed-model-output")["events"][1]["content"]
        elif stage == "chairman":
            content = copy.deepcopy(BASE["gateway"]["events"][6]["content"])
        elif stage == "critique":
            content = copy.deepcopy(self.data["gateway"]["events"][3]["content"])
            labels = [c["candidate_label"] for c in user["candidates"]]
            content["evaluations"] = [
                e for e in content["evaluations"] if e["candidate_label"] in labels
            ]
            if "ranking" in content:
                content["ranking"] = [
                    label for label in content["ranking"] if label in labels
                ]
        else:
            content = copy.deepcopy(
                self.data["gateway"]["events"][
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
    def test_default_non_code_council_has_no_rankings_or_chairman(self):
        runner, transport = council("non-code")
        manifest = fixture("non-code-decision")["manifest"]
        manifest["approval"] = {"authority": "autonomous"}
        result = validate_result(runner.run(manifest))
        self.assertEqual(len(transport.posts), 6)
        self.assertEqual(result["stages"]["chairman"]["status"], "skipped")
        self.assertEqual(result["ranking"]["order"], [])
        self.assertTrue(all(c["weighted_score"] is None for c in result["ranking"]["candidates"].values()))
        self.assertEqual(result["execution_report"]["backend"], "litellm")
        self.assertEqual(result["execution_report"]["round_number"], 1)
        for request in transport.posts:
            packet = json.loads(request["messages"][1]["content"])
            self.assertEqual(packet["decision_criteria"], manifest["decision_criteria"])
            if "evaluations" in packet["response_schema"]:
                self.assertNotIn("ranking", packet["response_schema"])
                self.assertNotIn("scores", packet["response_schema"]["evaluations"][0])
        rendered = render_result(result)
        self.assertIn("bounded pilot", rendered)
        self.assertIn("fixture-model-a", rendered)
        self.assertIn("E1", rendered)
        self.assertIn("Host must record accepted/rejected", rendered)
        self.assertNotIn("Benchmark-only run", rendered)
        self.assertIn("Demand for each format", rendered)
        self.assertNotIn("$0.00", rendered)

    def test_interactive_waits_but_autonomous_uses_existing_authority(self):
        runner, transport = council("non-code")
        manifest = fixture("non-code-decision")["manifest"]
        self.assertTrue(runner.preflight(manifest)["approval_required"])
        with self.assertRaises(ValidationError):
            runner.run(manifest)
        self.assertEqual(transport.posts, [])
        manifest["approval"] = {"authority": "autonomous"}
        preview = runner.preflight(manifest)
        self.assertFalse(preview["approval_required"])
        result = runner.run(manifest)
        self.assertEqual(result["execution_report"]["plan_fingerprint"], preview["plan_fingerprint"])
        self.assertEqual(result["status"], "complete")

    def test_approved_interactive_run_needs_no_second_approval(self):
        runner, _ = council()
        manifest = approved(runner)
        self.assertFalse(runner.preflight(manifest)["approval_required"])
        self.assertEqual(runner.run(manifest)["status"], "complete")

    def test_autonomous_does_not_accept_stale_fingerprint(self):
        runner, transport = council()
        manifest = approved(runner)
        manifest["approval"]["authority"] = "autonomous"
        manifest["question"] += " changed"
        with self.assertRaises(PreflightBlocked):
            runner.run(manifest)
        self.assertEqual(transport.posts, [])

    def test_three_round_ceiling_and_no_reset_after_changes(self):
        runner, transport = council("non-code")
        manifest = fixture("non-code-decision")["manifest"]
        manifest["approval"] = {"authority": "autonomous"}
        for number in (1, 2, 3):
            manifest["execution"] = {"round_number": number}
            result = runner.run(manifest)
            self.assertEqual(result["execution_report"]["round_number"], number)
        self.assertIn("limit reached", result["execution_report"]["stop_reason"])
        self.assertEqual(len(transport.posts), 18)
        manifest["execution"]["round_number"] = 4
        with self.assertRaises(ValidationError):
            runner.run(manifest)
        manifest["execution"]["round_number"] = 1
        manifest["context"]["evidence"][0]["content"] += " Revised staffing evidence."
        with self.assertRaises(BudgetBlocked):
            runner.run(manifest)
        self.assertEqual(len(transport.posts), 18)

    def test_models_discovery_requires_no_manifest_and_makes_no_inference(self):
        transport = Transport()
        def factory(**kwargs):
            return LiteLLMClient(**kwargs, transport=transport)
        output = io.StringIO()
        code = main(["models"], stdin=io.StringIO(""), stdout=output,
                    environ={"LLM_COUNCIL_BASE_URL": "https://gateway.example", "LLM_COUNCIL_API_KEY": "fixture-key"},
                    client_factory=factory)
        self.assertEqual(code, 0)
        self.assertEqual(json.loads(output.getvalue())["accessible_models"], sorted(BASE["gateway"]["models"]))
        self.assertEqual(transport.gets, 1)
        self.assertEqual(transport.posts, [])
        self.assertNotIn("fixture-key", output.getvalue())

    def test_validation_is_offline_without_gateway_environment(self):
        def forbidden_factory(**kwargs):
            self.fail("Offline validation tried to access a gateway")
        output = io.StringIO()
        code = main(["validate"], stdin=io.StringIO(json.dumps(fixture("non-code-decision")["manifest"])),
                    stdout=output, environ={}, client_factory=forbidden_factory)
        self.assertEqual(code, 0)
        self.assertEqual(json.loads(output.getvalue())["base_inference_requests"], 6)

    def test_named_model_unavailability_and_auth_block_before_fanout(self):
        for mode in ("unavailable", "auth"):
            runner, transport = council(mode)
            manifest = approved(runner)
            manifest["approval"]["authority"] = "autonomous"
            with self.assertRaises(PreflightBlocked):
                runner.run(manifest)
            self.assertEqual(transport.posts, [])
            self.assertEqual(transport.gets, 1)

    def test_multiple_distinct_models_in_one_provider_are_allowed(self):
        manifest = fixture("non-code-decision")["manifest"]
        for member in manifest["roster"]["members"]:
            member["family"] = "same-provider"
        normalized = validate_manifest(manifest)
        self.assertEqual(planned_base_requests(normalized), 6)
        runner, _ = council()
        self.assertEqual(runner.preflight(normalized)["diversity"]["status"], "perspective-only")
        manifest["policy"] = {"require_distinct_families": True}
        with self.assertRaises(ValidationError):
            validate_manifest(manifest)
        manifest["policy"] = {"allow_perspective_only_fallback": True}
        manifest["roster"]["members"][1]["model"] = manifest["roster"]["members"][0]["model"]
        with self.assertRaises(ValidationError):
            validate_manifest(manifest)

    def test_failed_critics_keep_independent_evidence_and_disclose_failure(self):
        runner, _ = council("critics")
        result = runner.run(approved(runner))
        self.assertEqual(result["status"], "partial")
        self.assertEqual(len(result["candidates"]), 3)
        self.assertEqual(result["stages"]["cross_critique"]["status"], "failed")
        self.assertEqual(result["ranking"]["order"], [])
        self.assertEqual(len(result["failures"]), 3)
        self.assertIn("cross_critique: LiteLLM request timed out", render_result(result))

    def test_exhausted_budget_returns_dossier_without_extra_synthesis_call(self):
        runner, transport = council("malformed")
        manifest = approved(runner)
        # Seven base calls, but the one authorized format repair consumes the
        # remaining synthesis slot. Evidence must survive the hard stop.
        manifest["budget"]["max_inference_requests"] = 7
        manifest["approval"]["plan_fingerprint"] = compute_plan_fingerprint(manifest, runner.client.base_url)
        result = validate_result(runner.run(manifest))
        self.assertEqual(len(transport.posts), 7)
        self.assertEqual(result["usage"]["inference_requests"], 7)
        self.assertEqual(result["status"], "partial")
        self.assertEqual(len(result["candidates"]), 3)
        self.assertEqual(result["chairman"]["status"], "failed")
        self.assertTrue(any(f["category"] == "budget" for f in result["failures"]))
        self.assertNotIn(("chairman", "fixture-chair"), transport.attempts)

    def test_independent_stage_never_receives_other_assessments(self):
        runner, transport = council()
        runner.run(approved(runner))
        independent = [r for r in transport.posts if "independent member" in r["messages"][0]["content"]]
        self.assertEqual(len(independent), 3)
        for request in independent:
            packet = json.loads(request["messages"][1]["content"])
            self.assertNotIn("candidates", packet)
            self.assertNotIn("critiques", packet)
            self.assertNotIn("ranking", packet)
            self.assertIn("launch other councils", request["messages"][0]["content"])

    def test_cli_autonomous_round_runs_without_manual_fingerprint_step(self):
        transport = Transport("non-code")
        def factory(**kwargs):
            return LiteLLMClient(**kwargs, transport=transport)
        manifest = fixture("non-code-decision")["manifest"]
        manifest["approval"] = {"authority": "autonomous"}
        output = io.StringIO()
        code = main(["run"], stdin=io.StringIO(json.dumps(manifest)), stdout=output, stderr=io.StringIO(),
                    environ={"LLM_COUNCIL_BASE_URL": "https://gateway.example", "LLM_COUNCIL_API_KEY": "fixture-key"},
                    client_factory=factory)
        self.assertEqual(code, 0)
        result = validate_result(json.loads(output.getvalue()))
        self.assertEqual(result["status"], "complete")
        self.assertEqual(len(transport.posts), 6)
        self.assertNotIn("fixture-key", output.getvalue())

    def test_independent_benchmark_uses_identical_inputs_without_extra_judge(self):
        runner, transport = council("non-code")
        manifest = fixture("non-code-decision")["manifest"]
        manifest["mode"] = "benchmark"
        manifest["policy"] = {"run_cross_critique": False, "run_addy_revision": True}
        manifest["approval"] = {"authority": "autonomous"}
        result = runner.run(manifest)
        self.assertEqual(len(transport.posts), 3)
        self.assertEqual(result["stages"]["revision"]["status"], "skipped")
        self.assertEqual(result["stages"]["chairman"]["status"], "skipped")
        messages = [request["messages"] for request in transport.posts]
        self.assertTrue(all(message == messages[0] for message in messages))

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
