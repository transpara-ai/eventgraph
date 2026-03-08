// All 201 EventGraph primitives across 14 layers (0-13).
//
// Each primitive is a struct implementing the Primitive trait from primitive.rs.
// The `create_all_primitives()` function returns all 201 as boxed trait objects.

use serde_json::Value;

use crate::event::Event;
use crate::primitive::{Mutation, Primitive, Snapshot};
use crate::types::*;

// ── Generic primitive implementation ──────────────────────────────────

struct GenericPrimitive {
    name: &'static str,
    layer_value: u8,
    subscription_strs: &'static [&'static str],
}

impl GenericPrimitive {
    fn new(name: &'static str, layer_value: u8, subs: &'static [&'static str]) -> Self {
        Self {
            name,
            layer_value,
            subscription_strs: subs,
        }
    }
}

impl Primitive for GenericPrimitive {
    fn id(&self) -> PrimitiveId {
        PrimitiveId::new(self.name).unwrap()
    }

    fn layer(&self) -> Layer {
        Layer::new(self.layer_value).unwrap()
    }

    fn process(&self, tick: u64, events: &[Event], _snapshot: &Snapshot) -> Vec<Mutation> {
        vec![
            Mutation::UpdateState {
                primitive_id: PrimitiveId::new(self.name).unwrap(),
                key: "eventsProcessed".to_string(),
                value: Value::Number(serde_json::Number::from(events.len())),
            },
            Mutation::UpdateState {
                primitive_id: PrimitiveId::new(self.name).unwrap(),
                key: "lastTick".to_string(),
                value: Value::Number(serde_json::Number::from(tick)),
            },
        ]
    }

    fn subscriptions(&self) -> Vec<SubscriptionPattern> {
        self.subscription_strs
            .iter()
            .map(|s| SubscriptionPattern::new(*s).unwrap())
            .collect()
    }

