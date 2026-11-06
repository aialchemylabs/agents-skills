#!/usr/bin/env python3
"""Inspect → propose → explicitly select → run. No imports of optional model libraries here."""

import argparse
import json
import sys
from pathlib import Path

sys.dont_write_bytecode = True

from forecasting.data import inspect
from forecasting.runtime import (
    SKILL,
    MANIFEST,
    external,
    host,
    propose,
    verify,
    execute,
    fetch_files,
    check_environment,
)


def main():
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest="command", required=True)
    for cmd in ("inspect", "propose", "run", "_worker", "_fetch"):
        s = sub.add_parser(cmd)
        s.add_argument("input", type=Path)
        if cmd == "propose":
            s.add_argument("--output", required=True)
        if cmd in ("run", "_worker", "_fetch"):
            s.add_argument("--selected", required=True)
    sub.add_parser("_check-env")
    a = p.parse_args()
    if a.command == "_check-env":
        check_environment(SKILL / json.loads(MANIFEST.read_text())["lock"])
        return
    obj = json.loads(a.input.read_text())
    if a.command == "inspect":
        result = {"data": inspect(obj["data"]), "host": host()}
    elif a.command == "propose":
        result = propose(obj)
        out = external(a.output)
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(json.dumps(result, indent=2) + "\n")
        result = {
            "fingerprint": result["fingerprint"],
            "proposal": str(out),
            "authorization": result["authorization"],
        }
    elif a.command == "run":
        result = execute(obj, a.selected)
    elif a.command == "_fetch":
        r = verify(obj, a.selected)
        if obj["scope"]["checkpoint"]:
            fetch_files(obj["scope"]["checkpoint"], r["cache"], r["offline"])
        return
    else:
        # Recheck the scope before optional model imports, even through the internal worker.
        verify(obj, a.selected)
        from forecasting.engine import run

        result = run(obj, a.selected)
    print(json.dumps(result, indent=2, allow_nan=False))


if __name__ == "__main__":
    try:
        main()
    except (ValueError, KeyError, OSError) as exc:
        print(f"forecast: {exc}", file=sys.stderr)
        sys.exit(2)
