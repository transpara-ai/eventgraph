"""Intelligence provider module — LLM abstraction with OpenAI-compatible and Claude CLI backends.

Provides a unified Provider protocol for interacting with language models.
Uses only stdlib (urllib.request, subprocess) — no external dependencies.

Supported providers:
- OpenAI, xAI/Grok, Groq, Together, Ollama, Azure (via OpenAI-compatible API)
- Claude CLI (shells out to the `claude` command)
"""

from __future__ import annotations

import json
import os
import subprocess
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from typing import Protocol, runtime_checkable

from .decision import Response
from .errors import EventGraphError, ValidationError
from .event import Event
from .types import Score


# ── Errors ───────────────────────────────────────────────────────────────

class IntelligenceError(EventGraphError):
    """Base for intelligence-related errors."""


class ProviderError(IntelligenceError):
    """An error from the underlying LLM provider."""

    def __init__(self, provider: str, message: str) -> None:
        self.provider = provider
        super().__init__(f"{provider}: {message}")


class ConfigError(IntelligenceError):
    """Invalid provider configuration."""


# ── Well-known base URLs ─────────────────────────────────────────────────

OPENAI_BASE_URL = "https://api.openai.com/v1"
XAI_BASE_URL = "https://api.x.ai/v1"
GROQ_BASE_URL = "https://api.groq.com/openai/v1"
TOGETHER_BASE_URL = "https://api.together.xyz/v1"
OLLAMA_BASE_URL = "http://localhost:11434/v1"

# ── Provider Protocol ────────────────────────────────────────────────────


@runtime_checkable
class Provider(Protocol):
    """Anything that can reason over a prompt and event history, with metadata."""

    def reason(self, prompt: str, history: list[Event]) -> Response:
        """Send a prompt with event history context and get a response."""
        ...

    @property
    def name(self) -> str:
        """Provider identifier (e.g., 'openai', 'claude-cli', 'ollama')."""
        ...

    @property
    def model_id(self) -> str:
        """Model identifier (e.g., 'gpt-4o', 'claude-sonnet-4-6')."""
        ...


# ── Config ───────────────────────────────────────────────────────────────

@dataclass(frozen=True, slots=True)
class Config:
    """Configuration for creating a Provider.

    Attributes:
        provider: Provider name — 'openai', 'openai-compatible', 'xai', 'groq',
                  'together', 'ollama', 'claude-cli'.
        model: Model identifier (provider-specific).
        api_key: API key for authentication. Falls back to environment variables
                 if empty (e.g., OPENAI_API_KEY).
        base_url: Override the default API endpoint. Useful for proxies, Ollama,
                  Azure, etc.
        max_tokens: Cap on response length. Defaults to 1024 if zero.
        temperature: Controls randomness. Zero means provider default.
        system_prompt: Prepended to every reason() call.
    """

    provider: str
    model: str = ""
    api_key: str = ""
    base_url: str = ""
    max_tokens: int = 1024
    temperature: float = 0.0
    system_prompt: str = ""


# ── Factory ──────────────────────────────────────────────────────────────

_SUPPORTED_PROVIDERS = frozenset({
    "openai", "openai-compatible", "xai", "groq",
    "together", "ollama", "azure", "claude-cli",
})


def new_provider(config: Config) -> Provider:
    """Create a Provider from the given Config.

    Raises:
        ConfigError: If the provider is unknown or required fields are missing.
    """
    if config.provider not in _SUPPORTED_PROVIDERS:
        raise ConfigError(
            f"unknown provider: {config.provider!r} "
            f"(supported: {', '.join(sorted(_SUPPORTED_PROVIDERS))})"
        )

    max_tokens = config.max_tokens if config.max_tokens > 0 else 1024

    if config.provider == "claude-cli":
        return _new_claude_cli_provider(config, max_tokens)

    # All other providers use the OpenAI-compatible path.
    return _new_openai_compatible_provider(config, max_tokens)


# ── Helper Functions ─────────────────────────────────────────────────────

