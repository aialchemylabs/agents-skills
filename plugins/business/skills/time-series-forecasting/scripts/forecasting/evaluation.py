"""Common-observation evaluation, training-only scales and uncalibrated intervals."""

import math


def metrics(actual, predicted, training, season, quantiles=None):
    indices = [i for i, y in enumerate(actual) if y is not None]
    if not indices:
        return {"n": 0, "mae": None, "rmse": None, "mase": None}
    errors = [predicted[i] - actual[i] for i in indices]
    mae = sum(abs(e) for e in errors) / len(errors)
    differences = [
        abs(training[i] - training[i - season]) for i in range(season, len(training))
    ]
    scale = sum(differences) / len(differences) if differences else 0
    result = {
        "n": len(indices),
        "mae": mae,
        "rmse": math.sqrt(sum(e * e for e in errors) / len(errors)),
        "mase": mae / scale if scale else None,
        "training_scale": scale or None,
    }
    if quantiles:
        result["pinball"] = {
            q: sum(
                max(
                    float(q) * (actual[i] - vs[i]), (float(q) - 1) * (actual[i] - vs[i])
                )
                for i in indices
            )
            / len(indices)
            for q, vs in quantiles.items()
        }
        if "0.1" in quantiles and "0.9" in quantiles:
            count = sum(
                quantiles["0.1"][i] <= actual[i] <= quantiles["0.9"][i] for i in indices
            )
            result["interval"] = {
                "nominal": 0.8,
                "covered": count,
                "n": len(indices),
                "coverage": count / len(indices),
                "label": "native nominal, not calibrated",
            }
    return result


def compare(folds, methods, weights, primary):
    keys = [
        {
            (f["origin"], f["series"])
            for f in folds
            if f["method"] == m and f["status"] == "ok"
        }
        for m in methods
    ]
    common = set.intersection(*keys) if keys else set()
    summaries = {}
    for method in methods:
        per_series = {}
        for key in weights:
            rows = [
                f
                for f in folds
                if f["method"] == method
                and f["series"] == key
                and (f["origin"], key) in common
                and f["metrics"][primary] is not None
            ]
            n = sum(f["metrics"]["n"] for f in rows)
            per_series[key] = {
                "n": n,
                primary: sum(f["metrics"][primary] * f["metrics"]["n"] for f in rows)
                / n
                if n
                else None,
            }
        valid = all(v[primary] is not None for v in per_series.values())
        summaries[method] = {
            "series": per_series,
            "pooled_primary": sum(
                weights[k] * v[primary] for k, v in per_series.items()
            )
            / sum(weights.values())
            if valid
            else None,
        }
    scores = {
        m: v["pooled_primary"]
        for m, v in summaries.items()
        if v["pooled_primary"] is not None
    }
    best = [m for m, v in scores.items() if v == min(scores.values())] if scores else []
    return {
        "best_on_primary": best,
        "common_series_folds": len(common),
        "methods": summaries,
        "excluded_folds": [
            {k: f[k] for k in ("origin", "series", "method", "status")}
            for f in folds
            if f["status"] != "ok" or (f["origin"], f["series"]) not in common
        ],
        "interpretation": "Exploratory rolling-origin comparison; baseline superiority is a completed finding.",
    }
