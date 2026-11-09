"""Deterministic advisory rubric and Borda aggregation."""

from __future__ import annotations

from typing import Any, Iterable, Mapping

from .errors import ValidationError
from .schemas import (
    CRITIQUE_CRITERIA,
    DEFAULT_RUBRIC_WEIGHTS,
    validate_critique_response,
)


def aggregate_rankings(
    critiques: Iterable[Mapping[str, Any]],
    candidate_labels: Iterable[str],
    weights: Mapping[str, float] | None = None,
) -> dict[str, Any]:
    """Aggregate valid critiques without turning their ranking into the final verdict."""
    labels = list(candidate_labels)
    if len(labels) < 2 or len(set(labels)) != len(labels):
        raise ValidationError("candidate_labels must contain unique labels")
    rubric = dict(DEFAULT_RUBRIC_WEIGHTS if weights is None else weights)
    if set(rubric) != set(CRITIQUE_CRITERIA):
        raise ValidationError("rubric weights must define every critique criterion")
    if abs(sum(float(value) for value in rubric.values()) - 1.0) > 1e-9:
        raise ValidationError("rubric weights must sum to 1.0")

    validated: list[dict[str, Any]] = []
    for critique in critiques:
        allow_partial = bool(critique.get("partial", False))
        validated.append(
            validate_critique_response(critique, labels, allow_partial=allow_partial)
        )

    weighted_values: dict[str, list[float]] = {label: [] for label in labels}
    rank_points: dict[str, float] = {label: 0.0 for label in labels}
    rank_possible: dict[str, float] = {label: 0.0 for label in labels}
    evaluation_coverage: dict[str, int] = {label: 0 for label in labels}
    ranking_coverage: dict[str, int] = {label: 0 for label in labels}

    for critique in validated:
        for evaluation in critique["evaluations"]:
            label = evaluation["candidate_label"]
            evaluation_coverage[label] += 1
            if "scores" not in evaluation:
                continue
            score = sum(
                float(evaluation["scores"][criterion]) * float(rubric[criterion])
                for criterion in CRITIQUE_CRITERIA
            )
            weighted_values[label].append((score - 1.0) / 4.0)
        ranking = critique.get("ranking", [])
        number_ranked = len(ranking)
        for index, label in enumerate(ranking):
            rank_points[label] += number_ranked - index - 1
            rank_possible[label] += max(number_ranked - 1, 1)
            ranking_coverage[label] += 1

    candidates: dict[str, dict[str, Any]] = {}
    for label in labels:
        weighted = weighted_values[label]
        candidates[label] = {
            "weighted_score": (
                round(sum(weighted) / len(weighted), 6) if weighted else None
            ),
            "rank_score": (
                round(rank_points[label] / rank_possible[label], 6)
                if rank_possible[label]
                else None
            ),
            "evaluation_coverage": evaluation_coverage[label],
            "ranking_coverage": ranking_coverage[label],
        }

    has_scores = any(weighted_values.values()) or any(ranking_coverage.values())
    order = sorted(
        labels,
        key=lambda label: (
            -(
                candidates[label]["weighted_score"]
                if candidates[label]["weighted_score"] is not None
                else -1.0
            ),
            -(
                candidates[label]["rank_score"]
                if candidates[label]["rank_score"] is not None
                else -1.0
            ),
            label,
        ),
    ) if has_scores else []
    return {
        "advisory": True,
        "subjective": True,
        "critic_count": len(validated),
        "expected_critic_count": len(labels),
        "coverage": (
            round(len(validated) / len(labels), 6) if labels else 0.0
        ),
        "weights": rubric,
        "candidates": candidates,
        "order": order,
        "limitations": (
            []
            if len(validated) == len(labels)
            else ["Critique coverage is partial; scores are not fully comparable."]
        ),
    }