def _events_to_messages(events: list[Event]) -> str:
    """Convert event history into text context for an LLM.

    Limits to the most recent 20 events to avoid overwhelming the context window.
    """
    if not events:
        return ""

    parts = ["Event history:"]
    for i, ev in enumerate(events):
        if i >= 20:
            parts.append(f"... and {len(events) - 20} more events")
            break
        parts.append(f"- [{ev.type.value}] {ev.id.value} by {ev.source.value}")

    return "\n".join(parts) + "\n"


def _parse_confidence(tokens_used: int) -> Score:
    """Derive a confidence score from token usage.

    This is a heuristic — real confidence requires model introspection.
    Default to 0.7; future versions may use log-probs or self-assessment.
    """
    return Score(0.7)


def _detect_provider() -> tuple[str, str, str]:
    """Check environment variables to auto-detect which OpenAI-compatible provider to use.

    Returns:
        Tuple of (api_key, base_url, provider_name).
    """
    if key := os.environ.get("XAI_API_KEY"):
        return key, XAI_BASE_URL, "xai"
    if key := os.environ.get("GROQ_API_KEY"):
        return key, GROQ_BASE_URL, "groq"
    if key := os.environ.get("TOGETHER_API_KEY"):
        return key, TOGETHER_BASE_URL, "together"
    if key := os.environ.get("OPENAI_API_KEY"):
        return key, OPENAI_BASE_URL, "openai"
    if os.environ.get("OLLAMA_HOST"):
        return "", os.environ["OLLAMA_HOST"] + "/v1", "ollama"
    return "", "", "openai-compatible"


def _infer_provider_name(base_url: str) -> str:
    """Derive a friendly provider name from the base URL."""
    lower = base_url.lower()
    if "azure" in lower:
        return "azure"
    if "openai.com" in lower:
        return "openai"
    if "x.ai" in lower:
        return "xai"
    if "groq.com" in lower:
        return "groq"
    if "together.xyz" in lower:
        return "together"
    if "localhost" in lower or "127.0.0.1" in lower:
        return "ollama"
    if "fireworks" in lower:
        return "fireworks"
    return "openai-compatible"


# ── OpenAI-Compatible Provider ───────────────────────────────────────────

