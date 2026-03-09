"""Tests for the intelligence provider module."""

from __future__ import annotations

import json
import os
import subprocess

import pytest

from eventgraph.intelligence import (
    ClaudeCliProvider,
    Config,
    ConfigError,
    IntelligenceError,
    OpenAICompatibleProvider,
    Provider,
    ProviderError,
    _detect_provider,
    _events_to_messages,
    _infer_provider_name,
    _parse_confidence,
    new_claude_cli_config,
    new_provider,
)
from eventgraph.types import Score


# ── Config validation ────────────────────────────────────────────────────

class TestConfig:
    """Tests for Config dataclass."""

    def test_defaults(self):
        cfg = Config(provider="openai", model="gpt-4o")
        assert cfg.provider == "openai"
        assert cfg.model == "gpt-4o"
        assert cfg.api_key == ""
        assert cfg.base_url == ""
        assert cfg.max_tokens == 1024
        assert cfg.temperature == 0.0
        assert cfg.system_prompt == ""

    def test_frozen(self):
        cfg = Config(provider="openai", model="gpt-4o")
        with pytest.raises(AttributeError):
            cfg.provider = "xai"  # type: ignore[misc]

    def test_all_fields(self):
        cfg = Config(
            provider="xai",
            model="grok-2",
            api_key="sk-test",
            base_url="https://custom.api.com/v1",
            max_tokens=2048,
            temperature=0.5,
            system_prompt="You are helpful.",
        )
        assert cfg.provider == "xai"
        assert cfg.model == "grok-2"
        assert cfg.api_key == "sk-test"
        assert cfg.base_url == "https://custom.api.com/v1"
        assert cfg.max_tokens == 2048
        assert cfg.temperature == 0.5
        assert cfg.system_prompt == "You are helpful."


# ── new_provider factory ─────────────────────────────────────────────────

class TestNewProvider:
    """Tests for the new_provider factory function."""

    def test_unknown_provider_raises(self):
        cfg = Config(provider="unknown_provider", model="test")
        with pytest.raises(ConfigError, match="unknown provider"):
            new_provider(cfg)

    def test_openai_creates_compatible_provider(self):
        cfg = Config(provider="openai", model="gpt-4o", api_key="sk-test")
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        assert p.name == "openai"
        assert p.model_id == "gpt-4o"

    def test_xai_creates_compatible_provider(self):
        cfg = Config(provider="xai", model="grok-2", api_key="xai-test")
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        assert p.name == "xai"

    def test_groq_creates_compatible_provider(self):
        cfg = Config(provider="groq", model="llama-3", api_key="gsk-test")
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        assert p.name == "groq"

    def test_together_creates_compatible_provider(self):
        cfg = Config(provider="together", model="mixtral", api_key="tog-test")
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        assert p.name == "together"

    def test_ollama_creates_compatible_provider(self):
        cfg = Config(provider="ollama", model="llama3")
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        assert p.name == "ollama"

    def test_claude_cli_creates_cli_provider(self):
        cfg = Config(provider="claude-cli", model="sonnet")
        p = new_provider(cfg)
        assert isinstance(p, ClaudeCliProvider)
        assert p.name == "claude-cli"
        assert p.model_id == "sonnet"

    def test_claude_cli_default_model(self):
        cfg = Config(provider="claude-cli")
        p = new_provider(cfg)
        assert isinstance(p, ClaudeCliProvider)
        assert p.model_id == "sonnet"

    def test_openai_compatible_requires_model(self):
        cfg = Config(provider="openai-compatible")
        with pytest.raises(ConfigError, match="requires a model"):
            new_provider(cfg)

    def test_azure_requires_base_url(self):
        cfg = Config(provider="azure", model="gpt-4o")
        with pytest.raises(ConfigError, match="requires a base_url"):
            new_provider(cfg)

    def test_azure_with_base_url(self):
        cfg = Config(
            provider="azure",
            model="gpt-4o",
            base_url="https://myinstance.openai.azure.com/v1",
            api_key="az-key",
        )
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        assert p.name == "azure"

    def test_max_tokens_defaults_to_1024(self):
        cfg = Config(provider="openai", model="gpt-4o", api_key="sk-test", max_tokens=0)
        p = new_provider(cfg)
        assert isinstance(p, OpenAICompatibleProvider)
        # The provider was created successfully with defaulted max_tokens.

    def test_satisfies_protocol(self):
        """Verify that concrete providers satisfy the Provider protocol."""
        cfg = Config(provider="openai", model="gpt-4o", api_key="sk-test")
        p = new_provider(cfg)
        assert isinstance(p, Provider)

        cli_cfg = Config(provider="claude-cli", model="sonnet")
        cp = new_provider(cli_cfg)
        assert isinstance(cp, Provider)


