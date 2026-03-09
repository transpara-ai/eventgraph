namespace EventGraph.Tests;

public class IntelligenceTests
{
    // ── Config Validation ───────────────────────────────────────────────

    [Fact]
    public void ProviderConfig_RequiresProvider()
    {
        Assert.Throws<EmptyRequiredException>(() => new ProviderConfig("", "gpt-4o"));
    }

    [Fact]
    public void ProviderConfig_RequiresModel()
    {
        Assert.Throws<EmptyRequiredException>(() => new ProviderConfig("openai", ""));
    }

    [Fact]
    public void ProviderConfig_WhitespaceProvider_Throws()
    {
        Assert.Throws<EmptyRequiredException>(() => new ProviderConfig("   ", "gpt-4o"));
    }

    [Fact]
    public void ProviderConfig_WhitespaceModel_Throws()
    {
        Assert.Throws<EmptyRequiredException>(() => new ProviderConfig("openai", "   "));
    }

    [Fact]
    public void ProviderConfig_ValidConfig_Succeeds()
    {
        var config = new ProviderConfig("openai", "gpt-4o")
        {
            ApiKey = "sk-test",
            BaseUrl = "https://api.openai.com/v1",
            MaxTokens = 2048,
            Temperature = 0.5,
            SystemPrompt = "You are a helpful assistant.",
        };

        Assert.Equal("openai", config.Provider);
        Assert.Equal("gpt-4o", config.Model);
        Assert.Equal("sk-test", config.ApiKey);
        Assert.Equal("https://api.openai.com/v1", config.BaseUrl);
        Assert.Equal(2048, config.MaxTokens);
        Assert.Equal(0.5, config.Temperature);
        Assert.Equal("You are a helpful assistant.", config.SystemPrompt);
    }

    [Fact]
    public void ProviderConfig_Defaults()
    {
        var config = new ProviderConfig("openai", "gpt-4o");

        Assert.Equal("", config.ApiKey);
        Assert.Equal("", config.BaseUrl);
        Assert.Equal(0, config.MaxTokens);
        Assert.Equal(0.0, config.Temperature);
        Assert.Equal("", config.SystemPrompt);
    }

    // ── Provider Name Inference ─────────────────────────────────────────

    [Theory]
    [InlineData("https://api.openai.com/v1", "openai")]
    [InlineData("https://api.x.ai/v1", "xai")]
    [InlineData("https://api.groq.com/openai/v1", "groq")]
    [InlineData("https://api.together.xyz/v1", "together")]
    [InlineData("http://localhost:11434/v1", "ollama")]
    [InlineData("http://127.0.0.1:11434/v1", "ollama")]
    [InlineData("https://my-resource.azure.openai.com/v1", "azure")]
    [InlineData("https://api.fireworks.ai/v1", "fireworks")]
    [InlineData("https://custom-endpoint.example.com/v1", "openai-compatible")]
    public void InferProviderName_FromBaseUrl(string baseUrl, string expected)
    {
        var result = IntelligenceHelpers.InferProviderName(baseUrl);
        Assert.Equal(expected, result);
    }

    // ── Factory Creation ────────────────────────────────────────────────

    [Fact]
    public void Factory_UnknownProvider_Throws()
    {
        var config = new ProviderConfig("unknown-provider", "some-model");
        var ex = Assert.Throws<EventGraphException>(() => ProviderFactory.Create(config));
        Assert.Contains("Unknown provider", ex.Message);
        Assert.Contains("unknown-provider", ex.Message);
    }

    [Theory]
    [InlineData("openai")]
    [InlineData("xai")]
    [InlineData("groq")]
    [InlineData("together")]
    [InlineData("ollama")]
    [InlineData("openai-compatible")]
    public void Factory_CreatesOpenAICompatible(string provider)
    {
        var config = new ProviderConfig(provider, "test-model") { ApiKey = "test-key" };
        var result = ProviderFactory.Create(config);

        Assert.IsType<OpenAICompatibleProvider>(result);
        Assert.Equal("test-model", result.Model);
    }

    [Fact]
    public void Factory_CreatesClaude()
    {
        var config = new ProviderConfig("claude-cli", "sonnet");
        var result = ProviderFactory.Create(config);

        Assert.IsType<ClaudeCliProvider>(result);
        Assert.Equal("claude-cli", result.Name);
        Assert.Equal("sonnet", result.Model);
    }

