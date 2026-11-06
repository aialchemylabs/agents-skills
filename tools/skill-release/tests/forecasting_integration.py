#!/usr/bin/env python3
"""Explicit real-weight integration; never included in test*.py discovery.

Run only under user-selected checkpoint/setup scope. Outputs are synthetic.
"""

import argparse
import csv
import datetime as dt
import json
import math
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[3]
SCRIPTS = ROOT / "plugins/business/skills/time-series-forecasting/scripts"
sys.path.insert(0, str(SCRIPTS))
from forecasting.runtime import propose


def request(workspace, model, offline=True, cross_learning=False):
    workspace = Path(workspace)
    workspace.mkdir(parents=True, exist_ok=True)
    data = workspace / "synthetic.csv"
    start = dt.date(2026, 1, 1)
    dates = [(start + dt.timedelta(days=i)).isoformat() for i in range(168)]
    with data.open("w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["id", "date", "value", "calendar_signal"])
        for key, offset in [("reactor-a", -4), ("reactor-b", 6)]:
            for i, date in enumerate(dates):
                writer.writerow(
                    [
                        key,
                        date,
                        offset + 3 * math.sin(i * 2 * math.pi / 7) + i * 0.01,
                        math.sin(i * 2 * math.pi / 7),
                    ]
                )
    home = Path.home()
    r = {
        "data": {
            "path": str(data),
            "id": "id",
            "timestamp": "date",
            "target": "value",
            "timezone": "UTC",
            "frequency": "daily",
            "as_of": dates[-1],
            "units": {"reactor-a": "units", "reactor-b": "units"},
        },
        "model": model,
        "report": {
            "start": (start + dt.timedelta(days=170)).isoformat(),
            "end": (start + dt.timedelta(days=175)).isoformat(),
            "calendar_basis": "calendar",
        },
        "origins": [dates[143], dates[151]],
        "context": 128,
        "season": 7,
        "baselines": ["mean", "naive", "seasonal-naive", "drift"],
        "primary_loss": "mase",
        "pooling": "weighted-series-mean",
        "weights": {"reactor-a": 1, "reactor-b": 1},
        "semantics": "flow",
        "aggregation": "sum-point-rules",
        "covariates": [],
        "groups": [["reactor-a"], ["reactor-b"]],
        "cross_learning": cross_learning,
        "settings": {
            "device": "cpu",
            "precision": "float32",
            "compile": False,
            "clipping": "none",
            "quantile_crossing": "clamp-about-median",
            "truncation": "last-context",
            "horizon_extension": "reject",
        },
        "budgets": {
            "seconds": 900,
            "threads": 4,
            "batch_size": 8,
            "max_series": 2,
            "max_origins": 2,
            "max_horizon": 8,
            "planned_download_bytes": 10 * 1024**3,
            "retries": 0,
        },
        "output": str(workspace / "output"),
        "environment": str(
            home / "Library/Caches/aialchemy-forecasting/envs/macos-arm64-cpu-v1"
        ),
        "cache": str(home / ".cache/huggingface/aialchemy-forecasting"),
        "offline": offline,
        "download_estimate_bytes": 2_000_000_000,
    }
    if model == "chronos":
        r["covariates"] = [
            {
                "name": "calendar_signal",
                "kind": "known",
                "values": [
                    {
                        "id": k,
                        "date": (start + dt.timedelta(days=i)).isoformat(),
                        "value": math.sin(i * 2 * math.pi / 7),
                        "available_at": dates[0],
                    }
                    for k in ("reactor-a", "reactor-b")
                    for i in range(176)
                ],
            }
        ]
    if cross_learning:
        r["groups"] = [["reactor-a", "reactor-b"]]
    return r


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--workspace", required=True, type=Path)
    p.add_argument("--model", choices=["timesfm", "chronos", "naive"], required=True)
    p.add_argument("--network-deny", action="store_true")
    p.add_argument("--cross-learning", action="store_true")
    a = p.parse_args()
    r = request(a.workspace, a.model, cross_learning=a.cross_learning)
    proposal = propose(r)
    path = a.workspace / "selected-proposal.json"
    path.write_text(json.dumps(proposal, indent=2) + "\n")
    command = [
        sys.executable,
        str(SCRIPTS / "forecast.py"),
        "run",
        str(path),
        "--selected",
        proposal["fingerprint"],
    ]
    if a.network_deny:
        if sys.platform != "darwin":
            raise RuntimeError(
                "network-deny integration currently requires macOS sandbox-exec"
            )
        command = [
            "/usr/bin/sandbox-exec",
            "-p",
            "(version 1)(allow default)(deny network*)",
            *command,
        ]
    with (a.workspace / "transcript.log").open("w") as f:
        subprocess.run(command, check=True, stdout=f, stderr=f, timeout=930)
    status = json.loads((Path(r["output"]) / "status.json").read_text())
    folds = json.loads((Path(r["output"]) / "metrics.json").read_text())["folds"]
    assert status["status"] == "complete"
    assert all(f["status"] == "ok" for f in folds), folds
    assert status["comparison"]["common_series_folds"] == 4
    if a.model in ("timesfm", "chronos"):
        assert status["point_statistics"]["reactor-a"] == "median"
        assert all(status["model_semantics"][:2])
    print(
        json.dumps(
            {
                "model": a.model,
                "network_denied": a.network_deny,
                "cross_learning": a.cross_learning,
                "runtime_seconds": status["runtime_seconds"],
                "peak_rss_bytes": status["peak_rss_bytes"],
                "common_series_folds": status["comparison"]["common_series_folds"],
                "fingerprint": proposal["fingerprint"],
            },
            indent=2,
        )
    )


if __name__ == "__main__":
    main()
