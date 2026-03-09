//! Intelligence provider module — LLM provider abstraction.
//!
//! Ports the Go `intelligence` package. Provides a `Provider` trait extending
//! `Intelligence` with metadata, a `Config` struct, and concrete providers:
//! OpenAI-compatible (gated on `intelligence` feature) and Claude CLI.

use std::fmt;

use crate::decision::{Intelligence, Response};
use crate::errors::{EventGraphError, Result};
use crate::event::Event;
use crate::types::Score;

// ── Provider trait ────────────────────────────────────────────────────

/// Extends `Intelligence` with metadata about the backing LLM.
pub trait Provider: Intelligence {
    /// Returns the provider identifier (e.g., "openai", "claude-cli", "ollama").
    fn name(&self) -> &str;

    /// Returns the model identifier (e.g., "claude-sonnet-4-6", "gpt-4o").
    fn model(&self) -> &str;
}

// ── Config ────────────────────────────────────────────────────────────

/// Configuration for creating a `Provider`.
#[derive(Debug, Clone)]
pub struct Config {
    /// Provider name: "claude-cli", "openai-compatible", "openai", "xai",
    /// "groq", "together", "ollama", "azure".
    pub provider: String,

    /// Model identifier (provider-specific).
    pub model: String,

    /// API key for authentication. If empty, the provider may fall back to
    /// environment variables (e.g., `OPENAI_API_KEY`).
    pub api_key: String,

    /// Base URL overrides the default API endpoint.
    /// Useful for proxies, Ollama (`http://localhost:11434/v1`), Azure, etc.
    pub base_url: String,

    /// Maximum response tokens. Defaults to 1024 if zero.
    pub max_tokens: usize,

    /// Temperature controls randomness. Zero means provider default.
    pub temperature: f64,

    /// System prompt prepended to every `reason` call.
    pub system_prompt: String,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            provider: String::new(),
            model: String::new(),
            api_key: String::new(),
            base_url: String::new(),
            max_tokens: 0,
            temperature: 0.0,
            system_prompt: String::new(),
        }
    }
}

impl fmt::Display for Config {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Config(provider={}, model={})", self.provider, self.model)
    }
}

// ── Factory ───────────────────────────────────────────────────────────

/// Creates a `Provider` from the given `Config`.
pub fn new(mut cfg: Config) -> Result<Box<dyn Provider>> {
    if cfg.max_tokens == 0 {
        cfg.max_tokens = 1024;
    }

    match cfg.provider.as_str() {
        "claude-cli" => {
            let p = ClaudeCliProvider::new(cfg)?;
            Ok(Box::new(p))
        }
        #[cfg(feature = "intelligence")]
        "openai-compatible" | "openai" | "xai" | "groq" | "together" | "ollama" | "azure" => {
            let p = OpenAIProvider::new(cfg)?;
            Ok(Box::new(p))
        }
        #[cfg(not(feature = "intelligence"))]
        "openai-compatible" | "openai" | "xai" | "groq" | "together" | "ollama" | "azure" => {
            Err(EventGraphError::GrammarViolation {
                detail: format!(
                    "provider {:?} requires the \"intelligence\" feature (ureq dependency)",
                    cfg.provider
                ),
            })
        }
        other => Err(EventGraphError::GrammarViolation {
            detail: format!(
                "unknown provider: {:?} (supported: claude-cli, openai-compatible, openai, xai, groq, together, ollama, azure)",
                other
            ),
        }),
    }
}

/// Creates a `Config` for the Claude CLI provider (convenience).
pub fn new_claude_cli_config(model: &str) -> Config {
    let model = if model.is_empty() { "sonnet" } else { model };
    Config {
        provider: "claude-cli".to_string(),
        model: model.to_string(),
        ..Config::default()
    }
}

// ── Shared helpers ────────────────────────────────────────────────────

