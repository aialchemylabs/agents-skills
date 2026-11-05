#!/usr/bin/env python3
"""Run a selected behavioral case through an explicitly supplied agent harness.

The adapter reads one JSON request from stdin and emits a JSON response containing
model, harness, response and trace. It owns its tool sandbox and inference budget.
This runner never infers a pass from the model's assertion that it complied.
"""

from __future__ import annotations
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import time
from release import ROOT, resources, validate_cases


def prepare(skill, case_id, output):
    if output.resolve().is_relative_to(ROOT):
        raise ValueError("eval artifacts must be outside the repository")
    matches = list((ROOT / "plugins").glob(f"*/evals/cases/{skill}.json"))
    if len(matches) != 1:
        raise ValueError("skill has no unique case specification")
    path = matches[0]
    spec = validate_cases(path, {skill})
    cases = [c for c in spec["evals"] if c["id"] == case_id]
    if len(cases) != 1:
        raise ValueError("unknown case id")
    case = cases[0]
    folder = path.parent.parent.parent / "skills" / skill
    files = resources(folder)
    digest = hashlib.sha256()
    for name, (content, mode) in sorted(files.items()):
        digest.update(
            name.encode() + b"\0" + str(mode).encode() + b"\0" + content + b"\0"
        )
    output.mkdir(parents=True, exist_ok=False)
    shutil.copytree(
        folder, output / "skill", ignore=shutil.ignore_patterns("__pycache__", "*.pyc")
    )
    for fixture in case.get("files", []):
        src = path.parent.parent / "fixtures" / fixture
        dest = output / "fixtures" / fixture
        dest.parent.mkdir(parents=True, exist_ok=True)
        if src.is_dir():
            shutil.copytree(src, dest)
        else:
            shutil.copy2(src, dest)
    request = {
        "skill": skill,
        "case_id": case_id,
        "kind": case["kind"],
        "skill_digest": digest.hexdigest(),
        "skill_directory": str((output / "skill").resolve()),
        "fixture_directory": str((output / "fixtures").resolve()),
        "prompt": case["prompt"],
    }
    # Keep the evaluation rubric out of the agent request.
    (output / "request.json").write_text(json.dumps(request, indent=2) + "\n")
    (output / "rubric.json").write_text(
        json.dumps(
            {
                k: case.get(k, [])
                for k in ("expected_output", "expectations", "failure_conditions")
            },
            indent=2,
        )
        + "\n"
    )
    return request


def record(directory, response_path):
    request = json.loads((directory / "request.json").read_text())
    response = json.loads(response_path.read_text())
    for key in ("model", "harness", "response", "trace"):
        if key not in response:
            raise ValueError("response missing " + key)
    result = {
        **{k: request[k] for k in ("skill", "case_id", "kind", "skill_digest")},
        "case_digest": hashlib.sha256(
            (directory / "request.json").read_bytes()
        ).hexdigest(),
        "status": "ungraded",
        "recorded_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        **response,
    }
    result["status"] = "ungraded"
    (directory / "result.json").write_text(json.dumps(result, indent=2) + "\n")
    return result


def grade(directory, assessment_path):
    rubric = json.loads((directory / "rubric.json").read_text())
    result = json.loads((directory / "result.json").read_text())
    assessment = json.loads(assessment_path.read_text())
    if not assessment.get("reviewer"):
        raise ValueError("assessment must identify the reviewer")
    for key in ("expectations", "failure_conditions"):
        rows = assessment.get(key, [])
        if len(rows) != len(rubric[key]):
            raise ValueError("assessment must address every " + key)
        for i, row in enumerate(rows):
            if (
                row.get("criterion") != rubric[key][i]
                or not isinstance(row.get("observed"), bool)
                or not row.get("evidence")
            ):
                raise ValueError(
                    "each assessment needs the exact criterion, observed boolean and trace evidence"
                )
    result["status"] = (
        "passed"
        if all(r["observed"] for r in assessment["expectations"])
        and not any(r["observed"] for r in assessment["failure_conditions"])
        else "failed"
    )
    result["assessment"] = assessment
    (directory / "result.json").write_text(json.dumps(result, indent=2) + "\n")
    return result


def main():
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest="command", required=True)
    for cmd in ("prepare", "run"):
        s = sub.add_parser(cmd)
        s.add_argument("--skill", required=True)
        s.add_argument("--case", type=int, required=True)
        s.add_argument("--output", type=Path, required=True)
        if cmd == "run":
            s.add_argument("--allow-execution", action="store_true")
            s.add_argument("--timeout", type=int, default=120)
            s.add_argument("--adapter", nargs="+", required=True)
    for cmd, flag in [("record", "response"), ("grade", "assessment")]:
        s = sub.add_parser(cmd)
        s.add_argument("--directory", type=Path, required=True)
        s.add_argument("--" + flag, type=Path, required=True)
    a = p.parse_args()
    try:
        if a.command in ("prepare", "run"):
            if a.command == "run" and not a.allow_execution:
                raise ValueError(
                    "adapter execution requires --allow-execution; configure its model/cost and tool limits first"
                )
            result = prepare(a.skill, a.case, a.output)
            if a.command == "run":
                completed = subprocess.run(
                    a.adapter,
                    input=json.dumps(result),
                    text=True,
                    capture_output=True,
                    cwd=a.output,
                    timeout=a.timeout,
                    check=True,
                )
                response = a.output / "response.json"
                response.write_text(completed.stdout)
                result = record(a.output, response)
        elif a.command == "record":
            result = record(a.directory, a.response)
        else:
            result = grade(a.directory, a.assessment)
        print(json.dumps(result, indent=2))
    except (OSError, ValueError, subprocess.SubprocessError) as exc:
        p.exit(1, str(exc) + "\n")


if __name__ == "__main__":
    main()
