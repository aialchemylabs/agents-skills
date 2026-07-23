from __future__ import annotations

import copy
import io
import json
import tempfile
import unittest
from pathlib import Path
from typing import Any

from support import valid_manifest, valid_result

from llm_council.cli import main
from llm_council.errors import Cancelled


FAKE_ENVIRONMENT = {
    "LLM_COUNCIL_BASE_URL": "https://gateway.invalid",
    "LLM_COUNCIL_API_KEY": "unit-test-placeholder",
}


class FakeCLIClient:
    instances: list["FakeCLIClient"] = []

    def __init__(self, **kwargs: Any) -> None:
        self.kwargs = kwargs
        self.base_url = kwargs["base_url"]
        self.__class__.instances.append(self)


class FakeCLIOrchestrator:
    instances: list["FakeCLIOrchestrator"] = []

    def __init__(self, client: FakeCLIClient, **kwargs: Any) -> None:
        self.client = client
        self.kwargs = kwargs
        self.preflight_manifest: dict[str, Any] | None = None
        self.run_manifest: dict[str, Any] | None = None
        self.__class__.instances.append(self)

    def preflight(self, manifest: dict[str, Any]) -> dict[str, Any]:
        self.preflight_manifest = copy.deepcopy(manifest)
        return {
            "schema_version": 1,
            "request_id": manifest["request_id"],
            "status": "ready",
            "accessible_models": sorted(
                {
                    member["model"] for member in manifest["roster"]["members"]
                }
            ),
            "cost": {"cost_usd": None, "cost_status": "unknown"},
            "approval_required": True,
            "plan_fingerprint": "sha256:" + ("a" * 64),
        }

    def run(self, manifest: dict[str, Any]) -> dict[str, Any]:
        self.run_manifest = copy.deepcopy(manifest)
        return valid_result()


class CancelledCLIOrchestrator(FakeCLIOrchestrator):
    def run(self, manifest: dict[str, Any]) -> dict[str, Any]:
        raise Cancelled("Fixture cancellation; an in-flight call may still bill")


