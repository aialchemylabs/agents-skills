from __future__ import annotations

import unittest

from support import critique_response

from llm_council.errors import ValidationError
from llm_council.ranking import aggregate_rankings


LABELS = ["Candidate A", "Candidate B", "Candidate C"]


def scored_critique(
    ranking: list[str],
    scores: dict[str, tuple[int, int, int, int, int]],
    *,
    partial: bool = False,
) -> dict:
    fields = (
        "evidence_quality",
        "correctness",
        "feasibility",
        "risk_awareness",
        "instruction_fit",
    )
    return {
        "evaluations": [
            {
                "candidate_label": label,
                "strengths": ["Fixture strength."],
                "weaknesses": ["Fixture weakness."],
                "unsupported_claims": [],
                "scores": dict(zip(fields, values)),
            }
            for label, values in scores.items()
        ],
        "ranking": ranking,
        "revised_position": "Fixture revised position.",
        "confidence": 0.8,
        **({"partial": True} if partial else {}),
    }


class RankingTests(unittest.TestCase):
    def test_exact_weighted_and_borda_scores(self) -> None:
        critiques = [
            scored_critique(
                LABELS,
                {
                    "Candidate A": (5, 5, 5, 5, 5),
                    "Candidate B": (3, 3, 3, 3, 3),
                    "Candidate C": (1, 1, 1, 1, 1),
                },
            ),
            scored_critique(
                LABELS,
                {
                    "Candidate A": (5, 5, 5, 5, 5),
                    "Candidate B": (3, 3, 3, 3, 3),
                    "Candidate C": (1, 1, 1, 1, 1),
                },
            ),
        ]
        result = aggregate_rankings(critiques, LABELS)
        candidates = result["candidates"]
        self.assertAlmostEqual(candidates["Candidate A"]["weighted_score"], 1.0)
        self.assertAlmostEqual(candidates["Candidate B"]["weighted_score"], 0.5)
        self.assertAlmostEqual(candidates["Candidate C"]["weighted_score"], 0.0)
        self.assertAlmostEqual(candidates["Candidate A"]["rank_score"], 1.0)
        self.assertAlmostEqual(candidates["Candidate B"]["rank_score"], 0.5)
        self.assertAlmostEqual(candidates["Candidate C"]["rank_score"], 0.0)
        self.assertEqual(result["order"], LABELS)

    def test_default_rubric_weights_are_applied(self) -> None:
        critique = scored_critique(
            LABELS,
            {
                "Candidate A": (5, 1, 1, 1, 1),
                "Candidate B": (1, 5, 1, 1, 1),
                "Candidate C": (1, 1, 5, 5, 5),
            },
        )
        result = aggregate_rankings([critique], LABELS)
        scores = result["candidates"]
        self.assertAlmostEqual(scores["Candidate A"]["weighted_score"], 0.3)
        self.assertAlmostEqual(scores["Candidate B"]["weighted_score"], 0.25)
        self.assertAlmostEqual(scores["Candidate C"]["weighted_score"], 0.45)

    def test_ties_are_deterministic_by_neutral_label(self) -> None:
        first = scored_critique(
            ["Candidate A", "Candidate B", "Candidate C"],
            {label: (4, 4, 4, 4, 4) for label in LABELS},
        )
        second = scored_critique(
            ["Candidate B", "Candidate A", "Candidate C"],
            {label: (4, 4, 4, 4, 4) for label in LABELS},
        )
        result = aggregate_rankings([first, second], LABELS)
        self.assertEqual(result["order"][:2], ["Candidate A", "Candidate B"])

    def test_rejects_duplicate_missing_and_unknown_ranking_labels(self) -> None:
        invalid_rankings = [
            ["Candidate A", "Candidate A", "Candidate C"],
            ["Candidate A", "Candidate B"],
            ["Candidate A", "Candidate B", "Candidate Z"],
        ]
        for ranking in invalid_rankings:
            with self.subTest(ranking=ranking):
                critique = scored_critique(
                    ranking,
                    {label: (4, 4, 4, 4, 4) for label in LABELS},
                )
                with self.assertRaises(ValidationError):
                    aggregate_rankings([critique], LABELS)

    def test_rejects_unknown_evaluation_label_and_out_of_range_score(self) -> None:
        critique = critique_response(LABELS)
        critique["evaluations"][0]["candidate_label"] = "Candidate Z"
        with self.assertRaises(ValidationError):
            aggregate_rankings([critique], LABELS)

        critique = critique_response(LABELS)
        critique["evaluations"][0]["scores"]["correctness"] = 6
        with self.assertRaises(ValidationError):
            aggregate_rankings([critique], LABELS)

    def test_partial_critique_reports_coverage_without_overstating(self) -> None:
        partial = scored_critique(
            ["Candidate A", "Candidate B"],
            {
                "Candidate A": (5, 5, 5, 5, 5),
                "Candidate B": (4, 4, 4, 4, 4),
            },
            partial=True,
        )
        result = aggregate_rankings([partial], LABELS)
        self.assertLess(result["coverage"], 1.0)
        self.assertTrue(result["limitations"])
        self.assertIsNone(result["candidates"]["Candidate C"]["weighted_score"])

    def test_ranking_is_explicitly_advisory(self) -> None:
        result = aggregate_rankings([critique_response(LABELS)], LABELS)
        self.assertTrue(result["advisory"])
        self.assertNotIn("chairman", result)
        self.assertNotIn("recommendation", result)


if __name__ == "__main__":
    unittest.main()
