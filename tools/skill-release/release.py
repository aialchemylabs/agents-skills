#!/usr/bin/env python3
"""Validate skills and build deterministic, self-contained portable archives."""

from __future__ import annotations
import argparse
import hashlib
import io
import json
from pathlib import Path
import re
import stat
import subprocess
import zipfile
import yaml

ROOT = Path(__file__).resolve().parents[2]
NAME = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
LINK = re.compile(r"\[[^\]]*\]\(([^\s)]+)(?:\s+[^)]*)?\)")


def frontmatter(path):
    text = path.read_text(encoding="utf-8")
    parts = text.split("---", 2)
    if len(parts) != 3 or parts[0].strip():
        raise ValueError(f"{path}: missing YAML frontmatter")
    data = yaml.safe_load(parts[1])
    if not isinstance(data, dict):
        raise ValueError(f"{path}: frontmatter must be an object")
    return data


def resources(folder):
    result = {}
    for path in sorted(folder.rglob("*")):
        if "__pycache__" in path.parts or path.suffix == ".pyc":
            continue
        if path.is_symlink():
            raise ValueError(f"{path}: portable packages cannot contain symlinks")
        if path.is_dir():
            continue
        if not path.is_file():
            raise ValueError(f"{path}: unsupported file type")
        result[path.relative_to(folder).as_posix()] = (
            path.read_bytes(),
            0o755 if path.stat().st_mode & 0o111 else 0o644,
        )
    return result


def validate_skill(folder):
    files = resources(folder)
    data = frontmatter(folder / "SKILL.md")
    name, desc = data.get("name"), data.get("description")
    if (
        not isinstance(name, str)
        or not NAME.fullmatch(name)
        or len(name) > 64
        or name != folder.name
    ):
        raise ValueError(f"{folder}: invalid skill name")
    if not isinstance(desc, str) or not 1 <= len(desc.encode("utf-8")) <= 1024:
        raise ValueError(
            f"{folder}: description must fit the installer 1024-byte limit"
        )
    for path in folder.rglob("*.md"):
        # Code examples may describe paths that are intentionally fictional.
        text = re.sub(r"```.*?```", "", path.read_text(encoding="utf-8"), flags=re.S)
        for link in LINK.findall(text):
            if "://" in link or link.startswith(("#", "mailto:", "/")):
                continue
            target = (path.parent / link.split("#")[0]).resolve()
            if not target.is_relative_to(folder.resolve()) or not target.exists():
                raise ValueError(f"{path}: missing or escaping local reference {link}")
    meta = folder / "agents/openai.yaml"
    if meta.exists():
        ui = yaml.safe_load(meta.read_text())
        if not isinstance(ui, dict) or not isinstance(ui.get("interface"), dict):
            raise ValueError(f"{meta}: missing interface")
        prompt = ui["interface"].get("default_prompt", "")
        if f"${name}" not in prompt:
            raise ValueError(f"{meta}: default prompt must name the skill")
    return data, files


def validate_cases(path, skill_names):
    data = json.loads(path.read_text())
    if data.get("skill_name") != path.stem or path.stem not in skill_names:
        raise ValueError(f"{path}: unknown skill")
    trigger = data.get("trigger", {})
    if not trigger.get("positive") or not trigger.get("negative"):
        raise ValueError(f"{path}: positive and negative trigger cases required")
    ids = set()
    for case in data.get("evals", []):
        if case.get("id") in ids or not isinstance(case.get("id"), int):
            raise ValueError(f"{path}: duplicate or invalid case id")
        ids.add(case["id"])
        if (
            case.get("kind") not in ("dialogue", "execution", "fixture")
            or not case.get("prompt")
            or not case.get("expectations")
        ):
            raise ValueError(f"{path}: incomplete behavioral case {case.get('id')}")
        for fixture in case.get("files", []):
            p = path.parent.parent / "fixtures" / fixture
            if (
                not p.resolve().is_relative_to(
                    (path.parent.parent / "fixtures").resolve()
                )
                or not p.exists()
            ):
                raise ValueError(f"{path}: invalid fixture {fixture}")
    if not ids:
        raise ValueError(f"{path}: no behavioral cases")
    return data


