using System.Diagnostics;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace EventGraph;

// ── Provider Interface ──────────────────────────────────────────────────

/// <summary>Extends IIntelligence with metadata about the backing LLM.</summary>
public interface IProvider : IIntelligence
{
    /// <summary>The provider identifier (e.g., "openai", "anthropic", "ollama").</summary>
    string Name { get; }

    /// <summary>The model identifier (e.g., "gpt-4o", "claude-sonnet-4-6").</summary>
    string Model { get; }
}

// ── Config ──────────────────────────────────────────────────────────────

/// <summary>Configuration for creating an IProvider via ProviderFactory.</summary>
public sealed class ProviderConfig
{
    /// <summary>Provider name: "openai", "openai-compatible", "xai", "groq", "together", "ollama", "azure", "claude-cli".</summary>
    public string Provider { get; }

    /// <summary>Model identifier (provider-specific).</summary>
    public string Model { get; }

    /// <summary>API key for authentication. Falls back to environment variables if empty.</summary>
    public string ApiKey { get; init; } = "";

    /// <summary>Overrides the default API endpoint. Useful for proxies, Ollama, Azure, etc.</summary>
    public string BaseUrl { get; init; } = "";

    /// <summary>Caps the response length. Defaults to 1024 if zero.</summary>
    public int MaxTokens { get; init; }

    /// <summary>Controls randomness. Zero means provider default.</summary>
    public double Temperature { get; init; }

    /// <summary>System prompt prepended to every Reason call.</summary>
    public string SystemPrompt { get; init; } = "";

    /// <summary>Create a new ProviderConfig.</summary>
    /// <param name="provider">Provider name (e.g., "openai", "claude-cli").</param>
    /// <param name="model">Model identifier.</param>
    public ProviderConfig(string provider, string model)
    {
        if (string.IsNullOrWhiteSpace(provider))
            throw new EmptyRequiredException("ProviderConfig.Provider");
        if (string.IsNullOrWhiteSpace(model))
            throw new EmptyRequiredException("ProviderConfig.Model");
        Provider = provider;
        Model = model;
    }
}

// ── Factory ─────────────────────────────────────────────────────────────

/// <summary>Creates IProvider instances from ProviderConfig.</summary>
public static class ProviderFactory
{
    /// <summary>Create a provider from the given configuration.</summary>
    /// <exception cref="EventGraphException">If the provider name is unknown or config is invalid.</exception>
    public static IProvider Create(ProviderConfig config)
    {
        var maxTokens = config.MaxTokens > 0 ? config.MaxTokens : 1024;
        var effectiveConfig = new ProviderConfig(config.Provider, config.Model)
        {
            ApiKey = config.ApiKey,
            BaseUrl = config.BaseUrl,
            MaxTokens = maxTokens,
            Temperature = config.Temperature,
            SystemPrompt = config.SystemPrompt,
        };

        return effectiveConfig.Provider.ToLowerInvariant() switch
        {
            "claude-cli" => new ClaudeCliProvider(effectiveConfig),
            "openai-compatible" or "openai" or "xai" or "groq" or "together" or "ollama" or "azure"
                => new OpenAICompatibleProvider(effectiveConfig),
            _ => throw new EventGraphException(
                $"Unknown provider: \"{config.Provider}\" (supported: openai, openai-compatible, xai, groq, together, ollama, azure, claude-cli)"),
        };
    }
}

// ── OpenAI-Compatible Provider ──────────────────────────────────────────

/// <summary>
/// Implements IProvider using the OpenAI Chat Completions API.
/// Compatible with: OpenAI, xAI/Grok, Ollama, Together, Azure OpenAI,
/// Fireworks, Groq, and any OpenAI-compatible endpoint.
/// </summary>
public sealed class OpenAICompatibleProvider : IProvider
{
    private static readonly HttpClient SharedClient = new();

    private readonly HttpClient _client;
    private readonly string _baseUrl;
    private readonly string _apiKey;
    private readonly string _model;
    private readonly int _maxTokens;
    private readonly double _temperature;
    private readonly string _systemPrompt;
    private readonly string _providerName;

    /// <summary>Well-known OpenAI-compatible base URLs.</summary>
    internal const string OpenAIBaseUrl = "https://api.openai.com/v1";
    internal const string XaiBaseUrl = "https://api.x.ai/v1";
    internal const string GroqBaseUrl = "https://api.groq.com/openai/v1";
    internal const string TogetherBaseUrl = "https://api.together.xyz/v1";
    internal const string OllamaBaseUrl = "http://localhost:11434/v1";