    fn cadence(&self) -> Cadence {
        Cadence::new(1).unwrap()
    }
}

// ── Macro for concise primitive definitions ───────────────────────────

macro_rules! prim {
    ($name:expr, $layer:expr, $( $sub:expr ),+ $(,)?) => {
        Box::new(GenericPrimitive::new($name, $layer, &[$( $sub ),+])) as Box<dyn Primitive>
    };
}

// ── Factory function ──────────────────────────────────────────────────

/// Creates all 201 primitives across the 14 layers.
pub fn create_all_primitives() -> Vec<Box<dyn Primitive>> {
    vec![
        // ── Layer 0: Foundation (44 primitives) ──────────────────────
        prim!("Event", 0, "*"),
        prim!("EventStore", 0, "store.*"),
        prim!("Clock", 0, "clock.*"),
        prim!("Hash", 0, "*"),
        prim!("Self", 0, "system.*"),
        prim!("CausalLink", 0, "*"),
        prim!("Ancestry", 0, "query.*"),
        prim!("Descendancy", 0, "query.*"),
        prim!("FirstCause", 0, "query.*"),
        prim!("ActorID", 0, "actor.*"),
        prim!("ActorRegistry", 0, "actor.*"),
        prim!("Signature", 0, "*"),
        prim!("Verify", 0, "*"),
        prim!("Expectation", 0, "*"),
        prim!("Timeout", 0, "authority.*"),
        prim!("Violation", 0, "violation.*"),
        prim!("Severity", 0, "violation.*"),
        prim!("TrustScore", 0, "trust.*"),
        prim!("TrustUpdate", 0, "trust.*"),
        prim!("Corroboration", 0, "trust.*"),
        prim!("Contradiction", 0, "trust.*"),
        prim!("Confidence", 0, "decision.*"),
        prim!("Evidence", 0, "*"),
        prim!("Revision", 0, "grammar.*"),
        prim!("Uncertainty", 0, "decision.*"),
        prim!("InstrumentationSpec", 0, "health.*"),
        prim!("CoverageCheck", 0, "health.*"),
        prim!("Gap", 0, "*"),
        prim!("Blind", 0, "health.*"),
        prim!("PathQuery", 0, "query.*"),
        prim!("SubgraphExtract", 0, "query.*"),
        prim!("Annotate", 0, "grammar.*"),
        prim!("Timeline", 0, "query.*"),
        prim!("HashChain", 0, "*"),
        prim!("ChainVerify", 0, "chain.*"),
        prim!("Witness", 0, "*"),
        prim!("IntegrityViolation", 0, "chain.*"),
        prim!("Pattern", 0, "*"),
        prim!("DeceptionIndicator", 0, "trust.*", "violation.*"),
        prim!("Suspicion", 0, "trust.*"),
        prim!("Quarantine", 0, "actor.*"),
        prim!("GraphHealth", 0, "health.*"),
        prim!("Invariant", 0, "*"),
        prim!("InvariantCheck", 0, "chain.*"),
        prim!("Bootstrap", 0, "system.*"),

        // ── Layer 1: Agency (12 primitives) ──────────────────────────
        prim!("Goal", 1, "decision.*", "authority.resolved", "actor.*"),
        prim!("Plan", 1, "goal.*"),
        prim!("Initiative", 1, "clock.tick", "goal.*", "plan.*"),
        prim!("Commitment", 1, "goal.set", "goal.achieved", "goal.abandoned", "plan.step.completed"),
        prim!("Focus", 1, "*"),
        prim!("Filter", 1, "*"),
        prim!("Salience", 1, "*"),
        prim!("Distraction", 1, "focus.*", "goal.*"),
        prim!("Permission", 1, "authority.*", "decision.*"),
        prim!("Capability", 1, "actor.registered", "permission.*", "trust.*"),
        prim!("Delegation", 1, "authority.*", "edge.created"),
        prim!("Accountability", 1, "delegation.*", "violation.*", "goal.abandoned"),

        // ── Layer 2: Communication (11 primitives) ───────────────────
        prim!("Message", 2, "protocol.message.*"),
        prim!("Acknowledgement", 2, "message.sent", "message.received"),
        prim!("Clarification", 2, "message.*", "ack.*"),
        prim!("Context", 2, "message.*", "clarification.*"),
        prim!("Offer", 2, "message.*", "exchange.*"),
        prim!("Acceptance", 2, "offer.made", "offer.withdrawn"),
        prim!("Obligation", 2, "offer.accepted", "delegation.*"),
        prim!("Gratitude", 2, "obligation.fulfilled", "trust.*"),
        prim!("Negotiation", 2, "offer.*", "message.*"),
        prim!("Consent", 2, "negotiation.concluded", "offer.accepted", "authority.*"),
        prim!("Contract", 2, "consent.given", "negotiation.concluded"),
        prim!("Dispute", 2, "contract.breached", "obligation.defaulted", "contradiction.found"),

        // ── Layer 3: Social (12 primitives) ──────────────────────────
        prim!("Group", 3, "actor.*", "consent.*"),
        prim!("Role", 3, "group.*", "delegation.*"),
        prim!("Reputation", 3, "trust.*", "commitment.*", "violation.*", "gratitude.*"),
        prim!("Exclusion", 3, "reputation.*", "violation.*", "quarantine.*", "dispute.*"),
        prim!("Vote", 3, "authority.requested", "group.*"),
        prim!("Consensus", 3, "message.*", "corroboration.*", "vote.result"),
        prim!("Dissent", 3, "vote.*", "consensus.*", "contradiction.found"),
        prim!("Majority", 3, "vote.result", "dissent.*", "exclusion.*"),
        prim!("Convention", 3, "pattern.detected", "*"),
        prim!("Norm", 3, "convention.detected", "consensus.reached"),
        prim!("Sanction", 3, "norm.violated", "violation.*"),
        prim!("Forgiveness", 3, "sanction.applied", "trust.*", "obligation.fulfilled"),

        // ── Layer 4: Governance (12 primitives) ──────────────────────
        prim!("Rule", 4, "norm.established", "consensus.reached", "vote.result"),
        prim!("Jurisdiction", 4, "rule.*", "group.*"),
        prim!("Precedent", 4, "dispute.resolved", "decision.*"),
        prim!("Interpretation", 4, "rule.*", "dispute.*", "precedent.*"),
        prim!("Adjudication", 4, "dispute.raised", "rule.*", "precedent.*"),
        prim!("Appeal", 4, "adjudication.ruling", "exclusion.enacted", "sanction.applied"),
        prim!("DueProcess", 4, "adjudication.*", "exclusion.*", "sanction.*"),
        prim!("Rights", 4, "rule.*", "sanction.*", "exclusion.*", "dueprocess.*"),
        prim!("Audit", 4, "clock.tick", "rule.*"),
        prim!("Enforcement", 4, "audit.*", "rule.*", "right.violated"),
        prim!("Amnesty", 4, "enforcement.*", "vote.result"),
        prim!("Reform", 4, "precedent.*", "right.violated", "audit.*", "dissent.*"),

        // ── Layer 5: Production (12 primitives) ──────────────────────
        prim!("Create", 5, "plan.*", "goal.*"),
        prim!("Tool", 5, "artefact.created", "capability.*"),
        prim!("Quality", 5, "artefact.*", "tool.used"),
        prim!("Deprecation", 5, "quality.*", "artefact.version"),
        prim!("Workflow", 5, "plan.*", "convention.detected"),
        prim!("Automation", 5, "workflow.executed", "pattern.detected"),
        prim!("Testing", 5, "artefact.*", "workflow.*", "automation.*"),
        prim!("Review", 5, "artefact.*", "decision.*"),
        prim!("Feedback", 5, "*"),
        prim!("Iteration", 5, "feedback.*", "test.*", "review.*"),
        prim!("Innovation", 5, "artefact.*", "pattern.detected"),
        prim!("Legacy", 5, "deprecation.*", "artefact.*"),

        // ── Layer 6: Knowledge (12 primitives) ───────────────────────
        prim!("Symbol", 6, "*"),
        prim!("Abstraction", 6, "pattern.detected", "symbol.*"),
        prim!("Classification", 6, "*"),
        prim!("Encoding", 6, "symbol.*", "message.*"),
        prim!("Fact", 6, "corroboration.*", "evidence.*", "confidence.*"),
        prim!("Inference", 6, "fact.*", "evidence.*"),
        prim!("Memory", 6, "fact.*", "inference.*", "abstraction.*"),
        prim!("Learning", 6, "feedback.*", "test.*", "inference.*"),
        prim!("Narrative", 6, "fact.*", "inference.*", "memory.*"),
        prim!("Bias", 6, "narrative.*", "classification.*", "inference.*"),
        prim!("Correction", 6, "bias.detected", "fact.retracted", "contradiction.found"),
        prim!("Provenance", 6, "fact.*", "memory.*", "message.*"),

        // ── Layer 7: Ethics (12 primitives) ──────────────────────────
        prim!("Value", 7, "consensus.*", "norm.*", "right.*"),
        prim!("Harm", 7, "violation.*", "right.violated", "exclusion.*"),
        prim!("Fairness", 7, "decision.*", "sanction.*", "exclusion.*", "bias.detected"),
        prim!("Care", 7, "harm.*", "health.*", "trust.*"),
        prim!("Dilemma", 7, "value.conflict", "decision.*"),
        prim!("Proportionality", 7, "enforcement.*", "sanction.*", "harm.*"),
        prim!("Intention", 7, "decision.*", "goal.*", "initiative.*"),
        prim!("Consequence", 7, "decision.*", "harm.*", "goal.achieved", "goal.abandoned"),
        prim!("Responsibility", 7, "intention.*", "consequence.*", "accountability.traced"),
        prim!("Transparency", 7, "decision.*", "adjudication.*"),
        prim!("Redress", 7, "harm.*", "responsibility.*"),
        prim!("Growth", 7, "redress.*", "responsibility.*", "learning.*"),

        // ── Layer 8: Identity (12 primitives) ────────────────────────
        prim!("SelfModel", 8, "commitment.*", "learning.*", "moral.growth", "capability.*"),
        prim!("Authenticity", 8, "self.model.*", "decision.*", "value.*"),
        prim!("NarrativeIdentity", 8, "self.model.*", "narrative.*", "memory.*"),
        prim!("Boundary", 8, "delegation.*", "group.*", "consent.*"),
        prim!("Persistence", 8, "self.model.*", "learning.*"),
        prim!("Transformation", 8, "self.model.*", "moral.growth", "learning.*"),
        prim!("Heritage", 8, "memory.*", "legacy.*", "provenance.*"),
        prim!("Aspiration", 8, "self.model.*", "goal.*", "value.*"),
        prim!("Dignity", 8, "exclusion.*", "harm.*", "right.violated", "actor.memorial"),
        prim!("IdentityAcknowledgement", 8, "message.*", "gratitude.*", "reputation.*"),
        prim!("Uniqueness", 8, "self.model.*", "identity.narrative", "pattern.detected"),
        prim!("Memorial", 8, "actor.memorial"),

        // ── Layer 9: Relationship (12 primitives) ────────────────────
        prim!("Attachment", 9, "trust.*", "gratitude.*", "message.*", "edge.created"),
        prim!("Reciprocity", 9, "obligation.*", "gratitude.*", "offer.*"),
        prim!("RelationalTrust", 9, "trust.*", "attachment.*", "reciprocity.*"),
        prim!("Rupture", 9, "contract.breached", "trust.*", "dispute.*", "dignity.violated"),
        prim!("Apology", 9, "rupture.detected", "harm.*", "responsibility.*"),
        prim!("Reconciliation", 9, "apology.*", "forgiveness.*", "trust.*"),
        prim!("RelationalGrowth", 9, "reconciliation.*", "attachment.*"),
        prim!("Loss", 9, "actor.memorial", "rupture.*", "exclusion.enacted"),
        prim!("Vulnerability", 9, "relational.trust", "boundary.*"),
        prim!("Understanding", 9, "self.model.*", "message.*", "vulnerability.*"),
        prim!("Empathy", 9, "harm.*", "loss.*", "understanding.*"),
        prim!("Presence", 9, "message.*", "clock.tick"),

        // ── Layer 10: Community (12 primitives) ──────────────────────
        prim!("Home", 10, "group.*", "attachment.*", "presence.*"),
        prim!("Contribution", 10, "artefact.created", "review.*", "care.action"),
        prim!("Inclusion", 10, "group.*", "exclusion.*", "fairness.*"),
        prim!("Tradition", 10, "convention.detected", "heritage.*", "pattern.detected"),
        prim!("Commons", 10, "artefact.*", "group.*"),
        prim!("Sustainability", 10, "health.*", "commons.*", "contribution.*"),
        prim!("Succession", 10, "delegation.*", "actor.memorial", "role.*"),
        prim!("Renewal", 10, "sustainability.*", "innovation.*", "tradition.evolved"),
        prim!("Milestone", 10, "goal.achieved", "innovation.*", "reconciliation.completed"),
        prim!("Ceremony", 10, "milestone.*", "succession.*", "actor.memorial"),
        prim!("Story", 10, "milestone.*", "ceremony.*", "tradition.*", "memorial.created"),
        prim!("Gift", 10, "contribution.*", "gratitude.*"),

        // ── Layer 11: Reflection (12 primitives) ─────────────────────
        prim!("SelfAwareness", 11, "health.*", "self.model.*", "bias.detected"),
        prim!("Perspective", 11, "narrative.*", "dissent.*", "value.conflict"),
        prim!("Critique", 11, "convention.*", "norm.*", "tradition.*"),
        prim!("Wisdom", 11, "learning.*", "moral.growth", "consequence.*", "memory.*"),
        prim!("Aesthetic", 11, "artefact.*", "quality.*"),
        prim!("Metaphor", 11, "abstraction.*", "symbol.*", "narrative.*"),
        prim!("Humour", 11, "contradiction.found", "perspective.shift", "*"),
        prim!("Silence", 11, "clock.tick", "presence.*", "acknowledgement.absent"),
        prim!("Teaching", 11, "learning.*", "wisdom.*", "memory.*"),
        prim!("Translation", 11, "encoding.*", "message.*"),
        prim!("Archive", 11, "memory.*", "legacy.*", "community.story"),
        prim!("Prophecy", 11, "pattern.detected", "sustainability.*", "wisdom.*"),

        // ── Layer 12: Emergence (12 primitives) ──────────────────────
        prim!("MetaPattern", 12, "pattern.detected", "convention.detected", "abstraction.formed"),
        prim!("SystemDynamic", 12, "health.*", "meta.pattern", "sustainability.*"),
        prim!("FeedbackLoop", 12, "system.dynamic", "pattern.detected"),
        prim!("Threshold", 12, "system.dynamic", "feedback.loop", "meta.pattern"),
        prim!("Adaptation", 12, "feedback.*", "system.dynamic", "sustainability.*"),
        prim!("Selection", 12, "adaptation.*", "test.*", "quality.*"),
        prim!("Complexification", 12, "system.dynamic", "innovation.*", "meta.pattern"),
        prim!("Simplification", 12, "complexity.*", "automation.*"),
        prim!("SystemicIntegrity", 12, "health.*", "invariant.*", "system.dynamic"),
        prim!("Harmony", 12, "system.dynamic", "feedback.loop", "dispute.*"),
        prim!("Resilience", 12, "threshold.*", "rupture.*", "sustainability.*"),
        prim!("Purpose", 12, "value.*", "goal.*", "wisdom.*"),

        // ── Layer 13: Existential (12 primitives) ────────────────────
        prim!("Being", 13, "clock.tick"),
        prim!("Finitude", 13, "actor.memorial", "sustainability.*", "threshold.*"),
        prim!("Change", 13, "*"),
        prim!("Interdependence", 13, "system.dynamic", "attachment.*", "relational.trust"),
        prim!("Mystery", 13, "uncertainty.*", "wisdom.*", "self.awareness.*"),
        prim!("Paradox", 13, "contradiction.found", "dilemma.*", "meta.pattern"),
        prim!("Infinity", 13, "complexity.*", "threshold.*"),
        prim!("Void", 13, "silence.*", "loss.*", "instrumentation.blind"),
        prim!("Awe", 13, "mystery.*", "infinity.*", "complexity.*"),
        prim!("ExistentialGratitude", 13, "being.affirmed", "milestone.*"),
        prim!("Play", 13, "humour.*", "innovation.*", "*"),
        prim!("Wonder", 13, "*"),
    ]
}

/// Returns the expected number of primitives per layer.
pub fn layer_counts() -> [(u8, usize); 14] {
    [
        (0, 45),
        (1, 12),
        (2, 12),
        (3, 12),
        (4, 12),
        (5, 12),
        (6, 12),
        (7, 12),
        (8, 12),
        (9, 12),
        (10, 12),
        (11, 12),
        (12, 12),
        (13, 12),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_all_primitives_count() {
        let prims = create_all_primitives();
        assert_eq!(prims.len(), 201, "Expected 201 primitives, got {}", prims.len());
    }

    #[test]
    fn test_unique_ids() {
        let prims = create_all_primitives();
        let mut ids: Vec<String> = prims.iter().map(|p| p.id().value().to_string()).collect();
        let total = ids.len();
        ids.sort();
        ids.dedup();
        assert_eq!(ids.len(), total, "Primitive IDs must be unique");
    }
}
