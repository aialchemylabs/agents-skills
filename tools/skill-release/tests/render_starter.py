#!/usr/bin/env python3
"""Explicit integration check: install into a fresh temp project and render it."""

import argparse
import json
from pathlib import Path
import shutil
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[3]
SKILL = ROOT / "plugins/engineering/skills/video-production"


def run(command, cwd=None):
    subprocess.run([str(part) for part in command], check=True, cwd=cwd)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--workspace", type=Path, required=True)
    parser.add_argument("--browser-executable", type=Path, required=True)
    args = parser.parse_args()
    workspace = args.workspace.resolve()
    if workspace.is_relative_to(ROOT):
        parser.error("integration artifacts must stay outside the repository")
    workspace.mkdir(parents=True, exist_ok=False)
    project = workspace / "project"
    shutil.copytree(SKILL / "assets/remotion-starter", project)
    run(
        [
            "npm",
            "ci",
            "--ignore-scripts",
            "--no-audit",
            "--no-fund",
        ],
        cwd=project,
    )
    helper = [sys.executable, SKILL / "scripts/video.py"]
    run(
        helper
        + [
            "render",
            "--project",
            project,
            "--browser-executable",
            args.browser_executable,
            "--output",
            workspace / "landscape",
        ]
    )
    public = project / "public"
    public.mkdir(exist_ok=True)
    run(
        [
            "ffmpeg",
            "-nostdin",
            "-v",
            "error",
            "-f",
            "lavfi",
            "-i",
            "sine=frequency=440:sample_rate=48000",
            "-t",
            "6",
            public / "tone.wav",
        ]
    )
    run(
        [
            "ffmpeg",
            "-nostdin",
            "-v",
            "error",
            "-f",
            "lavfi",
            "-i",
            "color=c=teal:s=320x180:r=30",
            "-t",
            "2",
            "-c:v",
            "libx264",
            "-pix_fmt",
            "yuv420p",
            public / "clip.mp4",
        ]
    )
    props = json.loads((project / "props.json").read_text())
    props.update(width=720, height=1280, audio="tone.wav")
    props["shots"][0]["media"] = "clip.mp4"
    (project / "props.json").write_text(json.dumps(props, indent=2) + "\n")
    run(
        helper
        + [
            "render",
            "--project",
            project,
            "--browser-executable",
            args.browser_executable,
            "--output",
            workspace / "portrait",
        ]
    )
    for variant in ("landscape", "portrait"):
        run(
            helper
            + [
                "frames",
                workspace / variant / "video.mp4",
                "--times",
                "0,0.5,1.9,2.1,4.1,5.9666",
                "--output",
                workspace / (variant + "-frames"),
            ]
        )


if __name__ == "__main__":
    main()