def validate(root=ROOT):
    market = json.loads((root / ".claude-plugin/marketplace.json").read_text())
    names = set()
    folders = []
    for plugin in market["plugins"]:
        base = root / plugin["source"]
        manifest = json.loads((base / ".claude-plugin/plugin.json").read_text())
        if manifest["version"] != plugin["version"]:
            raise ValueError(f"{base}: marketplace/plugin version mismatch")
        for file in sorted((base / "skills").glob("*/SKILL.md")):
            data, _ = validate_skill(file.parent)
            if data["name"] in names:
                raise ValueError("duplicate skill " + data["name"])
            names.add(data["name"])
            folders.append(file.parent)
    for path in (root / "plugins").glob("*/evals/cases/*.json"):
        validate_cases(path, names)
    for folder in (root / "plugins/engineering/skills").iterdir():
        if (folder / "SKILL.md").exists():
            if not (folder / "agents/openai.yaml").exists():
                raise ValueError(f"{folder}: missing discovery metadata")
            if not (
                root / "plugins/engineering/evals/cases" / f"{folder.name}.json"
            ).exists():
                raise ValueError(f"{folder}: missing eval cases")
    return folders


def bundle(folder, revision, root=ROOT):
    data, files = validate_skill(folder)
    plugin = folder.parent.parent
    for notice in (plugin / "LICENSES").glob("*"):
        if notice.is_file():
            files.setdefault("LICENSES/" + notice.name, (notice.read_bytes(), 0o644))
    if (root / "LICENSE").exists():
        files.setdefault(
            "LICENSES/REPOSITORY-LICENSE.txt", ((root / "LICENSE").read_bytes(), 0o644)
        )
    manifest = json.loads((plugin / ".claude-plugin/plugin.json").read_text())
    provenance = {
        "skill": data["name"],
        "version": manifest["version"],
        "repository": "https://github.com/aialchemylabs/agents-skills",
        "path": folder.relative_to(root).as_posix(),
        "revision": revision,
    }
    files["package-provenance.json"] = (
        json.dumps(provenance, sort_keys=True, indent=2).encode() + b"\n",
        0o644,
    )
    if (
        len(files) > 999
        or sum(len(v[0]) for v in files.values()) > 50 << 20
        or any(len(v[0]) > 10 << 20 for v in files.values())
    ):
        raise ValueError(f"{folder}: package exceeds portable gateway limits")
    output = io.BytesIO()
    with zipfile.ZipFile(output, "w", compression=zipfile.ZIP_DEFLATED) as archive:
        for name, (content, mode) in sorted(files.items()):
            info = zipfile.ZipInfo(name, date_time=(1980, 1, 1, 0, 0, 0))
            info.create_system = 3
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = (stat.S_IFREG | mode) << 16
            archive.writestr(info, content)
    content = output.getvalue()
    if len(content) > 10 << 20:
        raise ValueError("compressed package exceeds 10 MiB")
    return content, provenance


def build(output, allow_dirty=False, root=ROOT):
    folders = validate(root)
    if output.resolve().is_relative_to(root.resolve()):
        raise ValueError("put release artifacts outside the repository")
    status = subprocess.check_output(
        ["git", "status", "--porcelain"], cwd=root, text=True
    )
    if status and not allow_dirty:
        raise ValueError(
            "commit reviewed changes before building a publishable release"
        )
    revision = subprocess.check_output(
        ["git", "rev-parse", "HEAD"], cwd=root, text=True
    ).strip()
    if status:
        revision = "working-tree-" + revision
    output.mkdir(parents=True, exist_ok=True)
    packages = []
    for folder in folders:
        content, provenance = bundle(folder, revision, root)
        name = folder.name + "-" + provenance["version"] + ".zip"
        (output / name).write_bytes(content)
        packages.append(
            {
                **provenance,
                "archive": name,
                "digest": "sha256:" + hashlib.sha256(content).hexdigest(),
            }
        )
    result = {
        "revision": revision,
        "publishable": not bool(status),
        "packages": packages,
    }
    (output / "release.json").write_text(json.dumps(result, indent=2) + "\n")
    return result


def main():
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest="command", required=True)
    sub.add_parser("validate")
    b = sub.add_parser("build")
    b.add_argument("--output", type=Path, required=True)
    b.add_argument("--allow-dirty", action="store_true")
    args = p.parse_args()
    try:
        result = (
            {
                "skills": len(validate()),
                "status": "valid",
                "behavioral_evals": "not run",
            }
            if args.command == "validate"
            else build(args.output, args.allow_dirty)
        )
        print(json.dumps(result, indent=2))
    except (OSError, ValueError, yaml.YAMLError) as exc:
        p.exit(1, str(exc) + "\n")


if __name__ == "__main__":
    main()
