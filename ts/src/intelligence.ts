/**
 * Intelligence provider module — factory + implementations for LLM providers.
 * Ported from the Go reference implementation (go/pkg/intelligence/).
 *
 * Supports: OpenAI-compatible (OpenAI, xAI, Groq, Together, Ollama, Azure),
 * and Claude CLI (shells out to `claude -p`).
 *
 * No external runtime dependencies — uses Node's built-in `fetch` and `child_process`.
 */
import { execSync } from "node:child_process";
import { Score } from "./types.js";
import { Response, type Intelligence } from "./decision.js";
import type { Event } from "./event.js";

// ── Provider interface ───────────────────────────────────────────────────

/** Extends Intelligence with metadata about the backing LLM. */
export interface Provider extends Intelligence {
  readonly name: string;
  readonly model: string;
}

// ── Config ───────────────────────────────────────────────────────────────

/** Configuration for creating a Provider via createProvider(). */
export interface Config {
  /** Provider name: "openai-compatible", "openai", "xai", "groq", "together", "ollama", "azure", "claude-cli" */
  provider: string;
  /** Model identifier (provider-specific, e.g. "gpt-4o", "claude-sonnet-4-6"). */
  model: string;
  /** API key for authentication. Falls back to env vars if omitted. */
  apiKey?: string;
  /** Override the default API endpoint. */
  baseUrl?: string;
  /** Maximum response tokens. Defaults to 1024. */
  maxTokens?: number;
  /** Temperature (0-2). Undefined means provider default. */
  temperature?: number;
  /** System prompt prepended to every reason() call. */
  systemPrompt?: string;
}

// ── Well-known base URLs ─────────────────────────────────────────────────

const OPENAI_BASE_URL = "https://api.openai.com/v1";
const XAI_BASE_URL = "https://api.x.ai/v1";
const GROQ_BASE_URL = "https://api.groq.com/openai/v1";
const TOGETHER_BASE_URL = "https://api.together.xyz/v1";
const OLLAMA_BASE_URL = "http://localhost:11434/v1";

// ── Factory ──────────────────────────────────────────────────────────────

/**
 * Creates a Provider from the given Config.
 * Throws on unknown provider or missing required fields.
 */
export function createProvider(config: Config): Provider {
  const maxTokens = config.maxTokens ?? 1024;

  switch (config.provider) {
    case "claude-cli":
      return createClaudeCliProvider(config, maxTokens);

    case "openai-compatible":
    case "openai":
    case "xai":
    case "groq":
    case "together":
    case "ollama":
    case "azure":
      return createOpenAICompatibleProvider(config, maxTokens);

    default:
      throw new Error(
        `Unknown provider: "${config.provider}" (supported: openai-compatible, openai, xai, groq, together, ollama, azure, claude-cli)`,
      );
  }
}

// ── Convenience config ───────────────────────────────────────────────────

/** Creates a Config for the Claude CLI provider. */
export function newClaudeCliConfig(model?: string): Config {
  return {
    provider: "claude-cli",
    model: model || "sonnet",
  };
}

// ── OpenAI-compatible provider ───────────────────────────────────────────

class OpenAICompatibleProvider implements Provider {
  readonly name: string;
  readonly model: string;

  private readonly baseUrl: string;
  private readonly apiKey: string;
  private readonly maxTokens: number;
  private readonly temperature: number | undefined;
  private readonly systemPrompt: string | undefined;

  constructor(
    name: string,
    model: string,
    baseUrl: string,
    apiKey: string,
    maxTokens: number,
    temperature: number | undefined,
    systemPrompt: string | undefined,
  ) {
    this.name = name;
    this.model = model;
    this.baseUrl = baseUrl.replace(/\/+$/, "");
    this.apiKey = apiKey;
    this.maxTokens = maxTokens;
    this.temperature = temperature;
    this.systemPrompt = systemPrompt;
    Object.freeze(this);
  }

  async reason(prompt: string, history: Event[]): Promise<Response> {
    const messages: OpenAIMessage[] = [];

    // System prompt.
    if (this.systemPrompt) {
      messages.push({ role: "system", content: this.systemPrompt });
    }

    // Event history as context.
    const historyText = eventsToMessages(history);
    if (historyText) {
      messages.push(
        { role: "user", content: historyText },
        { role: "assistant", content: "I understand the event history. What would you like me to reason about?" },
      );
    }

    // The actual prompt.
    messages.push({ role: "user", content: prompt });

    const reqBody: OpenAIRequest = {
      model: this.model,
      messages,
    };
    if (this.maxTokens > 0) {
      reqBody.max_tokens = this.maxTokens;
    }
    if (this.temperature !== undefined && this.temperature > 0) {
      reqBody.temperature = this.temperature;
    }

    const url = `${this.baseUrl}/chat/completions`;
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (this.apiKey) {
      headers["Authorization"] = `Bearer ${this.apiKey}`;
    }

    const resp = await fetch(url, {
      method: "POST",
      headers,
      body: JSON.stringify(reqBody),
    });

    const body = await resp.text();

    if (!resp.ok) {
      // Try to extract error message from JSON body.
      try {
        const errResp = JSON.parse(body) as OpenAIResponse;
        if (errResp.error?.message) {
          throw new Error(`OpenAI API error (HTTP ${resp.status}): ${errResp.error.message}`);
        }
      } catch (e) {
        if (e instanceof Error && e.message.startsWith("OpenAI API error")) {
          throw e;
        }
      }
      throw new Error(`OpenAI API error (HTTP ${resp.status}): ${body}`);
    }

    const result = JSON.parse(body) as OpenAIResponse;

    if (!result.choices || result.choices.length === 0) {
      throw new Error("OpenAI API returned no choices");
    }

    const content = result.choices[0].message.content;
    const tokensUsed = result.usage?.total_tokens ?? 0;
    const confidence = parseConfidence(tokensUsed);

    return new Response(content, confidence, tokensUsed);
  }
}

