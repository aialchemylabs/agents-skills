"""Explicit regular calendars, snapshot validation and fold-local preparation."""

import calendar
import csv
import datetime as dt
import hashlib
import json
import math
from pathlib import Path
from zoneinfo import ZoneInfo


def digest(path):
    h = hashlib.sha256()
    with Path(path).open("rb") as f:
        for block in iter(lambda: f.read(1024 * 1024), b""):
            h.update(block)
    return h.hexdigest()


def stamp(value, timezone="UTC"):
    d = dt.datetime.fromisoformat(value)
    if d.tzinfo is None:
        d = d.replace(tzinfo=ZoneInfo(timezone))
        # Ambiguous/nonexistent local hours require an explicit offset.
        if d.replace(fold=1).utcoffset() != d.utcoffset():
            raise ValueError(
                "ambiguous or nonexistent local timestamp: supply UTC offset"
            )
    return d.astimezone(dt.timezone.utc)


def advance(value, frequency, timezone="UTC"):
    try:
        dt.date.fromisoformat(value)
        date_only = True
    except ValueError:
        date_only = False
    if frequency == "hourly" and date_only:
        raise ValueError("hourly cadence requires timestamps with a time component")
    d = stamp(value, timezone).astimezone(ZoneInfo(timezone))
    if frequency == "hourly":
        d = (d.astimezone(dt.timezone.utc) + dt.timedelta(hours=1)).astimezone(
            ZoneInfo(timezone)
        )
    elif frequency in ("daily", "weekly"):
        d += dt.timedelta(days=1 if frequency == "daily" else 7)
    elif frequency == "monthly":
        if d.day != calendar.monthrange(d.year, d.month)[1]:
            raise ValueError(
                "monthly cadence requires month-end timestamps; otherwise supply a calendar"
            )
        year, month = (d.year + 1, 1) if d.month == 12 else (d.year, d.month + 1)
        d = d.replace(year=year, month=month, day=calendar.monthrange(year, month)[1])
    else:
        raise ValueError("frequency must be hourly/daily/weekly/monthly or supplied")
    return d.date().isoformat() if date_only else d.isoformat()


def dates_after(last, end, frequency, timezone, supplied=None):
    if supplied is not None:
        dates = [
            x
            for x in supplied
            if stamp(last, timezone) < stamp(x, timezone) <= stamp(end, timezone)
        ]
        if not dates or stamp(dates[-1], timezone) != stamp(end, timezone):
            raise ValueError("supplied calendar does not reach report end")
        return dates
    dates = []
    while stamp(last, timezone) < stamp(end, timezone):
        last = advance(last, frequency, timezone)
        dates.append(last)
        if len(dates) > 100000:
            raise ValueError("calendar exceeds hard limit")
    if stamp(last, timezone) != stamp(end, timezone):
        raise ValueError("report end is off cadence")
    return dates


