namespace EventGraph;

// ── TrustMetrics ────────────────────────────────────────────────────────

/// <summary>Immutable trust metrics snapshot for an actor.</summary>
public sealed class TrustMetrics
{
    public ActorId Actor { get; }
    public Score Overall { get; }
    private readonly Dictionary<string, Score> _byDomain;
    public IReadOnlyDictionary<string, Score> ByDomain => new Dictionary<string, Score>(_byDomain);
    public Score Confidence { get; }
    public Weight Trend { get; }
    private readonly List<EventId> _evidence;
    public IReadOnlyList<EventId> Evidence => _evidence.AsReadOnly();
    public long LastUpdatedNanos { get; }
    public Score DecayRate { get; }

    public TrustMetrics(
        ActorId actor,
        Score overall,
        IReadOnlyDictionary<string, Score>? byDomain,
        Score confidence,
        Weight trend,
        IReadOnlyList<EventId>? evidence,
        long lastUpdatedNanos,
        Score decayRate)
    {
        Actor = actor;
        Overall = overall;
        _byDomain = byDomain is not null
            ? new Dictionary<string, Score>(byDomain)
            : new Dictionary<string, Score>();
        Confidence = confidence;
        Trend = trend;
        _evidence = evidence is not null
            ? new List<EventId>(evidence)
            : new List<EventId>();
        LastUpdatedNanos = lastUpdatedNanos;
        DecayRate = decayRate;
    }
}

// ── TrustConfig ─────────────────────────────────────────────────────────

/// <summary>Configuration for the default trust model.</summary>
public record TrustConfig
{
    public Score InitialTrust { get; init; } = new(0.0);
    public Score DecayRate { get; init; } = new(0.01);
    public Weight MaxAdjustment { get; init; } = new(0.1);
    public double ObservedEventDelta { get; init; } = 0.01;
    public double TrendDecayRate { get; init; } = 0.01;
}

// ── ITrustModel ─────────────────────────────────────────────────────────

/// <summary>Calculates, updates, and decays trust.</summary>
public interface ITrustModel
{
    /// <summary>Get current trust metrics for an actor.</summary>
    TrustMetrics Score(Actor actor);

    /// <summary>Get trust metrics for an actor in a specific domain.</summary>
    TrustMetrics ScoreInDomain(Actor actor, DomainScope domain);

    /// <summary>Update trust for an actor based on evidence.</summary>
    TrustMetrics Update(Actor actor, Event evidence);

    /// <summary>Update directional trust from one actor toward another.</summary>
    TrustMetrics UpdateBetween(Actor from, Actor to, Event evidence);

    /// <summary>Apply time-based trust decay for an actor.</summary>
    TrustMetrics Decay(Actor actor, TimeSpan elapsed);

    /// <summary>Get directional trust from one actor toward another.</summary>
    TrustMetrics Between(Actor from, Actor to);
}

// ── DefaultTrustModel ───────────────────────────────────────────────────

/// <summary>Thread-safe default trust model with linear decay and equal weighting.</summary>
public sealed class DefaultTrustModel : ITrustModel
{
    private readonly TrustConfig _config;
    private readonly object _lock = new();
    private readonly Dictionary<string, TrustState> _scores = new();
    private readonly Dictionary<(string from, string to), TrustState> _directed = new();

    private sealed class TrustState
    {
        public Score Score { get; set; }
        public Dictionary<string, Score> ByDomain { get; } = new();
        public List<EventId> Evidence { get; } = new();
        public long LastUpdatedNanos { get; set; }
        public Weight Trend { get; set; }
    }

    public DefaultTrustModel() : this(new TrustConfig()) { }

    public DefaultTrustModel(TrustConfig config)
    {
        _config = config;
    }

    private static long NowNanos() =>
        DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() * 1_000_000;

    private TrustState GetOrCreate(string actorId)
    {
        if (_scores.TryGetValue(actorId, out var state))
            return state;

        state = new TrustState
        {
            Score = _config.InitialTrust,
            LastUpdatedNanos = NowNanos(),
            Trend = new Weight(0.0),
        };
        _scores[actorId] = state;
        return state;
    }