    /// <inheritdoc />
    public string Name => _providerName;

    /// <inheritdoc />
    public string Model => _model;

    internal OpenAICompatibleProvider(ProviderConfig config) : this(config, null) { }

    internal OpenAICompatibleProvider(ProviderConfig config, HttpClient? httpClient)
    {
        _client = httpClient ?? SharedClient;
        _model = config.Model;
        _maxTokens = config.MaxTokens;
        _temperature = config.Temperature;
        _systemPrompt = config.SystemPrompt;

        var apiKey = config.ApiKey;
        var baseUrl = config.BaseUrl;
        var providerName = config.Provider.ToLowerInvariant();

        // Map shorthand provider names to default base URLs and env vars.
        switch (providerName)
        {
            case "openai":
                if (string.IsNullOrEmpty(baseUrl)) baseUrl = OpenAIBaseUrl;
                if (string.IsNullOrEmpty(apiKey)) apiKey = Environment.GetEnvironmentVariable("OPENAI_API_KEY") ?? "";
                break;
            case "xai":
                if (string.IsNullOrEmpty(baseUrl)) baseUrl = XaiBaseUrl;
                if (string.IsNullOrEmpty(apiKey)) apiKey = Environment.GetEnvironmentVariable("XAI_API_KEY") ?? "";
                break;
            case "groq":
                if (string.IsNullOrEmpty(baseUrl)) baseUrl = GroqBaseUrl;
                if (string.IsNullOrEmpty(apiKey)) apiKey = Environment.GetEnvironmentVariable("GROQ_API_KEY") ?? "";
                break;
            case "together":
                if (string.IsNullOrEmpty(baseUrl)) baseUrl = TogetherBaseUrl;
                if (string.IsNullOrEmpty(apiKey)) apiKey = Environment.GetEnvironmentVariable("TOGETHER_API_KEY") ?? "";
                break;
            case "ollama":
                if (string.IsNullOrEmpty(baseUrl))
                {
                    var host = Environment.GetEnvironmentVariable("OLLAMA_HOST") ?? "http://localhost:11434";
                    baseUrl = host.TrimEnd('/') + "/v1";
                }
                break;
            case "openai-compatible":
                if (string.IsNullOrEmpty(apiKey) && string.IsNullOrEmpty(baseUrl))
                {
                    (apiKey, baseUrl, providerName) = IntelligenceHelpers.DetectProvider();
                }
                break;
        }

        if (string.IsNullOrEmpty(baseUrl))
            baseUrl = OpenAIBaseUrl;

        // Infer provider name from base URL if still generic.
        if (providerName == "openai-compatible")
            providerName = IntelligenceHelpers.InferProviderName(baseUrl);

        _apiKey = apiKey;
        _baseUrl = baseUrl.TrimEnd('/');
        _providerName = providerName;
    }