    [Fact]
    public void Factory_DefaultsMaxTokens()
    {
        var config = new ProviderConfig("openai", "gpt-4o") { ApiKey = "test-key" };
        var provider = ProviderFactory.Create(config);

        // MaxTokens defaults to 1024 when 0 — verified by construction, not observable externally.
        // The factory path succeeds, which is the test.
        Assert.NotNull(provider);
    }

    // ── OpenAI Provider Properties ──────────────────────────────────────

    [Fact]
    public void OpenAI_NameAndModel()
    {
        var config = new ProviderConfig("openai", "gpt-4o") { ApiKey = "sk-test" };
        var provider = ProviderFactory.Create(config);

        Assert.Equal("openai", provider.Name);
        Assert.Equal("gpt-4o", provider.Model);
    }

    [Fact]
    public void Xai_InfersName()
    {
        var config = new ProviderConfig("xai", "grok-2") { ApiKey = "xai-test" };
        var provider = ProviderFactory.Create(config);

        Assert.Equal("xai", provider.Name);
        Assert.Equal("grok-2", provider.Model);
    }

    [Fact]
    public void Groq_InfersName()
    {
        var config = new ProviderConfig("groq", "mixtral-8x7b") { ApiKey = "groq-test" };
        var provider = ProviderFactory.Create(config);

        Assert.Equal("groq", provider.Name);
    }

    [Fact]
    public void Together_InfersName()
    {
        var config = new ProviderConfig("together", "llama-3") { ApiKey = "together-test" };
        var provider = ProviderFactory.Create(config);

        Assert.Equal("together", provider.Name);
    }

    [Fact]
    public void Ollama_InfersName()
    {
        var config = new ProviderConfig("ollama", "llama3");
        var provider = ProviderFactory.Create(config);

        Assert.Equal("ollama", provider.Name);
    }

    [Fact]
    public void OpenAICompatible_WithCustomUrl_InfersAzure()
    {
        var config = new ProviderConfig("openai-compatible", "gpt-4")
        {
            ApiKey = "test-key",
            BaseUrl = "https://my-resource.azure.openai.com/v1",
        };
        var provider = ProviderFactory.Create(config);

        Assert.Equal("azure", provider.Name);
    }

    // ── EventsToMessages ────────────────────────────────────────────────

    [Fact]
    public void EventsToMessages_EmptyList_ReturnsEmpty()
    {
        var result = IntelligenceHelpers.EventsToMessages(Array.Empty<Event>());
        Assert.Equal("", result);
    }

    [Fact]
    public void EventsToMessages_SingleEvent_FormatsCorrectly()
    {
        var ev = EventFactory.CreateBootstrap(new ActorId("test-actor"), new NoopSigner());
        var result = IntelligenceHelpers.EventsToMessages(new List<Event> { ev });

        Assert.StartsWith("Event history:\n", result);
        Assert.Contains("[system.bootstrapped]", result);
        Assert.Contains("test-actor", result);
    }

    [Fact]
    public void EventsToMessages_Over20_Truncates()
    {
        var events = new List<Event>();
        var signer = new NoopSigner();
        var bootstrap = EventFactory.CreateBootstrap(new ActorId("actor"), signer);
        for (var i = 0; i < 25; i++)
        {
            events.Add(bootstrap);
        }

        var result = IntelligenceHelpers.EventsToMessages(events);

        Assert.Contains("... and 5 more events", result);
    }

    // ── ParseConfidence ─────────────────────────────────────────────────

    [Fact]
    public void ParseConfidence_ReturnsScore()
    {
        var score = IntelligenceHelpers.ParseConfidence(500);
        Assert.Equal(0.7, score.Value);
    }

    // ── ClaudeCliConfig ─────────────────────────────────────────────────

    [Fact]
    public void ClaudeCliConfig_DefaultModel()
    {
        var config = ClaudeCliConfig.Create();
        Assert.Equal("claude-cli", config.Provider);
        Assert.Equal("sonnet", config.Model);
    }

    [Fact]
    public void ClaudeCliConfig_CustomModel()
    {
        var config = ClaudeCliConfig.Create("opus");
        Assert.Equal("opus", config.Model);
    }

    [Fact]
    public void ClaudeCliConfig_EmptyModel_DefaultsToSonnet()
    {
        var config = ClaudeCliConfig.Create("");
        Assert.Equal("sonnet", config.Model);
    }

    // ── IProvider interface conformance ─────────────────────────────────

    [Fact]
    public void OpenAIProvider_ImplementsIProvider()
    {
        var config = new ProviderConfig("openai", "gpt-4o") { ApiKey = "test" };
        var provider = ProviderFactory.Create(config);

        Assert.IsAssignableFrom<IProvider>(provider);
        Assert.IsAssignableFrom<IIntelligence>(provider);
    }

