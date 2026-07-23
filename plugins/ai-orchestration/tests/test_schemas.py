from __future__ import annotations

import copy
import unittest

from support import (
    chairman_response,
    critique_response,
    member_response,
    revision_response,
    valid_manifest,
    valid_result,
)

from llm_council.schemas import (
    ValidationError,
    validate_chairman_response,
    validate_critique_response,
    validate_manifest,
    validate_member_response,
    validate_revision_response,
    validate_result,
)
from llm_council.orchestrator import compute_plan_fingerprint


BASE_URL = "https://gateway.invalid"


class ManifestSchemaTests(unittest.TestCase):
    def assert_invalid(self, manifest: dict, command: str = "validate") -> None:
        with self.assertRaises(ValidationError):
            validate_manifest(manifest, command=command)

    def test_accepts_full_valid_manifest_for_validate_and_preflight(self) -> None:
        manifest = valid_manifest()
        for command in ("validate", "preflight"):
            with self.subTest(command=command):
                normalized = validate_manifest(manifest, command=command)
                self.assertEqual(normalized["schema_version"], 1)
                self.assertEqual(len(normalized["roster"]["members"]), 3)

    def test_accepts_minimum_manifest_and_fills_non_secret_defaults(self) -> None:
        minimum = {
            "schema_version": 1,
            "request_id": "fixture-minimum",
            "mode": "decision",
            "question": "Choose the safer bounded option.",
            "decision_criteria": ["correctness"],
            "roster": {
                "members": [
                    {
                        "member_id": "first",
                        "model": "fixture-model-a",
                        "family": "family-a",
                        "angle": "correctness",
                    },
                    {
                        "member_id": "second",
                        "model": "fixture-model-b",
                        "family": "family-b",
                        "angle": "risk",
                    },
                ],
                "chairman": {
                    "model": "fixture-chair",
                    "family": "family-c",
                },
            },
        }
        normalized = validate_manifest(minimum, command="validate")
        self.assertIn("Read-only analysis", normalized["context"]["instructions"])
        self.assertEqual(normalized["context"]["evidence"], [])
        self.assertEqual(normalized["context"]["artifact_refs"], [])
        self.assertEqual(normalized["execution"]["max_concurrency"], 2)
        self.assertEqual(normalized["budget"]["max_inference_requests"], 10)
        self.assertEqual(normalized["policy"]["persist"], "none")
        self.assertFalse(normalized["approval"]["inference_approved"])
        self.assertIsNone(normalized["approval"]["plan_fingerprint"])
        self.assertIsNone(normalized["roster"]["fallback_chairman"])

    def test_rejects_empty_decision_criteria(self) -> None:
        manifest = valid_manifest()
        manifest["decision_criteria"] = []
        self.assert_invalid(manifest)

    def test_rejects_unknown_schema_version_and_mode(self) -> None:
        manifest = valid_manifest()
        manifest["schema_version"] = 2
        self.assert_invalid(manifest)
        manifest = valid_manifest()
        manifest["mode"] = "oracle"
        self.assert_invalid(manifest)

    def test_member_count_and_identity_constraints(self) -> None:
        manifest = valid_manifest()
        manifest["roster"]["members"] = manifest["roster"]["members"][:1]
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        member = copy.deepcopy(manifest["roster"]["members"][0])
        for index in range(3):
            extra = copy.deepcopy(member)
            extra["member_id"] = f"extra-{index}"
            extra["model"] = f"fixture-extra-{index}"
            extra["family"] = f"family-extra-{index}"
            manifest["roster"]["members"].append(extra)
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["roster"]["members"][1]["member_id"] = "architect"
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["roster"]["members"][0]["model"] = ""
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["roster"]["members"][0]["angle"] = " "
        self.assert_invalid(manifest)

    def test_rejects_embedded_credentials_and_unknown_critical_fields(self) -> None:
        for key in ("api_key", "authorization", "litellm_master_key"):
            with self.subTest(key=key):
                manifest = valid_manifest()
                manifest[key] = "unit-test-placeholder"
                self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["execution"]["max_concurency"] = 2
        self.assert_invalid(manifest)

    def test_sensitive_evidence_blocks_without_echoing_the_detected_value(self) -> None:
        manifest = valid_manifest()
        detected = "Authorization: Bearer unit-test-sensitive-placeholder"
        manifest["context"]["evidence"][0]["content"] = detected
        with self.assertRaises(ValidationError) as raised:
            validate_manifest(manifest, command="preflight")
        self.assertIn("$.context.evidence", str(raised.exception))
        self.assertNotIn(detected, str(raised.exception))
        self.assertNotIn("unit-test-sensitive-placeholder", str(raised.exception))

    def test_rejects_unsafe_or_incomplete_persistence(self) -> None:
        manifest = valid_manifest()
        manifest["policy"]["persist"] = "summary"
        self.assert_invalid(manifest)

        for unsafe in ("/", str(__import__("pathlib").Path.home()), "."):
            with self.subTest(path=unsafe):
                manifest = valid_manifest()
                manifest["policy"]["persist"] = "summary"
                manifest["policy"]["output_directory"] = unsafe
                self.assert_invalid(manifest)

    def test_rejects_hard_limit_violations(self) -> None:
        cases = (
            ("max_concurrency", 6),
            ("connect_timeout_seconds", 61),
            ("request_timeout_seconds", 601),
            ("max_transient_retries", 3),
        )
        for key, value in cases:
            with self.subTest(key=key):
                manifest = valid_manifest()
                manifest["execution"][key] = value
                self.assert_invalid(manifest)

        budget_cases = (
            ("max_inference_requests", 26),
            ("max_input_characters", 500001),
            ("max_output_tokens", 50001),
        )
        for key, value in budget_cases:
            with self.subTest(key=key):
                manifest = valid_manifest()
                manifest["budget"][key] = value
                self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["execution"]["stage_output_tokens"]["chairman"] = 8193
        self.assert_invalid(manifest)

    def test_rejects_invalid_concurrency_negative_values_and_large_evidence(self) -> None:
        manifest = valid_manifest()
        manifest["execution"]["max_concurrency"] = 4
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["execution"]["request_timeout_seconds"] = 0
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["execution"]["stage_output_tokens"]["critique"] = -1
        self.assert_invalid(manifest)

        manifest = valid_manifest()
        manifest["budget"]["max_input_characters"] = 10
        self.assert_invalid(manifest)

    def test_verified_diversity_policy_rejects_same_or_unknown_families(self) -> None:
        manifest = valid_manifest()
        for member in manifest["roster"]["members"]:
            member["family"] = "one-family"
        self.assert_invalid(manifest, command="preflight")

        manifest["policy"]["require_distinct_families"] = False
        manifest["policy"]["allow_perspective_only_fallback"] = True
        normalized = validate_manifest(manifest, command="preflight")
        self.assertFalse(normalized["policy"]["require_distinct_families"])

        manifest = valid_manifest()
        manifest["roster"]["members"][0]["family"] = "unverified"
        self.assert_invalid(manifest, command="preflight")

    def test_run_requires_approval_and_exact_fingerprint(self) -> None:
        manifest = valid_manifest()
        self.assert_invalid(manifest, command="run")

        manifest["approval"]["inference_approved"] = True
        manifest["approval"]["plan_fingerprint"] = compute_plan_fingerprint(
            manifest,
            BASE_URL,
            data_egress_classification="non-sensitive fixture evidence",
        )
        normalized = validate_manifest(manifest, command="run")
        self.assertTrue(normalized["approval"]["inference_approved"])


