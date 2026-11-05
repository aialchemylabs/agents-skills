#!/usr/bin/env python3
"""Publish reviewed portable packages and emit a pinned gateway-only catalog.

This is a release client, not a gateway server. API credentials come only from
AIALCHEMY_SKILLS_API_KEY. Nothing uploads unless the publish command is invoked.
"""

from __future__ import annotations
import argparse
import hashlib
import io
import json
import os
from pathlib import Path
import re
import urllib.error
import urllib.parse
import urllib.request
import uuid
import zipfile


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        raise ValueError("gateway redirects are disabled")


class Gateway:
    def __init__(self, base, key):
        u = urllib.parse.urlsplit(base)
        if (
            u.scheme != "https"
            or not u.netloc
            or u.username
            or u.password
            or u.query
            or u.fragment
        ):
            raise ValueError("gateway must be a credential-free HTTPS origin")
        if u.path not in ("", "/"):
            raise ValueError("gateway URL must be an origin, without a path")
        if not key:
            raise ValueError("set AIALCHEMY_SKILLS_API_KEY to a scoped publishing key")
        self.base = base.rstrip("/")
        self.key = key
        self.opener = urllib.request.build_opener(NoRedirect())

    def request(self, method, path, body=None, content_type=None):
        if not path.startswith("/") or path.startswith("//"):
            raise ValueError("request path must be relative to the gateway")
        headers = {"Authorization": "Bearer " + self.key, "Accept": "application/json"}
        if content_type:
            headers["Content-Type"] = content_type
        req = urllib.request.Request(
            self.base + path, data=body, headers=headers, method=method
        )
        try:
            with self.opener.open(req, timeout=60) as response:
                data = response.read((10 << 20) + 1)
                if len(data) > 10 << 20:
                    raise ValueError("gateway response exceeds size bound")
                return data
        except urllib.error.HTTPError as exc:
            raise ValueError(
                f"gateway {method} {path}: HTTP {exc.code}; no automatic retry was made"
            ) from None

    def json(self, method, path, body=None, content_type=None):
        return json.loads(self.request(method, path, body, content_type))

    def preflight(self):
        return self.json("GET", "/skill-hub/skills?limit=1")

    def upload(self, name, version, data, replace=False):
        boundary = "skills-" + uuid.uuid4().hex
        concurrency = b""
        if replace:
            current = self.json("GET", f"/skill-hub/skills/{name}")
            timestamp = current.get("updated_at")
            if not isinstance(timestamp, str) or not timestamp:
                raise ValueError("gateway did not supply the current update timestamp")
            concurrency = (
                f'--{boundary}\r\nContent-Disposition: form-data; name="expected_updated_at"\r\n\r\n{timestamp}\r\n'
            ).encode()
        body = concurrency + (
            (
                f'--{boundary}\r\nContent-Disposition: form-data; name="version"\r\n\r\n{version}\r\n'
                f'--{boundary}\r\nContent-Disposition: form-data; name="file"; filename="{name}.zip"\r\nContent-Type: application/zip\r\n\r\n'
            ).encode()
            + data
            + f"\r\n--{boundary}--\r\n".encode()
        )
        path = (
            f"/skill-hub/skills/{name}/upload"
            if replace
            else "/skill-hub/skills/upload"
        )
        return self.json(
            "PUT" if replace else "POST",
            path,
            body,
            "multipart/form-data; boundary=" + boundary,
        )


def checked_files(data):
    with zipfile.ZipFile(io.BytesIO(data)) as z:
        if (
            len(z.infolist()) > 1000
            or sum(f.file_size for f in z.infolist()) > 50 << 20
        ):
            raise ValueError("archive exceeds portable limits")
        result = {}
        for f in z.infolist():
            if f.is_dir():
                continue
            if f.file_size > 10 << 20 or f.filename in result:
                raise ValueError("invalid archive member")
            result[f.filename] = z.read(f)
        return result