/// Converts event history into a text summary for LLM context.
fn events_to_messages(events: &[Event]) -> String {
    if events.is_empty() {
        return String::new();
    }
    let mut buf = String::from("Event history:\n");
    for (i, ev) in events.iter().enumerate() {
        if i >= 20 {
            buf.push_str(&format!("... and {} more events\n", events.len() - 20));
            break;
        }
        buf.push_str(&format!(
            "- [{}] {} by {}\n",
            ev.event_type.value(),
            ev.id.value(),
            ev.source.value(),
        ));
    }
    buf
}

/// Heuristic confidence from token usage. Returns 0.7 — real confidence
/// requires model introspection.
fn parse_confidence(_tokens_used: usize) -> Score {
    Score::new(0.7).expect("0.7 is always valid")
}

// ── Claude CLI provider ───────────────────────────────────────────────

/// JSON output from `claude -p --output-format json`.
#[derive(Debug)]
struct ClaudeCliResult {
    result: String,
    is_error: bool,
    subtype: String,
    input_tokens: usize,
    output_tokens: usize,
}

/// Implements `Provider` by shelling out to the `claude` CLI.
/// Uses whatever authentication Claude Code already has.
pub struct ClaudeCliProvider {
    model: String,
    max_budget: f64,
    system_prompt: String,
    claude_path: String,
}

impl ClaudeCliProvider {
    fn new(cfg: Config) -> Result<Self> {
        let model = if cfg.model.is_empty() {
            "sonnet".to_string()
        } else {
            cfg.model
        };

        let claude_path = if cfg.base_url.is_empty() {
            "claude".to_string()
        } else {
            // BaseURL repurposed as path to claude binary for testing.
            cfg.base_url
        };

        // Verify claude is available.
        let check = std::process::Command::new("which")
            .arg(&claude_path)
            .stdout(std::process::Stdio::null())
            .stderr(std::process::Stdio::null())
            .status();

        // On Windows, try `where` if `which` fails.
        let found = match check {
            Ok(status) => status.success(),
            Err(_) => {
                match std::process::Command::new("where")
                    .arg(&claude_path)
                    .stdout(std::process::Stdio::null())
                    .stderr(std::process::Stdio::null())
                    .status()
                {
                    Ok(status) => status.success(),
                    Err(_) => false,
                }
            }
        };

        if !found {
            return Err(EventGraphError::GrammarViolation {
                detail: format!("claude CLI not found in PATH: {}", claude_path),
            });
        }

        let max_budget = if cfg.temperature > 0.0 {
            // Repurpose Temperature field as max budget hint (dollars).
            cfg.temperature
        } else {
            1.0 // default $1 per call
        };

        Ok(Self {
            model,
            max_budget,
            system_prompt: cfg.system_prompt,
            claude_path,
        })
    }

    fn parse_result(stdout: &[u8]) -> Result<ClaudeCliResult> {
        let parsed: serde_json::Value = serde_json::from_slice(stdout).map_err(|e| {
            EventGraphError::GrammarViolation {
                detail: format!(
                    "failed to parse claude CLI JSON output: {}\nraw: {}",
                    e,
                    String::from_utf8_lossy(stdout)
                ),
            }
        })?;

        Ok(ClaudeCliResult {
            result: parsed["result"].as_str().unwrap_or("").to_string(),
            is_error: parsed["is_error"].as_bool().unwrap_or(false),
            subtype: parsed["subtype"].as_str().unwrap_or("").to_string(),
            input_tokens: parsed["usage"]["input_tokens"].as_u64().unwrap_or(0) as usize,
            output_tokens: parsed["usage"]["output_tokens"].as_u64().unwrap_or(0) as usize,
        })
    }
}

impl fmt::Debug for ClaudeCliProvider {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("ClaudeCliProvider")
            .field("model", &self.model)
            .field("claude_path", &self.claude_path)
            .finish()
    }
}

