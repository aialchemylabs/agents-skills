#!/usr/bin/env python3
"""Track session resources and conservatively clean up exact owned paths."""

import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile


def git(path, *args):
    result = subprocess.run(
        ["git", "-C", str(path), *args], capture_output=True, text=True, timeout=30
    )
    if result.returncode:
        raise ValueError(result.stderr.strip() or "Git command failed")
    return result.stdout.strip()


def identity(path):
    info = path.lstat()
    if path.is_symlink() or not path.is_dir():
        raise ValueError("resource must be a real directory")
    return {"device": info.st_dev, "inode": info.st_ino}


def load(ledger, session):
    if ledger.is_symlink():
        raise ValueError("ledger cannot be a symlink")
    data = (
        json.loads(ledger.read_text())
        if ledger.exists()
        else {"session": session, "resources": []}
    )
    if data.get("session") != session:
        raise ValueError("ledger belongs to another session")
    return data


def save(ledger, data):
    with tempfile.NamedTemporaryFile(
        mode="w", dir=ledger.parent, prefix=".session-ledger-", delete=False
    ) as handle:
        json.dump(data, handle, indent=2)
        handle.write("\n")
    os.replace(handle.name, ledger)


def worktrees(path):
    blocks = git(path, "worktree", "list", "--porcelain").split("\n\n")
    return [
        dict(
            line.split(" ", 1) if " " in line else (line, "")
            for line in block.splitlines()
        )
        for block in blocks
    ]


def register(
    ledger, data, path, purpose="task checkout", evidence="caller ownership declaration",
    consumers=(), host_managed=False,
):
    if path.is_symlink():
        raise ValueError("worktree path cannot be a symlink")
    path = path.resolve()
    listed = worktrees(path)
    entry = next(
        (item for item in listed if Path(item["worktree"]).resolve() == path), None
    )
    if not entry or Path(listed[0]["worktree"]).resolve() == path or "locked" in entry:
        raise ValueError("only an exact unlocked linked worktree may be registered")
    if ledger.resolve().is_relative_to(path):
        raise ValueError("ledger must live outside the worktree")
    if any(item["path"] == str(path) for item in data["resources"]):
        raise ValueError("resource is already registered")
    data["resources"].append(
        {
            "kind": "worktree",
            "purpose": purpose,
            "creation_evidence": evidence,
            "path": str(path),
            "identity": identity(path),
            "consumers": consumer_records(consumers),
            "host_managed": host_managed,
            "common_git_dir": git(
                path, "rev-parse", "--path-format=absolute", "--git-common-dir"
            ),
        }
    )
    save(ledger, data)


def consumer_records(handles):
    if any(not handle.strip() for handle in handles) or len(set(handles)) != len(handles):
        raise ValueError("consumer handles must be nonempty and unique")
    return [{"handle": handle} for handle in handles]


def resource(data, path, allow_removed=False):
    if path is None:
        raise ValueError("an exact registered --path is required")
    matches = [item for item in data["resources"] if item["path"] == str(path.resolve())]
    if len(matches) != 1 or (matches[0].get("removed") and not allow_removed):
        raise ValueError("path is not a live registered resource")
    return matches[0]


def update_consumers(ledger, data, path, handles, release_evidence=None):
    """Record consumers before use; release only after caller-observed closure."""
    additions = consumer_records(handles)
    if not additions:
        raise ValueError("at least one --consumer handle is required")
    item = resource(data, path)
    consumers = item.get("consumers", [])
    existing = {entry["handle"]: entry for entry in consumers}
    if release_evidence is not None:
        if not release_evidence.strip() or any(
            handle not in existing for handle in handles
        ):
            raise ValueError(
                "release requires known consumers and nonempty closure evidence"
            )
        for handle in handles:
            existing[handle]["release_evidence"] = release_evidence
    else:
        if any(handle in existing for handle in handles):
            raise ValueError(
                "consumer handle is already recorded; use a new handle for new use"
            )
        consumers.extend(additions)
        item["consumers"] = consumers
    save(ledger, data)
    return item


def reason(item, ledger, retain, target, inactive):
    path = Path(item["path"])
    if item.get("removed"):
        return "already removed"
    if item.get("host_managed"):
        return "host-managed checkout requires retirement through its owning host"
    held = [
        entry["handle"] for entry in item.get("consumers", [])
        if not entry.get("release_evidence")
    ]
    if held:
        return "consumers have not been confirmed closed: " + ", ".join(held)
    if not path.exists() or path.is_symlink():
        return "missing or symlink resource"
    if identity(path) != item["identity"]:
        return "filesystem identity changed"
    path = path.resolve()
    if Path.cwd().resolve().is_relative_to(path):
        return "current working directory is inside the resource"
    if ledger.resolve().is_relative_to(path) or any(
        p.resolve().is_relative_to(path) or path.is_relative_to(p.resolve())
        for p in retain
    ):
        return "contains or overlaps a retained artifact or ledger"
    if item["kind"] == "worktree":
        if not inactive:
            return "worktree inactivity has not been confirmed"
        listed = worktrees(path)
        entry = next((i for i in listed if Path(i["worktree"]).resolve() == path), None)
        if (
            not entry
            or Path(listed[0]["worktree"]).resolve() == path
            or "locked" in entry
        ):
            return "not an unlocked linked worktree"
        if (
            git(path, "rev-parse", "--path-format=absolute", "--git-common-dir")
            != item["common_git_dir"]
        ):
            return "repository identity changed"
        if git(path, "status", "--porcelain", "--untracked-files=all", "--ignored"):
            return "tracked, untracked or ignored work remains"
        if not target:
            return "no verified merge target provided"
        head = git(path, "rev-parse", "HEAD")
        if git(path, "merge-base", "HEAD", target) != head:
            return "HEAD is not reachable from the merge target"
    elif item["kind"] == "temp":
        if path.parent != Path(
            tempfile.gettempdir()
        ).resolve() or not path.name.startswith("aialchemy-session-"):
            return "not a helper-created temporary directory"
    else:
        return "unknown resource kind"
    for root, directories, files in os.walk(path, followlinks=False):
        candidate = Path(root)
        if (
            (candidate / "HEAD").is_file()
            and (candidate / "objects").is_dir()
            and ((candidate / "refs").is_dir() or (candidate / "packed-refs").is_file())
        ):
            return "nested bare repository present"
        for name in directories + files:
            child = Path(root) / name
            if child.name == ".git" and not (
                item["kind"] == "worktree" and child.parent == path and child.is_file()
            ):
                return "nested repository present"
            if (
                child.is_symlink()
                or child.lstat().st_dev != item["identity"]["device"]
                or child.is_mount()
            ):
                return "symlink or mounted resource requires inspection"
    return None


