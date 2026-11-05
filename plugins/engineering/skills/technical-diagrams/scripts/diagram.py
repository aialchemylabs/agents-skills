#!/usr/bin/env python3
"""Static, self-contained adapter for the pinned Fireworks renderer."""

from __future__ import annotations
import argparse
import importlib
import importlib.util
import json
from pathlib import Path
import shutil
import sys
import xml.etree.ElementTree as ET
from diagram_ir import normalize_diagram
from interactive_html import build_interactive_html, sanitize_svg
from png_export import export_png
from validate_svg import run_check

ROOT = Path(__file__).resolve().parent
STYLES = {1, 4, 9, 10, 11, 12}


def read_source(path):
    data = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise ValueError("diagram source must be an object")
    mode = data.get("diagram_type", "architecture")
    if mode not in {
        "architecture",
        "flowchart",
        "sequence",
        "c4",
        "deployment",
        "event_stream",
        "observability",
    }:
        raise ValueError("unsupported diagram_type")
    diagram = normalize_diagram(data, mode)
    if diagram.style_index not in STYLES:
        raise ValueError("supported themes: 1, 4, 9, 10, 11, 12")
    if diagram.style_index == 4:
        overrides = data.setdefault("style_overrides", {})
        if not isinstance(overrides, dict):
            raise ValueError("style_overrides must be an object")
        overrides.setdefault("subtitle_fill", "#6b7280")
    return data, mode, diagram


def generator():
    spec = importlib.util.spec_from_file_location(
        "fireworks_static", ROOT / "generate-from-template.py"
    )
    module = importlib.util.module_from_spec(spec)
    sys.modules[spec.name] = module
    spec.loader.exec_module(module)
    return module


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)
    sub.add_parser("doctor")
    for name in ("validate", "render"):
        p = sub.add_parser(name)
        p.add_argument("source", type=Path)
        if name == "render":
            p.add_argument("output", type=Path)
            p.add_argument("--report", type=Path)
    p = sub.add_parser("check")
    p.add_argument("source", type=Path)
    for name in ("html", "png"):
        p = sub.add_parser(name)
        p.add_argument("source", type=Path)
        p.add_argument("output", type=Path)
    args = parser.parse_args(argv)
    try:
        if args.command == "doctor":
            raster = bool(shutil.which("rsvg-convert"))
            try:
                importlib.import_module("cairosvg")
                raster = True
            except (ImportError, OSError):
                pass
            result = {"svg": True, "html": True, "png": raster, "animation": False}
        elif args.command in ("validate", "render"):
            data, mode, diagram = read_source(args.source)
            result = {"valid": True, "semantics": dict(diagram.semantic_report)}
            if args.command == "render":
                svg, report = generator().build_svg_with_report(mode, data)
                svg = sanitize_svg(svg)
                args.output.parent.mkdir(parents=True, exist_ok=True)
                args.output.write_text(svg, encoding="utf-8")
                if args.report:
                    args.report.parent.mkdir(parents=True, exist_ok=True)
                    args.report.write_text(
                        json.dumps(report, indent=2) + "\n", encoding="utf-8"
                    )
                result = {
                    "svg": str(args.output),
                    "report": str(args.report) if args.report else None,
                }
        elif args.command == "check":
            result = {
                check: run_check(args.source, check)
                for check in ("xml", "markers", "collisions", "geometry", "composition")
            }
            print(json.dumps(result, indent=2))
            return 0 if all(v[0] for v in result.values()) else 1
        elif args.command == "html":
            output = build_interactive_html(
                args.source.read_text(encoding="utf-8"),
                args.source.stem,
                {"slug": args.source.stem},
            )
            args.output.parent.mkdir(parents=True, exist_ok=True)
            args.output.write_text(output, encoding="utf-8")
            result = {"html": str(args.output)}
        else:
            args.output.parent.mkdir(parents=True, exist_ok=True)
            result = export_png(args.source, args.output, 1920)
        print(json.dumps(result, indent=2))
        return 0
    except (OSError, ValueError, RuntimeError, ET.ParseError) as exc:
        print(json.dumps({"error": str(exc)}), file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