def load(spec):
    path = Path(spec["path"])
    if path.stat().st_size > spec.get("max_bytes", 20_000_000):
        raise ValueError("input exceeds byte budget")
    if path.suffix.lower() == ".csv":
        with path.open(newline="") as f:
            raw = list(csv.DictReader(f))
    elif path.suffix.lower() == ".json":
        raw = json.loads(path.read_text())
    else:
        raise ValueError(
            "only CSV or JSON; authorize an external export for other formats"
        )
    if not isinstance(raw, list) or not raw:
        raise ValueError("expected nonempty row array")
    tz = spec["timezone"]
    if spec.get("calendar") is not None:
        stamps = [stamp(v, tz) for v in spec["calendar"]]
        if stamps != sorted(set(stamps)):
            raise ValueError("supplied calendar must be strictly increasing and unique")
    groups = {}
    for r in raw:
        key = str(r[spec["id"]])
        date = r[spec["timestamp"]]
        when = stamp(date, tz)
        val = r.get(spec["target"])
        val = None if val in ("", None) else float(val)
        if val is not None and not math.isfinite(val):
            raise ValueError("non-finite target")
        if r.get("complete", True) in (False, "false", "False", "0"):
            raise ValueError(
                "incomplete periods must be explicitly excluded before proposal"
            )
        if spec.get("calendar") is not None and when not in {
            stamp(v, tz) for v in spec["calendar"]
        }:
            raise ValueError("observation absent from supplied calendar")
        available = stamp(r.get("available_at") or date, tz)
        if when > stamp(spec["as_of"], tz) or available > stamp(spec["as_of"], tz):
            raise ValueError("snapshot contains observations unavailable at data as_of")
        groups.setdefault(key, []).append(
            {"date": date, "when": when, "value": val, "available": available, "raw": r}
        )
    for rows in groups.values():
        rows.sort(key=lambda r: r["when"])
        if len({r["when"] for r in rows}) != len(rows):
            raise ValueError("duplicate series/timestamp")
        for left, right in zip(rows, rows[1:]):
            expected = dates_after(
                left["date"], right["date"], spec["frequency"], tz, spec.get("calendar")
            )
            if len(expected) != 1:
                raise ValueError(
                    "irregular calendar: supply explicit missing rows or authorized calendar"
                )
    schedules = [tuple(r["when"] for r in rows) for rows in groups.values()]
    if len(set(schedules)) != 1:
        raise ValueError(
            "related series require an aligned calendar and common as-of boundary"
        )
    return groups


def inspect(spec):
    groups = load(spec)
    return {
        "sha256": digest(spec["path"]),
        "series": {
            k: {
                "observations": len(rs),
                "first": rs[0]["date"],
                "last": rs[-1]["date"],
                "missing": sum(r["value"] is None for r in rs),
                "zeros": sum(r["value"] == 0 for r in rs),
                "negative": sum(r["value"] is not None and r["value"] < 0 for r in rs),
            }
            for k, rs in groups.items()
        },
        "limitations": [
            "No universal sample threshold establishes reliability.",
            "Publication/revision history is unknown unless available_at and vintage inputs are supplied.",
            "Observed sales may be censored by stockouts; no unconstrained demand is inferred.",
        ],
    }


def prepare(rows, cutoff, context):
    # Publication delays mask values, not calendar positions shared across series.
    training = [
        {**r, "value": r["value"] if r["available"] <= cutoff else None}
        for r in rows if r["when"] <= cutoff
    ][-context:]
    observed = [r["value"] for r in training if r["value"] is not None]
    if len(observed) < 2:
        raise ValueError("fewer than two available training values")
    mean = sum(observed) / len(observed)
    return training, [mean if r["value"] is None else r["value"] for r in training]


def covariates(request, key, training, future, cutoff):
    """Known-future values need an origin-visible vintage or an explicit scenario."""
    past, known = {}, {}
    for cov in request.get("covariates", []):
        name = cov["name"]
        values = [
            float(r["raw"][name]) if r.get("available", cutoff) <= cutoff else None
            for r in training
        ]
        observed = [v for v in values if v is not None]
        if not observed:
            raise ValueError("no origin-visible historical covariate: " + name)
        mean = sum(observed) / len(observed)
        past[name] = [mean if v is None else v for v in values]
        if cov["kind"] == "past":
            continue
        if cov["kind"] not in ("known", "vintage", "scenario"):
            raise ValueError("invalid covariate kind")
        vals = []
        for date in future:
            candidates = [
                v
                for v in cov["values"]
                if str(v["id"]) == key
                and stamp(v["date"], request["data"]["timezone"])
                == stamp(date, request["data"]["timezone"])
                and stamp(v["available_at"], request["data"]["timezone"]) <= cutoff
            ]
            if not candidates:
                raise ValueError(
                    "future covariate missing origin-visible value: " + name
                )
            v = max(
                candidates,
                key=lambda v: stamp(v["available_at"], request["data"]["timezone"]),
            )
            vals.append(float(v["value"]))
        known[name] = vals
    if any(
        not math.isfinite(v)
        for values in [*past.values(), *known.values()]
        for v in values
    ):
        raise ValueError("non-finite covariate")
    return past, known