# ── Provider name inference ──────────────────────────────────────────────

class TestInferProviderName:
    """Tests for _infer_provider_name."""

    def test_openai(self):
        assert _infer_provider_name("https://api.openai.com/v1") == "openai"

    def test_xai(self):
        assert _infer_provider_name("https://api.x.ai/v1") == "xai"

    def test_groq(self):
        assert _infer_provider_name("https://api.groq.com/openai/v1") == "groq"

    def test_together(self):
        assert _infer_provider_name("https://api.together.xyz/v1") == "together"

    def test_azure(self):
        assert _infer_provider_name("https://myinstance.openai.azure.com/v1") == "azure"

    def test_localhost(self):
        assert _infer_provider_name("http://localhost:11434/v1") == "ollama"

    def test_127_0_0_1(self):
        assert _infer_provider_name("http://127.0.0.1:11434/v1") == "ollama"

    def test_fireworks(self):
        assert _infer_provider_name("https://api.fireworks.ai/v1") == "fireworks"

    def test_unknown(self):
        assert _infer_provider_name("https://custom.llm.example.com/v1") == "openai-compatible"

    def test_case_insensitive(self):
        assert _infer_provider_name("https://API.OPENAI.COM/V1") == "openai"


# ── Detect provider ─────────────────────────────────────────────────────

class TestDetectProvider:
    """Tests for _detect_provider environment auto-detection."""

    def test_detects_openai(self, monkeypatch):
        monkeypatch.delenv("XAI_API_KEY", raising=False)
        monkeypatch.delenv("GROQ_API_KEY", raising=False)
        monkeypatch.delenv("TOGETHER_API_KEY", raising=False)
        monkeypatch.delenv("OLLAMA_HOST", raising=False)
        monkeypatch.setenv("OPENAI_API_KEY", "sk-test")
        key, url, name = _detect_provider()
        assert key == "sk-test"
        assert "openai.com" in url
        assert name == "openai"

    def test_detects_xai(self, monkeypatch):
        monkeypatch.setenv("XAI_API_KEY", "xai-test")
        key, url, name = _detect_provider()
        assert key == "xai-test"
        assert "x.ai" in url
        assert name == "xai"

    def test_detects_groq(self, monkeypatch):
        monkeypatch.delenv("XAI_API_KEY", raising=False)
        monkeypatch.setenv("GROQ_API_KEY", "gsk-test")
        key, url, name = _detect_provider()
        assert name == "groq"

    def test_detects_together(self, monkeypatch):
        monkeypatch.delenv("XAI_API_KEY", raising=False)
        monkeypatch.delenv("GROQ_API_KEY", raising=False)
        monkeypatch.setenv("TOGETHER_API_KEY", "tog-test")
        key, url, name = _detect_provider()
        assert name == "together"

    def test_detects_ollama(self, monkeypatch):
        monkeypatch.delenv("XAI_API_KEY", raising=False)
        monkeypatch.delenv("GROQ_API_KEY", raising=False)
        monkeypatch.delenv("TOGETHER_API_KEY", raising=False)
        monkeypatch.delenv("OPENAI_API_KEY", raising=False)
        monkeypatch.setenv("OLLAMA_HOST", "http://myhost:11434")
        key, url, name = _detect_provider()
        assert name == "ollama"
        assert "myhost" in url

    def test_falls_back_to_generic(self, monkeypatch):
        monkeypatch.delenv("XAI_API_KEY", raising=False)
        monkeypatch.delenv("GROQ_API_KEY", raising=False)
        monkeypatch.delenv("TOGETHER_API_KEY", raising=False)
        monkeypatch.delenv("OPENAI_API_KEY", raising=False)
        monkeypatch.delenv("OLLAMA_HOST", raising=False)
        key, url, name = _detect_provider()
        assert name == "openai-compatible"
        assert key == ""


