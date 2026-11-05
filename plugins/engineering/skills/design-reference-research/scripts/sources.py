#!/usr/bin/env python3
"""Prepare a deduplicated research queue; this command never visits a site."""

import argparse
import json
from pathlib import Path
from urllib.parse import urlsplit, urlunsplit

CATALOG = Path(__file__).resolve().parents[1] / "references/source-catalog.json"


def canonical(url):
    parts = urlsplit(url)
    if (
        parts.scheme not in ("http", "https")
        or not parts.hostname
        or parts.username
        or parts.password
    ):
        raise ValueError("references must be HTTP(S) URLs without embedded credentials")
    host = parts.hostname.lower().removeprefix("www.")
    port = parts.port
    if port and (parts.scheme, port) not in (("http", 80), ("https", 443)):
        host += f":{port}"
    return urlunsplit(
        (
            parts.scheme.lower(),
            host,
            parts.path.rstrip("/"),
            parts.query,
            parts.fragment,
        )
    )


def select(
    catalog, urls, platform="any", topic=None, only_user=False, include_candidates=False
):
    queued = []
    seen = set()
    known = {canonical(entry["url"]): entry for entry in catalog["sources"]}
    for url in urls:
        key = canonical(url)
        if key not in seen:
            queued.append(
                {
                    **known.get(key, {}),
                    "url": url,
                    "origin": "user",
                    "inspection": "not_examined",
                }
            )
            seen.add(key)
    if not only_user:
        for entry in catalog["sources"]:
            if platform != "any" and platform not in entry["platforms"]:
                continue
            if topic and topic.lower() not in entry["tags"]:
                continue
            if (
                entry["verification"]["status"] == "unverified"
                and not include_candidates
            ):
                continue
            key = canonical(entry["url"])
            if key in seen:
                continue
            queued.append({**entry, "origin": "built-in", "inspection": "not_examined"})
            seen.add(key)
    return {
        "status": "source_queue_only",
        "sources": [{**item, "id": f"S{i:02}"} for i, item in enumerate(queued, 1)],
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--url", action="append", default=[])
    parser.add_argument("--platform", choices=["any", "web", "mobile"], default="any")
    parser.add_argument(
        "--topic", help="Exact catalog tag; user URLs are always retained"
    )
    parser.add_argument("--only-user", action="store_true")
    parser.add_argument("--include-candidates", action="store_true")
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    try:
        result = select(
            json.loads(CATALOG.read_text()),
            args.url,
            args.platform,
            args.topic,
            args.only_user,
            args.include_candidates,
        )
        serialized = json.dumps(result, indent=2) + "\n"
        if args.output:
            with args.output.open("x") as file:
                file.write(serialized)
        else:
            print(serialized, end="")
    except (ValueError, OSError) as error:
        parser.exit(1, f"{error}\n")


if __name__ == "__main__":
    main()
