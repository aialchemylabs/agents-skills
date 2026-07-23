from __future__ import annotations

import copy
import unittest

from support import valid_result

from llm_council.render import render_result


class RenderTests(unittest.TestCase):
    def test_recommendation_leads_and_raw_transcripts_are_omitted(self) -> None:
        result = valid_result()
        markdown = render_result(result)
        recommendation_index = markdown.index("## Recommendation")
        benchmark_index = markdown.index("## Council quality and benchmark")
        self.assertLess(recommendation_index, benchmark_index)
        self.assertIn("Use the bounded hybrid.", markdown)
        self.assertNotIn('"rationale":', markdown)
        self.assertNotIn('"evaluations":', markdown)

    def test_partial_failures_and_diversity_limitations_are_visible(self) -> None:
        result = valid_result(status="partial")
        result["diversity"]["status"] = "perspective-only"
        result["diversity"]["limitations"] = [
            "All successful aliases are in one verified family."
        ]
        result["failures"] = [
            {
                "stage": "independent",
                "model": "fixture-model-b",
                "category": "provider_unavailable",
                "message": "The fixture route was unavailable.",
            }
        ]
        markdown = render_result(result)
        self.assertIn("## Limitations", markdown)
        self.assertIn("partial", markdown.lower())
        self.assertIn("perspective-only", markdown.lower())
        self.assertIn("fixture route was unavailable", markdown.lower())

    def test_unknown_cost_is_not_rendered_as_zero(self) -> None:
        markdown = render_result(valid_result())
        self.assertIn("unknown", markdown.lower())
        self.assertNotIn("$0.00", markdown)

    def test_model_identities_follow_disclosure_policy(self) -> None:
        shown = valid_result()
        shown_markdown = render_result(shown)
        self.assertIn("fixture\\|model-a", shown_markdown)

        hidden = copy.deepcopy(shown)
        hidden["benchmark"]["calls"][0]["model"] = None
        hidden["benchmark"]["calls"][0]["family"] = None
        hidden["candidates"][0].pop("model")
        hidden["candidates"][0].pop("family")
        hidden_markdown = render_result(hidden)
        self.assertNotIn("fixture-model-a", hidden_markdown)
        self.assertNotIn("family-a", hidden_markdown)

    def test_markdown_table_content_is_escaped(self) -> None:
        result = valid_result()
        result["benchmark"]["calls"][0]["model"] = "route|with\nnewline"
        markdown = render_result(result)
        self.assertIn("route\\|with newline", markdown)
        table_lines = [line for line in markdown.splitlines() if "route" in line]
        self.assertEqual(len(table_lines), 1)

    def test_chairman_failure_renders_labeled_host_boundary(self) -> None:
        result = valid_result(status="partial")
        result["chairman"] = {}
        result["stages"]["chairman"] = {"status": "failed"}
        result["failures"] = [
            {
                "stage": "chairman",
                "category": "ambiguous_timeout",
                "message": "Chairman unavailable.",
            }
        ]
        markdown = render_result(result)
        self.assertIn("chairman", markdown.lower())
        self.assertIn("unavailable", markdown.lower())
        self.assertNotIn("Chairman recommendation:", markdown)


if __name__ == "__main__":
    unittest.main()
