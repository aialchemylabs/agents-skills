"""Weight-free observable calendar, leakage, metric and forecast invariants."""

import datetime as dt
import importlib.util
import json
from pathlib import Path
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[3]
sys.path.insert(
    0, str(ROOT / "plugins/business/skills/time-series-forecasting/scripts")
)
from forecasting.data import dates_after, stamp, load, prepare, covariates
from forecasting.baselines import predict
from forecasting.evaluation import metrics, compare
from forecasting.engine import valid


class ForecastingCore(unittest.TestCase):
    def test_calendar_bridge_leap_month_week_and_dst(self):
        dates = dates_after("2027-02-28", "2027-06-30", "daily", "UTC")
        self.assertEqual((len(dates), sum(d >= "2027-04-01" for d in dates)), (122, 91))
        self.assertEqual(
            dates_after("2028-02-28", "2028-03-01", "daily", "UTC"),
            ["2028-02-29", "2028-03-01"],
        )
        self.assertEqual(
            dates_after("2028-01-31", "2028-03-31", "monthly", "UTC"),
            ["2028-02-29", "2028-03-31"],
        )
        self.assertEqual(
            dates_after("2027-01-01", "2027-01-15", "weekly", "UTC"),
            ["2027-01-08", "2027-01-15"],
        )
        hours = dates_after(
            "2026-11-01T00:00:00-04:00",
            "2026-11-01T03:00:00-05:00",
            "hourly",
            "America/New_York",
        )
        self.assertEqual(len(hours), 4)
        self.assertEqual(
            [h[11:16] for h in hours], ["01:00", "01:00", "02:00", "03:00"]
        )
        self.assertEqual(
            dates_after("2026-01-01 01:00:00+00:00", "2026-01-01 02:00:00+00:00", "hourly", "UTC"),
            ["2026-01-01T02:00:00+00:00"],
        )
        with self.assertRaisesRegex(ValueError, "hourly.*time"):
            dates_after("2026-01-01", "2026-01-02", "hourly", "UTC")
        with self.assertRaisesRegex(ValueError, "ambiguous"):
            stamp("2026-11-01T01:30:00", "America/New_York")
        with self.assertRaisesRegex(ValueError, "off cadence"):
            dates_after("2027-01-01", "2027-01-10", "weekly", "UTC")

    def test_missing_zeros_negative_and_fold_local_leakage(self):
        rows = [
            {
                "when": stamp(f"2027-01-0{i + 1}"),
                "available": stamp(f"2027-01-0{i + 1}"),
                "value": v,
            }
            for i, v in enumerate([-2, None, 0, 999999])
        ]
        selected, values = prepare(rows, stamp("2027-01-03"), 128)
        self.assertEqual(values, [-2, -1, 0])
        rows[-1]["value"] = -999999
        self.assertEqual(prepare(rows, stamp("2027-01-03"), 128)[1], values)
        rows[2]["available"] = stamp("2027-01-04")
        with self.assertRaisesRegex(ValueError, "fewer than two"):
            prepare(rows, stamp("2027-01-03"), 128)

    def test_publication_delays_preserve_grid_and_hide_unavailable_values(self):
        rows = [
            {"when": stamp(f"2027-01-0{i}"), "available": stamp(f"2027-01-0{i}"),
             "value": float(i), "raw": {"sensor": float(i * 10)}}
            for i in range(1, 7)
        ]
        cutoff = stamp("2027-01-05")
        peer, _ = prepare(rows, cutoff, 5)
        for i in (1, 3, 4):
            rows[i]["available"] = stamp("2027-01-06")
        training, values = prepare(rows, cutoff, 5)
        self.assertEqual([r["when"] for r in training], [r["when"] for r in peer])
        self.assertEqual(values, [1, 2, 3, 2, 2])
        self.assertEqual([r["value"] for r in training], [1, None, 3, None, None])
        self.assertEqual(predict(values, 1, "seasonal-naive", 3)["point"], [3])
        request = {"covariates": [{"name": "sensor", "kind": "past"}]}
        self.assertEqual(covariates(request, "a", training, [], cutoff), ({"sensor": [10, 20, 30, 20, 20]}, {}))
        for i in (1, 3, 4):
            rows[i]["value"] = 999999
            rows[i]["raw"]["sensor"] = 999999
        updated, after = prepare(rows, cutoff, 5)
        self.assertEqual(after, values)
        self.assertEqual(covariates(request, "a", updated, [], cutoff), covariates(request, "a", training, [], cutoff))
        with self.assertRaisesRegex(ValueError, "fewer than two"):
            prepare(rows, cutoff, 3)

    def test_data_rejects_irregular_duplicates_incomplete_and_unaligned(self):
        with tempfile.TemporaryDirectory() as folder:
            path = Path(folder) / "input.json"
            spec = {
                "path": str(path),
                "id": "id",
                "timestamp": "date",
                "target": "v",
                "frequency": "daily",
                "timezone": "UTC",
                "as_of": "2027-02-01",
            }
            cases = [
                (
                    [
                        {"id": "a", "date": "2027-01-01", "v": 1},
                        {"id": "a", "date": "2027-01-03", "v": 0},
                    ],
                    "irregular",
                ),
                ([{"id": "a", "date": "2027-01-01", "v": 1}] * 2, "duplicate"),
                (
                    [{"id": "a", "date": "2027-01-01", "v": 1, "complete": False}],
                    "incomplete",
                ),
                (
                    [
                        {"id": "a", "date": "2027-01-01", "v": 1},
                        {"id": "b", "date": "2027-01-02", "v": 1},
                    ],
                    "aligned",
                ),
            ]
            for rows, error in cases:
                path.write_text(json.dumps(rows))
                with (
                    self.subTest(error=error),
                    self.assertRaisesRegex(ValueError, error),
                ):
                    load(spec)
            path.write_text(
                json.dumps(
                    [
                        {"id": "a", "date": "2027-01-01", "v": 0},
                        {"id": "a", "date": "2027-01-02", "v": None},
                        {"id": "a", "date": "2027-01-03", "v": -1},
                    ]
                )
            )
            self.assertEqual([v["value"] for v in load(spec)["a"]], [0, None, -1])

    def test_origin_visible_covariate_vintages(self):
        r = {
            "data": {"timezone": "UTC"},
            "covariates": [
                {
                    "name": "temperature",
                    "kind": "vintage",
                    "values": [
                        {
                            "id": "a",
                            "date": "2027-01-04",
                            "value": 10,
                            "available_at": "2027-01-02",
                        },
                        {
                            "id": "a",
                            "date": "2027-01-04",
                            "value": 999,
                            "available_at": "2027-01-04",
                        },
                    ],
                }
            ],
        }
        history = [{"raw": {"temperature": 4}}]
        self.assertEqual(
            covariates(r, "a", history, ["2027-01-04"], stamp("2027-01-03")),
            ({"temperature": [4.0]}, {"temperature": [10.0]}),
        )
        r["covariates"][0]["values"].pop(0)
        with self.assertRaisesRegex(ValueError, "origin-visible"):
            covariates(r, "a", history, ["2027-01-04"], stamp("2027-01-03"))

    def test_baseline_and_hand_calculated_metrics(self):
        self.assertEqual(predict([0, 2, 4], 2, "drift", 1)["point"], [6, 8])
        self.assertEqual(
            predict([1, 2, 3, 4], 3, "seasonal-naive", 2)["point"], [3, 4, 3]
        )
        m = metrics(
            [0, 2, None],
            [1, 0, 999],
            [0, 2, 4],
            1,
            {"0.1": [-1, 1, 0], "0.9": [1, 3, 0]},
        )
        self.assertEqual(m["mae"], 1.5)
        self.assertEqual(m["mase"], 0.75)
        self.assertEqual(m["interval"]["covered"], 2)
        self.assertEqual(m["interval"]["n"], 2)
        self.assertIsNone(metrics([0], [1], [0, 0], 1)["mase"])
        self.assertEqual(metrics([None], [1], [0, 0], 1)["n"], 0)

    def test_common_folds_prevent_easy_fold_advantage_and_keep_losses_ties(self):
        folds = [
            {
                "origin": str(o),
                "series": "a",
                "method": m,
                "status": "ok",
                "metrics": {"mae": loss, "n": 2},
            }
            for o, m, loss in [
                (1, "model", 1),
                (1, "naive", 2),
                (2, "model", 100),
                (2, "naive", 0),
            ]
        ]
        folds[2] = {"origin": "2", "series": "a", "method": "model", "status": "failed"}
        c = compare(folds, ["model", "naive"], {"a": 1}, "mae")
        self.assertEqual(c["common_series_folds"], 1)
        self.assertEqual(c["methods"]["naive"]["pooled_primary"], 2)
        for model_loss, baseline_loss in [(3, 1), (1, 1), (0, 1)]:
            folds[0]["metrics"]["mae"] = model_loss
            folds[1]["metrics"]["mae"] = baseline_loss
            c = compare(folds, ["model", "naive"], {"a": 1}, "mae")
            self.assertEqual(c["methods"]["model"]["pooled_primary"], model_loss)

    def test_quantiles_negative_values_and_declared_correction(self):
        res = {
            "point": [-3.0],
            "statistic": "median",
            "quantiles": {"0.1": [-2.0], "0.5": [-3.0], "0.9": [-4.0]},
        }
        with self.assertRaisesRegex(ValueError, "crossed"):
            valid(res, 1)
        corrected = valid(res, 1, "clamp-about-median")
        self.assertEqual(corrected["point"], [-3.0])
        self.assertEqual(
            corrected["quantiles"], {"0.1": [-3.0], "0.5": [-3.0], "0.9": [-3.0]}
        )
        with self.assertRaisesRegex(ValueError, "invalid"):
            valid({"point": [float("nan")], "quantiles": {}}, 1)


if __name__ == "__main__":
    unittest.main()
