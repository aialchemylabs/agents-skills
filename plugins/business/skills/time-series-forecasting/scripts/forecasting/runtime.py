"""External scope binding, pinned artifacts and bounded child execution (stdlib)."""

import hashlib
import json
import math
import os
from pathlib import Path
import platform
import re
import shutil
import subprocess
import sys
import time
import urllib.request
from .data import digest, inspect, load, dates_after, stamp
from .baselines import METHODS

SKILL = Path(__file__).resolve().parents[2]
MANIFEST = SKILL / "runtime-manifest.json"


def canonical(value):
    return hashlib.sha256(
        json.dumps(
            value, sort_keys=True, separators=(",", ":"), allow_nan=False
        ).encode()
    ).hexdigest()


def external(path):
    p = Path(path).expanduser().resolve()
    roots = [SKILL]
    for parent in SKILL.parents:
        if (parent / ".git").exists():
            roots.append(parent)
            break
    if any(p == root or p.is_relative_to(root) for root in roots):
        raise ValueError(
            "runtime destinations must be outside source checkout and installed skill"
        )
    return p


def host():
    def memory_pages(key):
        try:
            return os.sysconf(key) * os.sysconf("SC_PAGE_SIZE")
        except (ValueError, OSError):
            return None

    return {
        "system": platform.system(),
        "architecture": platform.machine(),
        "python": platform.python_version(),
        "cpu_count": os.cpu_count(),
        "physical_memory_bytes": memory_pages("SC_PHYS_PAGES"),
        "available_memory_bytes": memory_pages("SC_AVPHYS_PAGES"),
        "memory_limit_note": "unavailable counters are null; observations are not enforceable RAM limits",
        "free_bytes": shutil.disk_usage(Path.home()).free,
        "accelerator": "not probed; only tested CPU path is eligible",
    }


