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
SKILLS = ROOT / "plugins/engineering/skills"


def module(name, skill, script):
    spec = importlib.util.spec_from_file_location(
        name, SKILLS / skill / "scripts" / script
    )
    loaded = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(loaded)
    return loaded


sources = module("design_sources", "design-reference-research", "sources.py")
video = module("production_video", "video-production", "video.py")
housekeeping = module("session_cleanup", "session-housekeeping", "housekeeping.py")


class WorkflowTests(unittest.TestCase):
    def test_user_sources_preserve_deep_links_and_override_defaults(self):
        catalog = json.loads(sources.CATALOG.read_text())
        supplied = [
            "https://www.mobbin.com/",
            "https://mobbin.com",
            "https://example.com/demo?mode=mobile#tasks",
        ]
        result = sources.select(catalog, supplied, "mobile")
        self.assertEqual(result["sources"][0]["origin"], "user")
        self.assertEqual(result["sources"][0]["name"], "Mobbin")
        self.assertIn("login", result["sources"][0]["verification"]["evidence"])
        self.assertEqual(result["sources"][1]["url"], supplied[2])
        self.assertEqual(
            sum(
                sources.canonical(s["url"]) == "https://mobbin.com"
                for s in result["sources"]
            ),
            1,
        )
        self.assertTrue(
            all(s["inspection"] == "not_examined" for s in result["sources"])
        )
        self.assertEqual(
            len(sources.select(catalog, supplied, only_user=True)["sources"]), 2
        )
        candidates = sources.select(catalog, [], include_candidates=True)["sources"]
        self.assertEqual(len(candidates), 7)
        self.assertTrue(
            any(s["verification"]["status"] == "unverified" for s in candidates)
        )
        with self.assertRaises(ValueError):
            sources.select(catalog, ["https://name:secret@example.com"])

    @unittest.skipUnless(
        shutil.which("ffmpeg") and shutil.which("ffprobe"),
        "media tools are not installed",
    )
    def test_real_media_contract_and_frame_extraction(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            media = root / "sample.mp4"
            video.run(
                [
                    "ffmpeg",
                    "-nostdin",
                    "-v",
                    "error",
                    "-f",
                    "lavfi",
                    "-i",
                    "color=c=blue:s=320x180:r=10",
                    "-t",
                    "1",
                    "-c:v",
                    "libx264",
                    "-pix_fmt",
                    "yuv420p",
                    media,
                ]
            )
            report = video.inspect(media, 320, 180, 10, 1)
            self.assertEqual(report["status"], "technical_checks_passed")
            self.assertEqual(report["visual_review"], "pending")
            self.assertEqual(report["motion_review"], "pending")
            for contract in [
                (640, 180, 10, 1, False),
                (320, 180, 30, 1, False),
                (320, 180, 10, 2, False),
                (320, 180, 10, 1, True),
            ]:
                with self.subTest(contract=contract), self.assertRaises(ValueError):
                    video.inspect(media, *contract)
            frames = video.frames(media, [0, 0.5], root / "frames")
            self.assertTrue(
                all(
                    Path(frame["file"]).read_bytes().startswith(b"\x89PNG")
                    for frame in frames["frames"]
                )
            )
            self.assertTrue(all(not frame["inspected"] for frame in frames["frames"]))
            for times in [[float("nan")], [float("inf")], [-1], [1]]:
                with self.subTest(times=times), self.assertRaises(ValueError):
                    video.frames(media, times, root / "invalid")
            self.assertFalse((root / "invalid").exists())

    def test_props_reject_missing_media_bad_timing_and_duplicate_shots(self):
        with tempfile.TemporaryDirectory() as directory:
            project = Path(directory)
            original = json.loads(
                (
                    SKILLS
                    / "video-production/assets/remotion-starter/props.json"
                ).read_text()
            )
            (project / "props.json").write_text(json.dumps(original))
            self.assertEqual(
                sum(s["durationFrames"] for s in video.project_props(project)["shots"]),
                180,
            )
            for field, value in [
                ("fps", 0),
                ("width", 321),
                ("audio", "../outside.wav"),
            ]:
                props = {**original, field: value}
                (project / "props.json").write_text(json.dumps(props))
                with self.subTest(field=field), self.assertRaises(ValueError):
                    video.project_props(project)
            original["shots"][1]["id"] = original["shots"][0]["id"]
            (project / "props.json").write_text(json.dumps(original))
            with self.assertRaises(ValueError):
                video.project_props(project)

    @unittest.skipUnless(
        all(shutil.which(tool) for tool in ("node", "ffmpeg", "ffprobe")),
        "local render prerequisites missing",
    )
    def test_failed_renderer_never_reports_old_output_as_success(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            project = root / "project"
            shutil.copytree(
                SKILLS / "video-production/assets/remotion-starter",
                project,
            )
            for package in ("remotion", "@remotion/cli", "react", "react-dom"):
                path = project / "node_modules" / package
                path.mkdir(parents=True)
                (path / "package.json").write_text('{"version":"1.0.0"}')
            cli = project / "node_modules/.bin/remotion"
            cli.parent.mkdir()
            cli.write_text(
                '#!/bin/sh\nprintf "fixture rendering failure" >&2\nexit 7\n'
            )
            cli.chmod(0o755)
            old = root / "old.mp4"
            old.write_bytes(b"old output")
            with self.assertRaises(ValueError):
                video.render(project, cli, root / "export")
            self.assertEqual(
                json.loads((root / "export/failure.json").read_text())["status"],
                "failed",
            )
            self.assertFalse((root / "export/inspection.json").exists())
            self.assertEqual(old.read_bytes(), b"old output")

    def test_temporary_cleanup_retains_deliverables_and_rejects_replaced_identity(self):
        with tempfile.TemporaryDirectory() as directory:
            ledger = Path(directory) / "ledger.json"
            result = subprocess.run(
                [
                    sys.executable,
                    str(SKILLS / "session-housekeeping/scripts/housekeeping.py"),
                    "temp",
                    "--ledger",
                    str(ledger),
                    "--session",
                    "fixture",
                    "--purpose",
                    "fixture scratch",
                ],
                capture_output=True,
                text=True,
                check=True,
            )
            path = Path(json.loads(result.stdout)["created"])
            try:
                data = housekeeping.load(ledger, "fixture")
                deliverable = path / "demo.mp4"
                deliverable.write_bytes(b"delivered")
                report = housekeeping.cleanup(
                    ledger, data, [deliverable], None, False, True
                )
                self.assertEqual(report["resources"][0]["status"], "retain")
                self.assertTrue(deliverable.exists())
                data["resources"][0]["identity"]["inode"] += 1
                self.assertIn(
                    "identity changed",
                    housekeeping.cleanup(ledger, data, [], None, False, True)[
                        "resources"
                    ][0]["reason"],
                )
                data = housekeeping.load(ledger, "fixture")
                self.assertEqual(
                    housekeeping.cleanup(ledger, data, [], None, False, False)[
                        "resources"
                    ][0]["status"],
                    "eligible",
                )
                self.assertTrue(path.exists())
                self.assertEqual(
                    housekeeping.cleanup(ledger, data, [], None, False, True)[
                        "resources"
                    ][0]["status"],
                    "removed",
                )
                self.assertFalse(path.exists())
                self.assertEqual(
                    housekeeping.cleanup(ledger, data, [], None, False, True)[
                        "resources"
                    ][0]["status"],
                    "already_removed",
                )
                with self.assertRaises(ValueError):
                    housekeeping.load(ledger, "other-session")
            finally:
                if path.exists():
                    shutil.rmtree(path)

    def test_temp_symlinks_and_nested_repositories_are_not_followed_or_removed(self):
        with (
            tempfile.TemporaryDirectory() as directory,
            tempfile.TemporaryDirectory(prefix="aialchemy-session-") as owned,
        ):
            path = Path(owned).resolve()
            ledger = Path(directory) / "ledger.json"
            external = Path(directory) / "keep"
            external.write_text("keep")
            item = {
                "kind": "temp",
                "path": str(path),
                "identity": housekeeping.identity(path),
            }
            (path / "link").symlink_to(external)
            self.assertIn("symlink", housekeeping.reason(item, ledger, [], None, False))
            (path / "link").unlink()
            (path / ".git").mkdir()
            self.assertIn(
                "nested repository", housekeeping.reason(item, ledger, [], None, False)
            )
            (path / ".git").rmdir()
            housekeeping.git(path, "init", "--bare", "history.git")
            self.assertIn(
                "bare repository", housekeeping.reason(item, ledger, [], None, False)
            )
            self.assertEqual(external.read_text(), "keep")

    def test_worktree_cleanup_checks_ownership_dirty_state_inactivity_and_ancestry(
        self,
    ):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            repo = root / "repo"
            repo.mkdir()
            housekeeping.git(repo, "init", "-b", "main")
            housekeeping.git(repo, "config", "user.name", "Fixture")
            housekeeping.git(repo, "config", "user.email", "fixture@example.com")
            (repo / "file").write_text("initial")
            (repo / ".gitignore").write_text("cache/\n")
            housekeeping.git(repo, "add", "file", ".gitignore")
            housekeeping.git(repo, "commit", "-m", "initial")
            worktree = root / "worktree"
            unrelated = root / "unrelated"
            housekeeping.git(repo, "worktree", "add", "-b", "task", str(worktree))
            housekeeping.git(repo, "worktree", "add", "-b", "other", str(unrelated))
            ledger = root / "ledger.json"
            data = housekeeping.load(ledger, "fixture")
            with self.assertRaises(ValueError):
                housekeeping.register(ledger, data, repo)
            housekeeping.register(ledger, data, worktree)
            item = data["resources"][0]
            self.assertIn(
                "inactivity", housekeeping.reason(item, ledger, [], "main", False)
            )
            with patch.object(Path, "cwd", return_value=worktree):
                self.assertIn(
                    "current working",
                    housekeeping.reason(item, ledger, [], "main", True),
                )
            (worktree / "cache").mkdir()
            (worktree / "cache/keep.txt").write_text("ignored local data")
            self.assertIn(
                "work remains", housekeeping.reason(item, ledger, [], "main", True)
            )
            (worktree / "cache/keep.txt").unlink()
            (worktree / "cache").rmdir()
            (worktree / "new").write_text("untracked")
            self.assertIn(
                "work remains", housekeeping.reason(item, ledger, [], "main", True)
            )
            housekeeping.git(worktree, "add", "new")
            housekeeping.git(worktree, "commit", "-m", "task")
            self.assertIn(
                "not reachable", housekeeping.reason(item, ledger, [], "main", True)
            )
            housekeeping.git(repo, "merge", "--ff-only", "task")
            self.assertIsNone(housekeeping.reason(item, ledger, [], "main", True))
            housekeeping.git(repo, "worktree", "lock", str(worktree))
            self.assertIn(
                "unlocked", housekeeping.reason(item, ledger, [], "main", True)
            )
            housekeeping.git(repo, "worktree", "unlock", str(worktree))
            report = housekeeping.cleanup(ledger, data, [], "main", True, True)
            self.assertEqual(report["resources"][0]["status"], "removed", report)
            self.assertFalse(worktree.exists())
            self.assertTrue(unrelated.exists())
            self.assertEqual(
                housekeeping.git(repo, "rev-parse", "task"),
                housekeeping.git(repo, "rev-parse", "main"),
            )


if __name__ == "__main__":
    unittest.main()
