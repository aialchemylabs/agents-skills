#!/usr/bin/env python3
"""Render the local starter and verify media without installing dependencies."""

import argparse
from fractions import Fraction
import json
import math
from pathlib import Path
import shutil
import subprocess


def run(command, cwd=None):
    result = subprocess.run(
        [str(x) for x in command], cwd=cwd, capture_output=True, text=True, timeout=600
    )
    if result.returncode:
        raise ValueError(
            f"{Path(str(command[0])).name} failed ({result.returncode}): {(result.stderr or result.stdout)[-4000:]}"
        )
    return result.stdout


def require(name):
    path = shutil.which(name)
    if not path:
        raise ValueError(f"{name} is unavailable; no installation was attempted")
    return path


def probe(path):
    return json.loads(
        run(
            [
                require("ffprobe"),
                "-v",
                "error",
                "-show_streams",
                "-show_format",
                "-of",
                "json",
                path,
            ]
        )
    )


def inspect(path, width, height, fps, duration, require_audio=False):
    for label, value in [
        ("width", width),
        ("height", height),
        ("fps", fps),
        ("duration", duration),
    ]:
        if not math.isfinite(value) or value <= 0:
            raise ValueError(f"{label} must be finite and positive")
    metadata = probe(path)
    streams = metadata.get("streams", [])
    video = next((s for s in streams if s.get("codec_type") == "video"), None)
    audio = any(s.get("codec_type") == "audio" for s in streams)
    if not video:
        raise ValueError("export contains no video stream")
    actual_fps = float(Fraction(video["avg_frame_rate"]))
    actual_duration = float(video.get("duration", metadata["format"]["duration"]))
    if (video["width"], video["height"]) != (width, height):
        raise ValueError("export dimensions differ from the production contract")
    if (
        abs(actual_fps - fps) > 0.001
        or abs(actual_duration - duration) > 1 / fps + 0.001
    ):
        raise ValueError("export timing differs from the production contract")
    if video.get("codec_name") != "h264" or video.get("pix_fmt") not in (
        "yuv420p",
        "yuvj420p",
    ):
        raise ValueError("expected an H.264 8-bit 4:2:0 delivery file")
    if require_audio and not audio:
        raise ValueError("required audio stream is absent")
    run(
        [
            require("ffmpeg"),
            "-nostdin",
            "-v",
            "error",
            "-xerror",
            "-i",
            path,
            "-f",
            "null",
            "-",
        ]
    )
    return {
        "status": "technical_checks_passed",
        "file": str(Path(path).resolve()),
        "width": width,
        "height": height,
        "fps": actual_fps,
        "duration": actual_duration,
        "audio_present": audio,
        "codec": video["codec_name"],
        "pixel_format": video["pix_fmt"],
        "color_range": video.get("color_range", "unspecified"),
        "visual_review": "pending",
        "motion_review": "pending",
        "listening_review": "pending" if audio else "not_applicable",
    }


def project_props(project):
    props = json.loads((project / "props.json").read_text())
    for key in ("width", "height", "fps"):
        if type(props.get(key)) is not int or props[key] <= 0:
            raise ValueError(f"{key} must be a positive integer")
    if props["width"] % 2 or props["height"] % 2:
        raise ValueError("H.264 yuv420p dimensions must be even")
    if not isinstance(props.get("title"), str) or not props.get("shots"):
        raise ValueError("a title and nonempty shots are required")
    ids = set()
    for shot in props["shots"]:
        if not isinstance(shot.get("id"), str) or not shot["id"] or shot["id"] in ids:
            raise ValueError("shot IDs must be nonempty and unique")
        ids.add(shot["id"])
        if type(shot.get("durationFrames")) is not int or shot["durationFrames"] < 1:
            raise ValueError("shot durations must be positive integer frames")
        if not all(
            isinstance(shot.get(key), str) for key in ("headline", "body")
        ) or not isinstance(shot.get("caption", ""), str):
            raise ValueError("shot copy must be text")
    media = [
        (
            props.get("audio"),
            sum(s["durationFrames"] for s in props["shots"]) / props["fps"],
            "audio",
        )
    ]
    media.extend(
        (s.get("media"), s["durationFrames"] / props["fps"], "video")
        for s in props["shots"]
    )
    public = (project / "public").resolve()
    for name, seconds, kind in media:
        if name is None:
            continue
        if not isinstance(name, str) or not name or Path(name).is_absolute():
            raise ValueError("media must name a local file relative to public/")
        path = (public / name).resolve()
        if not path.is_relative_to(public) or not path.is_file():
            raise ValueError(f"missing or escaping media: {name}")
        info = probe(path)
        stream = next(
            (s for s in info.get("streams", []) if s.get("codec_type") == kind), None
        )
        if (
            not stream
            or float(stream.get("duration", info["format"]["duration"]))
            + 1 / props["fps"]
            < seconds
        ):
            raise ValueError(f"{kind} media is shorter than its planned usage: {name}")
    return props


