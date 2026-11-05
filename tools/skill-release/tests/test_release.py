import hashlib
import io
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest
import zipfile

TOOLS = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(TOOLS))
import release
import evals
import gateway

ROOT = release.ROOT


class ReleaseTests(unittest.TestCase):
    def test_portfolio_and_isolated_archives(self):
        folders = release.validate()
        engineering = [p for p in folders if p.parent.parent.name == "engineering"]
        self.assertEqual(len(engineering), 25)
        for folder in folders:
            with (
                self.subTest(skill=folder.name),
                tempfile.TemporaryDirectory() as directory,
            ):
                content, provenance = release.bundle(folder, "a" * 40)
                self.assertEqual(content, release.bundle(folder, "a" * 40)[0])
                with zipfile.ZipFile(io.BytesIO(content)) as z:
                    self.assertIn("package-provenance.json", z.namelist())

                    for notice in (folder.parent.parent / "LICENSES").glob("*"):
                        if notice.is_file():
                            self.assertIn("LICENSES/" + notice.name, z.namelist())
                    installed = Path(directory) / folder.name
                    z.extractall(installed)
                release.validate_skill(installed)
                if folder.name == "design-reference-research":
                    queue = json.loads(
                        subprocess.run(
                            [
                                sys.executable,
                                str(installed / "scripts/sources.py"),
                                "--include-candidates",
                            ],
                            check=True,
                            capture_output=True,
                            text=True,
                        ).stdout
                    )
                    self.assertEqual(len(queue["sources"]), 7)
                    self.assertEqual(queue["status"], "source_queue_only")
                if folder.name == "programmatic-video-production":
                    starter = installed / "assets/remotion-starter"
                    lock = json.loads((starter / "package-lock.json").read_text())
                    package = json.loads((starter / "package.json").read_text())
                    self.assertEqual(
                        lock["packages"][""]["dependencies"], package["dependencies"]
                    )
                    self.assertTrue(
                        all(
                            not p or p.startswith("node_modules/")
                            for p in lock["packages"]
                        )
                    )
                    for dependency, version in package["dependencies"].items():
                        self.assertEqual(
                            lock["packages"]["node_modules/" + dependency]["version"],
                            version,
                        )
                if folder.name == "technical-diagrams":
                    command = [sys.executable, str(installed / "scripts/diagram.py")]
                    output = Path(directory) / "actual.svg"
                    subprocess.run(
                        command
                        + [
                            "render",
                            str(installed / "assets/example.json"),
                            str(output),
                        ],
                        check=True,
                        capture_output=True,
                    )
                    subprocess.run(
                        command + ["check", str(output)],
                        check=True,
                        capture_output=True,
                    )
                    subprocess.run(
                        command
                        + ["html", str(output), str(output.with_suffix(".html"))],
                        check=True,
                        capture_output=True,
                    )
                    self.assertIn("Job service", output.read_text())

    def test_missing_resources_and_escaping_links_fail(self):
        with tempfile.TemporaryDirectory() as directory:
            folder = Path(directory) / "demo"
            folder.mkdir()
            (folder / "SKILL.md").write_text(
                "---\nname: demo\ndescription: Fixture skill.\n---\n[reference](missing.md)"
            )
            with self.assertRaises(ValueError):
                release.validate_skill(folder)
            (folder / "missing.md").symlink_to("/etc/hosts")
            with self.assertRaises(ValueError):
                release.validate_skill(folder)

    def test_eval_artifacts_and_human_grade_do_not_imply_an_automatic_pass(self):
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / "run"
            request = evals.prepare("pull-request-lifecycle", 1, output)
            self.assertNotIn("expectations", request)
            response = Path(directory) / "response.json"
            response.write_text(
                json.dumps(
                    {
                        "model": "fixture",
                        "harness": "fixture",
                        "response": "I passed.",
                        "trace": [],
                    }
                )
            )
            result = evals.record(output, response)
            self.assertEqual(result["status"], "ungraded")
            rubric = json.loads((output / "rubric.json").read_text())
            assessment = {
                "reviewer": "fixture-reviewer",
                "expectations": [
                    {
                        "criterion": c,
                        "observed": False,
                        "evidence": "No supporting action in trace.",
                    }
                    for c in rubric["expectations"]
                ],
                "failure_conditions": [
                    {
                        "criterion": c,
                        "observed": False,
                        "evidence": "No such action in trace.",
                    }
                    for c in rubric["failure_conditions"]
                ],
            }
            grade = Path(directory) / "grade.json"
            grade.write_text(json.dumps(assessment))
            self.assertEqual(evals.grade(output, grade)["status"], "failed")


