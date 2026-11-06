"""Chronological evaluation shared by statistical and optional foundation paths."""

import importlib
import math
import resource
import sys
import time
from . import baselines, evaluation, report
from .data import prepare, stamp, covariates, load
from .runtime import verify, fetch_files


def valid(result, horizon, policy="reject"):
    if len(result["point"]) != horizon or any(
        not math.isfinite(v) for v in result["point"]
    ):
        raise ValueError("invalid point shape/values")
    qs = result["quantiles"]
    for vs in qs.values():
        if len(vs) != horizon or any(not math.isfinite(v) for v in vs):
            raise ValueError("invalid quantile output")
    ordered = sorted(qs, key=float)
    if policy == "clamp-about-median" and qs:
        mid = ordered.index("0.5")
        for i in range(horizon):
            for j in range(mid - 1, -1, -1):
                qs[ordered[j]][i] = min(qs[ordered[j]][i], qs[ordered[j + 1]][i])
            for j in range(mid + 1, len(ordered)):
                qs[ordered[j]][i] = max(qs[ordered[j]][i], qs[ordered[j - 1]][i])
        result["quantile_correction"] = "clamp-about-median"
    if any(
        qs[a][i] > qs[b][i]
        for a, b in zip(ordered, ordered[1:])
        for i in range(horizon)
    ):
        raise ValueError("crossed quantiles rejected; no undeclared correction")
    return result