function createOpenAICompatibleProvider(config: Config, maxTokens: number): OpenAICompatibleProvider {
  if (!config.model) {
    throw new Error("OpenAI-compatible provider requires a model");
  }

  let apiKey = config.apiKey ?? "";
  let baseUrl = config.baseUrl ?? "";
  let providerName = config.provider;

  // Map shorthand provider names to default base URLs and env vars.
  switch (providerName) {
    case "openai":
      if (!baseUrl) baseUrl = OPENAI_BASE_URL;
      if (!apiKey) apiKey = process.env["OPENAI_API_KEY"] ?? "";
      break;
    case "xai":
      if (!baseUrl) baseUrl = XAI_BASE_URL;
      if (!apiKey) apiKey = process.env["XAI_API_KEY"] ?? "";
      break;
    case "groq":
      if (!baseUrl) baseUrl = GROQ_BASE_URL;
      if (!apiKey) apiKey = process.env["GROQ_API_KEY"] ?? "";
      break;
    case "together":
      if (!baseUrl) baseUrl = TOGETHER_BASE_URL;
      if (!apiKey) apiKey = process.env["TOGETHER_API_KEY"] ?? "";
      break;
    case "ollama":
      if (!baseUrl) {
        const host = process.env["OLLAMA_HOST"] ?? "http://localhost:11434";
        baseUrl = `${host}/v1`;
      }
      break;
    case "azure":
      // Azure requires explicit baseUrl and apiKey.
      break;
    case "openai-compatible":
      if (!apiKey && !baseUrl) {
        const detected = detectProvider();
        apiKey = detected.apiKey;
        baseUrl = detected.baseUrl;
        providerName = detected.name;
      }
      break;
  }

  if (!baseUrl) {
    baseUrl = OPENAI_BASE_URL;
  }

  // Infer provider name from base URL if still generic.
  if (providerName === "openai-compatible") {
    providerName = inferProviderName(baseUrl);
  }

  return new OpenAICompatibleProvider(
    providerName,
    config.model,
    baseUrl,
    apiKey,
    maxTokens,
    config.temperature,
    config.systemPrompt,
  );
}

// ── Claude CLI provider ──────────────────────────────────────────────────

/** JSON output from `claude -p --output-format json`. */
interface ClaudeCliResult {
  type: string;
  subtype: string;
  is_error: boolean;
  result: string;
  usage: {
    input_tokens: number;
    output_tokens: number;
  };
  total_cost_usd: number;
  stop_reason: string;
}

class ClaudeCliProvider implements Provider {
  readonly name = "claude-cli";
  readonly model: string;

  private readonly maxBudget: number;
  private readonly systemPrompt: string | undefined;
  private readonly claudePath: string;

  constructor(
    model: string,
    maxBudget: number,
    systemPrompt: string | undefined,
    claudePath: string,
  ) {
    this.model = model;
    this.maxBudget = maxBudget;
    this.systemPrompt = systemPrompt;
    this.claudePath = claudePath;
    Object.freeze(this);
  }

  reason(prompt: string, history: Event[]): Response {
    // Build the full prompt with history context.
    let fullPrompt = "";
    const historyText = eventsToMessages(history);
    if (historyText) {
      fullPrompt += historyText + "\n---\n\n";
    }
    fullPrompt += prompt;

    // Build command args.
    const args = [
      "-p",
      "--output-format", "json",
      "--model", this.model,
      "--max-budget-usd", this.maxBudget.toFixed(2),
      "--no-session-persistence",
    ];
    if (this.systemPrompt) {
      args.push("--system-prompt", this.systemPrompt);
    }

    const cmd = `${this.claudePath} ${args.map(escapeArg).join(" ")}`;

    // Remove CLAUDECODE env var for nested invocation.
    const env = { ...process.env };
    delete env["CLAUDECODE"];

    let stdout: string;
    try {
      stdout = execSync(cmd, {
        input: fullPrompt,
        env,
        encoding: "utf-8",
        maxBuffer: 10 * 1024 * 1024,
        stdio: ["pipe", "pipe", "pipe"],
      });
    } catch (err: unknown) {
      // Check if we got JSON output despite non-zero exit.
      const execErr = err as { stdout?: string; stderr?: string };
      if (execErr.stdout) {
        try {
          const result = JSON.parse(execErr.stdout) as ClaudeCliResult;
          if (result.result) {
            return this.resultToResponse(result);
          }
        } catch {
          // JSON parse failed, fall through to error.
        }
      }
      const stderr = execErr.stderr ?? "";
      throw new Error(`Claude CLI error: ${(err as Error).message}\nstderr: ${stderr}`);
    }

    const result = JSON.parse(stdout) as ClaudeCliResult;

    if (result.is_error) {
      throw new Error(`Claude CLI returned error: ${result.result} (subtype: ${result.subtype})`);
    }

    return this.resultToResponse(result);
  }