    private TrustState GetOrDefault(string actorId)
    {
        if (_scores.TryGetValue(actorId, out var state))
            return state;

        // Return a default without storing it
        return new TrustState
        {
            Score = _config.InitialTrust,
            LastUpdatedNanos = NowNanos(),
            Trend = new Weight(0.0),
        };
    }

    public TrustMetrics Score(Actor actor)
    {
        lock (_lock)
        {
            var state = GetOrDefault(actor.Id.Value);
            return BuildMetrics(actor.Id, state);
        }
    }

    public TrustMetrics ScoreInDomain(Actor actor, DomainScope domain)
    {
        lock (_lock)
        {
            var state = GetOrDefault(actor.Id.Value);

            // If domain-specific score exists, return metrics with that score
            if (state.ByDomain.TryGetValue(domain.Value, out var domainScore))
                return BuildDomainMetrics(actor.Id, state, domainScore);

            // Fall back to global score with halved confidence
            var evidenceCount = state.Evidence.Count;
            var globalConfidence = Math.Min(1.0, evidenceCount / 50.0);
            return new TrustMetrics(
                actor.Id,
                state.Score,
                state.ByDomain,
                new Score(globalConfidence * 0.5),
                state.Trend,
                state.Evidence,
                state.LastUpdatedNanos,
                _config.DecayRate);
        }
    }

    public TrustMetrics Update(Actor actor, Event evidence)
    {
        lock (_lock)
        {
            var state = GetOrCreate(actor.Id.Value);

            // Deduplicate
            foreach (var id in state.Evidence)
            {
                if (id == evidence.Id)
                    return BuildMetrics(actor.Id, state);
            }

            var delta = ExtractDeltaForScore(evidence, state.Score.Value);

            // Clamp to MaxAdjustment
            var maxAdj = _config.MaxAdjustment.Value;
            delta = Math.Clamp(delta, -maxAdj, maxAdj);

            // Apply delta, clamp to [0, 1]
            var newScore = Math.Clamp(state.Score.Value + delta, 0.0, 1.0);
            state.Score = new Score(newScore);

            // Update domain-specific score if evidence carries domain info
            UpdateDomainFromEvidence(state, evidence, delta);

            // Update trend
            if (delta > 0)
                state.Trend = new Weight(Math.Min(1.0, state.Trend.Value + 0.1));
            else if (delta < 0)
                state.Trend = new Weight(Math.Max(-1.0, state.Trend.Value - 0.1));

            // Track evidence, cap at 100
            state.Evidence.Add(evidence.Id);
            if (state.Evidence.Count > 100)
                state.Evidence.RemoveRange(0, state.Evidence.Count - 100);

            state.LastUpdatedNanos = NowNanos();

            return BuildMetrics(actor.Id, state);
        }
    }

    public TrustMetrics UpdateBetween(Actor from, Actor to, Event evidence)
    {
        lock (_lock)
        {
            var key = (from: from.Id.Value, to: to.Id.Value);
            if (!_directed.TryGetValue(key, out var state))
            {
                state = new TrustState
                {
                    Score = _config.InitialTrust,
                    LastUpdatedNanos = NowNanos(),
                    Trend = new Weight(0.0),
                };
                _directed[key] = state;
            }

            // Deduplicate
            foreach (var id in state.Evidence)
            {
                if (id == evidence.Id)
                    return BuildMetrics(to.Id, state);
            }

            var delta = ExtractDeltaForScore(evidence, state.Score.Value);

            var maxAdj = _config.MaxAdjustment.Value;
            delta = Math.Clamp(delta, -maxAdj, maxAdj);

            var newScore = Math.Clamp(state.Score.Value + delta, 0.0, 1.0);
            state.Score = new Score(newScore);

            UpdateDomainFromEvidence(state, evidence, delta);

            if (delta > 0)
                state.Trend = new Weight(Math.Min(1.0, state.Trend.Value + 0.1));
            else if (delta < 0)
                state.Trend = new Weight(Math.Max(-1.0, state.Trend.Value - 0.1));

            state.Evidence.Add(evidence.Id);
            if (state.Evidence.Count > 100)
                state.Evidence.RemoveRange(0, state.Evidence.Count - 100);

            state.LastUpdatedNanos = NowNanos();

            return BuildMetrics(to.Id, state);
        }
    }