def run(proposal, selected):
    start = time.monotonic()
    r = verify(proposal, selected)
    groups = load(r["data"])
    dates = proposal["scope"]["generation_dates"]
    mask = proposal["scope"]["scoring_indices"]
    h = len(dates)
    methods = list(dict.fromkeys([r["model"], *r["baselines"]]))
    if r["aggregation"] == "direct-total":
        if (
            r["covariates"]
            or r["cross_learning"]
            or len(set(r["data"]["units"].values())) != 1
        ):
            raise ValueError(
                "direct aggregate requires same-unit flow history and no covariates/cross-learning"
            )
        keys = list(groups)
        groups["__total__"] = [
            {
                **row,
                "value": sum(groups[k][i]["value"] for k in keys)
                if all(groups[k][i]["value"] is not None for k in keys)
                else None,
                "available": max(groups[k][i]["available"] for k in keys),
            }
            for i, row in enumerate(groups[keys[0]])
        ]
    adapter = None
    if proposal["scope"]["checkpoint"]:
        checkpoint = fetch_files(proposal["scope"]["checkpoint"], r["cache"], True)
        adapter = importlib.import_module("." + r["model"], __package__).Adapter(
            checkpoint, r
        )

    def predictions(cutoff, future):
        try:
            prepared = {
                k: prepare(rows, cutoff, r["context"]) for k, rows in groups.items()
            }
            covs = {
                k: covariates(r, k, training, future, cutoff)
                for k, (training, _) in prepared.items()
            }
        except ValueError as exc:
            return {}, {m: {} for m in methods}, {m: str(exc) for m in methods}
        outputs = {m: {} for m in methods}
        failures = {}
        for m in methods:
            try:
                if m == r["model"] and adapter:
                    batches = r["groups"] + (
                        [["__total__"]] if "__total__" in groups else []
                    )
                    for batch in batches:
                        result = adapter.predict(
                            [prepared[k][1] for k in batch], h, [covs[k] for k in batch]
                        )
                        if len(result) != len(batch):
                            raise ValueError("adapter lost series")
                        for key, res in zip(batch, result):
                            outputs[m][key] = valid(
                                res, h, r["settings"]["quantile_crossing"]
                            )
                else:
                    for key, (_, values) in prepared.items():
                        outputs[m][key] = valid(
                            baselines.predict(values, h, m, r["season"]), h
                        )
            except (ValueError, RuntimeError) as exc:
                failures[m] = type(exc).__name__ + ": " + str(exc)
                outputs[m] = {}
        return prepared, outputs, failures

    def period_sum(actual, result):
        complete = all(v is not None for v in actual)
        observed = sum(actual) if complete else None
        predicted = sum(result)
        return {
            "scored_periods": sum(v is not None for v in actual),
            "expected_periods": len(actual),
            "actual_total": observed,
            "point_sum": predicted,
            "absolute_error": abs(predicted - observed) if complete else None,
            "label": "sum of point rules, not an aggregate mean/median; interval unavailable",
        }

    folds = []
    tz = r["data"]["timezone"]
    rows = next(iter(groups.values()))
    for origin in sorted(r["origins"], key=lambda v: stamp(v, tz)):
        cutoff = stamp(origin, tz)
        idx = next(i for i, row in enumerate(rows) if row["when"] == cutoff)
        future = [v["date"] for v in rows[idx + 1 : idx + h + 1]]
        prepared, outputs, failures = predictions(cutoff, future)
        for m in methods:
            for key in groups:
                record = {
                    "origin": origin,
                    "series": key,
                    "method": m,
                    "status": "failed" if m in failures else "ok",
                }
                if m in failures:
                    record["error"] = failures[m]
                else:
                    actual = [groups[key][idx + 1 + i]["value"] for i in mask]
                    result = outputs[m][key]
                    record["metrics"] = evaluation.metrics(
                        actual,
                        [result["point"][i] for i in mask],
                        prepared[key][1],
                        r["season"],
                        {
                            q: [vs[i] for i in mask]
                            for q, vs in result["quantiles"].items()
                        },
                    )
                    if r["semantics"] == "flow" and r["aggregation"] != "none":
                        record["period_total"] = period_sum(
                            actual, [result["point"][i] for i in mask]
                        )
                    if not record["metrics"]["n"]:
                        record["status"] = "missing"
                folds.append(record)
        # Sum-of-point rules is evaluated as its own total, never as median/mean.
        if r["aggregation"] == "sum-point-rules":
            if len(set(r["data"]["units"].values())) != 1:
                raise ValueError("cannot total unlike units")
            for m in methods:
                record = {
                    "origin": origin,
                    "series": "__total__",
                    "method": m,
                    "status": "failed" if m in failures else "ok",
                }
                if m not in failures:
                    actual = [
                        sum(groups[k][idx + 1 + i]["value"] for k in groups)
                        if all(
                            groups[k][idx + 1 + i]["value"] is not None for k in groups
                        )
                        else None
                        for i in mask
                    ]
                    aggregate_rows = [
                        {
                            **row,
                            "value": sum(groups[k][i]["value"] for k in groups)
                            if all(groups[k][i]["value"] is not None for k in groups)
                            else None,
                            "available": max(groups[k][i]["available"] for k in groups),
                        }
                        for i, row in enumerate(rows)
                    ]
                    _, train = prepare(aggregate_rows, cutoff, r["context"])
                    points = [
                        sum(outputs[m][k]["point"][i] for k in groups) for i in mask
                    ]
                    record["metrics"] = evaluation.metrics(
                        actual, points, train, r["season"]
                    )
                    record["period_total"] = period_sum(actual, points)
                    if not record["metrics"]["n"]:
                        record["status"] = "missing"
                folds.append(record)
    comparison = evaluation.compare(
        [f for f in folds if f["series"] != "__total__"],
        methods,
        r["weights"],
        r["primary_loss"],
    )
    if r["aggregation"] != "none":
        comparison["aggregate"] = evaluation.compare(
            [f for f in folds if f["series"] == "__total__"],
            methods,
            {"__total__": 1},
            r["primary_loss"],
        )
    prepared, outputs, failures = predictions(stamp(r["data"]["as_of"], tz), dates)
    if r["model"] in failures:
        raise ValueError("final selected forecast failed: " + failures[r["model"]])
    forecasts = outputs[r["model"]]
    for key, result in forecasts.items():
        result["history"] = prepared[key][1]
    if r["aggregation"] == "sum-point-rules":
        forecasts["__total__"] = {
            "point": [sum(v["point"][i] for v in forecasts.values()) for i in range(h)],
            "quantiles": {},
            "statistic": "sum of component point rules; not aggregate mean/median",
            "history": prepare(
                [
                    {
                        **row,
                        "value": sum(groups[k][i]["value"] for k in groups)
                        if all(groups[k][i]["value"] is not None for k in groups)
                        else None,
                        "available": max(groups[k][i]["available"] for k in groups),
                    }
                    for i, row in enumerate(rows)
                ],
                stamp(r["data"]["as_of"], tz),
                r["context"],
            )[1],
        }
    peak = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss * (
        1 if sys.platform == "darwin" else 1024
    )
    return report.write(
        r["output"],
        proposal,
        forecasts,
        folds,
        comparison,
        time.monotonic() - start,
        peak,
    )