    /// <inheritdoc />
    public Response Reason(string prompt, IReadOnlyList<Event> history)
    {
        var messages = new List<OpenAIMessage>();

        if (!string.IsNullOrEmpty(_systemPrompt))
        {
            messages.Add(new OpenAIMessage { Role = "system", Content = _systemPrompt });
        }

        var historyText = IntelligenceHelpers.EventsToMessages(history);
        if (!string.IsNullOrEmpty(historyText))
        {
            messages.Add(new OpenAIMessage { Role = "user", Content = historyText });
            messages.Add(new OpenAIMessage { Role = "assistant", Content = "I understand the event history. What would you like me to reason about?" });
        }

        messages.Add(new OpenAIMessage { Role = "user", Content = prompt });

        var reqBody = new OpenAIRequest
        {
            Model = _model,
            Messages = messages,
        };
        if (_maxTokens > 0) reqBody.MaxTokens = _maxTokens;
        if (_temperature > 0) reqBody.Temperature = _temperature;

        var json = JsonSerializer.Serialize(reqBody, OpenAIJsonContext.Default.OpenAIRequest);
        var content = new StringContent(json, Encoding.UTF8, "application/json");

        var request = new HttpRequestMessage(HttpMethod.Post, _baseUrl + "/chat/completions")
        {
            Content = content,
        };
        request.Headers.Add("Accept", "application/json");
        if (!string.IsNullOrEmpty(_apiKey))
        {
            request.Headers.Add("Authorization", "Bearer " + _apiKey);
        }

        HttpResponseMessage response;
        try
        {
            response = _client.Send(request);
        }
        catch (Exception ex)
        {
            throw new EventGraphException($"OpenAI API request failed: {ex.Message}", ex);
        }

        var responseBody = response.Content.ReadAsStringAsync().GetAwaiter().GetResult();

        if (!response.IsSuccessStatusCode)
        {
            // Try to extract error message from JSON body.
            try
            {
                var errResp = JsonSerializer.Deserialize(responseBody, OpenAIJsonContext.Default.OpenAIResponse);
                if (errResp?.Error is not null)
                {
                    throw new EventGraphException($"OpenAI API error (HTTP {(int)response.StatusCode}): {errResp.Error.Message}");
                }
            }
            catch (JsonException) { }
            throw new EventGraphException($"OpenAI API error (HTTP {(int)response.StatusCode}): {responseBody}");
        }

        var result = JsonSerializer.Deserialize(responseBody, OpenAIJsonContext.Default.OpenAIResponse)
            ?? throw new EventGraphException("Failed to parse OpenAI response: null result");

        if (result.Choices is null || result.Choices.Count == 0)
        {
            throw new EventGraphException("OpenAI API returned no choices");
        }

        var responseContent = result.Choices[0].Message?.Content ?? "";
        var tokensUsed = result.Usage?.TotalTokens ?? 0;
        var confidence = IntelligenceHelpers.ParseConfidence(tokensUsed);

        return new Response(responseContent, confidence, tokensUsed);
    }
}

// ── Claude CLI Provider ─────────────────────────────────────────────────

/// <summary>
/// Implements IProvider by shelling out to the claude CLI.
/// Uses whatever authentication Claude Code already has (OAuth, API key, etc.)
/// without requiring separate credentials.
/// </summary>
public sealed class ClaudeCliProvider : IProvider
{
    private readonly string _model;
    private readonly double _maxBudget;
    private readonly string _systemPrompt;
    private readonly string _claudePath;

    /// <inheritdoc />
    public string Name => "claude-cli";

    /// <inheritdoc />
    public string Model => _model;

    internal ClaudeCliProvider(ProviderConfig config)
    {
        _model = string.IsNullOrEmpty(config.Model) ? "sonnet" : config.Model;
        _claudePath = string.IsNullOrEmpty(config.BaseUrl) ? "claude" : config.BaseUrl;
        _systemPrompt = config.SystemPrompt;

        // Repurpose Temperature as max budget hint (dollars). Default $1 per call.
        _maxBudget = config.Temperature > 0 ? config.Temperature : 1.0;
    }

