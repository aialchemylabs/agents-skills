from __future__ import annotations

import json
import stat
import tempfile
import unittest
from pathlib import Path
from unittest import mock

from support import valid_manifest, valid_result

from llm_council.errors import PersistenceFailure, ValidationError
from llm_council.persistence import persist_result


class PersistenceTests(unittest.TestCase):
    def test_none_mode_writes_nothing(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-test-") as temporary:
            output = Path(temporary) / "output"
            manifest = valid_manifest()
            persisted = persist_result(manifest, valid_result())
            self.assertEqual(persisted["mode"], "none")
            self.assertEqual(persisted["paths"], [])
            self.assertFalse(output.exists())

    def test_summary_omits_raw_evidence_candidates_and_critiques(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-test-") as temporary:
            output = Path(temporary) / "output"
            output.mkdir()
            manifest = valid_manifest()
            manifest["policy"]["persist"] = "summary"
            manifest["policy"]["output_directory"] = str(output)
            result = valid_result()
            result["candidates"][0]["response"][
                "recommendation"
            ] = "RAW MEMBER TRANSCRIPT MUST NOT PERSIST"
            result["critiques"][0]["response"][
                "revised_position"
            ] = "RAW CRITIQUE TRANSCRIPT MUST NOT PERSIST"
            persisted = persist_result(manifest, result)
            paths = [Path(path) for path in persisted["paths"]]
            self.assertTrue(paths)
            payload = "\n".join(path.read_text(encoding="utf-8") for path in paths)
            self.assertNotIn(manifest["context"]["evidence"][0]["content"], payload)
            self.assertNotIn('"critic_member_id"', payload)
            self.assertNotIn(
                result["candidates"][0]["response"]["recommendation"],
                payload,
            )
            self.assertNotIn(
                result["critiques"][0]["response"]["revised_position"],
                payload,
            )
            self.assertIn('"chairman"', payload)
            self.assertIn('"benchmark"', payload)

    def test_full_includes_approved_records_without_credentials(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-test-") as temporary:
            output = Path(temporary) / "output"
            output.mkdir()
            manifest = valid_manifest()
            manifest["policy"]["persist"] = "full"
            manifest["policy"]["output_directory"] = str(output)
            result = valid_result()
            persisted = persist_result(manifest, result)
            paths = [Path(path) for path in persisted["paths"]]
            payloads = [json.loads(path.read_text(encoding="utf-8")) for path in paths]
            serialized = json.dumps(payloads)
            self.assertIn("Candidate A", serialized)
            self.assertIn("fixture/architecture.md", serialized)
            self.assertNotIn("Authorization", serialized)
            self.assertNotIn("api_key", serialized.lower())

    def test_directory_and_file_modes_are_private(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-test-") as temporary:
            output = Path(temporary) / "output"
            output.mkdir()
            manifest = valid_manifest()
            manifest["policy"]["persist"] = "summary"
            manifest["policy"]["output_directory"] = str(output)
            persisted = persist_result(manifest, valid_result())
            paths = [Path(path) for path in persisted["paths"]]
            self.assertTrue(paths)
            run_directory = paths[0].parent
            self.assertEqual(stat.S_IMODE(run_directory.stat().st_mode), 0o700)
            for path in paths:
                self.assertEqual(stat.S_IMODE(path.stat().st_mode), 0o600)

    def test_existing_run_directory_is_not_overwritten(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-test-") as temporary:
            output = Path(temporary) / "output"
            output.mkdir()
            manifest = valid_manifest()
            manifest["policy"]["persist"] = "summary"
            manifest["policy"]["output_directory"] = str(output)
            result = valid_result()
            first = persist_result(manifest, result)
            snapshot = {
                path: Path(path).read_bytes()
                for path in first["paths"]
            }
            with self.assertRaises(PersistenceFailure):
                persist_result(manifest, result)
            for path, content in snapshot.items():
                self.assertEqual(Path(path).read_bytes(), content)

    def test_unsafe_targets_are_rejected(self) -> None:
        for unsafe in (Path("/"), Path.home(), Path(__file__).resolve().parents[3]):
            with self.subTest(unsafe=unsafe):
                manifest = valid_manifest()
                manifest["policy"]["persist"] = "summary"
                manifest["policy"]["output_directory"] = str(unsafe)
                with self.assertRaises((ValidationError, PersistenceFailure)):
                    persist_result(manifest, valid_result())

    def test_failed_atomic_replace_leaves_no_final_json(self) -> None:
        with tempfile.TemporaryDirectory(prefix="llm-council-test-") as temporary:
            output = Path(temporary) / "output"
            output.mkdir()
            manifest = valid_manifest()
            manifest["policy"]["persist"] = "summary"
            manifest["policy"]["output_directory"] = str(output)
            with mock.patch(
                "llm_council.persistence.os.replace",
                side_effect=OSError("fixture replace failure"),
            ):
                with self.assertRaises(PersistenceFailure):
                    persist_result(manifest, valid_result())
            self.assertEqual(list(output.rglob("*.json")), [])


if __name__ == "__main__":
    unittest.main()
