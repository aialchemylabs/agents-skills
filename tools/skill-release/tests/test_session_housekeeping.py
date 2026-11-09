"""Offline lifecycle checks for the session housekeeping helper."""

import importlib.util
import json
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[3]
SCRIPT = ROOT / "plugins/engineering/skills/session-housekeeping/scripts/housekeeping.py"
SPEC = importlib.util.spec_from_file_location("housekeeping_lifecycle", SCRIPT)
housekeeping = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(housekeeping)


class LifecycleTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        self.ledger = self.root / "ledger.json"

    def command(self, *args, success=True):
        result = subprocess.run(
            [sys.executable, str(SCRIPT), *args, "--ledger", str(self.ledger),
             "--session", "fixture"],
            capture_output=True, text=True, timeout=30,
        )
        self.assertEqual(result.returncode == 0, success, result.stderr)
        return json.loads(result.stdout) if success else result.stderr

    def scratch(self, *consumers):
        args = ["temp", "--purpose", "fixture scratch"]
        for handle in consumers:
            args.extend(["--consumer", handle])
        path = Path(self.command(*args)["created"])
        self.addCleanup(lambda: shutil.rmtree(path) if path.exists() else None)
        return path

    def release(self, path, handle):
        return self.command(
            "release", "--path", str(path), "--consumer", handle,
            "--release-evidence", "owning workflow confirms this instance closed",
        )

    def cleanup_path(self, path, *args):
        return self.command("cleanup", "--path", str(path), "--apply", *args)["resources"]

    def repository(self):
        repo = self.root / "repo"
        repo.mkdir()
        housekeeping.git(repo, "init", "-b", "main")
        housekeeping.git(repo, "config", "user.name", "Fixture")
        housekeeping.git(repo, "config", "user.email", "fixture@example.com")
        (repo / "file").write_text("initial")
        housekeeping.git(repo, "add", "file")
        housekeeping.git(repo, "commit", "-m", "initial")
        checkout = self.root / "checkout"
        housekeeping.git(repo, "worktree", "add", "-b", "task", str(checkout))
        return repo, checkout

    def test_scratch_waits_for_every_consumer_and_preserves_output_until_export(self):
        path = self.scratch("recorder:demo")
        other = self.scratch()
        self.command("hold", "--path", str(path), "--consumer", "watcher:export")
        output = path / "demo.mp4"
        output.write_bytes(b"fixture export")
        first = self.cleanup_path(path, "--inactive-worktrees")[0]
        self.assertEqual(first["status"], "retain")
        self.assertIn("recorder:demo", first["reason"])
        self.assertIn("watcher:export", first["reason"])
        self.release(path, "recorder:demo")
        remaining = self.cleanup_path(path)[0]
        self.assertNotIn("recorder:demo", remaining["reason"])
        self.assertIn("watcher:export", remaining["reason"])
        self.release(path, "watcher:export")
        self.assertEqual(
            self.cleanup_path(path, "--retain", str(output))[0]["status"], "retain"
        )
        durable = self.root / "delivered.mp4"
        shutil.copyfile(output, durable)
        self.assertEqual(durable.read_bytes(), output.read_bytes())
        # Closure plus verified output preservation permits prompt exact-path cleanup.
        self.assertEqual(self.cleanup_path(path)[0]["status"], "removed")
        self.assertFalse(path.exists())
        self.assertTrue(other.exists())
        self.assertEqual(durable.read_bytes(), b"fixture export")
        self.assertEqual(self.cleanup_path(path)[0]["status"], "already_removed")
        self.assertTrue(other.exists())

    def test_unknown_or_empty_release_never_clears_a_hold(self):
        path = self.scratch("browser:demo", "process:export")
        before = self.ledger.read_bytes()
        for handles, evidence in [
            (["browser:demo", "browser:unknown"], "closed"),
            (["browser:demo"], " "),
            (["browser:demo"], None),
        ]:
            args = ["release", "--path", str(path)]
            for handle in handles:
                args.extend(["--consumer", handle])
            if evidence is not None:
                args.extend(["--release-evidence", evidence])
            self.command(*args, success=False)
            self.assertEqual(self.ledger.read_bytes(), before)
        self.command("hold", "--path", str(path), "--consumer", "browser:demo", success=False)
        self.assertEqual(self.ledger.read_bytes(), before)
        self.assertEqual(self.cleanup_path(path)[0]["status"], "retain")

    def test_release_is_local_to_the_exact_resource_not_global_to_handle(self):
        first = self.scratch("watcher:shared")
        second = self.scratch("watcher:shared")
        self.release(first, "watcher:shared")
        self.assertEqual(self.cleanup_path(first)[0]["status"], "removed")
        self.assertEqual(self.cleanup_path(second)[0]["status"], "retain")
        self.assertTrue(second.exists())
        self.command("cleanup", "--path", str(self.root / "unknown"), "--apply", success=False)
        self.assertTrue(second.exists())

    def test_new_use_after_closure_records_a_fresh_handle(self):
        path = self.scratch("browser:first")
        self.release(path, "browser:first")
        self.command(
            "hold", "--path", str(path), "--consumer", "browser:first", success=False
        )
        self.command("hold", "--path", str(path), "--consumer", "browser:second")
        report = self.cleanup_path(path)[0]
        self.assertEqual(report["status"], "retain")
        self.assertIn("browser:second", report["reason"])
        self.assertNotIn("browser:first", report["reason"])
        self.assertTrue(path.exists())

    def test_host_checkout_needs_host_retirement_even_after_child_closes(self):
        _, checkout = self.repository()
        self.command(
            "register-worktree", "--path", str(checkout), "--purpose", "child checkout",
            "--created-by-session", "--creation-evidence", "host creation receipt",
            "--consumer", "task:review", "--host-managed",
        )
        self.release(checkout, "task:review")
        report = self.cleanup_path(checkout, "--inactive-worktrees", "--merged-into", "main")[0]
        self.assertEqual(report["status"], "retain")
        self.assertIn("host-managed", report["reason"])
        self.assertTrue(checkout.exists())

    def test_inactivity_flag_does_not_override_consumers_or_squash_ancestry(self):
        repo, checkout = self.repository()
        data = housekeeping.load(self.ledger, "fixture")
        housekeeping.register(
            self.ledger, data, checkout, consumers=["process:log-writer"]
        )
        report = self.cleanup_path(checkout, "--inactive-worktrees", "--merged-into", "main")[0]
        self.assertEqual(report["status"], "retain")
        self.assertIn("process:log-writer", report["reason"])
        (checkout / "file").write_text("task change")
        housekeeping.git(checkout, "commit", "-am", "task change")
        housekeeping.git(repo, "merge", "--squash", "task")
        housekeeping.git(repo, "commit", "-m", "squashed task")
        self.release(checkout, "process:log-writer")
        report = self.cleanup_path(checkout, "--inactive-worktrees", "--merged-into", "main")[0]
        self.assertEqual(report["status"], "retain")
        self.assertIn("not reachable", report["reason"])
        self.assertTrue(checkout.exists())

    def test_failed_ownership_save_does_not_leave_untracked_scratch(self):
        allocated = self.root / "scratch"
        allocated.mkdir()
        with (
            patch.object(sys, "argv", [str(SCRIPT), "temp", "--ledger", str(self.ledger),
                                      "--session", "fixture", "--purpose", "scratch"]),
            patch.object(housekeeping.tempfile, "mkdtemp", return_value=str(allocated)),
            patch.object(housekeeping, "save", side_effect=OSError("fixture write failure")),
            patch.object(sys, "stderr"),
        ):
            with self.assertRaises(SystemExit):
                housekeeping.main()
        self.assertFalse(allocated.exists())
        self.assertFalse(self.ledger.exists())


if __name__ == "__main__":
    unittest.main()
