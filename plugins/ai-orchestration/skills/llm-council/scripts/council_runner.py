#!/usr/bin/env python3
"""Minimal executable wrapper for the bundled LLM council package."""

from __future__ import annotations

import sys
from pathlib import Path

SCRIPT_DIRECTORY = Path(__file__).resolve().parent
if str(SCRIPT_DIRECTORY) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIRECTORY))

from llm_council.cli import main


if __name__ == "__main__":
    raise SystemExit(main())