def validate(request):
    expected = {
        "data",
        "model",
        "report",
        "origins",
        "context",
        "season",
        "baselines",
        "primary_loss",
        "pooling",
        "weights",
        "semantics",
        "aggregation",
        "covariates",
        "groups",
        "cross_learning",
        "settings",
        "budgets",
        "output",
        "environment",
        "cache",
        "offline",
        "download_estimate_bytes",
    }
    if set(request) != expected:
        raise ValueError(
            "request keys differ from contract: " + str(set(request) ^ expected)
        )
    for value in [
        request["data"]["path"],
        *[request[k] for k in ("output", "environment", "cache")],
    ]:
        if not Path(value).is_absolute():
            raise ValueError("absolute input/runtime paths required")
    for k in ("output", "environment", "cache"):
        external(request[k])
    if (
        Path(request["data"]["path"])
        .resolve()
        .is_relative_to(external(request["output"]))
    ):
        raise ValueError("input data must be outside generated output directory")
    roots = [external(request[k]) for k in ("output", "environment", "cache")]
    if any(
        a == b or a.is_relative_to(b) or b.is_relative_to(a)
        for i, a in enumerate(roots)
        for b in roots[i + 1 :]
    ):
        raise ValueError("output, cache and environment must be separate roots")
    if request["model"] not in (*METHODS, "timesfm", "chronos"):
        raise ValueError("unsupported model")
    if not request["baselines"] or any(m not in METHODS for m in request["baselines"]):
        raise ValueError("select statistical baselines")
    settings = request["settings"]
    if settings.get("quantile_crossing") not in ("reject", "clamp-about-median"):
        raise ValueError("unsupported quantile crossing policy")
    if {**settings, "quantile_crossing": "reject"} != {
        "device": "cpu",
        "precision": "float32",
        "compile": False,
        "clipping": "none",
        "quantile_crossing": "reject",
        "truncation": "last-context",
        "horizon_extension": "reject",
    }:
        raise ValueError("unsupported effective settings")
    if not isinstance(request["context"], int) or not 2 <= request["context"] <= 16384:
        raise ValueError("invalid context")
    if (
        not isinstance(request["season"], int)
        or not 1 <= request["season"] <= request["context"]
    ):
        raise ValueError("invalid selected season")
    if request["semantics"] not in ("flow", "stock", "ratio"):
        raise ValueError("specify flow/stock/ratio semantics")
    if request["aggregation"] not in ("none", "direct-total", "sum-point-rules"):
        raise ValueError("unsupported aggregation; marginal intervals cannot be summed")
    if request["aggregation"] != "none" and request["semantics"] != "flow":
        raise ValueError(
            "stocks/ratios cannot use flow totals; authorize transformed input separately"
        )
    if (
        request["primary_loss"] not in ("mae", "rmse", "mase")
        or request["pooling"] != "weighted-series-mean"
    ):
        raise ValueError("unsupported primary loss/pooling")
    b = request["budgets"]
    if set(b) != {
        "seconds",
        "threads",
        "batch_size",
        "max_series",
        "max_origins",
        "max_horizon",
        "planned_download_bytes",
        "retries",
    }:
        raise ValueError("invalid budget keys")
    if (
        any(not isinstance(v, int) or v < 0 for v in b.values())
        or not 1 <= b["seconds"] <= 86400
        or not 1 <= b["threads"] <= 16
        or not 1 <= b["batch_size"] <= 256
        or not 0 <= b["retries"] <= 2
    ):
        raise ValueError("finite budgets outside limits")
    if not 1 <= len(request["origins"]) <= b["max_origins"] <= 100:
        raise ValueError("select bounded explicit chronological origins")
    if len(set(request["origins"])) != len(request["origins"]):
        raise ValueError("duplicate origins")
    groups = load(request["data"])
    ids = set(groups)
    if "__total__" in ids:
        raise ValueError("reserved aggregate series ID")
    if set(request["data"]["units"]) != ids:
        raise ValueError("units required for every series")
    if not 1 <= len(ids) <= b["max_series"] <= 1000:
        raise ValueError("series budget exceeded")
    if set(request["weights"]) != ids or any(
        not math.isfinite(v) or v <= 0 for v in request["weights"].values()
    ):
        raise ValueError("explicit positive weights required for every series")
    if (
        request["primary_loss"] != "mase"
        and len({request["data"]["units"][k] for k in ids}) > 1
    ):
        raise ValueError("unlike units require scaled primary metric")
    if sorted(x for g in request["groups"] for x in g) != sorted(ids) or any(
        not g for g in request["groups"]
    ):
        raise ValueError("groups must partition series exactly")
    if not isinstance(request["cross_learning"], bool):
        raise ValueError("cross_learning must be explicit")
    if request["model"] != "chronos" and (
        request["covariates"] or request["cross_learning"]
    ):
        raise ValueError("only Chronos supports selected covariates/cross-learning")
    if (
        request["aggregation"] != "none"
        and len(set(request["data"]["units"].values())) != 1
    ):
        raise ValueError("totals require same units")
    if request["aggregation"] == "direct-total" and (
        request["covariates"] or request["cross_learning"]
    ):
        raise ValueError("direct aggregate excludes covariates and cross-learning")
    width = 1 + len(request["covariates"])
    if request["cross_learning"] and any(
        len(g) * width > b["batch_size"] for g in request["groups"]
    ):
        raise ValueError("cross-learning group must fit unchanged in one batch")
    for c in request["covariates"]:
        if c["name"] == request["data"]["target"] or c["kind"] not in (
            "past",
            "known",
            "vintage",
            "scenario",
        ):
            raise ValueError("target/unsupported future actual covariate")
    tz = request["data"]["timezone"]
    rows = next(iter(groups.values()))
    last = rows[-1]["date"]
    gen = dates_after(
        last,
        request["report"]["end"],
        request["data"]["frequency"],
        tz,
        request["data"].get("calendar"),
    )
    mask = [
        i
        for i, d in enumerate(gen)
        if stamp(d, tz) >= stamp(request["report"]["start"], tz)
    ]
    if not mask or request["report"]["calendar_basis"] not in (
        "calendar",
        "fiscal",
        "supplied",
    ):
        raise ValueError("explicit reporting window/calendar basis required")
    if stamp(request["report"]["start"], tz) != stamp(gen[mask[0]], tz):
        raise ValueError("report start off cadence")
    if len(gen) > b["max_horizon"] or len(gen) > 1024:
        raise ValueError(
            "horizon budget/native horizon exceeded; no automatic extension"
        )
    for origin in request["origins"]:
        idx = next(
            (i for i, r in enumerate(rows) if stamp(origin, tz) == r["when"]), None
        )
        if idx is None or idx + len(gen) >= len(rows):
            raise ValueError("origin lacks full bridge-matched holdout")
    return groups, gen, mask


