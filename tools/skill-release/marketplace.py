#!/usr/bin/env python3
"""Build pinned LiteLLM plugin registration payloads from a reviewed checkout."""

from __future__ import annotations
import argparse
import json
from pathlib import Path, PurePosixPath
import re
import subprocess

import release
import yaml

REPOSITORY = "https://github.com/aialchemylabs/agents-skills.git"


def revision(root):
    if subprocess.check_output(
        ["git", "status", "--porcelain"], cwd=root, text=True
    ):
        raise ValueError("commit reviewed changes before building a marketplace release")
    sha = subprocess.check_output(
        ["git", "rev-parse", "HEAD"], cwd=root, text=True
    ).strip()
    if not re.fullmatch(r"[0-9a-f]{40}", sha):
        raise ValueError("marketplace sources require a full Git commit SHA")
    return sha


def plugin_path(root, source):
    if not isinstance(source, str):
        raise ValueError("marketplace plugin sources must be local directory paths")
    normalized = source.removeprefix("./")
    parts = normalized.split("/")
    if (
        not normalized
        or any(part in ("", ".", "..") for part in parts)
        or "\\" in normalized
        or any(ord(char) < 32 for char in normalized)
        or PurePosixPath(normalized).is_absolute()
    ):
        raise ValueError(f"invalid plugin source path: {source!r}")
    folder = root
    for part in parts:
        folder = folder / part
        if folder.is_symlink():
            raise ValueError("marketplace plugin paths cannot contain symlinks")
    if not folder.is_dir():
        raise ValueError(f"plugin directory does not exist: {source!r}")
    return folder, normalized


def build(output, root=release.ROOT):
    root = root.resolve()
    if output.resolve().is_relative_to(root):
        raise ValueError("put marketplace release artifacts outside the repository")
    sha = revision(root)
    market = json.loads((root / ".claude-plugin/marketplace.json").read_text())
    entries = market.get("plugins")
    if not isinstance(entries, list) or not entries:
        raise ValueError("marketplace must contain at least one plugin")
    license_file = root / "LICENSE"
    if not license_file.is_file() or license_file.is_symlink():
        raise ValueError("repository LICENSE must be a regular file")
    license_bytes = license_file.read_bytes()
    plugins = []
    names = set()
    for entry in entries:
        folder, path = plugin_path(root, entry.get("source"))
        manifest = json.loads((folder / ".claude-plugin/plugin.json").read_text())
        name = manifest.get("name")
        if (
            not isinstance(name, str)
            or not release.NAME.fullmatch(name)
            or len(name) > 64
            or name != entry.get("name")
            or name in names
        ):
            raise ValueError(f"{folder}: invalid, duplicate or mismatched plugin name")
        names.add(name)
        version = manifest.get("version")
        if (
            not isinstance(version, str)
            or not re.fullmatch(r"\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?", version)
            or version != entry.get("version")
        ):
            raise ValueError(f"{folder}: invalid or mismatched plugin version")
        if manifest.get("skills") != "./skills/":
            raise ValueError(f"{folder}: marketplace release requires ./skills/")
        if not list((folder / "skills").glob("*/SKILL.md")):
            raise ValueError(f"{folder}: plugin has no skills")
        plugin_license = folder / "LICENSE"
        if (
            not plugin_license.is_file()
            or plugin_license.is_symlink()
            or plugin_license.read_bytes() != license_bytes
        ):
            raise ValueError(f"{folder}: plugin LICENSE must match repository LICENSE")
        description = manifest.get("description")
        if not isinstance(description, str) or not description:
            raise ValueError(f"{folder}: plugin description is required")
        author = manifest.get("author")
        if (
            not isinstance(author, dict)
            or not isinstance(author.get("name"), str)
            or not author["name"]
            or ("email" in author and not isinstance(author["email"], str))
        ):
            raise ValueError(f"{folder}: plugin author.name is required")
        payload = {
            "name": name,
            "version": version,
            "description": description,
            "author": {key: author[key] for key in ("name", "email") if key in author},
            "source": {
                "source": "git-subdir",
                "url": REPOSITORY,
                "path": path,
                "sha": sha,
            },
        }
        for key in ("homepage", "keywords"):
            if key in manifest:
                value = manifest[key]
                if key == "keywords":
                    valid = isinstance(value, list) and all(
                        isinstance(word, str) for word in value
                    )
                else:
                    valid = isinstance(value, str)
                if not valid:
                    raise ValueError(f"{folder}: invalid plugin {key}")
                payload[key] = manifest[key]
        if "category" in entry:
            if not isinstance(entry["category"], str):
                raise ValueError(f"{folder}: invalid plugin category")
            payload["category"] = entry["category"]
        plugins.append(payload)
    folders = release.validate(root)
    if revision(root) != sha:
        raise ValueError("checkout changed while building the marketplace release")
    result = {
        "revision": sha,
        "publishable": True,
        "skill_count": len(folders),
        "plugins": sorted(plugins, key=lambda item: item["name"]),
    }
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(result, sort_keys=True, indent=2) + "\n")
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    command = commands.add_parser("build")
    command.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    try:
        result = build(args.output)
        print(json.dumps(
            {
                "revision": result["revision"],
                "publishable": result["publishable"],
                "plugins": len(result["plugins"]),
                "skills": result["skill_count"],
                "output": str(args.output),
            },
            indent=2,
        ))
    except (OSError, ValueError, subprocess.CalledProcessError, yaml.YAMLError) as exc:
        parser.exit(1, str(exc) + "\n")


if __name__ == "__main__":
    main()
