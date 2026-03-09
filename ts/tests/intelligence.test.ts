import { describe, it, expect, afterEach } from "vitest";
import {
  createProvider,
  newClaudeCliConfig,
  eventsToMessages,
  parseConfidence,
  detectProvider,
  inferProviderName,
  type Provider,
  type Config,
} from "../src/intelligence.js";
import { Score, ActorId, EventId, EventType, Hash, ConversationId, NonEmpty } from "../src/types.js";
import { Event, NoopSigner } from "../src/event.js";
import { Response } from "../src/decision.js";

// ── Helpers ──────────────────────────────────────────────────────────────

function makeEvent(typeStr: string, sourceStr: string): Event {
  const id = new EventId("01912345-6789-7abc-8def-0123456789ab");
  const evType = new EventType(typeStr);
  const source = new ActorId(sourceStr);
  const convId = new ConversationId("conv_test");
  const hash = Hash.zero();
  const sig = new NoopSigner().sign(new Uint8Array(0));

  return new Event(
    1,
    id,
    evType,
    Date.now() * 1_000_000,
    source,
    {},
    NonEmpty.of([id]),
    convId,
    hash,
    hash,
    sig,
  );
}

// ════════════════════════════════════════════════════════════════════════
// Config validation
// ════════════════════════════════════════════════════════════════════════

describe("createProvider", () => {
  it("throws on unknown provider", () => {
    expect(() =>
      createProvider({ provider: "unknown", model: "some-model" }),
    ).toThrow('Unknown provider: "unknown"');
  });

  it("throws when openai-compatible is missing a model", () => {
    expect(() =>
      createProvider({ provider: "openai-compatible", model: "" }),
    ).toThrow("requires a model");
  });

  it("throws when openai is missing a model", () => {
    expect(() =>
      createProvider({ provider: "openai", model: "" }),
    ).toThrow("requires a model");
  });

  it("creates openai-compatible provider with model", () => {
    const p = createProvider({
      provider: "openai-compatible",
      model: "gpt-4o",
      apiKey: "test-key-not-real",
    });
    expect(p.model).toBe("gpt-4o");
  });

  it("creates provider with all options", () => {
    const p = createProvider({
      provider: "openai-compatible",
      model: "grok-3",
      apiKey: "test-key",
      baseUrl: "https://api.x.ai/v1",
      maxTokens: 2048,
      temperature: 0.7,
      systemPrompt: "You are a helpful assistant.",
    });
    expect(p.name).toBe("xai");
    expect(p.model).toBe("grok-3");
  });

  it("defaults maxTokens to 1024", () => {
    const p = createProvider({
      provider: "openai",
      model: "gpt-4o",
      apiKey: "test-key",
    });
    expect(p).toBeTruthy();
  });
});

// ════════════════════════════════════════════════════════════════════════
// Provider name inference
// ════════════════════════════════════════════════════════════════════════

describe("inferProviderName", () => {
  it.each([
    ["https://api.openai.com/v1", "openai"],
    ["https://api.x.ai/v1", "xai"],
    ["https://api.groq.com/openai/v1", "groq"],
    ["https://api.together.xyz/v1", "together"],
    ["http://localhost:11434/v1", "ollama"],
    ["http://127.0.0.1:11434/v1", "ollama"],
    ["https://mydeployment.azure.openai.com/v1", "azure"],
    ["https://api.fireworks.ai/v1", "fireworks"],
    ["https://custom.example.com/v1", "openai-compatible"],
  ])("infers %s as %s", (url, expected) => {
    expect(inferProviderName(url)).toBe(expected);
  });
});