def propose(request):
    groups, gen, mask = validate(request)
    manifest = json.loads(MANIFEST.read_text())
    model = request["model"]
    foundation = model in manifest["models"]
    if foundation:
        if platform.system() != "Darwin" or platform.machine() != "arm64":
            raise ValueError(
                "only macOS arm64 CPU validated; other platforms ineligible"
            )
        if request["context"] > manifest["models"][model]["max_context"]:
            raise ValueError("context exceeds native model support")
        footprint = (
            sum(f["bytes"] for f in manifest["models"][model]["files"])
            + request["download_estimate_bytes"]
        )
        if request["download_estimate_bytes"] < 1_500_000_000:
            raise ValueError(
                "include conservative Python, wheels, source/build/bootstrap estimate >=1.5GB"
            )
    else:
        footprint = 0
    if footprint > request["budgets"]["planned_download_bytes"]:
        raise ValueError("planned artifact footprint exceeds download budget")
    if (
        foundation
        and model == "timesfm"
        and (
            (request["context"] + 31) // 32 * 32 + (len(gen) + 127) // 128 * 128 > 16384
        )
    ):
        raise ValueError("TimesFM padded context+horizon exceeds native limit")
    lock = SKILL / manifest["lock"]
    if digest(SKILL / manifest["build_lock"]) != manifest["build_lock_sha256"]:
        raise ValueError("build lock digest mismatch")
    if digest(lock) != manifest["lock_sha256"]:
        raise ValueError("dependency lock digest mismatch")
    scope = {
        "schema": 1,
        "request": request,
        "destinations": {k: str(external(request[k])) for k in ("output", "environment", "cache")},
        "snapshot": digest(request["data"]["path"]),
        "generation_dates": gen,
        "scoring_indices": mask,
        "checkpoint": manifest["models"].get(model),
        "lock_sha256": manifest["lock_sha256"] if foundation else None,
        "build_lock_sha256": manifest["build_lock_sha256"] if foundation else None,
        "implementation_sha256": canonical(
            {
                str(p.relative_to(SKILL)): digest(p)
                for p in sorted((SKILL / "scripts").rglob("*.py"))
            }
        ),
        "planned_artifact_bytes": footprint,
        "download_limit_kind": "preflight footprint, not an enforceable total network traffic limit",
    }
    return {
        "scope": scope,
        "fingerprint": canonical(scope),
        "authorization": "Fingerprint detects drift; it is NOT proof of human selection.",
    }


def verify(proposal, selected):
    if (
        not selected
        or selected != proposal["fingerprint"]
        or canonical(proposal["scope"]) != selected
    ):
        raise ValueError(
            "explicit selected proposal fingerprint required; cached checkpoints do not bypass selection"
        )
    current = propose(proposal["scope"]["request"])
    if current["fingerprint"] != selected:
        raise ValueError(
            "scope/data/code/lock drift; present a new option for selection"
        )
    return {**current["scope"]["request"], **current["scope"]["destinations"]}


def check_environment(lock):
    import importlib.metadata

    if platform.python_version() != json.loads(MANIFEST.read_text())["python"]:
        raise ValueError("locked environment requires Python 3.12")
    for name, version in re.findall(
        r"^([a-zA-Z0-9_-]+)==([^\s\\]+)", lock.read_text(), re.M
    ):
        if importlib.metadata.version(name) != version:
            raise ValueError("environment version mismatch: " + name)
    for name, revision in [
        ("timesfm", "e51928e27119cb17bebc005be2696b75e0a9e688"),
        ("chronos-forecasting", "1224fafa91e1ecd6674f5204066b159e40b5e4a1"),
    ]:
        direct = json.loads(
            importlib.metadata.distribution(name).read_text("direct_url.json")
        )
        if revision not in direct["url"]:
            raise ValueError("environment source revision mismatch")