impl Intelligence for ClaudeCliProvider {
    fn reason(&self, prompt: &str, history: &[Event]) -> Result<Response> {
        let mut full_prompt = String::new();
        let history_text = events_to_messages(history);
        if !history_text.is_empty() {
            full_prompt.push_str(&history_text);
            full_prompt.push_str("\n---\n\n");
        }
        full_prompt.push_str(prompt);

        let mut args = vec![
            "-p".to_string(),
            "--output-format".to_string(),
            "json".to_string(),
            "--model".to_string(),
            self.model.clone(),
            "--max-budget-usd".to_string(),
            format!("{:.2}", self.max_budget),
            "--no-session-persistence".to_string(),
        ];
        if !self.system_prompt.is_empty() {
            args.push("--system-prompt".to_string());
            args.push(self.system_prompt.clone());
        }

        let mut cmd = std::process::Command::new(&self.claude_path);
        cmd.args(&args);
        cmd.stdin(std::process::Stdio::piped());
        cmd.stdout(std::process::Stdio::piped());
        cmd.stderr(std::process::Stdio::piped());

        // Unset CLAUDECODE to allow nested invocation.
        cmd.env_remove("CLAUDECODE");

        use std::io::Write;
        let mut child = cmd.spawn().map_err(|e| EventGraphError::GrammarViolation {
            detail: format!("failed to spawn claude CLI: {}", e),
        })?;

        if let Some(ref mut stdin) = child.stdin {
            let _ = stdin.write_all(full_prompt.as_bytes());
        }
        // Close stdin so claude reads EOF.
        drop(child.stdin.take());

        let output = child.wait_with_output().map_err(|e| {
            EventGraphError::GrammarViolation {
                detail: format!("claude CLI error: {}", e),
            }
        })?;

        if !output.status.success() {
            // Check if we got JSON output despite non-zero exit.
            if !output.stdout.is_empty() {
                if let Ok(result) = Self::parse_result(&output.stdout) {
                    if !result.result.is_empty() {
                        let tokens_used = result.input_tokens + result.output_tokens;
                        let confidence = parse_confidence(tokens_used);
                        return Ok(Response {
                            content: result.result,
                            confidence,
                            tokens_used,
                        });
                    }
                }
            }
            return Err(EventGraphError::GrammarViolation {
                detail: format!(
                    "claude CLI error (exit {})\nstderr: {}",
                    output.status,
                    String::from_utf8_lossy(&output.stderr),
                ),
            });
        }

        let result = Self::parse_result(&output.stdout)?;

        if result.is_error {
            return Err(EventGraphError::GrammarViolation {
                detail: format!(
                    "claude CLI returned error: {} (subtype: {})",
                    result.result, result.subtype
                ),
            });
        }

        let tokens_used = result.input_tokens + result.output_tokens;
        let confidence = parse_confidence(tokens_used);

        Ok(Response {
            content: result.result,
            confidence,
            tokens_used,
        })
    }
}

impl Provider for ClaudeCliProvider {
    fn name(&self) -> &str { "claude-cli" }
    fn model(&self) -> &str { &self.model }
}

// ── OpenAI-compatible provider ────────────────────────────────────────

/// Well-known OpenAI-compatible base URLs.
#[cfg(feature = "intelligence")]
const OPENAI_BASE_URL: &str = "https://api.openai.com/v1";
#[cfg(feature = "intelligence")]
const XAI_BASE_URL: &str = "https://api.x.ai/v1";
#[cfg(feature = "intelligence")]
const GROQ_BASE_URL: &str = "https://api.groq.com/openai/v1";
#[cfg(feature = "intelligence")]
const TOGETHER_BASE_URL: &str = "https://api.together.xyz/v1";
#[cfg(feature = "intelligence")]
#[allow(dead_code)]
const OLLAMA_BASE_URL: &str = "http://localhost:11434/v1";

