namespace EventGraph;

/// <summary>
/// A concrete primitive implementation that processes matching events
/// and emits tracking mutations (eventsProcessed, lastTick).
/// </summary>
public sealed class ConcretePrimitive : IPrimitive
{
    public PrimitiveId Id { get; }
    public Layer Layer { get; }
    public List<SubscriptionPattern> Subscriptions { get; }
    public Cadence Cadence { get; }

    public ConcretePrimitive(string name, int layer, List<SubscriptionPattern> subscriptions, int cadence = 1)
    {
        Id = new PrimitiveId(name);
        Layer = new Layer(layer);
        Subscriptions = subscriptions;
        Cadence = new Cadence(cadence);
    }

    public List<Mutation> Process(int tick, List<Event> events, Snapshot snapshot)
    {
        return new List<Mutation>
        {
            new UpdateStateMutation(Id, "eventsProcessed", events.Count),
            new UpdateStateMutation(Id, "lastTick", tick),
        };
    }
}

/// <summary>
/// Factory that creates all 201 primitives across 14 layers (0-13).
/// </summary>
public static class PrimitiveFactory
{
    private static List<SubscriptionPattern> Subs(params string[] patterns) =>
        patterns.Select(p => new SubscriptionPattern(p)).ToList();

    /// <summary>Create all 201 primitives in layer order.</summary>
    public static List<IPrimitive> CreateAll()
    {
        return new List<IPrimitive>
        {
            // ── Layer 0: Foundation (44 primitives) ─────────────────────────
            new ConcretePrimitive("Event", 0, Subs("*")),
            new ConcretePrimitive("EventStore", 0, Subs("store.*")),
            new ConcretePrimitive("Clock", 0, Subs("clock.*")),
            new ConcretePrimitive("Hash", 0, Subs("*")),
            new ConcretePrimitive("Self", 0, Subs("system.*")),
            new ConcretePrimitive("CausalLink", 0, Subs("*")),
            new ConcretePrimitive("Ancestry", 0, Subs("query.*")),
            new ConcretePrimitive("Descendancy", 0, Subs("query.*")),
            new ConcretePrimitive("FirstCause", 0, Subs("query.*")),
            new ConcretePrimitive("ActorID", 0, Subs("actor.*")),
            new ConcretePrimitive("ActorRegistry", 0, Subs("actor.*")),
            new ConcretePrimitive("Signature", 0, Subs("*")),
            new ConcretePrimitive("Verify", 0, Subs("*")),
            new ConcretePrimitive("Expectation", 0, Subs("*")),
            new ConcretePrimitive("Timeout", 0, Subs("authority.*")),
            new ConcretePrimitive("Violation", 0, Subs("violation.*")),
            new ConcretePrimitive("Severity", 0, Subs("violation.*")),
            new ConcretePrimitive("TrustScore", 0, Subs("trust.*")),
            new ConcretePrimitive("TrustUpdate", 0, Subs("trust.*")),
            new ConcretePrimitive("Corroboration", 0, Subs("trust.*")),
            new ConcretePrimitive("Contradiction", 0, Subs("trust.*")),
            new ConcretePrimitive("Confidence", 0, Subs("decision.*")),
            new ConcretePrimitive("Evidence", 0, Subs("*")),
            new ConcretePrimitive("Revision", 0, Subs("grammar.*")),
            new ConcretePrimitive("Uncertainty", 0, Subs("decision.*")),
            new ConcretePrimitive("InstrumentationSpec", 0, Subs("health.*")),
            new ConcretePrimitive("CoverageCheck", 0, Subs("health.*")),
            new ConcretePrimitive("Gap", 0, Subs("*")),
            new ConcretePrimitive("Blind", 0, Subs("health.*")),
            new ConcretePrimitive("PathQuery", 0, Subs("query.*")),
            new ConcretePrimitive("SubgraphExtract", 0, Subs("query.*")),
            new ConcretePrimitive("Annotate", 0, Subs("grammar.*")),
            new ConcretePrimitive("Timeline", 0, Subs("query.*")),
            new ConcretePrimitive("HashChain", 0, Subs("*")),
            new ConcretePrimitive("ChainVerify", 0, Subs("chain.*")),
            new ConcretePrimitive("Witness", 0, Subs("*")),
            new ConcretePrimitive("IntegrityViolation", 0, Subs("chain.*")),
            new ConcretePrimitive("Pattern", 0, Subs("*")),
            new ConcretePrimitive("DeceptionIndicator", 0, Subs("trust.*", "violation.*")),
            new ConcretePrimitive("Suspicion", 0, Subs("trust.*")),
            new ConcretePrimitive("Quarantine", 0, Subs("actor.*")),
            new ConcretePrimitive("GraphHealth", 0, Subs("health.*")),
            new ConcretePrimitive("Invariant", 0, Subs("*")),
            new ConcretePrimitive("InvariantCheck", 0, Subs("chain.*")),
            new ConcretePrimitive("Bootstrap", 0, Subs("system.*")),

            // ── Layer 1: Agency (12 primitives) ──────────────────────────────
            // Volition
            new ConcretePrimitive("Value", 1, Subs("decision.*", "consensus.*", "norm.*")),
            new ConcretePrimitive("Intent", 1, Subs("decision.*", "goal.*")),
            new ConcretePrimitive("Choice", 1, Subs("decision.*", "value.*", "intent.*")),
            new ConcretePrimitive("Risk", 1, Subs("decision.*", "uncertainty.*", "consequence.*")),
            // Action
            new ConcretePrimitive("Act", 1, Subs("choice.*", "intent.*", "plan.*")),
            new ConcretePrimitive("Consequence", 1, Subs("act.*", "decision.*", "harm.*")),
            new ConcretePrimitive("Capacity", 1, Subs("actor.*", "capability.*", "resource.*")),
            new ConcretePrimitive("Resource", 1, Subs("capacity.*", "budget.*", "health.*")),
            // Communication
            new ConcretePrimitive("Signal", 1, Subs("message.*", "act.*")),
            new ConcretePrimitive("Reception", 1, Subs("signal.*", "message.*")),
            new ConcretePrimitive("Acknowledgment", 1, Subs("reception.*", "signal.*")),
            new ConcretePrimitive("Commitment", 1, Subs("choice.*", "acknowledgment.*", "obligation.*")),

            // ── Layer 2: Exchange (12 primitives) ────────────────────────────
            new ConcretePrimitive("Term", 2, Subs("negotiation.*", "contract.*")),
            new ConcretePrimitive("Protocol", 2, Subs("term.*", "norm.*", "standard.*")),
            new ConcretePrimitive("Offer", 2, Subs("signal.*", "term.*", "exchange.*")),
            new ConcretePrimitive("Acceptance", 2, Subs("offer.*", "choice.*")),
            new ConcretePrimitive("Agreement", 2, Subs("acceptance.*", "term.*", "commitment.*")),
            new ConcretePrimitive("Obligation", 2, Subs("agreement.*", "commitment.*")),
            new ConcretePrimitive("Fulfillment", 2, Subs("obligation.*", "act.*")),
            new ConcretePrimitive("Breach", 2, Subs("obligation.*", "violation.*")),
            new ConcretePrimitive("Exchange", 2, Subs("fulfillment.*", "offer.*", "acceptance.*")),
            new ConcretePrimitive("Accountability", 2, Subs("breach.*", "obligation.*", "consequence.*")),
            new ConcretePrimitive("Debt", 2, Subs("breach.*", "obligation.*", "exchange.*")),
            new ConcretePrimitive("Reciprocity", 2, Subs("exchange.*", "fulfillment.*", "obligation.*")),

            // ── Layer 3: Society (12 primitives) ─────────────────────────────
            new ConcretePrimitive("Group", 3, Subs("actor.*", "agreement.*")),
            new ConcretePrimitive("Membership", 3, Subs("group.*", "acceptance.*", "actor.*")),
            new ConcretePrimitive("Role", 3, Subs("group.*", "membership.*", "capacity.*")),
            new ConcretePrimitive("Consent", 3, Subs("choice.*", "membership.*", "authority.*")),
            new ConcretePrimitive("Norm", 3, Subs("pattern.detected", "consensus.*", "group.*")),
            new ConcretePrimitive("Reputation", 3, Subs("trust.*", "fulfillment.*", "violation.*")),
            new ConcretePrimitive("Sanction", 3, Subs("norm.*", "violation.*", "breach.*")),
            new ConcretePrimitive("Authority", 3, Subs("role.*", "consent.*", "group.*")),
            new ConcretePrimitive("Property", 3, Subs("exchange.*", "agreement.*", "group.*")),
            new ConcretePrimitive("Commons", 3, Subs("property.*", "group.*", "norm.*")),
            new ConcretePrimitive("Governance", 3, Subs("authority.*", "norm.*", "consent.*")),
            new ConcretePrimitive("CollectiveAct", 3, Subs("group.*", "consent.*", "act.*")),

            // ── Layer 4: Legal (12 primitives) ───────────────────────────────
            new ConcretePrimitive("Law", 4, Subs("norm.*", "governance.*", "authority.*")),
            new ConcretePrimitive("Right", 4, Subs("law.*", "norm.*", "consent.*")),
            new ConcretePrimitive("Contract", 4, Subs("agreement.*", "law.*", "obligation.*")),
            new ConcretePrimitive("Liability", 4, Subs("breach.*", "harm.*", "contract.*")),
            new ConcretePrimitive("DueProcess", 4, Subs("law.*", "right.*", "sanction.*")),
            new ConcretePrimitive("Adjudication", 4, Subs("dispute.*", "law.*", "evidence.*")),
            new ConcretePrimitive("Remedy", 4, Subs("adjudication.*", "liability.*", "harm.*")),
            new ConcretePrimitive("Precedent", 4, Subs("adjudication.*", "decision.*")),
            new ConcretePrimitive("Jurisdiction", 4, Subs("law.*", "group.*", "authority.*")),
            new ConcretePrimitive("Sovereignty", 4, Subs("jurisdiction.*", "authority.*", "governance.*")),
            new ConcretePrimitive("Legitimacy", 4, Subs("consent.*", "authority.*", "law.*")),
            new ConcretePrimitive("Treaty", 4, Subs("agreement.*", "sovereignty.*", "jurisdiction.*")),

            // ── Layer 5: Technology (12 primitives) ──────────────────────────
            new ConcretePrimitive("Method", 5, Subs("plan.*", "pattern.detected")),
            new ConcretePrimitive("Measurement", 5, Subs("method.*", "evidence.*", "health.*")),
            new ConcretePrimitive("Knowledge", 5, Subs("measurement.*", "evidence.*", "corroboration.*")),
            new ConcretePrimitive("Model", 5, Subs("knowledge.*", "pattern.detected", "abstraction.*")),
            new ConcretePrimitive("Tool", 5, Subs("method.*", "capability.*")),
            new ConcretePrimitive("Technique", 5, Subs("method.*", "tool.*", "knowledge.*")),
            new ConcretePrimitive("Invention", 5, Subs("technique.*", "knowledge.*", "pattern.detected")),
            new ConcretePrimitive("Abstraction", 5, Subs("pattern.detected", "model.*", "knowledge.*")),
            new ConcretePrimitive("Infrastructure", 5, Subs("tool.*", "technique.*", "commons.*")),
            new ConcretePrimitive("Standard", 5, Subs("norm.*", "measurement.*", "infrastructure.*")),
            new ConcretePrimitive("Efficiency", 5, Subs("measurement.*", "method.*", "resource.*")),
            new ConcretePrimitive("Automation", 5, Subs("technique.*", "pattern.detected", "efficiency.*")),

            // ── Layer 6: Information (12 primitives) ─────────────────────────
            new ConcretePrimitive("Symbol", 6, Subs("signal.*", "pattern.detected")),
            new ConcretePrimitive("Language", 6, Subs("symbol.*", "norm.*", "protocol.*")),
            new ConcretePrimitive("Encoding", 6, Subs("symbol.*", "language.*", "signal.*")),
            new ConcretePrimitive("Record", 6, Subs("encoding.*", "store.*", "evidence.*")),
            new ConcretePrimitive("Channel", 6, Subs("signal.*", "protocol.*", "infrastructure.*")),
            new ConcretePrimitive("Copy", 6, Subs("record.*", "encoding.*")),
            new ConcretePrimitive("Noise", 6, Subs("channel.*", "signal.*", "reception.*")),
            new ConcretePrimitive("Redundancy", 6, Subs("copy.*", "record.*", "channel.*")),
            new ConcretePrimitive("Data", 6, Subs("record.*", "measurement.*", "encoding.*")),
            new ConcretePrimitive("Computation", 6, Subs("data.*", "method.*", "automation.*")),
            new ConcretePrimitive("Algorithm", 6, Subs("computation.*", "method.*", "pattern.detected")),
            new ConcretePrimitive("Entropy", 6, Subs("noise.*", "data.*", "uncertainty.*")),

            // ── Layer 7: Ethics (12 primitives) ──────────────────────────────
            new ConcretePrimitive("MoralStatus", 7, Subs("actor.*", "right.*", "dignity.*")),
            new ConcretePrimitive("Dignity", 7, Subs("moral.status.*", "right.*", "harm.*")),
            new ConcretePrimitive("Autonomy", 7, Subs("choice.*", "consent.*", "right.*")),
            new ConcretePrimitive("Flourishing", 7, Subs("value.*", "capacity.*", "care.*")),
            new ConcretePrimitive("Duty", 7, Subs("obligation.*", "moral.status.*", "right.*")),
            new ConcretePrimitive("Harm", 7, Subs("violation.*", "right.*", "consequence.*")),
            new ConcretePrimitive("Care", 7, Subs("harm.*", "flourishing.*", "trust.*")),
            new ConcretePrimitive("Justice", 7, Subs("right.*", "adjudication.*", "fairness.*")),
            new ConcretePrimitive("Conscience", 7, Subs("value.*", "duty.*", "harm.*")),
            new ConcretePrimitive("Virtue", 7, Subs("value.*", "flourishing.*", "conscience.*")),
            new ConcretePrimitive("Responsibility", 7, Subs("duty.*", "consequence.*", "accountability.*")),
            new ConcretePrimitive("Motive", 7, Subs("intent.*", "value.*", "conscience.*")),

            // ── Layer 8: Identity (12 primitives) ────────────────────────────
            new ConcretePrimitive("Narrative", 8, Subs("memory.*", "record.*", "meaning.*")),
            new ConcretePrimitive("SelfConcept", 8, Subs("narrative.*", "value.*", "reflection.*")),
            new ConcretePrimitive("Reflection", 8, Subs("self.concept.*", "conscience.*", "consequence.*")),
            new ConcretePrimitive("Memory", 8, Subs("record.*", "narrative.*", "data.*")),
            new ConcretePrimitive("Purpose", 8, Subs("value.*", "intent.*", "self.concept.*")),
            new ConcretePrimitive("Aspiration", 8, Subs("purpose.*", "value.*", "flourishing.*")),
            new ConcretePrimitive("Authenticity", 8, Subs("self.concept.*", "value.*", "choice.*")),
            new ConcretePrimitive("Expression", 8, Subs("signal.*", "authenticity.*", "language.*")),
            new ConcretePrimitive("Growth", 8, Subs("reflection.*", "aspiration.*", "learning.*")),
            new ConcretePrimitive("Continuity", 8, Subs("memory.*", "narrative.*", "self.concept.*")),
            new ConcretePrimitive("Integration", 8, Subs("self.concept.*", "continuity.*", "growth.*")),
            new ConcretePrimitive("Crisis", 8, Subs("self.concept.*", "rupture.*", "value.*")),

            // ── Layer 9: Relationship (12 primitives) ────────────────────────
            new ConcretePrimitive("Bond", 9, Subs("trust.*", "attachment.*", "commitment.*")),
            new ConcretePrimitive("Attachment", 9, Subs("bond.*", "care.*", "presence.*")),
            new ConcretePrimitive("Recognition", 9, Subs("signal.*", "acknowledgment.*", "dignity.*")),
            new ConcretePrimitive("Intimacy", 9, Subs("bond.*", "trust.*", "vulnerability.*")),
            new ConcretePrimitive("Attunement", 9, Subs("reception.*", "empathy.*", "care.*")),
            new ConcretePrimitive("Rupture", 9, Subs("breach.*", "trust.*", "harm.*")),
            new ConcretePrimitive("Repair", 9, Subs("rupture.*", "care.*", "trust.*")),
            new ConcretePrimitive("Loyalty", 9, Subs("bond.*", "commitment.*", "trust.*")),
            new ConcretePrimitive("MutualConstitution", 9, Subs("bond.*", "recognition.*", "attunement.*")),
            new ConcretePrimitive("RelationalObligation", 9, Subs("bond.*", "duty.*", "care.*")),
            new ConcretePrimitive("Grief", 9, Subs("loss.*", "rupture.*", "actor.memorial")),
            new ConcretePrimitive("Forgiveness", 9, Subs("rupture.*", "repair.*", "trust.*")),

            // ── Layer 10: Community (12 primitives) ──────────────────────────
            new ConcretePrimitive("Culture", 10, Subs("norm.*", "tradition.*", "group.*")),
            new ConcretePrimitive("SharedNarrative", 10, Subs("narrative.*", "group.*", "culture.*")),
            new ConcretePrimitive("Ethos", 10, Subs("value.*", "culture.*", "norm.*")),
            new ConcretePrimitive("Sacred", 10, Subs("value.*", "ethos.*", "tradition.*")),
            new ConcretePrimitive("Tradition", 10, Subs("pattern.detected", "culture.*", "norm.*")),
            new ConcretePrimitive("Ritual", 10, Subs("tradition.*", "sacred.*", "practice.*")),
            new ConcretePrimitive("Practice", 10, Subs("technique.*", "norm.*", "culture.*")),
            new ConcretePrimitive("Place", 10, Subs("group.*", "culture.*", "commons.*")),
            new ConcretePrimitive("Belonging", 10, Subs("membership.*", "culture.*", "bond.*")),
            new ConcretePrimitive("Solidarity", 10, Subs("belonging.*", "collective.act.*", "care.*")),
            new ConcretePrimitive("Voice", 10, Subs("expression.*", "consent.*", "group.*")),
            new ConcretePrimitive("Welcome", 10, Subs("membership.*", "belonging.*", "recognition.*")),

            // ── Layer 11: Culture (12 primitives) ────────────────────────────
            new ConcretePrimitive("Reflexivity", 11, Subs("reflection.*", "self.concept.*", "culture.*")),
            new ConcretePrimitive("Encounter", 11, Subs("recognition.*", "signal.*", "culture.*")),
            new ConcretePrimitive("Translation", 11, Subs("encoding.*", "language.*", "encounter.*")),
            new ConcretePrimitive("Pluralism", 11, Subs("culture.*", "voice.*", "encounter.*")),
            new ConcretePrimitive("Creativity", 11, Subs("invention.*", "expression.*", "imagination.*")),
            new ConcretePrimitive("Aesthetic", 11, Subs("expression.*", "value.*", "creativity.*")),
            new ConcretePrimitive("Interpretation", 11, Subs("symbol.*", "language.*", "narrative.*")),
            new ConcretePrimitive("Dialogue", 11, Subs("signal.*", "reception.*", "encounter.*")),
            new ConcretePrimitive("Syncretism", 11, Subs("culture.*", "encounter.*", "translation.*")),
            new ConcretePrimitive("Critique", 11, Subs("reflexivity.*", "norm.*", "tradition.*")),
            new ConcretePrimitive("Hegemony", 11, Subs("authority.*", "culture.*", "norm.*")),
            new ConcretePrimitive("CulturalEvolution", 11, Subs("culture.*", "creativity.*", "critique.*")),

            // ── Layer 12: Emergence (12 primitives) ──────────────────────────
            new ConcretePrimitive("Emergence", 12, Subs("pattern.detected", "complexity.*", "system.*")),
            new ConcretePrimitive("SelfOrganization", 12, Subs("emergence.*", "pattern.detected", "group.*")),
            new ConcretePrimitive("Feedback", 12, Subs("consequence.*", "measurement.*", "signal.*")),
            new ConcretePrimitive("Complexity", 12, Subs("emergence.*", "pattern.detected", "system.*")),
            new ConcretePrimitive("Consciousness", 12, Subs("self.concept.*", "reflexivity.*", "emergence.*")),
            new ConcretePrimitive("Recursion", 12, Subs("pattern.detected", "feedback.*", "self.organization.*")),
            new ConcretePrimitive("Paradox", 12, Subs("contradiction.*", "recursion.*", "complexity.*")),
            new ConcretePrimitive("Incompleteness", 12, Subs("knowledge.*", "paradox.*", "uncertainty.*")),
            new ConcretePrimitive("PhaseTransition", 12, Subs("emergence.*", "complexity.*", "feedback.*")),
            new ConcretePrimitive("DownwardCausation", 12, Subs("emergence.*", "self.organization.*", "norm.*")),
            new ConcretePrimitive("Autopoiesis", 12, Subs("self.organization.*", "recursion.*", "emergence.*")),
            new ConcretePrimitive("CoEvolution", 12, Subs("feedback.*", "emergence.*", "cultural.evolution.*")),

            // ── Layer 13: Existence (12 primitives) ──────────────────────────
            new ConcretePrimitive("Being", 13, Subs("clock.tick", "consciousness.*")),
            new ConcretePrimitive("Nothingness", 13, Subs("being.*", "absence.*", "void.*")),
            new ConcretePrimitive("Finitude", 13, Subs("actor.memorial", "being.*", "continuity.*")),
            new ConcretePrimitive("Contingency", 13, Subs("uncertainty.*", "risk.*", "being.*")),
            new ConcretePrimitive("Wonder", 13, Subs("mystery.*", "emergence.*", "being.*")),
            new ConcretePrimitive("ExistentialAcceptance", 13, Subs("finitude.*", "contingency.*", "conscience.*")),
            new ConcretePrimitive("Presence", 13, Subs("being.*", "clock.tick", "attunement.*")),
            new ConcretePrimitive("Gratitude", 13, Subs("being.*", "presence.*", "fulfillment.*")),
            new ConcretePrimitive("Mystery", 13, Subs("uncertainty.*", "incompleteness.*", "wonder.*")),
            new ConcretePrimitive("Transcendence", 13, Subs("being.*", "emergence.*", "consciousness.*")),
            new ConcretePrimitive("Groundlessness", 13, Subs("contingency.*", "nothingness.*", "paradox.*")),
            new ConcretePrimitive("Return", 13, Subs("recursion.*", "being.*", "acceptance.*")),
        };
    }

    /// <summary>Create all primitives and register them in a PrimitiveRegistry.</summary>
    public static PrimitiveRegistry CreateRegistry()
    {
        var registry = new PrimitiveRegistry();
        foreach (var p in CreateAll())
            registry.Register(p);
        return registry;
    }
}