    public TrustMetrics Decay(Actor actor, TimeSpan elapsed)
    {
        lock (_lock)
        {
            // Guard against negative durations
            if (elapsed <= TimeSpan.Zero)
            {
                var defaultState = GetOrDefault(actor.Id.Value);
                return BuildMetrics(actor.Id, defaultState);
            }

            var days = elapsed.TotalDays;
            var decayAmount = _config.DecayRate.Value * days;
            var trendDecay = _config.TrendDecayRate * days;

            // Decay undirected trust
            var hasUndirected = _scores.TryGetValue(actor.Id.Value, out var state);
            if (hasUndirected)
                DecayState(state!, decayAmount, trendDecay);

            // Decay directed trust where actor is "from"
            foreach (var kvp in _directed)
            {
                if (kvp.Key.from == actor.Id.Value)
                    DecayState(kvp.Value, decayAmount, trendDecay);
            }

            if (!hasUndirected)
            {
                return new TrustMetrics(
                    actor.Id,
                    _config.InitialTrust,
                    null,
                    new Score(0.0),
                    new Weight(0.0),
                    null,
                    NowNanos(),
                    _config.DecayRate);
            }

            return BuildMetrics(actor.Id, state!);
        }
    }

    public TrustMetrics Between(Actor from, Actor to)
    {
        lock (_lock)
        {
            var key = (from: from.Id.Value, to: to.Id.Value);
            if (_directed.TryGetValue(key, out var state))
                return BuildMetrics(to.Id, state);

            // No direct trust relationship
            return new TrustMetrics(
                to.Id,
                _config.InitialTrust,
                null,
                new Score(0.0),
                new Weight(0.0),
                null,
                NowNanos(),
                _config.DecayRate);
        }
    }

    // ── Private helpers ────────────────────────────────────────────────

    private double ExtractDeltaForScore(Event ev, double currentScore)
    {
        var content = ev.Content;
        if (content.TryGetValue("current", out var currentObj) && currentObj is not null)
        {
            var current = Convert.ToDouble(currentObj);
            return current - currentScore;
        }
        return _config.ObservedEventDelta;
    }

    private void UpdateDomainFromEvidence(TrustState state, Event evidence, double delta)
    {
        var content = evidence.Content;
        if (content.TryGetValue("domain", out var domainObj) && domainObj is string domainStr)
        {
            double domainScore;
            if (state.ByDomain.TryGetValue(domainStr, out var existing))
                domainScore = Math.Clamp(existing.Value + delta, 0.0, 1.0);
            else
                domainScore = Math.Clamp(_config.InitialTrust.Value + delta, 0.0, 1.0);

            state.ByDomain[domainStr] = new Score(domainScore);
        }
    }

    private static void DecayState(TrustState state, double decayAmount, double trendDecay)
    {
        state.Score = new Score(Math.Max(0.0, state.Score.Value - decayAmount));

        foreach (var domain in state.ByDomain.Keys.ToList())
        {
            state.ByDomain[domain] = new Score(Math.Max(0.0, state.ByDomain[domain].Value - decayAmount));
        }

        if (state.Trend.Value > 0)
            state.Trend = new Weight(Math.Max(0.0, state.Trend.Value - trendDecay));
        else if (state.Trend.Value < 0)
            state.Trend = new Weight(Math.Min(0.0, state.Trend.Value + trendDecay));

        state.LastUpdatedNanos = NowNanos();
    }

    private TrustMetrics BuildMetrics(ActorId actorId, TrustState state)
    {
        var evidenceCount = state.Evidence.Count;
        var confidence = Math.Min(1.0, evidenceCount / 50.0);

        return new TrustMetrics(
            actorId,
            state.Score,
            state.ByDomain,
            new Score(confidence),
            state.Trend,
            state.Evidence,
            state.LastUpdatedNanos,
            _config.DecayRate);
    }

    private TrustMetrics BuildDomainMetrics(ActorId actorId, TrustState state, Score domainScore)
    {
        var evidenceCount = state.Evidence.Count;
        var confidence = Math.Min(1.0, evidenceCount / 50.0);

        return new TrustMetrics(
            actorId,
            domainScore,
            state.ByDomain,
            new Score(confidence),
            state.Trend,
            state.Evidence,
            state.LastUpdatedNanos,
            _config.DecayRate);
    }
}
