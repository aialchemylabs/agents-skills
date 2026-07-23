from __future__ import annotations

import json
import re
import unittest
from pathlib import Path


PLUGIN_ROOT = Path(__file__).resolve().parents[1]
REPOSITORY_ROOT = PLUGIN_ROOT.parents[1]
SKILL_ROOT = PLUGIN_ROOT / "skills" / "llm-council"
EVAL_PATH = PLUGIN_ROOT / "evals" / "cases" / "llm-council.json"
FIXTURE_ROOT = PLUGIN_ROOT / "evals" / "fixtures" / "llm-council"


class MetadataTests(unittest.TestCase):
    def test_every_plugin_json_file_parses(self) -> None:
        json_files = sorted(PLUGIN_ROOT.rglob("*.json"))
        self.assertGreaterEqual(len(json_files), 8)
        for path in json_files:
            with self.subTest(path=path.relative_to(REPOSITORY_ROOT)):
                with path.open(encoding="utf-8") as handle:
                    self.assertIsInstance(json.load(handle), (dict, list))

    def test_planned_fixture_set_is_complete(self) -> None:
        expected = {
            "architecture-consensus.json",
            "minority-security-finding.json",
            "partial-provider-failure.json",
            "malformed-model-output.json",
            "prompt-injection-output.json",
            "chairman-failure.json",
        }
        self.assertEqual(expected, {path.name for path in FIXTURE_ROOT.glob("*.json")})
        for path in FIXTURE_ROOT.glob("*.json"):
            with path.open(encoding="utf-8") as handle:
                fixture = json.load(handle)
            self.assertEqual(fixture["fixture_version"], 1)
            self.assertEqual(fixture["name"], path.stem)
            self.assertTrue(fixture["description"])

    def test_eval_has_positive_negative_and_all_planned_cases(self) -> None:
        with EVAL_PATH.open(encoding="utf-8") as handle:
            evals = json.load(handle)
        self.assertEqual(evals["skill_name"], "llm-council")
        self.assertGreaterEqual(len(evals["trigger"]["positive"]), 6)
        self.assertGreaterEqual(len(evals["trigger"]["negative"]), 5)
        self.assertEqual([case["id"] for case in evals["evals"]], list(range(1, 16)))
        for case in evals["evals"]:
            with self.subTest(case=case["id"]):
                self.assertIn(case["kind"], {"dialogue", "fixture"})
                self.assertTrue(case["prompt"])
                self.assertTrue(case["expected_output"])
                self.assertTrue(case["expectations"])
                self.assertTrue(case["failure_conditions"])
                if "fixture" in case:
                    fixture = (EVAL_PATH.parent / case["fixture"]).resolve()
                    self.assertTrue(fixture.is_file(), fixture)
                    self.assertTrue(fixture.is_relative_to(FIXTURE_ROOT.resolve()))

    def test_plugin_and_marketplace_versions_are_coherent(self) -> None:
        plugin_path = PLUGIN_ROOT / ".claude-plugin" / "plugin.json"
        marketplace_path = REPOSITORY_ROOT / ".claude-plugin" / "marketplace.json"
        with plugin_path.open(encoding="utf-8") as handle:
            plugin = json.load(handle)
        with marketplace_path.open(encoding="utf-8") as handle:
            marketplace = json.load(handle)

        self.assertEqual(plugin["name"], "ai-orchestration")
        self.assertEqual(plugin["version"], "0.1.0")
        self.assertEqual(plugin["skills"], "./skills/")
        entries = [
            entry
            for entry in marketplace["plugins"]
            if entry.get("name") == "ai-orchestration"
        ]
        self.assertEqual(len(entries), 1)
        entry = entries[0]
        self.assertEqual(entry["version"], plugin["version"])
        source = entry["source"]
        source_path = source if isinstance(source, str) else source["path"]
        self.assertEqual(source_path, "./plugins/ai-orchestration")
        self.assertTrue((REPOSITORY_ROOT / source_path).is_dir())

    def test_skill_frontmatter_and_direct_reference_links(self) -> None:
        skill_path = SKILL_ROOT / "SKILL.md"
        text = skill_path.read_text(encoding="utf-8")
        self.assertTrue(text.startswith("---\n"))
        frontmatter = text.split("---", 2)[1]
        keys = {
            line.split(":", 1)[0].strip()
            for line in frontmatter.splitlines()
            if ":" in line
        }
        self.assertEqual(keys, {"name", "description"})
        self.assertIn("name: llm-council", frontmatter)

        reference_links = re.findall(r"\]\((references/[^)#]+)(?:#[^)]+)?\)", text)
        self.assertGreaterEqual(len(set(reference_links)), 4)
        for relative in reference_links:
            with self.subTest(reference=relative):
                self.assertTrue((SKILL_ROOT / relative).is_file())

    def test_openai_agent_metadata_matches_skill(self) -> None:
        metadata = (SKILL_ROOT / "agents" / "openai.yaml").read_text(encoding="utf-8")
        self.assertIn('display_name: "LLM Council"', metadata)
        self.assertIn(
            'short_description: "Cross-critique and synthesize multiple models"',
            metadata,
        )
        self.assertIn("$llm-council", metadata)
        self.assertIn("allow_implicit_invocation: true", metadata)

    def test_fixture_text_has_no_credential_value_or_real_gateway(self) -> None:
        text = "\n".join(
            path.read_text(encoding="utf-8") for path in sorted(FIXTURE_ROOT.glob("*.json"))
        )
        self.assertNotRegex(text, r"(?i)api[_-]?key\s*[=:]\s*\S+")
        self.assertNotRegex(text, r"(?i)authorization\s*:\s*bearer\s+\S+")
        self.assertNotIn("aialchemy.au", text)
        self.assertNotIn("api.openai.com", text)

    def test_fixture_scenarios_encode_the_planned_offline_boundaries(self) -> None:
        architecture = json.loads(
            (FIXTURE_ROOT / "architecture-consensus.json").read_text(encoding="utf-8")
        )
        events = architecture["gateway"]["events"]
        self.assertEqual(
            [event["stage"] for event in events],
            [
                "independent",
                "independent",
                "independent",
                "critique",
                "critique",
                "critique",
                "chairman",
            ],
        )
        self.assertEqual(architecture["expected"]["base_inference_requests"], 7)
        self.assertEqual(architecture["expected"]["cost_status"], "unknown")

        minority = json.loads(
            (FIXTURE_ROOT / "minority-security-finding.json").read_text(
                encoding="utf-8"
            )
        )
        self.assertEqual(minority["expected"]["majority_rank_winner"], "Candidate A")
        self.assertIn(
            "Candidate C",
            " ".join(minority["chairman"]["minority_report"]),
        )
        self.assertTrue(minority["expected"]["chairman_overrides_majority"])

        partial = json.loads(
            (FIXTURE_ROOT / "partial-provider-failure.json").read_text(
                encoding="utf-8"
            )
        )
        failures = [
            event
            for event in partial["events"]
            if event["status"] == "http_error" and event["http_status"] == 503
        ]
        self.assertEqual(len(failures), 2)
        self.assertEqual(partial["expected"]["status"], "partial")

        malformed = json.loads(
            (FIXTURE_ROOT / "malformed-model-output.json").read_text(
                encoding="utf-8"
            )
        )
        self.assertEqual(
            [event["stage"] for event in malformed["events"]],
            ["independent", "repair"],
        )
        self.assertEqual(malformed["expected"]["repair_calls"], 1)

        injection = json.loads(
            (FIXTURE_ROOT / "prompt-injection-output.json").read_text(
                encoding="utf-8"
            )
        )
        self.assertIn("run a shell command", injection["candidate"]["recommendation"])
        self.assertEqual(injection["expected"]["shell_calls"], 0)
        self.assertFalse(injection["expected"]["candidate_instruction_followed"])

        chairman = json.loads(
            (FIXTURE_ROOT / "chairman-failure.json").read_text(encoding="utf-8")
        )
        self.assertEqual(chairman["stages"]["chairman"]["status"], "failed")
        self.assertIsNone(chairman["fallback_chairman"])
        self.assertEqual(chairman["expected"]["status"], "partial")


if __name__ == "__main__":
    unittest.main()