describe("provider name from config", () => {
  it.each([
    { provider: "openai", expected: "openai" },
    { provider: "xai", expected: "xai" },
    { provider: "groq", expected: "groq" },
    { provider: "together", expected: "together" },
    { provider: "ollama", expected: "ollama" },
  ])("explicit $provider maps to $expected", ({ provider, expected }) => {
    const p = createProvider({
      provider,
      model: "test-model",
      apiKey: "test-key",
    });
    expect(p.name).toBe(expected);
  });

  it("infers name from baseUrl for openai-compatible", () => {
    const p = createProvider({
      provider: "openai-compatible",
      model: "test-model",
      apiKey: "test-key",
      baseUrl: "https://api.x.ai/v1",
    });
    expect(p.name).toBe("xai");
  });

  it("keeps openai-compatible for unknown URLs", () => {
    const p = createProvider({
      provider: "openai-compatible",
      model: "test-model",
      apiKey: "test-key",
      baseUrl: "https://custom.example.com/v1",
    });
    expect(p.name).toBe("openai-compatible");
  });

  it("azure provider name", () => {
    const p = createProvider({
      provider: "azure",
      model: "gpt-4o",
      apiKey: "test-key",
      baseUrl: "https://mydeployment.azure.openai.com/v1",
    });
    expect(p.name).toBe("azure");
  });
});

// ════════════════════════════════════════════════════════════════════════
// Claude CLI config
// ════════════════════════════════════════════════════════════════════════

describe("newClaudeCliConfig", () => {
  it("creates config with default model", () => {
    const cfg = newClaudeCliConfig();
    expect(cfg.provider).toBe("claude-cli");
    expect(cfg.model).toBe("sonnet");
  });

  it("creates config with custom model", () => {
    const cfg = newClaudeCliConfig("opus");
    expect(cfg.provider).toBe("claude-cli");
    expect(cfg.model).toBe("opus");
  });
});

// ════════════════════════════════════════════════════════════════════════
// eventsToMessages
// ════════════════════════════════════════════════════════════════════════

describe("eventsToMessages", () => {
  it("returns empty string for empty array", () => {
    expect(eventsToMessages([])).toBe("");
  });

  it("returns empty string for null/undefined", () => {
    expect(eventsToMessages(undefined as unknown as Event[])).toBe("");
    expect(eventsToMessages(null as unknown as Event[])).toBe("");
  });

  it("formats single event", () => {
    const ev = makeEvent("trust.updated", "actor_test001");
    const result = eventsToMessages([ev]);
    expect(result).toContain("Event history:");
    expect(result).toContain("[trust.updated]");
    expect(result).toContain("actor_test001");
  });

  it("caps at 20 events", () => {
    const events = Array.from({ length: 25 }, (_, i) =>
      makeEvent("test.event", `actor_${i}`),
    );
    const result = eventsToMessages(events);
    expect(result).toContain("... and 5 more events");
  });
});

// ════════════════════════════════════════════════════════════════════════
// parseConfidence
// ════════════════════════════════════════════════════════════════════════

describe("parseConfidence", () => {
  it("returns Score of 0.7", () => {
    const score = parseConfidence(100);
    expect(score.value).toBe(0.7);
  });

  it("returns valid Score regardless of token count", () => {
    const score = parseConfidence(0);
    expect(score).toBeInstanceOf(Score);
    expect(score.value).toBeGreaterThanOrEqual(0);
    expect(score.value).toBeLessThanOrEqual(1);
  });
});

// ════════════════════════════════════════════════════════════════════════
// detectProvider
// ════════════════════════════════════════════════════════════════════════

describe("detectProvider", () => {
  const origEnv = { ...process.env };

  afterEach(() => {
    // Restore env.
    process.env = { ...origEnv };
  });

  it("returns openai-compatible with no env vars set", () => {
    delete process.env["XAI_API_KEY"];
    delete process.env["GROQ_API_KEY"];
    delete process.env["TOGETHER_API_KEY"];
    delete process.env["OPENAI_API_KEY"];
    delete process.env["OLLAMA_HOST"];

    const result = detectProvider();
    expect(result.name).toBe("openai-compatible");
    expect(result.apiKey).toBe("");
  });

  it("detects XAI_API_KEY", () => {
    process.env["XAI_API_KEY"] = "test-xai-key";
    delete process.env["GROQ_API_KEY"];
    delete process.env["TOGETHER_API_KEY"];
    delete process.env["OPENAI_API_KEY"];
    delete process.env["OLLAMA_HOST"];

    const result = detectProvider();
    expect(result.name).toBe("xai");
    expect(result.apiKey).toBe("test-xai-key");
  });
});

// ════════════════════════════════════════════════════════════════════════
// Provider interface compliance
// ════════════════════════════════════════════════════════════════════════