/// Infers a friendly provider name from the base URL.
#[cfg_attr(not(feature = "intelligence"), allow(dead_code))]
fn infer_provider_name(base_url: &str) -> &'static str {
    let lower = base_url.to_lowercase();
    if lower.contains("azure") {
        "azure"
    } else if lower.contains("openai.com") {
        "openai"
    } else if lower.contains("x.ai") {
        "xai"
    } else if lower.contains("groq.com") {
        "groq"
    } else if lower.contains("together.xyz") {
        "together"
    } else if lower.contains("localhost") || lower.contains("127.0.0.1") {
        "ollama"
    } else if lower.contains("fireworks") {
        "fireworks"
    } else {
        "openai-compatible"
    }
}

/// Auto-detect from environment variables which OpenAI-compatible provider to use.
#[cfg(feature = "intelligence")]
fn detect_openai_provider() -> (String, String, &'static str) {
    if let Ok(key) = std::env::var("XAI_API_KEY") {
        if !key.is_empty() {
            return (key, XAI_BASE_URL.to_string(), "xai");
        }
    }
    if let Ok(key) = std::env::var("GROQ_API_KEY") {
        if !key.is_empty() {
            return (key, GROQ_BASE_URL.to_string(), "groq");
        }
    }
    if let Ok(key) = std::env::var("TOGETHER_API_KEY") {
        if !key.is_empty() {
            return (key, TOGETHER_BASE_URL.to_string(), "together");
        }
    }
    if let Ok(key) = std::env::var("OPENAI_API_KEY") {
        if !key.is_empty() {
            return (key, OPENAI_BASE_URL.to_string(), "openai");
        }
    }
    if let Ok(host) = std::env::var("OLLAMA_HOST") {
        if !host.is_empty() {
            return (String::new(), format!("{}/v1", host), "ollama");
        }
    }
    (String::new(), String::new(), "openai-compatible")
}

#[cfg(feature = "intelligence")]
pub struct OpenAIProvider {
    base_url: String,
    api_key: String,
    model: String,
    max_tokens: usize,
    temperature: f64,
    system_prompt: String,
    provider_name: String,
}

#[cfg(feature = "intelligence")]
impl OpenAIProvider {
    fn new(cfg: Config) -> Result<Self> {
        if cfg.model.is_empty() {
            return Err(EventGraphError::GrammarViolation {
                detail: "openai-compatible provider requires a model".to_string(),
            });
        }

        let mut api_key = cfg.api_key;
        let mut base_url = cfg.base_url;
        let mut provider_name = cfg.provider.clone();

        // Map shorthand provider names to default base URLs and env vars.
        match provider_name.as_str() {
            "openai" => {
                if base_url.is_empty() {
                    base_url = OPENAI_BASE_URL.to_string();
                }
                if api_key.is_empty() {
                    api_key = std::env::var("OPENAI_API_KEY").unwrap_or_default();
                }
            }
            "xai" => {
                if base_url.is_empty() {
                    base_url = XAI_BASE_URL.to_string();
                }
                if api_key.is_empty() {
                    api_key = std::env::var("XAI_API_KEY").unwrap_or_default();
                }
            }
            "groq" => {
                if base_url.is_empty() {
                    base_url = GROQ_BASE_URL.to_string();
                }
                if api_key.is_empty() {
                    api_key = std::env::var("GROQ_API_KEY").unwrap_or_default();
                }
            }
            "together" => {
                if base_url.is_empty() {
                    base_url = TOGETHER_BASE_URL.to_string();
                }
                if api_key.is_empty() {
                    api_key = std::env::var("TOGETHER_API_KEY").unwrap_or_default();
                }
            }
            "ollama" => {
                if base_url.is_empty() {
                    let host = std::env::var("OLLAMA_HOST")
                        .unwrap_or_else(|_| "http://localhost:11434".to_string());
                    base_url = format!("{}/v1", host);
                }
            }
            "openai-compatible" => {
                if api_key.is_empty() && base_url.is_empty() {
                    let (detected_key, detected_url, detected_name) = detect_openai_provider();
                    api_key = detected_key;
                    base_url = detected_url;
                    provider_name = detected_name.to_string();
                }
            }
            "azure" => {
                // Azure uses the provided base_url and api_key directly.
            }
            _ => {}
        }

        if base_url.is_empty() {
            base_url = OPENAI_BASE_URL.to_string();
        }

        // Infer provider name from base URL if still generic.
        if provider_name == "openai-compatible" {
            provider_name = infer_provider_name(&base_url).to_string();
        }

        let base_url = base_url.trim_end_matches('/').to_string();

        Ok(Self {
            base_url,
            api_key,
            model: cfg.model,
            max_tokens: cfg.max_tokens,
            temperature: cfg.temperature,
            system_prompt: cfg.system_prompt,
            provider_name,
        })
    }
}

