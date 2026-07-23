"""Typed, secret-safe runtime errors and stable process exit codes."""

from __future__ import annotations


class CouncilError(Exception):
    """Base class for expected errors that are safe to show to callers."""

    exit_code = 2
    category = "council_error"

    def __init__(self, message: str, *, details: list[str] | None = None) -> None:
        super().__init__(message)
        self.message = message
        self.details = details or []


class ValidationError(CouncilError):
    exit_code = 2
    category = "invalid_input"


class PreflightBlocked(CouncilError):
    exit_code = 3
    category = "preflight_blocked"


class AuthenticationBlocked(PreflightBlocked):
    category = "authentication"


class BudgetBlocked(CouncilError):
    exit_code = 4
    category = "budget"


class InsufficientMembers(CouncilError):
    exit_code = 5
    category = "insufficient_members"


class ChairmanUnavailable(CouncilError):
    exit_code = 6
    category = "chairman_unavailable"


class PersistenceFailure(CouncilError):
    exit_code = 7
    category = "persistence"


class Cancelled(CouncilError):
    exit_code = 130
    category = "cancelled"


class ClientError(CouncilError):
    """A sanitized gateway/client failure."""

    exit_code = 3
    retryable = False

    def __init__(
        self,
        message: str,
        *,
        category: str = "client",
        status: int | None = None,
        retryable: bool = False,
        retry_after: float | None = None,
    ) -> None:
        super().__init__(message)
        self.category = category
        self.status = status
        self.retryable = retryable
        self.retry_after = retry_after


class SemanticError(ClientError):
    def __init__(self, message: str = "Gateway response was not semantically valid") -> None:
        super().__init__(message, category="semantic", retryable=False)