class OpenAICompatibleProvider:
    """Provider using the OpenAI Chat Completions API.

    Compatible with: OpenAI, xAI/Grok, Ollama, Together, Azure OpenAI,
    Fireworks, Groq, and any OpenAI-compatible endpoint.

    Uses urllib.request from stdlib — no external HTTP dependencies.
    """

    def __init__(
        self,
        base_url: str,
        api_key: str,
        model: str,
        max_tokens: int,
        temperature: float,
        system_prompt: str,
        provider_name: str,
    ) -> None:
        self._base_url = base_url.rstrip("/")
        self._api_key = api_key
        self._model = model
        self._max_tokens = max_tokens
        self._temperature = temperature
        self._system_prompt = system_prompt
        self._provider_name = provider_name

    @property
    def name(self) -> str:
        """Provider identifier."""
        return self._provider_name

    @property
    def model_id(self) -> str:
        """Model identifier."""
        return self._model

    def reason(self, prompt: str, history: list[Event]) -> Response:
        """Send a prompt to the OpenAI-compatible API and return a Response.

        Raises:
            ProviderError: On HTTP errors or malformed responses.
        """
        messages: list[dict[str, str]] = []

        if self._system_prompt:
            messages.append({"role": "system", "content": self._system_prompt})

        history_text = _events_to_messages(history)
        if history_text:
            messages.append({"role": "user", "content": history_text})
            messages.append({
                "role": "assistant",
                "content": "I understand the event history. What would you like me to reason about?",
            })

        messages.append({"role": "user", "content": prompt})

        req_body: dict = {
            "model": self._model,
            "messages": messages,
        }
        if self._max_tokens > 0:
            req_body["max_tokens"] = self._max_tokens
        if self._temperature > 0:
            req_body["temperature"] = self._temperature

        body_bytes = json.dumps(req_body).encode("utf-8")
        url = self._base_url + "/chat/completions"

        req = urllib.request.Request(
            url,
            data=body_bytes,
            headers={
                "Content-Type": "application/json",
            },
            method="POST",
        )
        if self._api_key:
            req.add_header("Authorization", f"Bearer {self._api_key}")

        try:
            with urllib.request.urlopen(req) as resp:
                resp_body = resp.read()
        except urllib.error.HTTPError as e:
            error_body = e.read().decode("utf-8", errors="replace")
            # Try to extract structured error message.
            try:
                err_json = json.loads(error_body)
                if "error" in err_json and "message" in err_json["error"]:
                    msg = err_json["error"]["message"]
                    raise ProviderError(
                        self._provider_name,
                        f"API error (HTTP {e.code}): {msg}",
                    ) from e
            except (json.JSONDecodeError, KeyError):
                pass
            raise ProviderError(
                self._provider_name,
                f"API error (HTTP {e.code}): {error_body}",
            ) from e
        except urllib.error.URLError as e:
            raise ProviderError(
                self._provider_name,
                f"connection error: {e.reason}",
            ) from e

        try:
            result = json.loads(resp_body)
        except json.JSONDecodeError as e:
            raise ProviderError(
                self._provider_name,
                f"failed to parse response JSON: {e}",
            ) from e

        if "error" in result and result["error"]:
            msg = result["error"].get("message", str(result["error"]))
            raise ProviderError(self._provider_name, f"API error: {msg}")

        choices = result.get("choices", [])
        if not choices:
            raise ProviderError(self._provider_name, "API returned no choices")

        content = choices[0].get("message", {}).get("content", "")
        usage = result.get("usage", {})
        tokens_used = usage.get("total_tokens", 0)
        confidence = _parse_confidence(tokens_used)

        return Response(content=content, confidence=confidence, tokens_used=tokens_used)


def _new_openai_compatible_provider(config: Config, max_tokens: int) -> OpenAICompatibleProvider:
    """Create an OpenAICompatibleProvider from Config."""
    if not config.model:
        raise ConfigError("openai-compatible provider requires a model")

    api_key = config.api_key
    base_url = config.base_url
    provider_name = config.provider

    # Map shorthand provider names to default base URLs and env vars.
    if provider_name == "openai":
        if not base_url:
            base_url = OPENAI_BASE_URL
        if not api_key:
            api_key = os.environ.get("OPENAI_API_KEY", "")
    elif provider_name == "xai":
        if not base_url:
            base_url = XAI_BASE_URL
        if not api_key:
            api_key = os.environ.get("XAI_API_KEY", "")
    elif provider_name == "groq":
        if not base_url:
            base_url = GROQ_BASE_URL
        if not api_key:
            api_key = os.environ.get("GROQ_API_KEY", "")
    elif provider_name == "together":
        if not base_url:
            base_url = TOGETHER_BASE_URL
        if not api_key:
            api_key = os.environ.get("TOGETHER_API_KEY", "")
    elif provider_name == "ollama":
        if not base_url:
            host = os.environ.get("OLLAMA_HOST", "http://localhost:11434")
            base_url = host + "/v1"
    elif provider_name == "azure":
        if not base_url:
            raise ConfigError("azure provider requires a base_url")
        if not api_key:
            api_key = os.environ.get("AZURE_OPENAI_API_KEY", "")
    elif provider_name == "openai-compatible":
        if not api_key and not base_url:
            api_key, base_url, provider_name = _detect_provider()

    if not base_url:
        base_url = OPENAI_BASE_URL

    # Infer provider name from base URL if still generic.
    if provider_name == "openai-compatible":
        provider_name = _infer_provider_name(base_url)

    return OpenAICompatibleProvider(
        base_url=base_url,
        api_key=api_key,
        model=config.model,
        max_tokens=max_tokens,
        temperature=config.temperature,
        system_prompt=config.system_prompt,
        provider_name=provider_name,
    )


# ── Claude CLI Provider ──────────────────────────────────────────────────

