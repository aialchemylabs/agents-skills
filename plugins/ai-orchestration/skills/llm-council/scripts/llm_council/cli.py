"""JSON-on-stdin command line interface with stable exit codes."""

from __future__ import annotations

import argparse
import json
import os
import signal
import sys
import threading
from typing import Any, Callable, Mapping, TextIO

from .client import LiteLLMClient
from .config import apply_runtime_policy, load_runtime_config
from .errors import Cancelled, CouncilError, PersistenceFailure, ValidationError
from .orchestrator import CouncilOrchestrator
from .persistence import persist_result
from .render import render_result
from .schemas import (
    planned_base_requests,
    planned_output_tokens,
    planned_stages,
    validate_manifest,
)


def _parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="council_runner.py",
        description="Validate, preflight, run, or render one LLM council request.",
    )
    parser.add_argument("subcommand", choices=("validate", "preflight", "run", "render"))
    return parser


def _read_json(stream: TextIO) -> dict[str, Any]:
    try:
        value = json.load(stream)
    except json.JSONDecodeError as exc:
        raise ValidationError(
            f"stdin must contain one JSON object (line {exc.lineno}, column {exc.colno})"
        ) from exc
    if not isinstance(value, dict):
        raise ValidationError("stdin must contain one JSON object")
    return value


def _write_json(stream: TextIO, value: Mapping[str, Any]) -> None:
    json.dump(value, stream, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
    stream.write("\n")
    stream.flush()


def main(
    argv: list[str] | None = None,
    *,
    stdin: TextIO | None = None,
    stdout: TextIO | None = None,
    stderr: TextIO | None = None,
    environ: Mapping[str, str] | None = None,
    client_factory: Callable[..., LiteLLMClient] = LiteLLMClient,
    orchestrator_factory: Callable[..., CouncilOrchestrator] = CouncilOrchestrator,
) -> int:
    input_stream = stdin or sys.stdin
    output_stream = stdout or sys.stdout
    error_stream = stderr or sys.stderr
    environment = os.environ if environ is None else environ
    cancellation_event = threading.Event()
    previous_sigint: Any = None

    def on_sigint(_signum: int, _frame: Any) -> None:
        cancellation_event.set()

    try:
        if threading.current_thread() is threading.main_thread():
            previous_sigint = signal.signal(signal.SIGINT, on_sigint)
        args = _parser().parse_args(argv)
        payload = _read_json(input_stream)
        if args.subcommand == "render":
            output_stream.write(render_result(payload))
            output_stream.flush()
            return 0

        if args.subcommand == "validate":
            normalized = validate_manifest(payload, command=args.subcommand)
            _write_json(
                output_stream,
                {
                    "schema_version": 1,
                    "request_id": normalized["request_id"],
                    "status": "valid",
                    "mode": normalized["mode"],
                    "stages": planned_stages(normalized),
                    "base_inference_requests": planned_base_requests(normalized),
                    "planned_output_tokens": planned_output_tokens(normalized),
                    "approval_required_for_run": True,
                },
            )
            return 0

        runtime = load_runtime_config(
            environment,
            require_key=True,
        )
        effective_payload = apply_runtime_policy(payload, runtime)
        normalized = validate_manifest(effective_payload, command=args.subcommand)
        client = client_factory(
            base_url=runtime.base_url,
            api_key=runtime.api_key,
            connect_timeout_seconds=normalized["execution"]["connect_timeout_seconds"],
            request_timeout_seconds=normalized["execution"]["request_timeout_seconds"],
            max_transient_retries=normalized["execution"]["max_transient_retries"],
        )

        def progress(message: str) -> None:
            print(message, file=error_stream, flush=True)

        orchestrator = orchestrator_factory(
            client,
            cancellation_event=cancellation_event,
            progress=progress,
        )
        if args.subcommand == "preflight":
            _write_json(output_stream, orchestrator.preflight(normalized))
            return 0

        result = orchestrator.run(normalized)
        try:
            result["persistence"] = persist_result(normalized, result)
        except PersistenceFailure as exc:
            # The in-memory result is still safe to return; disclose the failed opt-in write.
            result["status"] = "partial"
            result["failures"].append(
                {
                    "stage": "persistence",
                    "subject": "local_output",
                    "category": exc.category,
                    "message": str(exc),
                }
            )
            result["persistence"] = {
                "mode": normalized["policy"]["persist"],
                "status": "failed",
                "paths": [],
            }
        _write_json(output_stream, result)
        return 0
    except Cancelled as exc:
        _write_json(
            output_stream,
            {
                "schema_version": 1,
                "status": "cancelled",
                "error": {"category": exc.category, "message": str(exc)},
            },
        )
        return exc.exit_code
    except CouncilError as exc:
        _write_json(
            output_stream,
            {
                "schema_version": 1,
                "status": "blocked",
                "error": {
                    "category": exc.category,
                    "message": str(exc),
                    "details": exc.details,
                },
            },
        )
        return exc.exit_code
    except KeyboardInterrupt:
        _write_json(
            output_stream,
            {
                "schema_version": 1,
                "status": "cancelled",
                "error": {
                    "category": "cancelled",
                    "message": "Council run cancelled; in-flight calls may still complete or bill",
                },
            },
        )
        return 130
    finally:
        if previous_sigint is not None:
            signal.signal(signal.SIGINT, previous_sigint)
