"""Opt-in, permission-restricted, atomic council result persistence."""

from __future__ import annotations

import copy
import json
import os
import re
import tempfile
from pathlib import Path
from typing import Any, Mapping

from .errors import PersistenceFailure
from .schemas import validate_manifest, validate_result


def _atomic_json(path: Path, payload: Mapping[str, Any]) -> None:
    descriptor = -1
    temporary_path: str | None = None
    try:
        descriptor, temporary_path = tempfile.mkstemp(
            prefix=f".{path.name}.",
            suffix=".tmp",
            dir=path.parent,
        )
        os.fchmod(descriptor, 0o600)
        with os.fdopen(descriptor, "w", encoding="utf-8") as handle:
            descriptor = -1
            json.dump(payload, handle, ensure_ascii=False, sort_keys=True, indent=2)
            handle.write("\n")
            handle.flush()
            os.fsync(handle.fileno())
        os.replace(temporary_path, path)
        temporary_path = None
        os.chmod(path, 0o600)
    except OSError as exc:
        raise PersistenceFailure("Could not atomically persist the council result") from exc
    finally:
        if descriptor >= 0:
            os.close(descriptor)
        if temporary_path is not None:
            try:
                os.unlink(temporary_path)
            except OSError:
                pass


def _summary(result: Mapping[str, Any]) -> dict[str, Any]:
    return {
        "schema_version": 1,
        "request_id": result["request_id"],
        "run_id": result["run_id"],
        "status": result["status"],
        "mode": result["mode"],
        "diversity": copy.deepcopy(result["diversity"]),
        "stages": copy.deepcopy(result["stages"]),
        "ranking": copy.deepcopy(result["ranking"]),
        "chairman": copy.deepcopy(result["chairman"]),
        "benchmark": copy.deepcopy(result["benchmark"]),
        "usage": copy.deepcopy(result["usage"]),
        "failures": copy.deepcopy(result["failures"]),
        "persistence": {"mode": "summary"},
    }


def persist_result(
    manifest: Mapping[str, Any],
    result: Mapping[str, Any],
) -> dict[str, Any]:
    """Persist a summary or full record and return non-secret artifact paths."""
    normalized = validate_manifest(manifest, command="validate")
    validated_result = validate_result(result)
    mode = normalized["policy"]["persist"]
    if mode == "none":
        return {"mode": "none", "status": "not_requested", "paths": []}
    output_directory = normalized["policy"]["output_directory"]
    if not output_directory:
        raise PersistenceFailure("Persistence requires an explicit output directory")
    target = Path(output_directory).resolve()
    home = Path.home().resolve()
    cwd = Path.cwd().resolve()
    if target in {Path("/"), home, cwd} or len(target.parts) < 3:
        raise PersistenceFailure("Refusing unsafe broad persistence target")
    try:
        target.mkdir(mode=0o700, parents=True, exist_ok=True)
        os.chmod(target, 0o700)
    except OSError as exc:
        raise PersistenceFailure("Could not prepare the persistence output directory") from exc
    run_id = validated_result["run_id"]
    if not isinstance(run_id, str) or not re.fullmatch(r"[A-Za-z0-9_-]{1,128}", run_id):
        raise PersistenceFailure("Run ID is not safe for persistence")
    run_directory = target / f"council-{run_id}"
    try:
        run_directory.mkdir(mode=0o700, exist_ok=False)
        os.chmod(run_directory, 0o700)
    except FileExistsError as exc:
        raise PersistenceFailure("Refusing to overwrite an existing run directory") from exc
    except OSError as exc:
        raise PersistenceFailure("Could not create the run persistence directory") from exc

    paths: list[str] = []
    try:
        if mode == "summary":
            summary_path = run_directory / "summary.json"
            _atomic_json(summary_path, _summary(validated_result))
            paths.append(str(summary_path))
        else:
            manifest_path = run_directory / "manifest.json"
            result_path = run_directory / "result.json"
            # Manifest validation has already rejected secret-shaped fields and values.
            _atomic_json(manifest_path, normalized)
            _atomic_json(result_path, validated_result)
            paths.extend((str(manifest_path), str(result_path)))
    except PersistenceFailure:
        # Leave no valid-looking final artifact after a failed atomic write.
        for path in run_directory.iterdir():
            try:
                path.unlink()
            except OSError:
                pass
        try:
            run_directory.rmdir()
        except OSError:
            pass
        raise
    return {"mode": mode, "status": "complete", "paths": paths}