# ── Helper functions ─────────────────────────────────────────────────────

class TestHelpers:
    """Tests for helper functions."""

    def test_events_to_messages_empty(self):
        assert _events_to_messages([]) == ""

    def test_parse_confidence_returns_score(self):
        score = _parse_confidence(500)
        assert isinstance(score, Score)
        assert score.value == 0.7

    def test_parse_confidence_always_0_7(self):
        assert _parse_confidence(0).value == 0.7
        assert _parse_confidence(10000).value == 0.7


# ── Claude CLI config convenience ────────────────────────────────────────

class TestClaudeCliConfig:
    """Tests for new_claude_cli_config convenience function."""

    def test_default_model(self):
        cfg = new_claude_cli_config()
        assert cfg.provider == "claude-cli"
        assert cfg.model == "sonnet"

    def test_custom_model(self):
        cfg = new_claude_cli_config("opus")
        assert cfg.model == "opus"


# ── Integration tests (skipped without env vars) ────────────────────────

@pytest.mark.skipif(
    not os.environ.get("OPENAI_API_KEY"),
    reason="OPENAI_API_KEY not set",
)
class TestOpenAIIntegration:
    """Integration tests requiring a live OpenAI API key."""

    def test_simple_reason(self):
        cfg = Config(
            provider="openai",
            model="gpt-4o-mini",
            max_tokens=50,
        )
        p = new_provider(cfg)
        response = p.reason("Say 'hello' and nothing else.", [])
        assert response.content
        assert response.tokens_used > 0
        assert 0.0 <= response.confidence.value <= 1.0

    def test_invalid_api_key_returns_error(self):
        cfg = Config(
            provider="openai",
            model="gpt-4o-mini",
            api_key="sk-invalid-key-for-testing",
            max_tokens=50,
        )
        p = new_provider(cfg)
        with pytest.raises(ProviderError, match="API error"):
            p.reason("hello", [])


@pytest.mark.skipif(
    not os.environ.get("EVENTGRAPH_TEST_CLAUDE_CLI"),
    reason="EVENTGRAPH_TEST_CLAUDE_CLI not set",
)
class TestClaudeCliIntegration:
    """Integration tests requiring the claude CLI."""

    def test_simple_reason(self):
        cfg = Config(
            provider="claude-cli",
            model="sonnet",
        )
        p = new_provider(cfg)
        response = p.reason("Say 'hello' and nothing else.", [])
        assert response.content
        assert response.tokens_used > 0

    def test_with_system_prompt(self):
        cfg = Config(
            provider="claude-cli",
            model="sonnet",
            system_prompt="You are a helpful assistant. Respond with one word only.",
        )
        p = new_provider(cfg)
        response = p.reason("What color is the sky?", [])
        assert response.content


# ── Error hierarchy ──────────────────────────────────────────────────────

class TestErrorHierarchy:
    """Verify error types are properly structured."""

    def test_intelligence_error_is_eventgraph_error(self):
        from eventgraph.errors import EventGraphError
        assert issubclass(IntelligenceError, EventGraphError)

    def test_provider_error_is_intelligence_error(self):
        assert issubclass(ProviderError, IntelligenceError)

    def test_config_error_is_intelligence_error(self):
        assert issubclass(ConfigError, IntelligenceError)

    def test_provider_error_includes_name(self):
        err = ProviderError("openai", "something failed")
        assert "openai" in str(err)
        assert "something failed" in str(err)
        assert err.provider == "openai"