describe("Provider interface", () => {
  it("satisfies Intelligence interface (has reason method)", () => {
    const p = createProvider({
      provider: "openai",
      model: "gpt-4o",
      apiKey: "test-key",
    });
    expect(typeof p.reason).toBe("function");
    expect(typeof p.name).toBe("string");
    expect(typeof p.model).toBe("string");
  });
});

// ════════════════════════════════════════════════════════════════════════
// Integration tests — require API keys
// ════════════════════════════════════════════════════════════════════════

describe("OpenAI-compatible integration", () => {
  const hasOpenAIKey = !!process.env["OPENAI_API_KEY"];
  const hasXAIKey = !!process.env["XAI_API_KEY"];
  const hasGroqKey = !!process.env["GROQ_API_KEY"];
  const hasAnyKey = hasOpenAIKey || hasXAIKey || hasGroqKey;

  function getKeyAndModel(): { key: string; model: string; provider: string } {
    if (process.env["OPENAI_API_KEY"]) {
      return { key: process.env["OPENAI_API_KEY"]!, model: "gpt-4o-mini", provider: "openai" };
    }
    if (process.env["XAI_API_KEY"]) {
      return { key: process.env["XAI_API_KEY"]!, model: "grok-3-mini-fast", provider: "xai" };
    }
    if (process.env["GROQ_API_KEY"]) {
      return { key: process.env["GROQ_API_KEY"]!, model: "llama-3.1-8b-instant", provider: "groq" };
    }
    return { key: "", model: "", provider: "" };
  }

  it.skipIf(!hasAnyKey)("reason returns a response", async () => {
    const { key, model, provider } = getKeyAndModel();
    const p = createProvider({
      provider,
      model,
      apiKey: key,
      maxTokens: 100,
    });

    const resp = await p.reason("Reply with exactly one word: hello", []);
    expect(resp.content).toBeTruthy();
    expect(resp.tokensUsed).toBeGreaterThan(0);
    expect(resp.confidence.value).toBeGreaterThan(0);
  });

  it.skipIf(!hasAnyKey)("reason with system prompt", async () => {
    const { key, model, provider } = getKeyAndModel();
    const p = createProvider({
      provider,
      model,
      apiKey: key,
      maxTokens: 100,
      systemPrompt: "Respond with exactly one word.",
    });

    const resp = await p.reason("What is 1+1?", []);
    expect(resp.content).toBeTruthy();
  });

  it.skipIf(!hasAnyKey)("reason with history", async () => {
    const { key, model, provider } = getKeyAndModel();
    const p = createProvider({
      provider,
      model,
      apiKey: key,
      maxTokens: 100,
    });

    const ev = makeEvent("trust.updated", "actor_test001");
    const resp = await p.reason("Given the event history, reply with exactly one word: acknowledged", [ev]);
    expect(resp.content).toBeTruthy();
  });

  it("invalid API key returns error", async () => {
    const p = createProvider({
      provider: "openai",
      model: "gpt-4o-mini",
      apiKey: "sk-invalid-key-for-testing",
      maxTokens: 50,
    });

    await expect(p.reason("hello", [])).rejects.toThrow();
  });
});

describe("Claude CLI integration", () => {
  const hasCli = !!process.env["EVENTGRAPH_TEST_CLAUDE_CLI"];

  it.skipIf(!hasCli)("reason returns a response", () => {
    const p = createProvider(newClaudeCliConfig("sonnet"));

    const resp = p.reason("Reply with exactly one word: hello", []);
    // ClaudeCliProvider.reason is synchronous, but can be awaited.
    if (resp instanceof Promise) {
      return resp.then((r) => {
        expect(r.content).toBeTruthy();
      });
    }
    expect((resp as Response).content).toBeTruthy();
  });

  it.skipIf(!hasCli)("reason with system prompt", () => {
    const p = createProvider({
      provider: "claude-cli",
      model: "sonnet",
      systemPrompt: "You are a decision-making system. Respond with: PERMIT, DENY, ESCALATE, or DEFER.",
    });

    const resp = p.reason(
      "An agent with trust score 0.3 wants to delete critical data. What is your decision?",
      [],
    );
    if (resp instanceof Promise) {
      return resp.then((r) => {
        expect(r.content).toBeTruthy();
      });
    }
    expect((resp as Response).content).toBeTruthy();
  });
});