def doctor(project):
    installed = {}
    for package in ("remotion", "@remotion/cli", "react", "react-dom"):
        path = project / "node_modules" / package / "package.json"
        installed[package] = (
            json.loads(path.read_text())["version"] if path.exists() else None
        )
    return {
        "packages": installed,
        "tools": {name: shutil.which(name) for name in ("node", "ffmpeg", "ffprobe")},
        "local_cli": (project / "node_modules/.bin/remotion").is_file(),
        "browser": "pass an existing executable explicitly; no browser download is performed",
    }


def render(project, browser, output):
    project, browser = project.resolve(), browser.resolve()
    status = doctor(project)
    if (
        not all(status["tools"].values())
        or not status["local_cli"]
        or not all(status["packages"].values())
    ):
        raise ValueError("missing prerequisites: " + json.dumps(status))
    if status["packages"]["remotion"] != status["packages"]["@remotion/cli"]:
        raise ValueError("Remotion and its CLI must have matching versions")
    if not browser.is_file():
        raise ValueError("browser executable does not exist; no download was attempted")
    props = project_props(project)
    output.mkdir(parents=True, exist_ok=False)
    video = output.resolve() / "video.mp4"
    try:
        run(
            [
                project / "node_modules/.bin/remotion",
                "render",
                "src/index.jsx",
                "ProductStory",
                video,
                "--props",
                "props.json",
                "--codec",
                "h264",
                "--pixel-format",
                "yuv420p",
                "--concurrency",
                "1",
                "--browser-executable",
                browser,
            ],
            cwd=project,
        )
        report = inspect(
            video,
            props["width"],
            props["height"],
            props["fps"],
            sum(s["durationFrames"] for s in props["shots"]) / props["fps"],
            props.get("audio") is not None,
        )
        report["runtime"] = status
        (output / "inspection.json").write_text(json.dumps(report, indent=2) + "\n")
        return report
    except Exception as error:
        (output / "failure.json").write_text(
            json.dumps(
                {"status": "failed", "error": str(error), "output_verified": False},
                indent=2,
            )
            + "\n"
        )
        raise


def frames(path, times, output):
    info = probe(path)
    duration = float(info["format"]["duration"])
    if not times or any(not math.isfinite(t) or t < 0 or t >= duration for t in times):
        raise ValueError(
            "frame timestamps must be finite, nonnegative and within the video"
        )
    output.mkdir(parents=True, exist_ok=False)
    captured = []
    for index, timestamp in enumerate(times):
        target = output / f"frame-{index:03}.png"
        run(
            [
                require("ffmpeg"),
                "-nostdin",
                "-n",
                "-v",
                "error",
                "-ss",
                timestamp,
                "-i",
                path,
                "-frames:v",
                "1",
                target,
            ]
        )
        if not target.is_file() or not target.stat().st_size:
            raise ValueError(f"no frame produced at {timestamp}")
        captured.append(
            {"seconds": timestamp, "file": str(target.resolve()), "inspected": False}
        )
    result = {"source": str(path.resolve()), "frames": captured}
    (output / "frames.json").write_text(json.dumps(result, indent=2) + "\n")
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    check = commands.add_parser("doctor")
    check.add_argument("--project", type=Path, required=True)
    production = commands.add_parser("render")
    production.add_argument("--project", type=Path, required=True)
    production.add_argument("--browser-executable", type=Path, required=True)
    production.add_argument("--output", type=Path, required=True)
    inspect_parser = commands.add_parser("inspect")
    inspect_parser.add_argument("input", type=Path)
    for field, type_ in [
        ("width", int),
        ("height", int),
        ("fps", float),
        ("duration", float),
    ]:
        inspect_parser.add_argument("--" + field, type=type_, required=True)
    inspect_parser.add_argument("--require-audio", action="store_true")
    extract = commands.add_parser("frames")
    extract.add_argument("input", type=Path)
    extract.add_argument("--times", required=True)
    extract.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    try:
        if args.command == "doctor":
            result = doctor(args.project.resolve())
        elif args.command == "render":
            result = render(args.project, args.browser_executable, args.output)
        elif args.command == "inspect":
            result = inspect(
                args.input,
                args.width,
                args.height,
                args.fps,
                args.duration,
                args.require_audio,
            )
        else:
            result = frames(
                args.input, [float(t) for t in args.times.split(",")], args.output
            )
        print(json.dumps(result, indent=2))
    except (
        ValueError,
        OSError,
        KeyError,
        TypeError,
        subprocess.TimeoutExpired,
    ) as error:
        parser.exit(1, f"{error}\n")


if __name__ == "__main__":
    main()
