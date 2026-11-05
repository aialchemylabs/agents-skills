#!/usr/bin/env python3
"""Evaluate one skill-selection case without disclosing the expected owner."""

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
from release import ROOT, frontmatter


def request_for(skill, kind, case_id):
    matches = list((ROOT / "plugins").glob(f"*/evals/cases/{skill}.json"))
    if len(matches) != 1:
        raise ValueError("no unique case specification")
    spec = json.loads(matches[0].read_text())
    row = spec["trigger"][kind][case_id - 1]
    catalog = [
        frontmatter(p) for p in sorted((ROOT / "plugins").glob("*/skills/*/SKILL.md"))
    ]
    request = {
        "task": "Select up to three relevant skills for this user request. Return JSON with model, harness and selected (a ranked array of skill names).",
        "prompt": row["prompt"],
        "catalog": [
            {"name": s["name"], "description": s["description"]} for s in catalog
        ],
    }
    return request, row


def assess(skill, kind, row, response):
    selected = response.get("selected")
    if (
        not isinstance(selected, list)
        or len(selected) > 3
        or any(not isinstance(s, str) for s in selected)
        or len(set(selected)) != len(selected)
    ):
        raise ValueError("selected must be up to three unique skill names")
    if not response.get("model") or not response.get("harness"):
        raise ValueError("model and harness provenance required")
    if kind == "positive":
        return skill in selected[: row.get("top_k", 1)]
    return skill not in selected and (not row.get("owner") or row["owner"] in selected)


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--skill", required=True)
    p.add_argument("--kind", choices=["positive", "negative"], required=True)
    p.add_argument("--case", type=int, default=1)
    p.add_argument("--output", type=Path, required=True)
    p.add_argument("--response", type=Path)
    p.add_argument("--allow-execution", action="store_true")
    p.add_argument("--adapter", nargs="+")
    p.add_argument("--timeout", type=int, default=120)
    a = p.parse_args()
    try:
        if a.case < 1:
            raise ValueError("case is one-based")
        if a.output.resolve().is_relative_to(ROOT):
            raise ValueError("put eval results outside the repository")
        request, row = request_for(a.skill, a.kind, a.case)
        if a.adapter:
            if not a.allow_execution:
                raise ValueError(
                    "adapter execution requires --allow-execution and configured inference limits"
                )
            response = json.loads(
                subprocess.run(
                    a.adapter,
                    input=json.dumps(request),
                    text=True,
                    capture_output=True,
                    timeout=a.timeout,
                    check=True,
                ).stdout
            )
        elif a.response:
            response = json.loads(a.response.read_text())
        else:
            a.output.write_text(json.dumps(request, indent=2) + "\n")
            return
        result = {
            "skill": a.skill,
            "kind": a.kind,
            "case_id": a.case,
            "request_digest": hashlib.sha256(
                json.dumps(request, sort_keys=True).encode()
            ).hexdigest(),
            "status": "passed" if assess(a.skill, a.kind, row, response) else "failed",
            "response": response,
        }
        a.output.write_text(json.dumps(result, indent=2) + "\n")
        print(json.dumps(result, indent=2))
        if result["status"] == "failed":
            return 1
    except (OSError, ValueError, IndexError, subprocess.SubprocessError) as exc:
        p.exit(1, str(exc) + "\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
