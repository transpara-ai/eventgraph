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

            // ── Layer 1: Agency (12 primitives) ─────────────────────────────
            new ConcretePrimitive("Goal", 1, Subs("decision.*", "authority.resolved", "actor.*")),
            new ConcretePrimitive("Plan", 1, Subs("goal.*")),
            new ConcretePrimitive("Initiative", 1, Subs("clock.tick", "goal.*", "plan.*")),
            new ConcretePrimitive("Commitment", 1, Subs("goal.set", "goal.achieved", "goal.abandoned", "plan.step.completed")),
            new ConcretePrimitive("Focus", 1, Subs("*")),
            new ConcretePrimitive("Filter", 1, Subs("*")),
            new ConcretePrimitive("Salience", 1, Subs("*")),
            new ConcretePrimitive("Distraction", 1, Subs("focus.*", "goal.*")),
            new ConcretePrimitive("Permission", 1, Subs("authority.*", "decision.*")),
            new ConcretePrimitive("Capability", 1, Subs("actor.registered", "permission.*", "trust.*")),
            new ConcretePrimitive("Delegation", 1, Subs("authority.*", "edge.created")),
            new ConcretePrimitive("Accountability", 1, Subs("delegation.*", "violation.*", "goal.abandoned")),

            // ── Layer 2: Communication (12 primitives) ──────────────────────
            new ConcretePrimitive("Message", 2, Subs("protocol.message.*")),
            new ConcretePrimitive("Acknowledgement", 2, Subs("message.sent", "message.received")),
            new ConcretePrimitive("Clarification", 2, Subs("message.*", "ack.*")),
            new ConcretePrimitive("Context", 2, Subs("message.*", "clarification.*")),
            new ConcretePrimitive("Offer", 2, Subs("message.*", "exchange.*")),
            new ConcretePrimitive("Acceptance", 2, Subs("offer.made", "offer.withdrawn")),
            new ConcretePrimitive("Obligation", 2, Subs("offer.accepted", "delegation.*")),
            new ConcretePrimitive("Gratitude", 2, Subs("obligation.fulfilled", "trust.*")),
            new ConcretePrimitive("Negotiation", 2, Subs("offer.*", "message.*")),
            new ConcretePrimitive("Consent", 2, Subs("negotiation.concluded", "offer.accepted", "authority.*")),
            new ConcretePrimitive("Contract", 2, Subs("consent.given", "negotiation.concluded")),
            new ConcretePrimitive("Dispute", 2, Subs("contract.breached", "obligation.defaulted", "contradiction.found")),

            // ── Layer 3: Social (12 primitives) ─────────────────────────────
            new ConcretePrimitive("Group", 3, Subs("actor.*", "consent.*")),
            new ConcretePrimitive("Role", 3, Subs("group.*", "delegation.*")),
            new ConcretePrimitive("Reputation", 3, Subs("trust.*", "commitment.*", "violation.*", "gratitude.*")),
            new ConcretePrimitive("Exclusion", 3, Subs("reputation.*", "violation.*", "quarantine.*", "dispute.*")),
            new ConcretePrimitive("Vote", 3, Subs("authority.requested", "group.*")),
            new ConcretePrimitive("Consensus", 3, Subs("message.*", "corroboration.*", "vote.result")),
            new ConcretePrimitive("Dissent", 3, Subs("vote.*", "consensus.*", "contradiction.found")),
            new ConcretePrimitive("Majority", 3, Subs("vote.result", "dissent.*", "exclusion.*")),
            new ConcretePrimitive("Convention", 3, Subs("pattern.detected", "*")),
            new ConcretePrimitive("Norm", 3, Subs("convention.detected", "consensus.reached")),
            new ConcretePrimitive("Sanction", 3, Subs("norm.violated", "violation.*")),
            new ConcretePrimitive("Forgiveness", 3, Subs("sanction.applied", "trust.*", "obligation.fulfilled")),

            // ── Layer 4: Governance (12 primitives) ─────────────────────────
            new ConcretePrimitive("Rule", 4, Subs("norm.established", "consensus.reached", "vote.result")),
            new ConcretePrimitive("Jurisdiction", 4, Subs("rule.*", "group.*")),
            new ConcretePrimitive("Precedent", 4, Subs("dispute.resolved", "decision.*")),
            new ConcretePrimitive("Interpretation", 4, Subs("rule.*", "dispute.*", "precedent.*")),
            new ConcretePrimitive("Adjudication", 4, Subs("dispute.raised", "rule.*", "precedent.*")),
            new ConcretePrimitive("Appeal", 4, Subs("adjudication.ruling", "exclusion.enacted", "sanction.applied")),
            new ConcretePrimitive("DueProcess", 4, Subs("adjudication.*", "exclusion.*", "sanction.*")),
            new ConcretePrimitive("Rights", 4, Subs("rule.*", "sanction.*", "exclusion.*", "dueprocess.*")),
            new ConcretePrimitive("Audit", 4, Subs("clock.tick", "rule.*")),
            new ConcretePrimitive("Enforcement", 4, Subs("audit.*", "rule.*", "right.violated")),
            new ConcretePrimitive("Amnesty", 4, Subs("enforcement.*", "vote.result")),
            new ConcretePrimitive("Reform", 4, Subs("precedent.*", "right.violated", "audit.*", "dissent.*")),

            // ── Layer 5: Production (12 primitives) ─────────────────────────
            new ConcretePrimitive("Create", 5, Subs("plan.*", "goal.*")),
            new ConcretePrimitive("Tool", 5, Subs("artefact.created", "capability.*")),
            new ConcretePrimitive("Quality", 5, Subs("artefact.*", "tool.used")),
            new ConcretePrimitive("Deprecation", 5, Subs("quality.*", "artefact.version")),
            new ConcretePrimitive("Workflow", 5, Subs("plan.*", "convention.detected")),
            new ConcretePrimitive("Automation", 5, Subs("workflow.executed", "pattern.detected")),
            new ConcretePrimitive("Testing", 5, Subs("artefact.*", "workflow.*", "automation.*")),
            new ConcretePrimitive("Review", 5, Subs("artefact.*", "decision.*")),
            new ConcretePrimitive("Feedback", 5, Subs("*")),
            new ConcretePrimitive("Iteration", 5, Subs("feedback.*", "test.*", "review.*")),
            new ConcretePrimitive("Innovation", 5, Subs("artefact.*", "pattern.detected")),
            new ConcretePrimitive("Legacy", 5, Subs("deprecation.*", "artefact.*")),

            // ── Layer 6: Knowledge (12 primitives) ──────────────────────────
            new ConcretePrimitive("Symbol", 6, Subs("*")),
            new ConcretePrimitive("Abstraction", 6, Subs("pattern.detected", "symbol.*")),
            new ConcretePrimitive("Classification", 6, Subs("*")),
            new ConcretePrimitive("Encoding", 6, Subs("symbol.*", "message.*")),
            new ConcretePrimitive("Fact", 6, Subs("corroboration.*", "evidence.*", "confidence.*")),
            new ConcretePrimitive("Inference", 6, Subs("fact.*", "evidence.*")),
            new ConcretePrimitive("Memory", 6, Subs("fact.*", "inference.*", "abstraction.*")),
            new ConcretePrimitive("Learning", 6, Subs("feedback.*", "test.*", "inference.*")),
            new ConcretePrimitive("Narrative", 6, Subs("fact.*", "inference.*", "memory.*")),
            new ConcretePrimitive("Bias", 6, Subs("narrative.*", "classification.*", "inference.*")),
            new ConcretePrimitive("Correction", 6, Subs("bias.detected", "fact.retracted", "contradiction.found")),
            new ConcretePrimitive("Provenance", 6, Subs("fact.*", "memory.*", "message.*")),

            // ── Layer 7: Ethics (12 primitives) ─────────────────────────────
            new ConcretePrimitive("Value", 7, Subs("consensus.*", "norm.*", "right.*")),
            new ConcretePrimitive("Harm", 7, Subs("violation.*", "right.violated", "exclusion.*")),
            new ConcretePrimitive("Fairness", 7, Subs("decision.*", "sanction.*", "exclusion.*", "bias.detected")),
            new ConcretePrimitive("Care", 7, Subs("harm.*", "health.*", "trust.*")),
            new ConcretePrimitive("Dilemma", 7, Subs("value.conflict", "decision.*")),
            new ConcretePrimitive("Proportionality", 7, Subs("enforcement.*", "sanction.*", "harm.*")),
            new ConcretePrimitive("Intention", 7, Subs("decision.*", "goal.*", "initiative.*")),
            new ConcretePrimitive("Consequence", 7, Subs("decision.*", "harm.*", "goal.achieved", "goal.abandoned")),
            new ConcretePrimitive("Responsibility", 7, Subs("intention.*", "consequence.*", "accountability.traced")),
            new ConcretePrimitive("Transparency", 7, Subs("decision.*", "adjudication.*")),
            new ConcretePrimitive("Redress", 7, Subs("harm.*", "responsibility.*")),
            new ConcretePrimitive("Growth", 7, Subs("redress.*", "responsibility.*", "learning.*")),

            // ── Layer 8: Identity (12 primitives) ───────────────────────────
            new ConcretePrimitive("SelfModel", 8, Subs("commitment.*", "learning.*", "moral.growth", "capability.*")),
            new ConcretePrimitive("Authenticity", 8, Subs("self.model.*", "decision.*", "value.*")),
            new ConcretePrimitive("NarrativeIdentity", 8, Subs("self.model.*", "narrative.*", "memory.*")),
            new ConcretePrimitive("Boundary", 8, Subs("delegation.*", "group.*", "consent.*")),
            new ConcretePrimitive("Persistence", 8, Subs("self.model.*", "learning.*")),
            new ConcretePrimitive("Transformation", 8, Subs("self.model.*", "moral.growth", "learning.*")),
            new ConcretePrimitive("Heritage", 8, Subs("memory.*", "legacy.*", "provenance.*")),
            new ConcretePrimitive("Aspiration", 8, Subs("self.model.*", "goal.*", "value.*")),
            new ConcretePrimitive("Dignity", 8, Subs("exclusion.*", "harm.*", "right.violated", "actor.memorial")),
            new ConcretePrimitive("IdentityAcknowledgement", 8, Subs("message.*", "gratitude.*", "reputation.*")),
            new ConcretePrimitive("Uniqueness", 8, Subs("self.model.*", "identity.narrative", "pattern.detected")),
            new ConcretePrimitive("Memorial", 8, Subs("actor.memorial")),

            // ── Layer 9: Relational (12 primitives) ─────────────────────────
            new ConcretePrimitive("Attachment", 9, Subs("trust.*", "gratitude.*", "message.*", "edge.created")),
            new ConcretePrimitive("Reciprocity", 9, Subs("obligation.*", "gratitude.*", "offer.*")),
            new ConcretePrimitive("RelationalTrust", 9, Subs("trust.*", "attachment.*", "reciprocity.*")),
            new ConcretePrimitive("Rupture", 9, Subs("contract.breached", "trust.*", "dispute.*", "dignity.violated")),
            new ConcretePrimitive("Apology", 9, Subs("rupture.detected", "harm.*", "responsibility.*")),
            new ConcretePrimitive("Reconciliation", 9, Subs("apology.*", "forgiveness.*", "trust.*")),
            new ConcretePrimitive("RelationalGrowth", 9, Subs("reconciliation.*", "attachment.*")),
            new ConcretePrimitive("Loss", 9, Subs("actor.memorial", "rupture.*", "exclusion.enacted")),
            new ConcretePrimitive("Vulnerability", 9, Subs("relational.trust", "boundary.*")),
            new ConcretePrimitive("Understanding", 9, Subs("self.model.*", "message.*", "vulnerability.*")),
            new ConcretePrimitive("Empathy", 9, Subs("harm.*", "loss.*", "understanding.*")),
            new ConcretePrimitive("Presence", 9, Subs("message.*", "clock.tick")),

            // ── Layer 10: Community (12 primitives) ─────────────────────────
            new ConcretePrimitive("Home", 10, Subs("group.*", "attachment.*", "presence.*")),
            new ConcretePrimitive("Contribution", 10, Subs("artefact.created", "review.*", "care.action")),
            new ConcretePrimitive("Inclusion", 10, Subs("group.*", "exclusion.*", "fairness.*")),
            new ConcretePrimitive("Tradition", 10, Subs("convention.detected", "heritage.*", "pattern.detected")),
            new ConcretePrimitive("Commons", 10, Subs("artefact.*", "group.*")),
            new ConcretePrimitive("Sustainability", 10, Subs("health.*", "commons.*", "contribution.*")),
            new ConcretePrimitive("Succession", 10, Subs("delegation.*", "actor.memorial", "role.*")),
            new ConcretePrimitive("Renewal", 10, Subs("sustainability.*", "innovation.*", "tradition.evolved")),
            new ConcretePrimitive("Milestone", 10, Subs("goal.achieved", "innovation.*", "reconciliation.completed")),
            new ConcretePrimitive("Ceremony", 10, Subs("milestone.*", "succession.*", "actor.memorial")),
            new ConcretePrimitive("Story", 10, Subs("milestone.*", "ceremony.*", "tradition.*", "memorial.created")),
            new ConcretePrimitive("Gift", 10, Subs("contribution.*", "gratitude.*")),

            // ── Layer 11: Reflective (12 primitives) ────────────────────────
            new ConcretePrimitive("SelfAwareness", 11, Subs("health.*", "self.model.*", "bias.detected")),
            new ConcretePrimitive("Perspective", 11, Subs("narrative.*", "dissent.*", "value.conflict")),
            new ConcretePrimitive("Critique", 11, Subs("convention.*", "norm.*", "tradition.*")),
            new ConcretePrimitive("Wisdom", 11, Subs("learning.*", "moral.growth", "consequence.*", "memory.*")),
            new ConcretePrimitive("Aesthetic", 11, Subs("artefact.*", "quality.*")),
            new ConcretePrimitive("Metaphor", 11, Subs("abstraction.*", "symbol.*", "narrative.*")),
            new ConcretePrimitive("Humour", 11, Subs("contradiction.found", "perspective.shift", "*")),
            new ConcretePrimitive("Silence", 11, Subs("clock.tick", "presence.*", "acknowledgement.absent")),
            new ConcretePrimitive("Teaching", 11, Subs("learning.*", "wisdom.*", "memory.*")),
            new ConcretePrimitive("Translation", 11, Subs("encoding.*", "message.*")),
            new ConcretePrimitive("Archive", 11, Subs("memory.*", "legacy.*", "community.story")),
            new ConcretePrimitive("Prophecy", 11, Subs("pattern.detected", "sustainability.*", "wisdom.*")),

            // ── Layer 12: Systemic (12 primitives) ──────────────────────────
            new ConcretePrimitive("MetaPattern", 12, Subs("pattern.detected", "convention.detected", "abstraction.formed")),
            new ConcretePrimitive("SystemDynamic", 12, Subs("health.*", "meta.pattern", "sustainability.*")),
            new ConcretePrimitive("FeedbackLoop", 12, Subs("system.dynamic", "pattern.detected")),
            new ConcretePrimitive("Threshold", 12, Subs("system.dynamic", "feedback.loop", "meta.pattern")),
            new ConcretePrimitive("Adaptation", 12, Subs("feedback.*", "system.dynamic", "sustainability.*")),
            new ConcretePrimitive("Selection", 12, Subs("adaptation.*", "test.*", "quality.*")),
            new ConcretePrimitive("Complexification", 12, Subs("system.dynamic", "innovation.*", "meta.pattern")),
            new ConcretePrimitive("Simplification", 12, Subs("complexity.*", "automation.*")),
            new ConcretePrimitive("SystemicIntegrity", 12, Subs("health.*", "invariant.*", "system.dynamic")),
            new ConcretePrimitive("Harmony", 12, Subs("system.dynamic", "feedback.loop", "dispute.*")),
            new ConcretePrimitive("Resilience", 12, Subs("threshold.*", "rupture.*", "sustainability.*")),
            new ConcretePrimitive("Purpose", 12, Subs("value.*", "goal.*", "wisdom.*")),

            // ── Layer 13: Existential (12 primitives) ───────────────────────
            new ConcretePrimitive("Being", 13, Subs("clock.tick")),
            new ConcretePrimitive("Finitude", 13, Subs("actor.memorial", "sustainability.*", "threshold.*")),
            new ConcretePrimitive("Change", 13, Subs("*")),
            new ConcretePrimitive("Interdependence", 13, Subs("system.dynamic", "attachment.*", "relational.trust")),
            new ConcretePrimitive("Mystery", 13, Subs("uncertainty.*", "wisdom.*", "self.awareness.*")),
            new ConcretePrimitive("Paradox", 13, Subs("contradiction.found", "dilemma.*", "meta.pattern")),
            new ConcretePrimitive("Infinity", 13, Subs("complexity.*", "threshold.*")),
            new ConcretePrimitive("Void", 13, Subs("silence.*", "loss.*", "instrumentation.blind")),
            new ConcretePrimitive("Awe", 13, Subs("mystery.*", "infinity.*", "complexity.*")),
            new ConcretePrimitive("ExistentialGratitude", 13, Subs("being.affirmed", "milestone.*")),
            new ConcretePrimitive("Play", 13, Subs("humour.*", "innovation.*", "*")),
            new ConcretePrimitive("Wonder", 13, Subs("*")),
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