class CLITests(unittest.TestCase):
    def setUp(self) -> None:
        FakeCLIClient.instances.clear()
        FakeCLIOrchestrator.instances.clear()
        CancelledCLIOrchestrator.instances.clear()

    def invoke(
        self,
        subcommand: str,
        payload: dict[str, Any] | str,
        *,
        orchestrator_factory: Any = FakeCLIOrchestrator,
        environ: dict[str, str] | None = None,
    ) -> tuple[int, str, str]:
        input_text = payload if isinstance(payload, str) else json.dumps(payload)
        stdout = io.StringIO()
        stderr = io.StringIO()
        code = main(
            [subcommand],
            stdin=io.StringIO(input_text),
            stdout=stdout,
            stderr=stderr,
            environ=FAKE_ENVIRONMENT if environ is None else environ,
            client_factory=FakeCLIClient,
            orchestrator_factory=orchestrator_factory,
        )
        return code, stdout.getvalue(), stderr.getvalue()

    def test_validate_is_offline_and_emits_call_plan(self) -> None:
        code, stdout, stderr = self.invoke("validate", valid_manifest(), environ={})
        self.assertEqual(code, 0)
        result = json.loads(stdout)
        self.assertEqual(result["status"], "valid")
        self.assertEqual(result["base_inference_requests"], 7)
        self.assertEqual(result["planned_output_tokens"], 4000)
        self.assertTrue(result["approval_required_for_run"])
        self.assertEqual(FakeCLIClient.instances, [])
        self.assertEqual(FakeCLIOrchestrator.instances, [])
        self.assertEqual(stderr, "")

    def test_preflight_uses_injected_fake_and_makes_no_real_network_call(self) -> None:
        code, stdout, stderr = self.invoke("preflight", valid_manifest())
        self.assertEqual(code, 0)
        preview = json.loads(stdout)
        self.assertEqual(preview["status"], "ready")
        self.assertEqual(preview["cost"]["cost_status"], "unknown")
        self.assertRegex(preview["plan_fingerprint"], r"^sha256:[0-9a-f]{64}$")
        self.assertEqual(len(FakeCLIClient.instances), 1)
        self.assertEqual(len(FakeCLIOrchestrator.instances), 1)
        client_kwargs = FakeCLIClient.instances[0].kwargs
        self.assertEqual(client_kwargs["base_url"], "https://gateway.invalid")
        self.assertEqual(client_kwargs["api_key"], "unit-test-placeholder")
        self.assertNotIn("unit-test-placeholder", stdout)
        self.assertNotIn("unit-test-placeholder", stderr)

    def test_full_run_uses_injected_fake_and_emits_versioned_result(self) -> None:
        manifest = valid_manifest()
        manifest["approval"]["inference_approved"] = True
        manifest["approval"]["plan_fingerprint"] = "sha256:" + ("b" * 64)
        code, stdout, stderr = self.invoke("run", manifest)
        self.assertEqual(code, 0)
        result = json.loads(stdout)
        self.assertEqual(result["schema_version"], 1)
        self.assertEqual(result["status"], "complete")
        self.assertEqual(result["persistence"]["mode"], "none")
        self.assertEqual(len(FakeCLIOrchestrator.instances), 1)
        self.assertIsNotNone(FakeCLIOrchestrator.instances[0].run_manifest)
        self.assertNotIn("unit-test-placeholder", stdout)
        self.assertNotIn("unit-test-placeholder", stderr)

    def test_render_is_offline_and_leads_with_recommendation(self) -> None:
        code, stdout, stderr = self.invoke("render", valid_result(), environ={})
        self.assertEqual(code, 0)
        self.assertTrue(stdout.startswith("## Recommendation\n"))
        self.assertIn("## Council quality and benchmark", stdout)
        self.assertEqual(FakeCLIClient.instances, [])
        self.assertEqual(FakeCLIOrchestrator.instances, [])
        self.assertEqual(stderr, "")

    def test_invalid_json_maps_to_stable_input_exit(self) -> None:
        code, stdout, _ = self.invoke("validate", "{not json", environ={})
        self.assertEqual(code, 2)
        result = json.loads(stdout)
        self.assertEqual(result["status"], "blocked")
        self.assertEqual(result["error"]["category"], "invalid_input")
        self.assertNotIn("{not json", stdout)

    def test_preflight_without_key_is_blocked_before_client_creation(self) -> None:
        code, stdout, _ = self.invoke(
            "preflight",
            valid_manifest(),
            environ={
                "LLM_COUNCIL_BASE_URL": "https://gateway.invalid",
                "LITELLM_MASTER_KEY": "unit-test-master-placeholder",
            },
        )
        self.assertEqual(code, 3)
        result = json.loads(stdout)
        self.assertEqual(result["status"], "blocked")
        self.assertEqual(result["error"]["category"], "authentication")
        self.assertEqual(FakeCLIClient.instances, [])
        self.assertNotIn("unit-test-master-placeholder", stdout)

    def test_preflight_applies_config_defaults_before_orchestration(self) -> None:
        manifest = valid_manifest()
        manifest["execution"].pop("max_concurrency")
        with tempfile.TemporaryDirectory(prefix="llm-council-cli-config-") as temporary:
            config_path = Path(temporary) / "policy.json"
            config_path.write_text(
                json.dumps(
                    {
                        "schema_version": 1,
                        "policy": {
                            "defaults": {"execution": {"max_concurrency": 2}}
                        },
                    }
                ),
                encoding="utf-8",
            )
            environment = {
                **FAKE_ENVIRONMENT,
                "LLM_COUNCIL_CONFIG": str(config_path),
            }
            code, _, _ = self.invoke(
                "preflight",
                manifest,
                environ=environment,
            )
        self.assertEqual(code, 0)
        effective = FakeCLIOrchestrator.instances[0].preflight_manifest
        self.assertIsNotNone(effective)
        self.assertEqual(effective["execution"]["max_concurrency"], 2)

    def test_preflight_rejects_deployment_maximum_violation_before_client(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-cli-config-") as temporary:
            config_path = Path(temporary) / "policy.json"
            config_path.write_text(
                json.dumps(
                    {
                        "schema_version": 1,
                        "policy": {"maximums": {"max_concurrency": 2}},
                    }
                ),
                encoding="utf-8",
            )
            environment = {
                **FAKE_ENVIRONMENT,
                "LLM_COUNCIL_CONFIG": str(config_path),
            }
            code, stdout, _ = self.invoke(
                "preflight",
                valid_manifest(),
                environ=environment,
            )
        self.assertEqual(code, 2)
        self.assertEqual(json.loads(stdout)["error"]["category"], "invalid_input")
        self.assertEqual(FakeCLIClient.instances, [])

    def test_cancellation_maps_to_130_and_bounded_json(self) -> None:
        manifest = valid_manifest()
        manifest["approval"]["inference_approved"] = True
        manifest["approval"]["plan_fingerprint"] = "sha256:" + ("c" * 64)
        code, stdout, _ = self.invoke(
            "run",
            manifest,
            orchestrator_factory=CancelledCLIOrchestrator,
        )
        self.assertEqual(code, 130)
        result = json.loads(stdout)
        self.assertEqual(result["status"], "cancelled")
        self.assertEqual(result["error"]["category"], "cancelled")


if __name__ == "__main__":
    unittest.main()