def cleanup(ledger, data, retain, target, inactive, apply, only_path=None):
    selected = (
        resource(data, only_path, allow_removed=True) if only_path is not None else None
    )
    report = []
    for item in data["resources"]:
        if selected is not None and item is not selected:
            continue
        if item.get("removed"):
            report.append(
                {
                    "path": item["path"],
                    "status": "already_removed",
                    "reason": "recorded successful removal",
                }
            )
            continue
        try:
            blocked = reason(item, ledger, retain, target, inactive)
            status = "retain" if blocked else "eligible"
            if not blocked and apply:
                blocked = reason(item, ledger, retain, target, inactive)
                if blocked:
                    status = "retain"
                else:
                    path = Path(item["path"])
                    if item["kind"] == "worktree":
                        git(path, "worktree", "remove", str(path))
                    else:
                        shutil.rmtree(path)
                    item["removed"] = True
                    save(ledger, data)
                    status = "removed"
            report.append({"path": item["path"], "status": status, "reason": blocked})
        except (OSError, ValueError, subprocess.TimeoutExpired) as error:
            report.append(
                {"path": item["path"], "status": "retain", "reason": str(error)}
            )
    return {"applied": apply, "resources": report}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "command", choices=["temp", "register-worktree", "hold", "release", "cleanup"]
    )
    parser.add_argument("--ledger", type=Path, required=True)
    parser.add_argument("--session", required=True)
    parser.add_argument("--path", type=Path)
    parser.add_argument("--purpose")
    parser.add_argument("--creation-evidence")
    parser.add_argument("--created-by-session", action="store_true")
    parser.add_argument(
        "--consumer", action="append", default=[],
        help="exact owning-workflow consumer handle; repeat for shared use",
    )
    parser.add_argument("--release-evidence", help="observed owning-workflow closure result")
    parser.add_argument(
        "--host-managed", action="store_true",
        help="checkout must be retired by its owning host, never this helper",
    )
    parser.add_argument("--merged-into")
    parser.add_argument("--inactive-worktrees", action="store_true")
    parser.add_argument("--retain", type=Path, action="append", default=[])
    parser.add_argument("--apply", action="store_true")
    args = parser.parse_args()
    try:
        ledger = args.ledger.absolute()
        data = load(ledger, args.session)
        if args.command in ("temp", "register-worktree") and not args.purpose:
            raise ValueError("resource creation/registration requires --purpose")
        if args.command == "temp":
            consumers = consumer_records(args.consumer)
            path = Path(tempfile.mkdtemp(prefix="aialchemy-session-")).resolve()
            data["resources"].append(
                {
                    "kind": "temp",
                    "purpose": args.purpose,
                    "creation_evidence": "created by helper tempfile.mkdtemp in this session",
                    "path": str(path),
                    "identity": identity(path),
                    "consumers": consumers,
                }
            )
            try:
                save(ledger, data)
            except OSError:
                path.rmdir()  # Creation is not complete until ownership is recorded.
                raise
            result = {"created": str(path)}
        elif args.command == "register-worktree":
            if (
                not args.path
                or not args.created_by_session
                or not args.creation_evidence
            ):
                raise ValueError(
                    "registration needs --path, --created-by-session and --creation-evidence"
                )
            register(
                ledger, data, args.path, args.purpose, args.creation_evidence,
                args.consumer, args.host_managed,
            )
            result = {"registered": str(args.path.resolve())}
        elif args.command in ("hold", "release"):
            if args.command == "release" and not args.release_evidence:
                raise ValueError(
                    "release requires --release-evidence from the owning workflow"
                )
            item = update_consumers(
                ledger, data, args.path, args.consumer,
                args.release_evidence if args.command == "release" else None,
            )
            result = {"path": item["path"], "consumers": item["consumers"]}
        else:
            if not ledger.exists():
                raise ValueError(
                    "no session ledger exists; ownership cannot be inferred"
                )
            result = cleanup(
                ledger,
                data,
                args.retain,
                args.merged_into,
                args.inactive_worktrees,
                args.apply,
                args.path,
            )
        print(json.dumps(result, indent=2))
    except (OSError, ValueError, KeyError, subprocess.TimeoutExpired) as error:
        parser.exit(1, f"{error}\n")


if __name__ == "__main__":
    main()