#[cfg(feature = "intelligence")]
impl fmt::Debug for OpenAIProvider {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("OpenAIProvider")
            .field("provider_name", &self.provider_name)
            .field("model", &self.model)
            .field("base_url", &self.base_url)
            .finish()
    }
}

#[cfg(feature = "intelligence")]
impl Intelligence for OpenAIProvider {
    fn reason(&self, prompt: &str, history: &[Event]) -> Result<Response> {
        let mut messages = Vec::new();

        // System prompt.
        if !self.system_prompt.is_empty() {
            messages.push(serde_json::json!({
                "role": "system",
                "content": self.system_prompt,
            }));
        }

        // Event history as context.
        let history_text = events_to_messages(history);
        if !history_text.is_empty() {
            messages.push(serde_json::json!({
                "role": "user",
                "content": history_text,
            }));
            messages.push(serde_json::json!({
                "role": "assistant",
                "content": "I understand the event history. What would you like me to reason about?",
            }));
        }

        // The actual prompt.
        messages.push(serde_json::json!({
            "role": "user",
            "content": prompt,
        }));

        let mut req_body = serde_json::json!({
            "model": self.model,
            "messages": messages,
        });

        if self.max_tokens > 0 {
            req_body["max_tokens"] = serde_json::json!(self.max_tokens);
        }
        if self.temperature > 0.0 {
            req_body["temperature"] = serde_json::json!(self.temperature);
        }

        let url = format!("{}/chat/completions", self.base_url);

        let mut request = ureq::post(&url)
            .header("Content-Type", "application/json");

        if !self.api_key.is_empty() {
            request = request.header("Authorization", &format!("Bearer {}", self.api_key));
        }

        let mut resp = request
            .send_json(&req_body)
            .map_err(|e| EventGraphError::GrammarViolation {
                detail: format!("openai API request error: {}", e),
            })?;

        let body: serde_json::Value = resp.body_mut().read_json().map_err(|e| {
            EventGraphError::GrammarViolation {
                detail: format!("openai API response parse error: {}", e),
            }
        })?;

        // Check for API error in body.
        if let Some(err) = body.get("error") {
            let msg = err["message"].as_str().unwrap_or("unknown error");
            return Err(EventGraphError::GrammarViolation {
                detail: format!("openai API error: {}", msg),
            });
        }

        let choices = body["choices"].as_array().ok_or_else(|| {
            EventGraphError::GrammarViolation {
                detail: "openai API returned no choices".to_string(),
            }
        })?;

        if choices.is_empty() {
            return Err(EventGraphError::GrammarViolation {
                detail: "openai API returned no choices".to_string(),
            });
        }

        let content = choices[0]["message"]["content"]
            .as_str()
            .unwrap_or("")
            .to_string();

        let tokens_used = body["usage"]["total_tokens"].as_u64().unwrap_or(0) as usize;
        let confidence = parse_confidence(tokens_used);

        Ok(Response {
            content,
            confidence,
            tokens_used,
        })
    }
}