class ClaudeCliProvider:
    """Provider that shells out to the `claude` CLI.

    Uses whatever authentication Claude Code already has (OAuth, API key, etc.)
    without requiring separate credentials.

    The claude CLI is invoked with `-p --output-format json` for structured output.
    """

    def __init__(
        self,
        model: str,
        max_budget: float,
        system_prompt: str,
        claude_path: str,
    ) -> None:
        self._model = model
        self._max_budget = max_budget
        self._system_prompt = system_prompt
        self._claude_path = claude_path

    @property
    def name(self) -> str:
        """Provider identifier."""
        return "claude-cli"

    @property
    def model_id(self) -> str:
        """Model identifier."""
        return self._model

    def reason(self, prompt: str, history: list[Event]) -> Response:
        """Run the claude CLI with the prompt and return a Response.

        Raises:
            ProviderError: On CLI errors or malformed output.
        """
        # Build the full prompt with history context.
        full_prompt_parts: list[str] = []
        history_text = _events_to_messages(history)
        if history_text:
            full_prompt_parts.append(history_text)
            full_prompt_parts.append("---\n")
        full_prompt_parts.append(prompt)
        full_prompt = "\n".join(full_prompt_parts)

        # Build command args.
        args = [
            self._claude_path,
            "-p",
            "--output-format", "json",
            "--model", self._model,
            "--max-budget-usd", f"{self._max_budget:.2f}",
            "--no-session-persistence",
        ]
        if self._system_prompt:
            args.extend(["--system-prompt", self._system_prompt])

        # Remove CLAUDECODE env var to allow nested invocation.
        env = {k: v for k, v in os.environ.items() if k != "CLAUDECODE"}

        try:
            result = subprocess.run(
                args,
                input=full_prompt,
                capture_output=True,
                text=True,
                env=env,
            )
        except FileNotFoundError as e:
            raise ProviderError("claude-cli", f"claude CLI not found: {e}") from e
        except OSError as e:
            raise ProviderError("claude-cli", f"failed to run claude CLI: {e}") from e

        # Try to parse JSON even on non-zero exit (budget exceeded but got result).
        if result.stdout:
            try:
                parsed = json.loads(result.stdout)
            except json.JSONDecodeError:
                parsed = None

            if parsed and isinstance(parsed, dict):
                if parsed.get("is_error"):
                    raise ProviderError(
                        "claude-cli",
                        f"CLI error: {parsed.get('result', 'unknown')} "
                        f"(subtype: {parsed.get('subtype', 'unknown')})",
                    )

                content = parsed.get("result", "")
                if content:
                    usage = parsed.get("usage", {})
                    tokens_used = (
                        usage.get("input_tokens", 0) + usage.get("output_tokens", 0)
                    )
                    confidence = _parse_confidence(tokens_used)
                    return Response(
                        content=content,
                        confidence=confidence,
                        tokens_used=tokens_used,
                    )

        if result.returncode != 0:
            raise ProviderError(
                "claude-cli",
                f"CLI exited with code {result.returncode}: {result.stderr}",
            )

        raise ProviderError("claude-cli", "CLI produced no parseable output")


def _new_claude_cli_provider(config: Config, max_tokens: int) -> ClaudeCliProvider:
    """Create a ClaudeCliProvider from Config."""
    model = config.model if config.model else "sonnet"

    claude_path = "claude"
    if config.base_url:
        # base_url repurposed as path to claude binary for testing.
        claude_path = config.base_url

    max_budget = 1.0  # default $1 per call
    if config.temperature > 0:
        # Repurpose temperature field as max budget hint (dollars).
        max_budget = config.temperature

    return ClaudeCliProvider(
        model=model,
        max_budget=max_budget,
        system_prompt=config.system_prompt,
        claude_path=claude_path,
    )


def new_claude_cli_config(model: str = "") -> Config:
    """Create a Config for the Claude CLI provider.

    This is a convenience function for the most common case.
    """
    if not model:
        model = "sonnet"
    return Config(provider="claude-cli", model=model)