    /// <inheritdoc />
    public Response Reason(string prompt, IReadOnlyList<Event> history)
    {
        var fullPrompt = new StringBuilder();
        var historyText = IntelligenceHelpers.EventsToMessages(history);
        if (!string.IsNullOrEmpty(historyText))
        {
            fullPrompt.Append(historyText);
            fullPrompt.Append("\n---\n\n");
        }
        fullPrompt.Append(prompt);

        var args = new StringBuilder();
        args.Append("-p --output-format json");
        args.Append($" --model {_model}");
        args.Append($" --max-budget-usd {_maxBudget:F2}");
        args.Append(" --no-session-persistence");
        if (!string.IsNullOrEmpty(_systemPrompt))
        {
            args.Append($" --system-prompt \"{_systemPrompt.Replace("\"", "\\\"")}\"");
        }

        var psi = new ProcessStartInfo
        {
            FileName = _claudePath,
            Arguments = args.ToString(),
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
            CreateNoWindow = true,
        };

        // Remove CLAUDECODE env var to allow nested invocation.
        RemoveEnvironmentVariable(psi, "CLAUDECODE");

        Process proc;
        try
        {
            proc = Process.Start(psi) ?? throw new EventGraphException("Failed to start claude CLI process");
        }
        catch (Exception ex)
        {
            throw new EventGraphException($"Claude CLI not found or failed to start: {ex.Message}", ex);
        }

        proc.StandardInput.Write(fullPrompt.ToString());
        proc.StandardInput.Close();

        var stdout = proc.StandardOutput.ReadToEnd();
        var stderr = proc.StandardError.ReadToEnd();
        proc.WaitForExit();

        if (proc.ExitCode != 0)
        {
            // Check if we got JSON output despite non-zero exit.
            if (!string.IsNullOrEmpty(stdout))
            {
                try
                {
                    var partial = JsonSerializer.Deserialize(stdout, ClaudeCliJsonContext.Default.ClaudeCliResult);
                    if (partial is not null && !string.IsNullOrEmpty(partial.Result))
                    {
                        return ResultToResponse(partial);
                    }
                }
                catch (JsonException) { }
            }
            throw new EventGraphException($"Claude CLI error (exit {proc.ExitCode}): {stderr}");
        }

        ClaudeCliResult result;
        try
        {
            result = JsonSerializer.Deserialize(stdout, ClaudeCliJsonContext.Default.ClaudeCliResult)
                ?? throw new EventGraphException($"Failed to parse claude CLI JSON output: null result");
        }
        catch (JsonException ex)
        {
            throw new EventGraphException($"Failed to parse claude CLI JSON output: {ex.Message}\nraw: {stdout}", ex);
        }

        if (result.IsError)
        {
            throw new EventGraphException($"Claude CLI returned error: {result.Result} (subtype: {result.Subtype})");
        }

        return ResultToResponse(result);
    }

    private static Response ResultToResponse(ClaudeCliResult result)
    {
        var tokensUsed = result.Usage.InputTokens + result.Usage.OutputTokens;
        var confidence = IntelligenceHelpers.ParseConfidence(tokensUsed);
        return new Response(result.Result ?? "", confidence, tokensUsed);
    }

    private static void RemoveEnvironmentVariable(ProcessStartInfo psi, string key)
    {
        // Accessing EnvironmentVariables copies the current process env.
        if (psi.EnvironmentVariables.ContainsKey(key))
        {
            psi.EnvironmentVariables.Remove(key);
        }
    }
}

// ── Helper Methods ──────────────────────────────────────────────────────

/// <summary>Shared helper methods for intelligence providers.</summary>
internal static class IntelligenceHelpers
{
    /// <summary>Convert event history into a simple text context for an LLM.</summary>
    internal static string EventsToMessages(IReadOnlyList<Event> events)
    {
        if (events.Count == 0) return "";

        var sb = new StringBuilder("Event history:\n");
        var limit = Math.Min(events.Count, 20);
        for (var i = 0; i < limit; i++)
        {
            var ev = events[i];
            sb.Append($"- [{ev.Type.Value}] {ev.Id.Value} by {ev.Source.Value}\n");
        }
        if (events.Count > 20)
        {
            sb.Append($"... and {events.Count - 20} more events\n");
        }
        return sb.ToString();
    }

    /// <summary>
    /// Derives a confidence score from token usage. This is a heuristic —
    /// real confidence requires model introspection.
    /// </summary>
    internal static Score ParseConfidence(int tokensUsed)
    {
        // Default to 0.7 — we can't truly measure confidence from outside the model.
        return new Score(0.7);
    }

    /// <summary>
    /// Checks environment variables to auto-detect which OpenAI-compatible provider to use.
    /// </summary>
    internal static (string apiKey, string baseUrl, string name) DetectProvider()
    {
        var xaiKey = Environment.GetEnvironmentVariable("XAI_API_KEY");
        if (!string.IsNullOrEmpty(xaiKey))
            return (xaiKey, OpenAICompatibleProvider.XaiBaseUrl, "xai");

        var groqKey = Environment.GetEnvironmentVariable("GROQ_API_KEY");
        if (!string.IsNullOrEmpty(groqKey))
            return (groqKey, OpenAICompatibleProvider.GroqBaseUrl, "groq");

        var togetherKey = Environment.GetEnvironmentVariable("TOGETHER_API_KEY");
        if (!string.IsNullOrEmpty(togetherKey))
            return (togetherKey, OpenAICompatibleProvider.TogetherBaseUrl, "together");

        var openaiKey = Environment.GetEnvironmentVariable("OPENAI_API_KEY");
        if (!string.IsNullOrEmpty(openaiKey))
            return (openaiKey, OpenAICompatibleProvider.OpenAIBaseUrl, "openai");

        var ollamaHost = Environment.GetEnvironmentVariable("OLLAMA_HOST");
        if (!string.IsNullOrEmpty(ollamaHost))
            return ("", ollamaHost.TrimEnd('/') + "/v1", "ollama");

        return ("", "", "openai-compatible");
    }

