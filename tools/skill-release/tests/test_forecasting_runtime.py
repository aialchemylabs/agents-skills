"""Selection, side-effect, workload and extracted-package tests; no model downloads."""

import copy
import importlib.util
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[3]
sys.path.insert(
    0, str(ROOT / "plugins/business/skills/time-series-forecasting/scripts")
)
from forecasting import runtime
from forecasting.data import inspect

spec = importlib.util.spec_from_file_location(
    "forecast_integration", Path(__file__).with_name("forecasting_integration.py")
)
helper = importlib.util.module_from_spec(spec)
spec.loader.exec_module(helper)


class ForecastingRuntime(unittest.TestCase):
    def test_inspection_proposal_and_denied_execution_have_no_model_side_effects(self):
        with tempfile.TemporaryDirectory() as folder:
            r = helper.request(folder, "naive")
            with (
                patch(
                    "subprocess.run",
                    side_effect=AssertionError("subprocess during inspection"),
                ),
                patch(
                    "urllib.request.urlopen",
                    side_effect=AssertionError("network during inspection"),
                ),
            ):
                summary = inspect(r["data"])
                self.assertEqual(len(summary["series"]), 2)
                p = runtime.propose(r)
                with self.assertRaisesRegex(ValueError, "explicit selected"):
                    runtime.execute(p, "")
                with self.assertRaisesRegex(ValueError, "explicit selected"):
                    runtime.execute(p, "forged")
                changed = copy.deepcopy(p)
                changed["scope"]["request"]["context"] = 64
                with self.assertRaisesRegex(ValueError, "explicit selected"):
                    runtime.execute(changed, p["fingerprint"])
            self.assertNotIn("torch", sys.modules)
            self.assertNotIn("numpy", sys.modules)
            self.assertFalse(Path(r["output"]).exists())

    def test_snapshot_code_lock_and_option_drift_are_detected(self):
        with tempfile.TemporaryDirectory() as folder:
            r = helper.request(folder, "naive")
            p = runtime.propose(r)
            self.assertEqual(runtime.verify(p, p["fingerprint"]), r)
            path = Path(r["data"]["path"])
            path.write_text(path.read_text().replace("-4.0", "-8.0", 1))
            with self.assertRaisesRegex(ValueError, "drift"):
                runtime.verify(p, p["fingerprint"])
            r = helper.request(folder, "naive")
            p = runtime.propose(r)
            with patch.object(runtime, "digest", return_value="bad"):
                with self.assertRaisesRegex(ValueError, "lock digest"):
                    runtime.verify(p, p["fingerprint"])

    def test_scope_rejects_bounds_unlike_units_and_unsupported_options(self):
        with tempfile.TemporaryDirectory() as folder:
            base = helper.request(folder, "naive")
            mutations = [
                (lambda r: r["budgets"].update(max_horizon=7), "horizon budget"),
                (lambda r: r["settings"].update(clipping="nonnegative"), "settings"),
                (lambda r: r.update(semantics="stock"), "stocks/ratios"),
                (lambda r: r.update(cross_learning=True), "only Chronos"),
                (
                    lambda r: r.update(primary_loss="mae")
                    or r["data"]["units"].update({"reactor-b": "seconds"}),
                    "unlike units|totals require",
                ),
                (lambda r: r["origins"].append("2026-06-01"), "origins"),
                (lambda r: r.update(model="unapproved"), "unsupported model"),
                (
                    lambda r: r.update(output=str(ROOT / "private-report")),
                    "outside source",
                ),
            ]
            for mutation, error in mutations:
                r = copy.deepcopy(base)
                mutation(r)
                with (
                    self.subTest(error=error),
                    self.assertRaisesRegex(ValueError, error),
                ):
                    runtime.propose(r)

    def test_retargeted_runtime_destinations_are_rejected_before_writes(self):
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder)
            for name in ("output", "environment", "cache"):
                with self.subTest(destination=name):
                    r = helper.request(root / name, "naive")
                    original = root / (name + "-original")
                    other = root / (name + "-other")
                    original.mkdir()
                    other.mkdir()
                    link = root / (name + "-link")
                    link.symlink_to(original, target_is_directory=True)
                    r[name] = str(link)
                    p = runtime.propose(r)
                    link.unlink()
                    link.symlink_to(other, target_is_directory=True)
                    with self.assertRaisesRegex(ValueError, "drift"):
                        runtime.execute(p, p["fingerprint"])
                    self.assertEqual(list(other.iterdir()), [])
                    self.assertEqual(list(original.iterdir()), [])
                    link.unlink()

    def test_checksum_offline_and_planned_download_budget(self):
        with tempfile.TemporaryDirectory() as folder:
            model = {
                "repo": "fictional/model",
                "revision": "a" * 40,
                "files": [{"name": "config.json", "bytes": 2, "sha256": "a" * 64}],
            }
            with self.assertRaisesRegex(ValueError, "offline cache"):
                runtime.fetch_files(model, folder, True)
            dest = Path(folder) / "fictional--model" / ("a" * 40)
            dest.mkdir(parents=True)
            (dest / "config.json").write_text("{}")
            with self.assertRaisesRegex(ValueError, "checksum"):
                runtime.fetch_files(model, folder, True)
            self.assertEqual((dest / "config.json").read_text(), "{}")
            r = helper.request(Path(folder) / "task", "timesfm")
            r["budgets"]["planned_download_bytes"] = 1
            with (
                patch("platform.system", return_value="Darwin"),
                patch("platform.machine", return_value="arm64"),
            ):
                with self.assertRaisesRegex(ValueError, "footprint"):
                    runtime.propose(r)
            with patch("platform.system", return_value="Linux"):
                with self.assertRaisesRegex(ValueError, "unverified|validated"):
                    runtime.propose(r)

    def test_installed_copy_runs_from_unrelated_cwd_and_sums_points_honestly(self):
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder)
            installed = root / "installed"
            shutil.copytree(
                runtime.SKILL,
                installed,
                ignore=shutil.ignore_patterns("__pycache__", "*.pyc"),
            )
            r = helper.request(root / "task", "naive")
            request = root / "request.json"
            request.write_text(json.dumps(r))
            cli = installed / "scripts/forecast.py"
            before = {str(p.relative_to(installed)): p.read_bytes() for p in installed.rglob("*") if p.is_file()}
            subprocess.run(
                [sys.executable, str(cli), "inspect", str(request)],
                cwd=root, check=True, capture_output=True,
                env={k: v for k, v in os.environ.items() if k not in ("PYTHONDONTWRITEBYTECODE", "PYTHONPYCACHEPREFIX")},
            )
            self.assertEqual(before, {str(p.relative_to(installed)): p.read_bytes() for p in installed.rglob("*") if p.is_file()})
            proposal = root / "proposal.json"
            subprocess.run(
                [
                    sys.executable,
                    str(cli),
                    "propose",
                    str(request),
                    "--output",
                    str(proposal),
                ],
                cwd=root,
                check=True,
                capture_output=True,
            )
            selected = json.loads(proposal.read_text())["fingerprint"]
            subprocess.run(
                [
                    sys.executable,
                    str(cli),
                    "run",
                    str(proposal),
                    "--selected",
                    selected,
                ],
                cwd=root,
                check=True,
                capture_output=True,
            )
            status = json.loads((Path(r["output"]) / "status.json").read_text())
            self.assertEqual(status["status"], "complete")
            self.assertIn("not an aggregate", status["period_sum_label"])
            self.assertEqual(
                status["comparison"]["aggregate"]["common_series_folds"], 2
            )
            self.assertTrue(
                (Path(r["output"]) / "series-0.svg").read_text().startswith("<svg")
            )
            self.assertFalse(any(installed.rglob("model.safetensors")))
            folds = json.loads((Path(r["output"]) / "metrics.json").read_text())[
                "folds"
            ]
            self.assertTrue(
                all(f["period_total"]["absolute_error"] is not None for f in folds)
            )
            # A short early fold cannot quietly disappear from the comparison.
            r["origins"][0] = "2026-01-01"
            r["aggregation"] = "direct-total"
            r["output"] = str(root / "direct-output")
            request.write_text(json.dumps(r))
            subprocess.run(
                [
                    sys.executable,
                    str(cli),
                    "propose",
                    str(request),
                    "--output",
                    str(root / "direct.json"),
                ],
                cwd=root,
                check=True,
                capture_output=True,
            )
            selected = json.loads((root / "direct.json").read_text())["fingerprint"]
            subprocess.run(
                [
                    sys.executable,
                    str(cli),
                    "run",
                    str(root / "direct.json"),
                    "--selected",
                    selected,
                ],
                cwd=root,
                check=True,
                capture_output=True,
            )
            status = json.loads((Path(r["output"]) / "status.json").read_text())
            self.assertEqual(status["comparison"]["common_series_folds"], 2)
            self.assertTrue(status["comparison"]["excluded_folds"])
            self.assertIn("__total__", status["point_statistics"])
            r["aggregation"] = "none"
            r["output"] = str(root / "no-totals-output")
            p = runtime.propose(r)
            status = runtime.execute(p, p["fingerprint"])
            self.assertIsNone(status["period_point_sums"])
            self.assertIsNone(status["period_sum_label"])
            self.assertNotIn("aggregate", status["comparison"])
            folds = json.loads((Path(r["output"]) / "metrics.json").read_text())["folds"]
            self.assertTrue(all("period_total" not in f for f in folds))


if __name__ == "__main__":
    unittest.main()