class FingerprintTests(unittest.TestCase):
    def fingerprint(self, manifest: dict) -> str:
        return compute_plan_fingerprint(
            manifest,
            BASE_URL,
            data_egress_classification="non-sensitive fixture evidence",
        )

    def test_fingerprint_is_stable_sha256_and_does_not_accept_credentials(self) -> None:
        manifest = valid_manifest()
        first = self.fingerprint(manifest)
        second = self.fingerprint(copy.deepcopy(manifest))
        self.assertEqual(first, second)
        self.assertRegex(first, r"^sha256:[0-9a-f]{64}$")
        self.assertNotIn("unit-test-placeholder", first)

    def test_execution_relevant_changes_invalidate_fingerprint(self) -> None:
        base = valid_manifest()
        fingerprint = self.fingerprint(base)
        mutations = []

        changed = valid_manifest()
        changed["question"] += " Include rollback."
        mutations.append(changed)

        changed = valid_manifest()
        changed["context"]["evidence"][0]["content"] += " Changed."
        mutations.append(changed)

        changed = valid_manifest()
        changed["roster"]["members"][0]["model"] = "fixture-model-z"
        mutations.append(changed)

        changed = valid_manifest()
        changed["roster"]["chairman"]["model"] = "fixture-chair-z"
        mutations.append(changed)

        changed = valid_manifest()
        changed["policy"]["expose_model_identities_after_synthesis"] = False
        mutations.append(changed)

        changed = valid_manifest()
        changed["budget"]["max_inference_requests"] = 9
        mutations.append(changed)

        changed = valid_manifest()
        changed["policy"]["persist"] = "summary"
        changed["policy"]["output_directory"] = "/tmp/llm-council-fixture-output"
        mutations.append(changed)

        for changed in mutations:
            with self.subTest(change=changed):
                self.assertNotEqual(fingerprint, self.fingerprint(changed))

        self.assertNotEqual(
            fingerprint,
            compute_plan_fingerprint(
                base,
                "https://other-gateway.invalid",
                data_egress_classification="non-sensitive fixture evidence",
            ),
        )
        self.assertNotEqual(
            fingerprint,
            compute_plan_fingerprint(
                base,
                BASE_URL,
                data_egress_classification="proprietary",
            ),
        )