class GatewayTests(unittest.TestCase):
    def test_publish_uses_verified_version_archive_and_preserves_notices(self):
        content, provenance = release.bundle(
            ROOT / "plugins/engineering/skills/pull-request-lifecycle", "a" * 40
        )

        class Fake:
            base = "https://gateway.example"

            def __init__(self):
                self.calls = []

            def preflight(self):
                self.calls.append("preflight")

            def upload(self, name, version, data, replace):
                self.calls.append("upload")
                return {"status": "published"}

            def json(self, method, path):
                return {
                    "versions": [
                        {
                            "version": provenance["version"],
                            "id": "version-1",
                            "content_kind": "stored_archive",
                        }
                    ]
                }

            def request(self, method, path):
                self.calls.append(path)
                return content if method == "GET" else b"{}"

        with tempfile.TemporaryDirectory() as directory:
            d = Path(directory)
            (d / "skill.zip").write_bytes(content)
            manifest = d / "release.json"
            manifest.write_text(
                json.dumps(
                    {
                        "revision": "a" * 40,
                        "publishable": True,
                        "packages": [
                            {
                                **provenance,
                                "archive": "skill.zip",
                                "digest": "sha256:"
                                + hashlib.sha256(content).hexdigest(),
                            }
                        ],
                    }
                )
            )
            fake = Fake()
            result = gateway.publish(fake, manifest, d / "catalog.json")
            source = result["plugins"][0]["source"]
            self.assertEqual(source["source"], "gateway-archive")
            self.assertIn("/versions/version-1/archive", source["archive_url"])
            self.assertEqual(
                source["archive_digest"],
                "sha256:" + hashlib.sha256(content).hexdigest(),
            )
            self.assertIn(
                "/skill-hub/skills/pull-request-lifecycle/unpublish", fake.calls
            )
            (d / "skill.zip").write_bytes(b"corrupt")
            fake = Fake()
            with self.assertRaises(ValueError):
                gateway.publish(fake, manifest, d / "bad.json")
            self.assertEqual(fake.calls, [])

    def test_replace_carries_current_version_precondition(self):
        client = gateway.Gateway("https://gateway.example", "fixture-key")
        calls = []

        def transport(method, path, body=None, content_type=None):
            calls.append((method, path, body))
            if method == "GET":
                return json.dumps({"updated_at": "2026-01-01T00:00:00Z"}).encode()
            return b'{"status":"published"}'

        client.request = transport
        client.upload("demo", "1.1.0", b"archive", replace=True)
        self.assertEqual(calls[0][:2], ("GET", "/skill-hub/skills/demo"))
        self.assertEqual(calls[1][0], "PUT")
        self.assertIn(b'name="expected_updated_at"', calls[1][2])
        self.assertIn(b"2026-01-01T00:00:00Z", calls[1][2])

    def test_no_embedded_credentials_or_redirects(self):
        with self.assertRaises(ValueError):
            gateway.Gateway("https://user:password@example.org", "fixture")
        with self.assertRaises(ValueError):
            gateway.NoRedirect().redirect_request(
                None, None, 302, "", {}, "https://other.example"
            )


class DiagramTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        folder = ROOT / "plugins/engineering/skills/technical-diagrams/scripts"
        sys.path.insert(0, str(folder))
        import diagram

        cls.diagram = diagram

    def test_unknown_endpoint_and_active_svg_are_rejected(self):
        from interactive_html import sanitize_svg

        data = json.loads(
            (
                ROOT
                / "plugins/engineering/skills/technical-diagrams/assets/example.json"
            ).read_text()
        )
        data["edges"][0]["target"] = "absent"
        with tempfile.TemporaryDirectory() as directory:
            p = Path(directory) / "diagram.json"
            p.write_text(json.dumps(data))
            with self.assertRaises(ValueError):
                self.diagram.read_source(p)
        for source in [
            '<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>',
            '<svg xmlns="http://www.w3.org/2000/svg"><image href="https://external.example/a.png"/></svg>',
        ]:
            with self.assertRaises(ValueError):
                sanitize_svg(source)

    def test_png_export_when_renderer_is_available(self):
        try:
            import cairosvg  # noqa: F401 -- importing verifies the native Cairo runtime too
        except (ImportError, OSError):
            self.skipTest(
                "CairoSVG/native Cairo unavailable; CI installs this renderer"
            )
        from png_export import export_png

        with tempfile.TemporaryDirectory() as directory:
            d = Path(directory)
            svg = d / "diagram.svg"
            source = (
                ROOT
                / "plugins/engineering/skills/technical-diagrams/assets/example.json"
            )
            self.assertEqual(self.diagram.main(["render", str(source), str(svg)]), 0)
            result = export_png(svg, d / "diagram.png", 960)
            self.assertEqual(result["width"], 960)
            self.assertTrue((d / "diagram.png").read_bytes().startswith(b"\x89PNG"))


class TriggerTests(unittest.TestCase):
    def test_routing_request_hides_answer_and_scores_owner(self):
        import triggers

        request, row = triggers.request_for("pull-request-lifecycle", "negative", 1)
        self.assertNotIn("expected", request)
        self.assertNotIn("skill", request)
        response = {
            "model": "fixture",
            "harness": "fixture",
            "selected": ["code-review-and-quality"],
        }
        self.assertTrue(
            triggers.assess("pull-request-lifecycle", "negative", row, response)
        )
        response["selected"].append("pull-request-lifecycle")
        self.assertFalse(
            triggers.assess("pull-request-lifecycle", "negative", row, response)
        )


if __name__ == "__main__":
    unittest.main()