def publish(gateway, release_path, output, visibility="private", replace=False):
    release = json.loads(release_path.read_text())
    if not release.get("publishable") or not re.fullmatch(
        "[0-9a-f]{40}", release.get("revision", "")
    ):
        raise ValueError("only a clean committed release can be published")
    # Read and verify every local artifact before making the first remote mutation.
    prepared = []
    for package in release["packages"]:
        name = package["skill"]
        if not re.fullmatch("[a-z0-9]+(?:-[a-z0-9]+)*", name) or len(name) > 64:
            raise ValueError("invalid package name")
        path = (release_path.parent / package["archive"]).resolve()
        if not path.is_relative_to(release_path.parent.resolve()):
            raise ValueError("archive path escapes release")
        data = path.read_bytes()
        if (
            len(data) > 10 << 20
            or "sha256:" + hashlib.sha256(data).hexdigest() != package["digest"]
        ):
            raise ValueError("release artifact digest mismatch")
        original = checked_files(data)
        provenance = json.loads(original["package-provenance.json"])
        if any(
            provenance[k] != package[k]
            for k in ("skill", "version", "repository", "path", "revision")
        ):
            raise ValueError("package provenance mismatch")
        prepared.append((package, data, original))
    gateway.preflight()
    entries = []
    progress = output.with_suffix(".progress.json")
    output.parent.mkdir(parents=True, exist_ok=True)
    for package, data, original in prepared:
        name = package["skill"]
        result = gateway.upload(name, package["version"], data, replace)
        if result.get("status") == "pending_review":
            raise ValueError(
                f"{name} requires gateway administrator review; not published"
            )
        if visibility == "private":
            gateway.request("POST", f"/skill-hub/skills/{name}/unpublish")
        versions = gateway.json("GET", f"/skill-hub/skills/{name}/versions?limit=100")[
            "versions"
        ]
        matches = [
            v
            for v in versions
            if v.get("version") == package["version"]
            and v.get("content_kind") == "stored_archive"
        ]
        if len(matches) != 1:
            raise ValueError(f"{name}: no unique stored version for this release")
        version_id = urllib.parse.quote(str(matches[0]["id"]), safe="")
        archive_path = f"/skill-hub/skills/{name}/versions/{version_id}/archive"
        served = gateway.request("GET", archive_path)
        actual = checked_files(served)
        if any(actual.get(k) != v for k, v in original.items()):
            raise ValueError(f"{name}: gateway changed or lost a reviewed file")
        if any(
            k not in original and not k.startswith(".claude-plugin/") for k in actual
        ):
            raise ValueError(f"{name}: unexpected gateway-added files")
        entries.append(
            {
                "name": name,
                "version": package["version"],
                "source": {
                    "source": "gateway-archive",
                    "url": package["repository"],
                    "path": package["path"],
                    "sha": package["revision"],
                    "archive_url": gateway.base + archive_path,
                    "archive_digest": "sha256:" + hashlib.sha256(served).hexdigest(),
                },
            }
        )
        progress.write_text(
            json.dumps({"published": entries, "complete": False}, indent=2) + "\n"
        )
    catalog = {"name": "aialchemy-reviewed-skills", "plugins": entries}
    temporary = output.with_suffix(".tmp")
    temporary.write_text(json.dumps(catalog, indent=2) + "\n")
    temporary.replace(output)
    progress.write_text(
        json.dumps({"published": entries, "complete": True}, indent=2) + "\n"
    )
    return catalog


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("command", choices=["preflight", "publish"])
    p.add_argument("--gateway", required=True)
    p.add_argument("--release", type=Path)
    p.add_argument("--catalog-output", type=Path)
    p.add_argument("--visibility", choices=["private", "public"], default="private")
    p.add_argument(
        "--replace",
        action="store_true",
        help="replace existing versions using the update API",
    )
    a = p.parse_args()
    try:
        g = Gateway(a.gateway, os.environ.get("AIALCHEMY_SKILLS_API_KEY", ""))
        if a.command == "preflight":
            result = g.preflight()
        else:
            if not a.release or not a.catalog_output:
                raise ValueError("publish requires --release and --catalog-output")
            result = publish(g, a.release, a.catalog_output, a.visibility, a.replace)
        print(json.dumps(result, indent=2))
    except (OSError, ValueError, KeyError, zipfile.BadZipFile) as exc:
        p.exit(1, str(exc) + "\n")


if __name__ == "__main__":
    main()