class ModelOutputSchemaTests(unittest.TestCase):
    def test_accepts_all_valid_stage_objects(self) -> None:
        validate_member_response(member_response(), evidence_ids={"E1"})
        validate_critique_response(
            critique_response(["Candidate A", "Candidate B"]),
            candidate_labels=["Candidate A", "Candidate B"],
        )
        validate_revision_response(revision_response())
        validate_chairman_response(chairman_response())

    def test_rejects_member_with_missing_field_bad_evidence_or_confidence(self) -> None:
        invalid = member_response()
        invalid.pop("risks")
        with self.assertRaises(ValidationError):
            validate_member_response(invalid, evidence_ids={"E1"})

        invalid = member_response()
        invalid["evidence"][0]["evidence_id"] = "E99"
        with self.assertRaises(ValidationError):
            validate_member_response(invalid, evidence_ids={"E1"})

        invalid = member_response()
        invalid["confidence"] = 1.1
        with self.assertRaises(ValidationError):
            validate_member_response(invalid, evidence_ids={"E1"})

    def test_rejects_critique_with_unknown_label_duplicate_rank_or_bad_score(self) -> None:
        labels = ["Candidate A", "Candidate B"]
        invalid = critique_response(labels)
        invalid["evaluations"][0]["candidate_label"] = "Candidate Z"
        with self.assertRaises(ValidationError):
            validate_critique_response(invalid, candidate_labels=labels)

        invalid = critique_response(labels)
        invalid["ranking"] = ["Candidate A", "Candidate A"]
        with self.assertRaises(ValidationError):
            validate_critique_response(invalid, candidate_labels=labels)

        invalid = critique_response(labels)
        invalid["evaluations"][0]["scores"]["evidence_quality"] = 0
        with self.assertRaises(ValidationError):
            validate_critique_response(invalid, candidate_labels=labels)

    def test_rejects_invalid_revision_and_chairman(self) -> None:
        invalid_revision = revision_response()
        invalid_revision["accepted_feedback"][0].pop("response")
        with self.assertRaises(ValidationError):
            validate_revision_response(invalid_revision)

        invalid_chairman = chairman_response()
        invalid_chairman["confidence"] = "high"
        with self.assertRaises(ValidationError):
            validate_chairman_response(invalid_chairman)


class ResultSchemaTests(unittest.TestCase):
    def test_accepts_complete_versioned_result(self) -> None:
        validated = validate_result(valid_result())
        self.assertEqual(validated["schema_version"], 1)
        self.assertEqual(validated["usage"]["requested_output_tokens"], 4000)
        self.assertFalse(validated["chairman"]["participated_as_member"])

    def test_deep_result_validation_rejects_malformed_nested_contracts(self) -> None:
        mutations: list[tuple[str, callable]] = [
            (
                "candidate response",
                lambda result: result["candidates"][0]["response"].pop("risks"),
            ),
            (
                "unknown critique label",
                lambda result: result["critiques"][0]["response"]["evaluations"][0].update(
                    {"candidate_label": "Candidate Z"}
                ),
            ),
            (
                "stage status",
                lambda result: result["stages"]["independent"].update(
                    {"status": "done"}
                ),
            ),
            (
                "chairman response",
                lambda result: result["chairman"]["response"].update(
                    {"confidence": "certain"}
                ),
            ),
            (
                "unknown benchmark cost",
                lambda result: result["benchmark"]["calls"][0].update(
                    {"cost_usd": 0.0}
                ),
            ),
            (
                "benchmark quality",
                lambda result: result["benchmark"]["calls"][0]["quality"].update(
                    {"peer_score": 1.2}
                ),
            ),
            (
                "negative usage",
                lambda result: result["usage"].update(
                    {"requested_output_tokens": -1}
                ),
            ),
            (
                "persistence status",
                lambda result: result["persistence"].update({"status": "saved"}),
            ),
        ]
        for name, mutate in mutations:
            with self.subTest(contract=name):
                result = valid_result()
                mutate(result)
                with self.assertRaises(ValidationError):
                    validate_result(result)

    def test_identity_mapping_and_transitions_must_reference_valid_state(self) -> None:
        result = valid_result()
        result["identity_mapping"] = {"architect": "Candidate Z"}
        with self.assertRaises(ValidationError):
            validate_result(result)

        result = valid_result()
        result["transitions"] = [
            {
                "stage": "independent",
                "status": "complete",
                "monotonic_ms": -1,
                "reason": None,
                "call_ids": [],
            }
        ]
        with self.assertRaises(ValidationError):
            validate_result(result)


if __name__ == "__main__":
    unittest.main()