    /// <summary>Derives a friendly provider name from the base URL.</summary>
    internal static string InferProviderName(string baseUrl)
    {
        var lower = baseUrl.ToLowerInvariant();
        if (lower.Contains("azure")) return "azure";
        if (lower.Contains("openai.com")) return "openai";
        if (lower.Contains("x.ai")) return "xai";
        if (lower.Contains("groq.com")) return "groq";
        if (lower.Contains("together.xyz")) return "together";
        if (lower.Contains("localhost") || lower.Contains("127.0.0.1")) return "ollama";
        if (lower.Contains("fireworks")) return "fireworks";
        return "openai-compatible";
    }
}

// ── JSON DTOs (OpenAI) ──────────────────────────────────────────────────

internal sealed class OpenAIMessage
{
    [JsonPropertyName("role")]
    public string Role { get; set; } = "";

    [JsonPropertyName("content")]
    public string Content { get; set; } = "";
}

internal sealed class OpenAIRequest
{
    [JsonPropertyName("model")]
    public string Model { get; set; } = "";

    [JsonPropertyName("messages")]
    public List<OpenAIMessage> Messages { get; set; } = new();

    [JsonPropertyName("max_tokens")]
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingDefault)]
    public int MaxTokens { get; set; }

    [JsonPropertyName("temperature")]
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingDefault)]
    public double Temperature { get; set; }
}

internal sealed class OpenAIChoice
{
    [JsonPropertyName("message")]
    public OpenAIMessage? Message { get; set; }

    [JsonPropertyName("finish_reason")]
    public string? FinishReason { get; set; }
}

internal sealed class OpenAIUsage
{
    [JsonPropertyName("prompt_tokens")]
    public int PromptTokens { get; set; }

    [JsonPropertyName("completion_tokens")]
    public int CompletionTokens { get; set; }

    [JsonPropertyName("total_tokens")]
    public int TotalTokens { get; set; }
}

internal sealed class OpenAIErrorDetail
{
    [JsonPropertyName("message")]
    public string? Message { get; set; }

    [JsonPropertyName("type")]
    public string? Type { get; set; }

    [JsonPropertyName("code")]
    public string? Code { get; set; }
}

internal sealed class OpenAIResponse
{
    [JsonPropertyName("id")]
    public string? Id { get; set; }

    [JsonPropertyName("choices")]
    public List<OpenAIChoice>? Choices { get; set; }

    [JsonPropertyName("usage")]
    public OpenAIUsage? Usage { get; set; }

    [JsonPropertyName("error")]
    public OpenAIErrorDetail? Error { get; set; }
}

// ── JSON DTOs (Claude CLI) ──────────────────────────────────────────────

internal sealed class ClaudeCliUsage
{
    [JsonPropertyName("input_tokens")]
    public int InputTokens { get; set; }

    [JsonPropertyName("output_tokens")]
    public int OutputTokens { get; set; }
}

internal sealed class ClaudeCliResult
{
    [JsonPropertyName("type")]
    public string? Type { get; set; }

    [JsonPropertyName("subtype")]
    public string? Subtype { get; set; }

    [JsonPropertyName("is_error")]
    public bool IsError { get; set; }

    [JsonPropertyName("result")]
    public string? Result { get; set; }

    [JsonPropertyName("usage")]
    public ClaudeCliUsage Usage { get; set; } = new();

    [JsonPropertyName("total_cost_usd")]
    public double TotalCostUsd { get; set; }

    [JsonPropertyName("stop_reason")]
    public string? StopReason { get; set; }
}

// ── Source-generated JSON contexts ──────────────────────────────────────

[JsonSerializable(typeof(OpenAIRequest))]
[JsonSerializable(typeof(OpenAIResponse))]
internal partial class OpenAIJsonContext : JsonSerializerContext { }

[JsonSerializable(typeof(ClaudeCliResult))]
internal partial class ClaudeCliJsonContext : JsonSerializerContext { }

// ── Convenience ─────────────────────────────────────────────────────────

/// <summary>Creates a ProviderConfig for the Claude CLI provider.</summary>
public static class ClaudeCliConfig
{
    /// <summary>Create a config for the Claude CLI provider with the given model.</summary>
    /// <param name="model">Model name (defaults to "sonnet" if empty).</param>
    public static ProviderConfig Create(string model = "sonnet")
    {
        if (string.IsNullOrEmpty(model)) model = "sonnet";
        return new ProviderConfig("claude-cli", model);
    }
}