#[cfg(feature = "intelligence")]
impl Provider for OpenAIProvider {
    fn name(&self) -> &str { &self.provider_name }
    fn model(&self) -> &str { &self.model }
}

// ── Tests ─────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;

    // ── Unit tests — no API calls ─────────────────────────────────────

    #[test]
    fn test_new_unknown_provider() {
        let result = new(Config {
            provider: "unknown".to_string(),
            model: "some-model".to_string(),
            ..Config::default()
        });
        assert!(result.is_err(), "expected error for unknown provider");
    }

    #[test]
    fn test_config_requires_model_for_openai_compatible() {
        // Without the intelligence feature, this returns a feature-required error.
        // With the feature, it returns a model-required error.
        // Either way, it should error.
        let result = new(Config {
            provider: "openai-compatible".to_string(),
            ..Config::default()
        });
        assert!(result.is_err(), "expected error when model is empty");
    }

    #[test]
    fn test_default_max_tokens() {
        let cfg = Config {
            provider: "claude-cli".to_string(),
            model: "sonnet".to_string(),
            ..Config::default()
        };
        assert_eq!(cfg.max_tokens, 0, "default max_tokens should be 0 before factory");
    }

    #[test]
    fn test_new_claude_cli_config() {
        let cfg = new_claude_cli_config("haiku");
        assert_eq!(cfg.provider, "claude-cli");
        assert_eq!(cfg.model, "haiku");
    }

    #[test]
    fn test_new_claude_cli_config_default_model() {
        let cfg = new_claude_cli_config("");
        assert_eq!(cfg.model, "sonnet");
    }

    #[test]
    fn test_events_to_messages_empty() {
        let result = events_to_messages(&[]);
        assert!(result.is_empty());
    }

    #[test]
    fn test_parse_confidence_returns_0_7() {
        let score = parse_confidence(100);
        assert!((score.value() - 0.7).abs() < f64::EPSILON);
    }

    #[test]
    fn test_infer_provider_name_from_urls() {
        assert_eq!(infer_provider_name("https://api.openai.com/v1"), "openai");
        assert_eq!(infer_provider_name("https://api.x.ai/v1"), "xai");
        assert_eq!(infer_provider_name("https://api.groq.com/openai/v1"), "groq");
        assert_eq!(infer_provider_name("https://api.together.xyz/v1"), "together");
        assert_eq!(infer_provider_name("http://localhost:11434/v1"), "ollama");
        assert_eq!(infer_provider_name("http://127.0.0.1:11434/v1"), "ollama");
        assert_eq!(
            infer_provider_name("https://mydeployment.azure.openai.com/v1"),
            "azure"
        );
        assert_eq!(
            infer_provider_name("https://custom.example.com/v1"),
            "openai-compatible"
        );
        assert_eq!(
            infer_provider_name("https://api.fireworks.ai/v1"),
            "fireworks"
        );
    }

    #[cfg(feature = "intelligence")]
    mod openai_tests {
        use super::*;

        #[test]
        fn test_new_openai_compatible_requires_model() {
            let result = new(Config {
                provider: "openai-compatible".to_string(),
                ..Config::default()
            });
            assert!(result.is_err(), "expected error when model is empty");
        }

        #[test]
        fn test_new_openai_compatible_success() {
            let p = new(Config {
                provider: "openai-compatible".to_string(),
                model: "gpt-4o".to_string(),
                api_key: "test-key-not-real".to_string(),
                ..Config::default()
            })
            .expect("should create provider");
            assert_eq!(p.model(), "gpt-4o");
        }

        #[test]
        fn test_openai_compatible_infers_provider_name() {
            let cases = vec![
                ("openai", "", "openai"),
                ("xai", "", "xai"),
                ("groq", "", "groq"),
                ("together", "", "together"),
                ("ollama", "", "ollama"),
                ("openai-compatible", "https://api.openai.com/v1", "openai"),
                ("openai-compatible", "https://api.x.ai/v1", "xai"),
                ("openai-compatible", "https://api.groq.com/openai/v1", "groq"),
                ("openai-compatible", "https://api.together.xyz/v1", "together"),
                ("openai-compatible", "http://localhost:11434/v1", "ollama"),
                (
                    "openai-compatible",
                    "https://mydeployment.azure.openai.com/v1",
                    "azure",
                ),
                (
                    "openai-compatible",
                    "https://custom.example.com/v1",
                    "openai-compatible",
                ),
            ];

            for (provider, base_url, want_name) in cases {
                let p = new(Config {
                    provider: provider.to_string(),
                    model: "test-model".to_string(),
                    api_key: "test-key".to_string(),
                    base_url: base_url.to_string(),
                    ..Config::default()
                })
                .unwrap_or_else(|e| {
                    panic!("unexpected error for provider={provider}, base_url={base_url}: {e}")
                });
                assert_eq!(
                    p.name(),
                    want_name,
                    "provider={provider}, base_url={base_url}: name={}, want={want_name}",
                    p.name()
                );
            }
        }

        #[test]
        fn test_openai_compatible_with_all_options() {
            let p = new(Config {
                provider: "openai-compatible".to_string(),
                model: "grok-3".to_string(),
                api_key: "test-key".to_string(),
                base_url: "https://api.x.ai/v1".to_string(),
                max_tokens: 2048,
                temperature: 0.7,
                system_prompt: "You are a helpful assistant.".to_string(),
            })
            .expect("should create provider");
            assert_eq!(p.name(), "xai");
            assert_eq!(p.model(), "grok-3");
        }
    }

    // ── Integration tests — gated on env vars ─────────────────────────

    #[cfg(feature = "intelligence")]
    #[test]
    fn test_integration_openai_compatible_reason() {
        let api_key = std::env::var("OPENAI_API_KEY").unwrap_or_default();
        if api_key.is_empty() {
            eprintln!("OPENAI_API_KEY not set — skipping integration test");
            return;
        }

        let p = new(Config {
            provider: "openai-compatible".to_string(),
            model: "gpt-4o-mini".to_string(),
            api_key,
            max_tokens: 100,
            ..Config::default()
        })
        .expect("should create provider");

        let resp = p.reason("Reply with exactly one word: hello", &[]);
        assert!(resp.is_ok(), "Reason failed: {:?}", resp.err());
        let resp = resp.unwrap();
        assert!(!resp.content.is_empty(), "response content is empty");
        assert!(resp.tokens_used > 0, "tokens used is 0");
    }

    #[cfg(feature = "intelligence")]
    #[test]
    fn test_integration_openai_compatible_invalid_key() {
        let p = new(Config {
            provider: "openai-compatible".to_string(),
            model: "gpt-4o-mini".to_string(),
            api_key: "sk-invalid-key-for-testing".to_string(),
            max_tokens: 50,
            ..Config::default()
        })
        .expect("should create provider");

        let resp = p.reason("hello", &[]);
        assert!(resp.is_err(), "expected error with invalid API key");
    }

    #[test]
    fn test_integration_claude_cli_reason() {
        if std::env::var("EVENTGRAPH_TEST_CLAUDE_CLI").unwrap_or_default().is_empty() {
            eprintln!("EVENTGRAPH_TEST_CLAUDE_CLI not set — skipping Claude CLI integration test");
            return;
        }

        let p = new(new_claude_cli_config("sonnet")).expect("should create provider");
        assert_eq!(p.name(), "claude-cli");
        assert_eq!(p.model(), "sonnet");

        let resp = p.reason("Reply with exactly one word: hello", &[]);
        assert!(resp.is_ok(), "Reason failed: {:?}", resp.err());
        let resp = resp.unwrap();
        assert!(!resp.content.is_empty(), "response content is empty");
    }
}