def fetch_files(model, root, offline):
    root = external(root) / model["repo"].replace("/", "--") / model["revision"]
    for file in model["files"]:
        dest = root / file["name"]
        if dest.exists():
            if digest(dest) != file["sha256"]:
                raise ValueError(
                    "cached artifact checksum mismatch; never overwrite cache"
                )
            continue
        if offline:
            raise ValueError("offline cache incomplete: " + file["name"])
        root.mkdir(parents=True, exist_ok=True)
        temp = dest.with_suffix(dest.suffix + ".partial")
        if temp.exists():
            raise ValueError(
                "partial download exists; inspect task-owned partial before retry"
            )
        try:
            with (
                urllib.request.urlopen(
                    f"https://huggingface.co/{model['repo']}/resolve/{model['revision']}/{file['name']}",
                    timeout=60,
                ) as response,
                temp.open("xb") as out,
            ):
                count = 0
                while block := response.read(1024 * 1024):
                    count += len(block)
                    if count > file["bytes"]:
                        raise ValueError("artifact exceeds manifest byte size")
                    out.write(block)
            if temp.stat().st_size != file["bytes"] or digest(temp) != file["sha256"]:
                raise ValueError("download size/checksum mismatch")
            temp.rename(dest)
        finally:
            if temp.exists():
                temp.unlink()
    return root


def execute(proposal, selected):
    request = verify(proposal, selected)
    out = external(request["output"])
    out.mkdir(parents=True, exist_ok=True)
    saved = out / "proposal.json"
    if saved.exists() and json.loads(saved.read_text())["fingerprint"] != selected:
        raise ValueError("output contains another scope")
    saved.write_text(json.dumps(proposal, indent=2) + "\n")
    start = time.monotonic()
    limit = request["budgets"]["seconds"]

    def child(args, env=None):
        with (out / "runtime.log").open("a") as log:
            subprocess.run(
                args,
                check=True,
                timeout=max(1, limit - (time.monotonic() - start)),
                env=env,
                stdout=log,
                stderr=log,
            )

    try:
        python = sys.executable
        if proposal["scope"]["checkpoint"]:
            envroot = external(request["environment"])
            python = str(envroot / "bin/python")
            if not Path(python).exists():
                if request["offline"]:
                    raise ValueError("offline environment missing")
                child(["uv", "venv", "--python", "3.12.14", str(envroot)])
                child(
                    [
                        "uv",
                        "pip",
                        "sync",
                        "--python",
                        python,
                        "--require-hashes",
                        "--require-build-hashes",
                        "--build-constraints",
                        str(SKILL / json.loads(MANIFEST.read_text())["build_lock"]),
                        str(SKILL / json.loads(MANIFEST.read_text())["lock"]),
                    ]
                )
            child([python, str(SKILL / "scripts/forecast.py"), "_check-env"])
        env = {
            **os.environ,
            "HF_HUB_OFFLINE": "1",
            "TRANSFORMERS_OFFLINE": "1",
            "HF_HUB_DISABLE_TELEMETRY": "1",
            "OMP_NUM_THREADS": str(request["budgets"]["threads"]),
            "MKL_NUM_THREADS": str(request["budgets"]["threads"]),
            "CUDA_VISIBLE_DEVICES": "",
            "PYTHONDONTWRITEBYTECODE": "1",
        }
        # Downloads run in a child too, so the approved wall-clock bound includes them.
        for attempt in range(request["budgets"]["retries"] + 1):
            try:
                child(
                    [
                        sys.executable,
                        str(SKILL / "scripts/forecast.py"),
                        "_fetch",
                        str(saved),
                        "--selected",
                        selected,
                    ]
                )
                child(
                    [
                        python,
                        str(SKILL / "scripts/forecast.py"),
                        "_worker",
                        str(saved),
                        "--selected",
                        selected,
                    ],
                    env,
                )
                return json.loads((out / "status.json").read_text())
            except subprocess.CalledProcessError:
                if attempt == request["budgets"]["retries"]:
                    raise
    except Exception as exc:
        status = {
            "status": "failed",
            "fingerprint": selected,
            "error": str(exc),
            "runtime_seconds": time.monotonic() - start,
        }
        (out / "status.json").write_text(json.dumps(status, indent=2) + "\n")
        raise
