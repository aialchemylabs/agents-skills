"""Small standard-library LiteLLM client with injectable offline transport."""

from __future__ import annotations

import json
import http.client
import inspect
import random
import socket
import ssl
import time
import urllib.error
from dataclasses import asdict, dataclass
from typing import Any, Callable, Mapping, Optional
from urllib.parse import urlsplit

from .errors import AuthenticationBlocked, ClientError, SemanticError


@dataclass(frozen=True)
class TransportResponse:
    status: int
    headers: dict[str, str]
    body: bytes


@dataclass(frozen=True)
class CallMetrics:
    model: str
    stage: str
    status: str
    latency_ms: int
    input_tokens: int | None
    output_tokens: int | None
    total_tokens: int | None
    cost_usd: float | None
    cost_status: str
    retry_count: int
    request_id: str | None
    http_status: int

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass(frozen=True)
class CallResult:
    content: str
    metrics: CallMetrics
    raw: dict[str, Any]


Transport = Callable[..., TransportResponse]
BeforeAttempt = Callable[..., None]


class _AmbiguousRequestFailure(Exception):
    """A request may have reached the provider; automatic retry is unsafe."""


def stdlib_transport(
    method: str,
    url: str,
    headers: Mapping[str, str],
    body: bytes | None,
    connect_timeout_seconds: float,
    request_timeout_seconds: float,
) -> TransportResponse:
    """Blocking HTTP transport with a distinct connect and remaining total timeout."""
    parsed = urlsplit(url)
    connection_type = (
        http.client.HTTPSConnection if parsed.scheme == "https" else http.client.HTTPConnection
    )
    connection = connection_type(
        parsed.hostname,
        parsed.port,
        timeout=connect_timeout_seconds,
    )
    started = time.monotonic()

    def remaining() -> float:
        value = request_timeout_seconds - (time.monotonic() - started)
        if value <= 0:
            raise TimeoutError("request timeout elapsed")
        return value

    path = parsed.path or "/"
    if parsed.query:
        path += "?" + parsed.query
    try:
        try:
            connection.connect()
        except ssl.SSLError as exc:
            raise ClientError(
                "LiteLLM TLS validation failed",
                category="tls",
                retryable=False,
            ) from exc
        except (TimeoutError, socket.timeout, OSError, http.client.HTTPException) as exc:
            raise ConnectionError("LiteLLM connection establishment failed") from exc
        if connection.sock is not None:
            connection.sock.settimeout(remaining())
        try:
            connection.request(method, path, body=body, headers=dict(headers))
            if connection.sock is not None:
                connection.sock.settimeout(remaining())
            response = connection.getresponse()
            if connection.sock is not None:
                connection.sock.settimeout(remaining())
            body_bytes = response.read(10_000_001)
            remaining()
        except (TimeoutError, socket.timeout):
            raise
        except (OSError, http.client.HTTPException) as exc:
            raise _AmbiguousRequestFailure from exc
        if len(body_bytes) > 10_000_000:
            raise ClientError(
                "LiteLLM response exceeded the bounded response size",
                category="semantic",
                retryable=False,
            )
        return TransportResponse(
            status=int(response.status),
            headers={key.lower(): value for key, value in response.getheaders()},
            body=body_bytes,
        )
    finally:
        connection.close()


# Backward-compatible import name; implementation no longer uses urllib.
urllib_transport = stdlib_transport


