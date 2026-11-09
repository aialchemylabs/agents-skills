"""Concise Markdown rendering for validated council results."""

from __future__ import annotations

from typing import Any, Mapping

from .schemas import validate_result


def _escape(value: Any) -> str:
    if value is None:
        return "unknown"
    return str(value).replace("\\", "\\\\").replace("|", "\\|").replace("\r", " ").replace("\n", " ")


def _bullets(values: list[Any], fallback: str) -> list[str]:
    return [f"- {value}" for value in values] if values else [f"- {fallback}"]


def render_result(result: Mapping[str, Any]) -> str:
    """Render outcome first, disclose partial evidence and avoid raw transcript dumps."""
    obj = validate_result(result)
    chairman = obj["chairman"]
    chairman_response = (
        chairman.get("response") if isinstance(chairman.get("response"), dict) else None
    )
    if chairman_response:
        recommendation = chairman_response["recommendation"]
        why = chairman_response.get("why", [])
        consensus = chairman_response.get("consensus", [])
        disagreements = chairman_response.get("disagreements", [])
        minority = chairman_response.get("minority_report", [])
        risks = chairman_response.get("risks", [])
        next_steps = chairman_response.get("next_steps", [])
        unknowns = chairman_response.get("material_unknowns", [])
        confidence = chairman_response.get("confidence")
    else:
        if chairman.get("status") == "skipped" and obj["mode"] == "benchmark":
            recommendation = "Benchmark-only run completed; no paid chairman or quality judge was requested."
        elif chairman.get("status") == "skipped":
            recommendation = "Council evidence is ready for host reconciliation; no separate synthesis model was requested."
        else:
            recommendation = "Chairman unavailable; use the structured council dossier for a labeled host synthesis."
        why = []
        consensus = []
        disagreements = []
        minority = []
        risks = []
        next_steps = []
        unknowns = []
        confidence = None

    lines = ["## Recommendation", "", recommendation, "", "## Why", ""]
    lines.extend(_bullets(why, "No chairman rationale was available."))
    if confidence is not None:
        lines.append(f"- Chairman confidence: {confidence:.2f} (subjective).")
    lines.extend(["", "## Consensus and disagreement", ""])
    lines.extend(_bullets(consensus, "No validated consensus was recorded."))
    for item in disagreements:
        positions = "; ".join(item.get("positions", []))
        lines.append(
            f"- {item.get('issue', 'Disagreement')}: {positions}. "
            f"Resolution: {item.get('resolution', 'unresolved')}"
        )
    lines.extend(["", "## Minority report", ""])
    lines.extend(_bullets(minority, "No credible minority report was recorded."))
    lines.extend(["", "## Risks and next steps", ""])
    for item in risks:
        owner = f" Owner: {item.get('owner')}." if item.get("owner") else ""
        lines.append(
            f"- Risk: {item.get('risk')}. Mitigation: {item.get('mitigation')}.{owner}"
        )
    lines.extend(f"- Next: {item}" for item in next_steps)
    if not risks and not next_steps:
        lines.append("- No validated chairman risk or next step was available.")

    report = obj.get("execution_report", {})
    if report:
        lines.extend(["", "## Execution", "",
                      f"Backend: {report['backend']}; round {report['round_number']} of at most {report['maximum_rounds']}.",
                      report["stop_reason"]])
        for member in report["roster"]:
            lines.append(f"- {_escape(member['member_id'])}: {_escape(member['model'])} ({_escape(member['family'])}); {_escape(member['angle'])}.")

    lines.extend(["", "## Findings and evidence", ""])
    for candidate in obj["candidates"]:
        response = candidate["response"]
        lines.append(f"- {candidate['candidate_label']}: {response['recommendation']}")
        for evidence in response["evidence"]:
            lines.append(f"  - {evidence['evidence_id']}: {evidence['claim']}")
        for risk in response["risks"]:
            lines.append(f"  - {risk['severity']}: {risk['risk']} Mitigation: {risk['mitigation']}")
        unknowns.extend(response["unknowns"])
    lines.extend(["", "## Feedback disposition", ""])
    revision = obj.get("revision")
    if revision:
        for item in revision["response"]["accepted_feedback"]:
            lines.append(f"- Accepted: {item['feedback']} Reason/response: {item['response']}")
        for item in revision["response"]["rejected_feedback"]:
            lines.append(f"- Rejected: {item['feedback']} Reason: {item['reason']}")
    else:
        lines.append("- Host must record accepted/rejected material feedback with evidence-based reasons; no proposer reply was requested.")
    for critique in obj["critiques"]:
        lines.append(f"- Reviewer position: {critique['response']['revised_position']}")

    lines.extend(
        [
            "",
            "## Benchmark" if obj["mode"] == "benchmark" else "## Call accounting",
            "",
            "| Member | Family | Stage | Status | Latency | Tokens | Cost | Peer score |",
            "|---|---|---|---|---:|---:|---:|---:|",
        ]
    )
    calls = obj["benchmark"].get("calls", [])
    for index, call in enumerate(calls):
        member = (
            call.get("candidate_label")
            or call.get("member_id")
            or call.get("model")
            or f"Call {index + 1}"
        )
        tokens = call.get("total_tokens")
        cost = (
            f"${call['cost_usd']:.6f}"
            if call.get("cost_status") in {"reported", "estimated"}
            and isinstance(call.get("cost_usd"), (int, float))
            else "unknown"
        )
        quality = call.get("quality", {})
        peer_score = quality.get("peer_score") if isinstance(quality, dict) else None
        lines.append(
            "| "
            + " | ".join(
                _escape(value)
                for value in (
                    member,
                    call.get("family"),
                    call.get("stage"),
                    call.get("status"),
                    f"{call.get('latency_ms', 0)} ms",
                    tokens,
                    cost,
                    peer_score,
                )
            )
            + " |"
        )
    if not calls:
        lines.append("| No successful calls | unknown | unknown | failed | unknown | unknown | unknown | unknown |")

    limitations: list[str] = []
    if obj["status"] != "complete":
        limitations.append(f"Run status is {obj['status']}.")
    if not chairman_response and obj["mode"] != "benchmark":
        limitations.append("This is a review dossier, not a reconciled verdict. Host reconciliation must preserve dissent and disclose final verification evidence.")
    if obj["diversity"].get("status") != "verified":
        limitations.append(
            f"Council diversity status: {obj['diversity'].get('status', 'unverified')}."
        )
    limitations.extend(obj["diversity"].get("limitations", []))
    limitations.extend(obj["ranking"].get("limitations", []))
    limitations.extend(unknowns)
    limitations.extend(
        f"{failure.get('stage')}: {failure.get('message')}"
        for failure in obj["failures"]
    )
    if obj["usage"].get("cost_status") == "unknown":
        limitations.append("At least one call lacked trustworthy cost metadata; cost is unknown, not zero.")
    limitations.extend(obj["benchmark"].get("caveats", []))
    lines.extend(["", "## Limitations", ""])
    lines.extend(_bullets(list(dict.fromkeys(limitations)), "No material runtime limitation was recorded."))
    return "\n".join(lines).rstrip() + "\n"