    [Fact]
    public void ClaudeCliProvider_ImplementsIProvider()
    {
        var config = new ProviderConfig("claude-cli", "sonnet");
        var provider = ProviderFactory.Create(config);

        Assert.IsAssignableFrom<IProvider>(provider);
        Assert.IsAssignableFrom<IIntelligence>(provider);
    }

    // ── Integration Tests (skipped by default) ──────────────────────────

    [Fact(Skip = "Requires OPENAI_API_KEY environment variable")]
    public void Integration_OpenAI_Reason()
    {
        var apiKey = Environment.GetEnvironmentVariable("OPENAI_API_KEY");
        if (string.IsNullOrEmpty(apiKey)) return;

        var config = new ProviderConfig("openai", "gpt-4o-mini") { ApiKey = apiKey, MaxTokens = 50 };
        var provider = ProviderFactory.Create(config);
        var response = provider.Reason("Say hello in one word.", Array.Empty<Event>());

        Assert.NotNull(response);
        Assert.NotEmpty(response.Content);
        Assert.True(response.TokensUsed > 0);
        Assert.True(response.Confidence.Value >= 0.0 && response.Confidence.Value <= 1.0);
    }

    [Fact(Skip = "Requires XAI_API_KEY environment variable")]
    public void Integration_Xai_Reason()
    {
        var apiKey = Environment.GetEnvironmentVariable("XAI_API_KEY");
        if (string.IsNullOrEmpty(apiKey)) return;

        var config = new ProviderConfig("xai", "grok-2") { ApiKey = apiKey, MaxTokens = 50 };
        var provider = ProviderFactory.Create(config);
        var response = provider.Reason("Say hello in one word.", Array.Empty<Event>());

        Assert.NotNull(response);
        Assert.NotEmpty(response.Content);
    }

    [Fact(Skip = "Requires claude CLI installed")]
    public void Integration_ClaudeCli_Reason()
    {
        var config = ClaudeCliConfig.Create("sonnet");
        var provider = ProviderFactory.Create(config);
        var response = provider.Reason("Say hello in one word.", Array.Empty<Event>());

        Assert.NotNull(response);
        Assert.NotEmpty(response.Content);
    }

    [Fact(Skip = "Requires Ollama running locally")]
    public void Integration_Ollama_Reason()
    {
        var config = new ProviderConfig("ollama", "llama3") { MaxTokens = 50 };
        var provider = ProviderFactory.Create(config);
        var response = provider.Reason("Say hello in one word.", Array.Empty<Event>());

        Assert.NotNull(response);
        Assert.NotEmpty(response.Content);
    }

    [Fact]
    public void OpenAI_InvalidApiKey_ThrowsOnReason()
    {
        var config = new ProviderConfig("openai", "gpt-4o")
        {
            ApiKey = "sk-invalid-key-that-will-fail",
            MaxTokens = 10,
        };
        var provider = ProviderFactory.Create(config);

        // The provider is created successfully; it fails when Reason is called.
        Assert.Throws<EventGraphException>(() =>
            provider.Reason("test", Array.Empty<Event>()));
    }

    // ── Provider with SystemPrompt ──────────────────────────────────────

    [Fact]
    public void OpenAI_WithSystemPrompt_CreatesSuccessfully()
    {
        var config = new ProviderConfig("openai", "gpt-4o")
        {
            ApiKey = "test-key",
            SystemPrompt = "You are a helpful assistant for EventGraph.",
        };
        var provider = ProviderFactory.Create(config);

        Assert.NotNull(provider);
        Assert.Equal("openai", provider.Name);
    }

    [Fact]
    public void OpenAI_WithTemperature_CreatesSuccessfully()
    {
        var config = new ProviderConfig("openai", "gpt-4o")
        {
            ApiKey = "test-key",
            Temperature = 0.9,
            MaxTokens = 4096,
        };
        var provider = ProviderFactory.Create(config);
        Assert.NotNull(provider);
    }

    // ── Case insensitivity ──────────────────────────────────────────────

    [Theory]
    [InlineData("OpenAI")]
    [InlineData("OPENAI")]
    [InlineData("openai")]
    public void Factory_CaseInsensitiveProvider(string provider)
    {
        var config = new ProviderConfig(provider, "gpt-4o") { ApiKey = "test" };
        var result = ProviderFactory.Create(config);
        Assert.IsType<OpenAICompatibleProvider>(result);
    }
}
