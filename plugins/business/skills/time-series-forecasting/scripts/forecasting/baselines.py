"""Original stdlib reference forecasts, fit only on supplied training values."""

import math

METHODS = ("mean", "naive", "seasonal-naive", "drift")


def predict(values, horizon, method, season):
    if method == "mean":
        point = [sum(values) / len(values)] * horizon
    elif method == "naive":
        point = [values[-1]] * horizon
    elif method == "seasonal-naive":
        if len(values) < season or season < 1:
            raise ValueError("seasonal-naive lacks a full selected season")
        point = [values[-season + i % season] for i in range(horizon)]
    elif method == "drift":
        slope = (values[-1] - values[0]) / (len(values) - 1)
        point = [values[-1] + slope * (i + 1) for i in range(horizon)]
    else:
        raise ValueError("unknown statistical method")
    if any(not math.isfinite(v) for v in point):
        raise ValueError("invalid statistical output")
    return {
        "point": point,
        "statistic": "reference point rule: " + method,
        "quantiles": {},
    }