  private resultToResponse(result: ClaudeCliResult): Response {
    const tokensUsed = (result.usage?.input_tokens ?? 0) + (result.usage?.output_tokens ?? 0);
    const confidence = parseConfidence(tokensUsed);
    return new Response(result.result, confidence, tokensUsed);
  }
}

function createClaudeCliProvider(config: Config, _maxTokens: number): ClaudeCliProvider {
  const model = config.model || "sonnet";
  const claudePath = config.baseUrl || "claude";

  // Repurpose temperature as max budget hint (dollars).
  const maxBudget = (config.temperature && config.temperature > 0) ? config.temperature : 1.0;

  return new ClaudeCliProvider(model, maxBudget, config.systemPrompt, claudePath);
}

function escapeArg(arg: string): string {
  // Wrap in quotes if it contains spaces or special characters.
  if (/["\s]/.test(arg)) {
    return `"${arg.replace(/"/g, '\\"')}"`;
  }
  return arg;
}

// ── OpenAI API types ─────────────────────────────────────────────────────

interface OpenAIMessage {
  role: string;
  content: string;
}

interface OpenAIRequest {
  model: string;
  messages: OpenAIMessage[];
  max_tokens?: number;
  temperature?: number;
}

interface OpenAIResponse {
  id?: string;
  choices?: {
    message: OpenAIMessage;
    finish_reason: string;
  }[];
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
  error?: {
    message: string;
    type: string;
    code: string;
  };
}

// ── Helper functions ─────────────────────────────────────────────────────

/**
 * Converts event history into a simple text context
 * suitable for passing as conversation history to an LLM.
 * Caps at 20 events for token efficiency.
 */
export function eventsToMessages(events: Event[]): string {
  if (!events || events.length === 0) {
    return "";
  }

  const lines: string[] = ["Event history:"];
  const limit = Math.min(events.length, 20);

  for (let i = 0; i < limit; i++) {
    const ev = events[i];
    lines.push(`- [${ev.type.value}] ${ev.id.value} by ${ev.source.value}`);
  }

  if (events.length > 20) {
    lines.push(`... and ${events.length - 20} more events`);
  }

  return lines.join("\n") + "\n";
}

/**
 * Extracts a confidence score from token usage.
 * Higher token usage with a stop reason of "end_turn" suggests higher confidence.
 * This is a heuristic — real confidence requires model introspection.
 */
export function parseConfidence(_tokensUsed: number): Score {
  // Default to 0.7 — we can't truly measure confidence from outside the model.
  // Future: use log-probs or model self-assessment.
  return new Score(0.7);
}

/**
 * Auto-detects which OpenAI-compatible provider to use
 * based on environment variables.
 */
export function detectProvider(): { apiKey: string; baseUrl: string; name: string } {
  // Check in order of specificity.
  const xaiKey = process.env["XAI_API_KEY"];
  if (xaiKey) {
    return { apiKey: xaiKey, baseUrl: XAI_BASE_URL, name: "xai" };
  }

  const groqKey = process.env["GROQ_API_KEY"];
  if (groqKey) {
    return { apiKey: groqKey, baseUrl: GROQ_BASE_URL, name: "groq" };
  }

  const togetherKey = process.env["TOGETHER_API_KEY"];
  if (togetherKey) {
    return { apiKey: togetherKey, baseUrl: TOGETHER_BASE_URL, name: "together" };
  }

  const openaiKey = process.env["OPENAI_API_KEY"];
  if (openaiKey) {
    return { apiKey: openaiKey, baseUrl: OPENAI_BASE_URL, name: "openai" };
  }

  // Ollama doesn't need an API key.
  const ollamaHost = process.env["OLLAMA_HOST"];
  if (ollamaHost) {
    return { apiKey: "", baseUrl: `${ollamaHost}/v1`, name: "ollama" };
  }

  return { apiKey: "", baseUrl: "", name: "openai-compatible" };
}

/**
 * Derives a friendly provider name from the base URL.
 */
export function inferProviderName(baseUrl: string): string {
  const lower = baseUrl.toLowerCase();

  if (lower.includes("azure")) return "azure";
  if (lower.includes("openai.com")) return "openai";
  if (lower.includes("x.ai")) return "xai";
  if (lower.includes("groq.com")) return "groq";
  if (lower.includes("together.xyz")) return "together";
  if (lower.includes("localhost") || lower.includes("127.0.0.1")) return "ollama";
  if (lower.includes("fireworks")) return "fireworks";

  return "openai-compatible";
}
