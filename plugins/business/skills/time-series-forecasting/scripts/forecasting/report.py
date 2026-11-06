"""Self-contained SVG and JSON/CSV reports without plotting dependencies."""

import csv
import html
import json
from pathlib import Path


def chart(history, forecast, label, destination):
    vals = [v for v in history + forecast if v is not None]
    low, high = min(vals), max(vals)
    span = high - low or 1
    n = len(history) + len(forecast)

    def line(values, offset, color):
        pts = " ".join(
            f"{40 + 700 * (i + offset) / max(1, n - 1):.2f},{210 - 170 * (v - low) / span:.2f}"
            for i, v in enumerate(values)
            if v is not None
        )
        return (
            f'<polyline points="{pts}" fill="none" stroke="{color}" stroke-width="2"/>'
        )

    svg = f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 780 260"><rect width="780" height="260" fill="white"/><text x="40" y="24">{html.escape(label)}</text><text x="40" y="245">History: blue; forecast: orange. Full generation window, including bridge.</text>'
    svg += (
        line(history, 0, "#2563eb") + line(forecast, len(history), "#c2410c") + "</svg>"
    )
    destination.write_text(svg)


def write(output, proposal, forecasts, folds, comparison, elapsed, peak_rss):
    output = Path(output)
    r = proposal["scope"]["request"]
    mask = proposal["scope"]["scoring_indices"]
    fields = [
        "series",
        "date",
        "in_reporting_window",
        "unit",
        "point",
        "statistic",
        *[f"q{q}" for q in range(10, 100, 10)],
    ]
    with (output / "forecast.csv").open("w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=fields)
        writer.writeheader()
        for key, result in forecasts.items():
            for i, date in enumerate(proposal["scope"]["generation_dates"]):
                writer.writerow(
                    {
                        "series": key,
                        "date": date,
                        "in_reporting_window": i in mask,
                        "unit": r["data"]["units"].get(key, "same-unit total"),
                        "point": result["point"][i],
                        "statistic": result["statistic"],
                        **{
                            f"q{q}": result["quantiles"].get(
                                str(q / 100), [""] * len(result["point"])
                            )[i]
                            for q in range(10, 100, 10)
                        },
                    }
                )
    (output / "metrics.json").write_text(
        json.dumps(
            {"folds": folds, "comparison": comparison}, indent=2, allow_nan=False
        )
        + "\n"
    )
    period_totals = (
        {k: sum(v["point"][i] for i in mask) for k, v in forecasts.items()}
        if r["semantics"] == "flow" and r["aggregation"] != "none"
        else None
    )
    status = {
        "status": "complete",
        "fingerprint": proposal["fingerprint"],
        "runtime_seconds": elapsed,
        "peak_rss_bytes": peak_rss,
        "quantile_correction": r["settings"]["quantile_crossing"],
        "point_statistics": {k: v["statistic"] for k, v in forecasts.items()},
        "period_point_sums": period_totals,
        "period_sum_label": "sum of point rules, not an aggregate mean/median; interval unavailable" if period_totals is not None else None,
        "comparison": comparison,
        "model_semantics": [v.get("runtime_semantics") for v in forecasts.values()],
        "limitations": [
            "Quantiles are nominal and uncalibrated, with the selected correction applied; no aggregate interval.",
            "Rolling-origin evidence is exploratory; origins may overlap.",
            "As-of availability cannot reconstruct unknown historical revisions.",
            "No causal, accounting or unconstrained-demand conclusion.",
        ],
    }
    (output / "status.json").write_text(
        json.dumps(status, indent=2, allow_nan=False) + "\n"
    )
    (output / "provenance.json").write_text(json.dumps(proposal, indent=2) + "\n")
    body = [
        '<!doctype html><meta charset="utf-8"><title>Forecast report</title><h1>Forecast report</h1>',
        f"<p>Selected method: {html.escape(r['model'])}. Complete; baseline wins and ties remain valid findings.</p>",
        "<p>Generation "
        + html.escape(proposal["scope"]["generation_dates"][0])
        + " to "
        + html.escape(r["report"]["end"])
        + "; reporting "
        + html.escape(r["report"]["start"])
        + " to "
        + html.escape(r["report"]["end"])
        + ".</p>",
    ]
    for i, (key, result) in enumerate(forecasts.items()):
        filename = f"series-{i}.svg"
        chart(
            result.pop("history"),
            result["point"],
            f"{key}: {result['statistic']}",
            output / filename,
        )
        body.append(
            f'<img alt="{html.escape(key, quote=True)} forecast" src="{filename}">'
        )
    body += [
        "<h2>Comparable evaluation</h2><pre>"
        + html.escape(json.dumps(comparison, indent=2))
        + "</pre>",
        "<h2>Totals and limitations</h2><pre>"
        + html.escape(
            json.dumps(
                {
                    k: status[k]
                    for k in ("period_point_sums", "period_sum_label", "limitations")
                },
                indent=2,
            )
        )
        + "</pre>",
    ]
    (output / "report.html").write_text("\n".join(body) + "\n")
    return status