class LiteLLMClient:
    """OpenAI-compatible model-list/chat client with safe bounded retries."""

    def __init__(
        self,
        base_url: str,
        api_key: str,
        transport: Transport | None = None,
        connect_timeout_seconds: float = 10,
        request_timeout_seconds: float = 120,
        max_transient_retries: int = 1,
        before_attempt: BeforeAttempt | None = None,
        sleep: Callable[[float], None] | None = None,
    ) -> None:
        if not api_key:
            raise AuthenticationBlocked("A scoped LLM_COUNCIL_API_KEY is required")
        self.base_url = base_url.rstrip("/")
        self._api_key = api_key
        self._transport = transport or stdlib_transport
        self.connect_timeout_seconds = connect_timeout_seconds
        self.request_timeout_seconds = request_timeout_seconds
        self.max_transient_retries = max_transient_retries
        self._before_attempt = before_attempt
        self._sleep = sleep or time.sleep
        self.last_diagnostic: dict[str, Any] | None = None

    def _notify_before_attempt(
        self,
        stage: str,
        model: str,
        requested_output_tokens: int,
    ) -> None:
        if self._before_attempt is None:
            return
        arguments = (stage, model, requested_output_tokens)
        try:
            signature = inspect.signature(self._before_attempt)
            signature.bind(*arguments)
        except (TypeError, ValueError):
            self._before_attempt(stage, model)
            return
        self._before_attempt(*arguments)

    def _invoke_transport(
        self,
        method: str,
        url: str,
        headers: Mapping[str, str],
        body: bytes | None,
    ) -> TransportResponse:
        """Call the new six-argument contract, retaining five-argument fake compatibility."""
        arguments = (
            method,
            url,
            headers,
            body,
            self.connect_timeout_seconds,
            self.request_timeout_seconds,
        )
        try:
            signature = inspect.signature(self._transport)
            signature.bind(*arguments)
        except (TypeError, ValueError):
            return self._transport(
                method,
                url,
                headers,
                body,
                self.request_timeout_seconds,
            )
        return self._transport(*arguments)

    def _headers(self, *, json_body: bool = False) -> dict[str, str]:
        headers = {
            "Authorization": f"Bearer {self._api_key}",
            "Accept": "application/json",
            "User-Agent": "aialchemylabs-llm-council/0.1.0",
        }
        if json_body:
            headers["Content-Type"] = "application/json"
        return headers

    @staticmethod
    def _sanitized_headers(headers: Mapping[str, str]) -> dict[str, str]:
        return {
            key: ("<redacted>" if key.lower() == "authorization" else value)
            for key, value in headers.items()
        }

    @staticmethod
    def _retry_after(headers: Mapping[str, str]) -> float | None:
        raw = next(
            (value for key, value in headers.items() if key.lower() == "retry-after"),
            None,
        )
        if raw is None:
            return None
        try:
            return max(0.0, min(float(raw), 60.0))
        except ValueError:
            return None

    def _classify_http(self, response: TransportResponse) -> ClientError:
        status = response.status
        if status in {401, 403}:
            return AuthenticationBlocked(
                "LiteLLM authentication or authorization failed; verify the scoped key"
            )
        if status == 429:
            return ClientError(
                "LiteLLM rate limit exceeded",
                category="rate_limit",
                status=status,
                retryable=True,
                retry_after=self._retry_after(response.headers),
            )
        if status in {502, 503, 504}:
            return ClientError(
                "LiteLLM or its provider is temporarily unavailable",
                category="transient_http",
                status=status,
                retryable=True,
            )
        category = "model_not_found" if status == 404 else "http"
        return ClientError(
            f"LiteLLM request failed with HTTP {status}",
            category=category,
            status=status,
            retryable=False,
        )

    def _request(
        self,
        method: str,
        path: str,
        *,
        body: dict[str, Any] | None,
        stage: str,
        model: str,
        requested_output_tokens: int = 0,
        allow_retry: bool = True,
    ) -> tuple[TransportResponse, int, int]:
        encoded = (
            json.dumps(body, ensure_ascii=False, separators=(",", ":")).encode("utf-8")
            if body is not None
            else None
        )
        headers = self._headers(json_body=body is not None)
        started = time.monotonic()
        retry_count = 0
        while True:
            self._notify_before_attempt(stage, model, requested_output_tokens)
            try:
                response = self._invoke_transport(
                    method,
                    self.base_url + path,
                    headers,
                    encoded,
                )
            except (TimeoutError, socket.timeout) as exc:
                # An ambiguous read timeout may have billed; never retry automatically.
                raise ClientError(
                    "LiteLLM request timed out; billing status may be unknown",
                    category="timeout",
                    retryable=False,
                ) from exc
            except _AmbiguousRequestFailure as exc:
                raise ClientError(
                    "LiteLLM request failed after transmission; billing status may be unknown",
                    category="ambiguous_transport",
                    retryable=False,
                ) from exc
            except (
                ConnectionError,
                urllib.error.URLError,
                http.client.HTTPException,
                OSError,
            ) as exc:
                if isinstance(exc, urllib.error.URLError) and isinstance(
                    exc.reason, (TimeoutError, socket.timeout)
                ):
                    raise ClientError(
                        "LiteLLM request timed out; billing status may be unknown",
                        category="timeout",
                        retryable=False,
                    ) from exc
                error = ClientError(
                    "Could not establish a connection to LiteLLM",
                    category="connection",
                    retryable=True,
                )
                if not allow_retry or retry_count >= self.max_transient_retries:
                    raise error from exc
                delay = min(0.25 * (2**retry_count) + random.uniform(0.0, 0.1), 2.0)
                retry_count += 1
                self._sleep(delay)
                continue

            self.last_diagnostic = {
                "method": method,
                "url": self.base_url + path,
                "headers": self._sanitized_headers(headers),
                "status": response.status,
            }
            if 200 <= response.status < 300:
                elapsed_ms = round((time.monotonic() - started) * 1000)
                return response, retry_count, elapsed_ms
            error = self._classify_http(response)
            if (
                not allow_retry
                or not getattr(error, "retryable", False)
                or retry_count >= self.max_transient_retries
            ):
                raise error
            delay = (
                error.retry_after
                if error.retry_after is not None
                else min(0.25 * (2**retry_count) + random.uniform(0.0, 0.1), 2.0)
            )
            retry_count += 1
            self._sleep(delay)

    @staticmethod
    def _json_body(response: TransportResponse) -> dict[str, Any]:
        try:
            value = json.loads(response.body.decode("utf-8"))
        except (UnicodeDecodeError, json.JSONDecodeError) as exc:
            raise SemanticError("LiteLLM returned malformed JSON") from exc
        if not isinstance(value, dict):
            raise SemanticError("LiteLLM response must be a JSON object")
        return value

    def list_models(self) -> list[str]:
        response, _, _ = self._request(
            "GET",
            "/v1/models",
            body=None,
            stage="preflight",
            model="<model-list>",
            requested_output_tokens=0,
            allow_retry=True,
        )
        payload = self._json_body(response)
        data = payload.get("data")
        if not isinstance(data, list):
            raise SemanticError("LiteLLM model list is missing data")
        models: list[str] = []
        for row in data:
            if isinstance(row, dict) and isinstance(row.get("id"), str) and row["id"].strip():
                models.append(row["id"])
        return sorted(set(models))

    @staticmethod
    def _usage(payload: Mapping[str, Any]) -> tuple[int | None, int | None, int | None]:
        usage = payload.get("usage")
        if not isinstance(usage, Mapping):
            return None, None, None

        def integer(*names: str) -> int | None:
            for name in names:
                value = usage.get(name)
                if isinstance(value, int) and not isinstance(value, bool) and value >= 0:
                    return value
            return None

        return (
            integer("prompt_tokens", "input_tokens"),
            integer("completion_tokens", "output_tokens"),
            integer("total_tokens"),
        )

    @staticmethod
    def _cost(payload: Mapping[str, Any]) -> tuple[float | None, str]:
        candidates: list[Any] = [payload.get("response_cost"), payload.get("cost")]
        hidden = payload.get("_hidden_params")
        if isinstance(hidden, Mapping):
            candidates.extend((hidden.get("response_cost"), hidden.get("cost")))
        for candidate in candidates:
            if (
                isinstance(candidate, (int, float))
                and not isinstance(candidate, bool)
                and candidate >= 0
            ):
                return float(candidate), "reported"
        return None, "unknown"

    def chat_completion(
        self,
        *,
        model: str,
        system_prompt: str,
        user_payload: dict[str, Any] | str,
        max_tokens: int,
        temperature: float,
        stage: str,
    ) -> CallResult:
        user_content = (
            user_payload
            if isinstance(user_payload, str)
            else json.dumps(user_payload, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
        )
        body = {
            "model": model,
            "messages": [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_content},
            ],
            "temperature": temperature,
            "max_tokens": max_tokens,
        }
        response, retry_count, elapsed_ms = self._request(
            "POST",
            "/v1/chat/completions",
            body=body,
            stage=stage,
            model=model,
            requested_output_tokens=max_tokens,
            allow_retry=True,
        )
        payload = self._json_body(response)
        choices = payload.get("choices")
        if not isinstance(choices, list) or not choices:
            raise SemanticError("LiteLLM completion is missing choices")
        first = choices[0]
        if not isinstance(first, Mapping):
            raise SemanticError("LiteLLM completion choice is invalid")
        message = first.get("message")
        if not isinstance(message, Mapping) or not isinstance(message.get("content"), str):
            raise SemanticError("LiteLLM completion is missing message content")
        content = message["content"]
        if not content.strip():
            raise SemanticError("LiteLLM completion content is empty")
        input_tokens, output_tokens, total_tokens = self._usage(payload)
        if total_tokens is None and input_tokens is not None and output_tokens is not None:
            total_tokens = input_tokens + output_tokens
        cost, cost_status = self._cost(payload)
        response_headers = {key.lower(): value for key, value in response.headers.items()}
        request_id = (
            response_headers.get("x-litellm-request-id")
            or response_headers.get("x-request-id")
            or (payload.get("id") if isinstance(payload.get("id"), str) else None)
        )
        return CallResult(
            content=content,
            metrics=CallMetrics(
                model=model,
                stage=stage,
                status="success",
                latency_ms=elapsed_ms,
                input_tokens=input_tokens,
                output_tokens=output_tokens,
                total_tokens=total_tokens,
                cost_usd=cost,
                cost_status=cost_status,
                retry_count=retry_count,
                request_id=request_id,
                http_status=response.status,
            ),
            raw=payload,
        )
