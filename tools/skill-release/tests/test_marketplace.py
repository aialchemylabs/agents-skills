import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

TOOLS = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(TOOLS))
import marketplace


class MarketplaceTests(unittest.TestCase):
    def fixture(self, directory):
        root = Path(directory) / "repository"
        plugin = root / "plugins/business"
        (plugin / ".claude-plugin").mkdir(parents=True)
        (root / ".claude-plugin").mkdir()
        (root / "plugins/engineering/skills").mkdir(parents=True)
        (root / "LICENSE").write_text("MIT fixture grant\n")
        (plugin / "LICENSE").write_bytes((root / "LICENSE").read_bytes())
        manifest = {
            "name": "business",
            "version": "1.2.3",
            "description": "Fixture workflows.",
            "author": {"name": "Example Authors", "url": "https://example.com"},
            "homepage": "https://example.com",
            "keywords": ["fixture"],
            "skills": "./skills/",
        }
        (plugin / ".claude-plugin/plugin.json").write_text(json.dumps(manifest))
        entry = {
            "name": "business", "version": "1.2.3",
            "source": "./plugins/business", "category": "business", "strict": True,
        }
        (root / ".claude-plugin/marketplace.json").write_text(
            json.dumps({"name": "fixture", "plugins": [entry]})
        )
        for name in ("alpha", "beta"):
            folder = plugin / "skills" / name
            folder.mkdir(parents=True)
            (folder / "SKILL.md").write_text(
                f"---\nname: {name}\ndescription: Fixture {name}.\n---\nUse this fixture.\n"
            )
        self.git(root, "init", "--quiet")
        self.commit(root)
        return root

    def git(self, root, *args):
        return subprocess.check_output(
            ["git", "-C", str(root), "-c", "core.hooksPath=/dev/null",
             "-c", "user.name=Example", "-c", "user.email=example@example.com", *args],
            text=True, stderr=subprocess.PIPE,
        ).strip()

    def commit(self, root):
        self.git(root, "add", "--all")
        self.git(root, "commit", "--quiet", "-m", "Fixture")

    def test_registration_is_deterministic_pinned_and_schema_filtered(self):
        with tempfile.TemporaryDirectory() as directory:
            root = self.fixture(directory)
            output = Path(directory) / "release/marketplace.json"
            result = marketplace.build(output, root)
            content = output.read_bytes()
            self.assertEqual(marketplace.build(output, root), result)
            self.assertEqual(output.read_bytes(), content)
            self.assertEqual(result["revision"], self.git(root, "rev-parse", "HEAD"))
            self.assertTrue(result["publishable"])
            self.assertEqual(result["skill_count"], 2)
            self.assertEqual(len(result["plugins"]), 1)
            plugin = result["plugins"][0]
            self.assertEqual(plugin["source"], {
                "source": "git-subdir", "url": marketplace.REPOSITORY,
                "path": "plugins/business", "sha": result["revision"],
            })
            self.assertEqual(plugin["author"], {"name": "Example Authors"})
            self.assertEqual(set(plugin), {
                "name", "version", "description", "author", "source",
                "homepage", "keywords", "category",
            })

    def test_dirty_checkout_and_repository_output_are_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            root = self.fixture(directory)
            output = Path(directory) / "marketplace.json"
            (root / "unreviewed.txt").write_text("Uncommitted work")
            with self.assertRaisesRegex(ValueError, "commit reviewed changes"):
                marketplace.build(output, root)
            self.assertFalse(output.exists())
            with self.assertRaisesRegex(ValueError, "outside the repository"):
                marketplace.build(root / "generated.json", root)

    def test_invalid_sources_names_and_missing_or_different_notices_are_rejected(self):
        cases = (
            "traversal", "absolute", "symlink", "name", "keywords",
            "missing-license", "license",
        )
        for case in cases:
            with self.subTest(case=case), tempfile.TemporaryDirectory() as directory:
                root = self.fixture(directory)
                plugin = root / "plugins/business"
                catalog = root / ".claude-plugin/marketplace.json"
                data = json.loads(catalog.read_text())
                if case == "traversal":
                    data["plugins"][0]["source"] = "../elsewhere"
                elif case == "absolute":
                    data["plugins"][0]["source"] = str(plugin)
                elif case == "symlink":
                    (root / "linked").symlink_to("plugins/business", target_is_directory=True)
                    data["plugins"][0]["source"] = "./linked"
                elif case == "name":
                    data["plugins"][0]["name"] = "wrong-name"
                elif case == "keywords":
                    path = plugin / ".claude-plugin/plugin.json"
                    manifest = json.loads(path.read_text())
                    manifest["keywords"] = ["fixture", 3]
                    path.write_text(json.dumps(manifest))
                elif case == "missing-license":
                    (plugin / "LICENSE").unlink()
                else:
                    (plugin / "LICENSE").write_text("Different grant\n")
                catalog.write_text(json.dumps(data))
                self.commit(root)
                output = Path(directory) / "marketplace.json"
                with self.assertRaises(ValueError):
                    marketplace.build(output, root)
                self.assertFalse(output.exists())


if __name__ == "__main__":
    unittest.main()
