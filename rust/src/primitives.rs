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
        // ── Layer 0: Foundation (45 primitives) ──────────────────────
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
        // Volition
        prim!("Value", 1, "decision.*", "actor.*"),
        prim!("Intent", 1, "value.*", "expectation.*"),
        prim!("Choice", 1, "intent.*", "value.*", "confidence.*"),
        prim!("Risk", 1, "intent.*", "uncertainty.*", "value.*"),
        // Action
        prim!("Act", 1, "choice.*", "intent.*"),
        prim!("Consequence", 1, "act.*", "violation.*"),
        prim!("Capacity", 1, "actor.*", "resource.*", "trust.*"),
        prim!("Resource", 1, "act.*", "budget.*"),
        // Communication
        prim!("Signal", 1, "act.*", "actor.*"),
        prim!("Reception", 1, "*"),
        prim!("Acknowledgment", 1, "signal.*"),
        prim!("Commitment", 1, "signal.*", "agreement.*", "intent.*"),

        // ── Layer 2: Exchange (12 primitives) ────────────────────────
        prim!("Term", 2, "signal.*", "commitment.*"),
        prim!("Protocol", 2, "term.*", "norm.*"),
        prim!("Offer", 2, "signal.*", "term.*", "value.*"),
        prim!("Acceptance", 2, "offer.*", "acknowledgment.*"),
        prim!("Agreement", 2, "offer.*", "acceptance.*", "commitment.*"),
        prim!("Obligation", 2, "agreement.*", "commitment.*"),
        prim!("Fulfillment", 2, "obligation.*", "act.*"),
        prim!("Breach", 2, "obligation.*", "violation.*"),
        prim!("Exchange", 2, "fulfillment.*", "reciprocity.*"),
        prim!("Accountability", 2, "obligation.*", "consequence.*"),
        prim!("Debt", 2, "obligation.*", "breach.*", "exchange.*"),
        prim!("Reciprocity", 2, "exchange.*", "obligation.*", "fulfillment.*"),

        // ── Layer 3: Society (12 primitives) ─────────────────────────
        prim!("Group", 3, "actor.*", "membership.*"),
        prim!("Membership", 3, "group.*", "actor.*", "consent.*"),
        prim!("Role", 3, "group.*", "membership.*", "capacity.*"),
        prim!("Consent", 3, "signal.*", "acknowledgment.*", "authority.*"),
        prim!("Norm", 3, "group.*", "agreement.*", "pattern.*"),
        prim!("Reputation", 3, "trust.*", "fulfillment.*", "violation.*"),
        prim!("Sanction", 3, "norm.*", "violation.*", "authority.*"),
        prim!("Authority", 3, "role.*", "trust.*", "consent.*"),
        prim!("Property", 3, "resource.*", "agreement.*", "actor.*"),
        prim!("Commons", 3, "property.*", "group.*", "norm.*"),
        prim!("Governance", 3, "authority.*", "norm.*", "group.*"),
        prim!("CollectiveAct", 3, "group.*", "consent.*", "act.*"),

        // ── Layer 4: Legal (12 primitives) ───────────────────────────
        prim!("Law", 4, "norm.*", "governance.*", "authority.*"),
        prim!("Right", 4, "law.*", "actor.*", "dignity.*"),
        prim!("Contract", 4, "agreement.*", "obligation.*", "law.*"),
        prim!("Liability", 4, "breach.*", "consequence.*", "law.*"),
        prim!("DueProcess", 4, "adjudication.*", "right.*", "law.*"),
        prim!("Adjudication", 4, "dispute.*", "law.*", "precedent.*"),
        prim!("Remedy", 4, "adjudication.*", "liability.*", "right.*"),
        prim!("Precedent", 4, "adjudication.*", "law.*", "decision.*"),
        prim!("Jurisdiction", 4, "law.*", "group.*", "authority.*"),
        prim!("Sovereignty", 4, "jurisdiction.*", "authority.*", "governance.*"),
        prim!("Legitimacy", 4, "consent.*", "authority.*", "law.*"),
        prim!("Treaty", 4, "agreement.*", "sovereignty.*", "obligation.*"),

        // ── Layer 5: Technology (12 primitives) ──────────────────────
        prim!("Method", 5, "act.*", "knowledge.*", "pattern.*"),
        prim!("Measurement", 5, "method.*", "evidence.*"),
        prim!("Knowledge", 5, "measurement.*", "evidence.*", "corroboration.*"),
        prim!("Model", 5, "knowledge.*", "abstraction.*", "pattern.*"),
        prim!("Tool", 5, "method.*", "capacity.*", "resource.*"),
        prim!("Technique", 5, "method.*", "tool.*", "knowledge.*"),
        prim!("Invention", 5, "technique.*", "knowledge.*", "pattern.*"),
        prim!("Abstraction", 5, "pattern.*", "model.*", "knowledge.*"),
        prim!("Infrastructure", 5, "tool.*", "resource.*", "standard.*"),
        prim!("Standard", 5, "measurement.*", "norm.*", "protocol.*"),
        prim!("Efficiency", 5, "method.*", "measurement.*", "resource.*"),
        prim!("Automation", 5, "technique.*", "tool.*", "pattern.*"),

        // ── Layer 6: Information (12 primitives) ─────────────────────
        prim!("Symbol", 6, "signal.*", "encoding.*"),
        prim!("Language", 6, "symbol.*", "protocol.*", "norm.*"),
        prim!("Encoding", 6, "symbol.*", "channel.*"),
        prim!("Record", 6, "encoding.*", "data.*", "evidence.*"),
        prim!("Channel", 6, "signal.*", "infrastructure.*"),
        prim!("Copy", 6, "record.*", "encoding.*"),
        prim!("Noise", 6, "channel.*", "signal.*"),
        prim!("Redundancy", 6, "copy.*", "noise.*", "record.*"),
        prim!("Data", 6, "record.*", "measurement.*", "encoding.*"),
        prim!("Computation", 6, "data.*", "algorithm.*", "tool.*"),
        prim!("Algorithm", 6, "method.*", "abstraction.*", "computation.*"),
        prim!("Entropy", 6, "noise.*", "data.*", "redundancy.*"),

        // ── Layer 7: Ethics (12 primitives) ──────────────────────────
        prim!("MoralStatus", 7, "actor.*", "dignity.*", "right.*"),
        prim!("Dignity", 7, "moral.status.*", "right.*", "autonomy.*"),
        prim!("Autonomy", 7, "choice.*", "capacity.*", "consent.*"),
        prim!("Flourishing", 7, "autonomy.*", "care.*", "value.*"),
        prim!("Duty", 7, "obligation.*", "moral.status.*", "norm.*"),
        prim!("Harm", 7, "violation.*", "consequence.*", "dignity.*"),
        prim!("Care", 7, "harm.*", "flourishing.*", "trust.*"),
        prim!("Justice", 7, "right.*", "fairness.*", "remedy.*"),
        prim!("Conscience", 7, "duty.*", "harm.*", "value.*"),
        prim!("Virtue", 7, "duty.*", "flourishing.*", "reputation.*"),
        prim!("Responsibility", 7, "consequence.*", "duty.*", "accountability.*"),
        prim!("Motive", 7, "intent.*", "value.*", "conscience.*"),

        // ── Layer 8: Identity (12 primitives) ────────────────────────
        prim!("Narrative", 8, "memory.*", "self.concept.*", "expression.*"),
        prim!("SelfConcept", 8, "reflection.*", "memory.*", "value.*"),
        prim!("Reflection", 8, "self.concept.*", "conscience.*", "knowledge.*"),
        prim!("Memory", 8, "record.*", "narrative.*", "continuity.*"),
        prim!("Purpose", 8, "aspiration.*", "value.*", "narrative.*"),
        prim!("Aspiration", 8, "self.concept.*", "value.*", "intent.*"),
        prim!("Authenticity", 8, "self.concept.*", "expression.*", "value.*"),
        prim!("Expression", 8, "signal.*", "self.concept.*", "narrative.*"),
        prim!("Growth", 8, "reflection.*", "aspiration.*", "learning.*"),
        prim!("Continuity", 8, "memory.*", "self.concept.*", "narrative.*"),
        prim!("Integration", 8, "self.concept.*", "narrative.*", "growth.*"),
        prim!("Crisis", 8, "self.concept.*", "rupture.*", "continuity.*"),

        // ── Layer 9: Relationship (12 primitives) ────────────────────
        prim!("Bond", 9, "trust.*", "attachment.*", "commitment.*"),
        prim!("Attachment", 9, "bond.*", "signal.*", "care.*"),
        prim!("Recognition", 9, "signal.*", "dignity.*", "identity.*"),
        prim!("Intimacy", 9, "bond.*", "trust.*", "vulnerability.*"),
        prim!("Attunement", 9, "signal.*", "reception.*", "care.*"),
        prim!("Rupture", 9, "bond.*", "breach.*", "harm.*"),
        prim!("Repair", 9, "rupture.*", "care.*", "trust.*"),
        prim!("Loyalty", 9, "bond.*", "commitment.*", "trust.*"),
        prim!("MutualConstitution", 9, "bond.*", "self.concept.*", "recognition.*"),
        prim!("RelationalObligation", 9, "bond.*", "duty.*", "care.*"),
        prim!("Grief", 9, "bond.*", "rupture.*", "loss.*"),
        prim!("Forgiveness", 9, "rupture.*", "repair.*", "trust.*"),

        // ── Layer 10: Community (12 primitives) ──────────────────────
        prim!("Culture", 10, "group.*", "norm.*", "narrative.*"),
        prim!("SharedNarrative", 10, "narrative.*", "group.*", "culture.*"),
        prim!("Ethos", 10, "culture.*", "value.*", "norm.*"),
        prim!("Sacred", 10, "ethos.*", "culture.*", "value.*"),
        prim!("Tradition", 10, "culture.*", "practice.*", "memory.*"),
        prim!("Ritual", 10, "tradition.*", "sacred.*", "practice.*"),
        prim!("Practice", 10, "norm.*", "method.*", "culture.*"),
        prim!("Place", 10, "group.*", "belonging.*", "culture.*"),
        prim!("Belonging", 10, "group.*", "bond.*", "culture.*"),
        prim!("Solidarity", 10, "belonging.*", "collective.act.*", "care.*"),
        prim!("Voice", 10, "expression.*", "belonging.*", "signal.*"),
        prim!("Welcome", 10, "belonging.*", "recognition.*", "care.*"),

        // ── Layer 11: Culture (12 primitives) ────────────────────────
        prim!("Reflexivity", 11, "reflection.*", "culture.*", "self.concept.*"),
        prim!("Encounter", 11, "recognition.*", "signal.*", "culture.*"),
        prim!("Translation", 11, "language.*", "encoding.*", "encounter.*"),
        prim!("Pluralism", 11, "culture.*", "encounter.*", "norm.*"),
        prim!("Creativity", 11, "expression.*", "invention.*", "aesthetic.*"),
        prim!("Aesthetic", 11, "expression.*", "value.*", "culture.*"),
        prim!("Interpretation", 11, "symbol.*", "narrative.*", "culture.*"),
        prim!("Dialogue", 11, "signal.*", "encounter.*", "recognition.*"),
        prim!("Syncretism", 11, "encounter.*", "translation.*", "culture.*"),
        prim!("Critique", 11, "reflexivity.*", "norm.*", "culture.*"),
        prim!("Hegemony", 11, "authority.*", "culture.*", "norm.*"),
        prim!("CulturalEvolution", 11, "culture.*", "creativity.*", "tradition.*"),

        // ── Layer 12: Emergence (12 primitives) ──────────────────────
        prim!("Emergence", 12, "pattern.*", "complexity.*", "self.organization.*"),
        prim!("SelfOrganization", 12, "feedback.*", "pattern.*", "autopoiesis.*"),
        prim!("Feedback", 12, "consequence.*", "signal.*", "act.*"),
        prim!("Complexity", 12, "emergence.*", "pattern.*", "self.organization.*"),
        prim!("Consciousness", 12, "self.concept.*", "reflection.*", "recursion.*"),
        prim!("Recursion", 12, "feedback.*", "self.organization.*", "pattern.*"),
        prim!("Paradox", 12, "contradiction.*", "recursion.*", "incompleteness.*"),
        prim!("Incompleteness", 12, "knowledge.*", "recursion.*", "complexity.*"),
        prim!("PhaseTransition", 12, "emergence.*", "complexity.*", "feedback.*"),
        prim!("DownwardCausation", 12, "emergence.*", "complexity.*", "pattern.*"),
        prim!("Autopoiesis", 12, "self.organization.*", "feedback.*", "emergence.*"),
        prim!("CoEvolution", 12, "feedback.*", "emergence.*", "complexity.*"),

        // ── Layer 13: Existence (12 primitives) ──────────────────────
        prim!("Being", 13, "clock.tick", "presence.*"),
        prim!("Nothingness", 13, "being.*", "absence.*"),
        prim!("Finitude", 13, "being.*", "actor.memorial", "continuity.*"),
        prim!("Contingency", 13, "being.*", "uncertainty.*", "risk.*"),
        prim!("Wonder", 13, "mystery.*", "being.*", "emergence.*"),
        prim!("ExistentialAcceptance", 13, "finitude.*", "contingency.*", "grief.*"),
        prim!("Presence", 13, "being.*", "clock.tick", "attunement.*"),
        prim!("Gratitude", 13, "presence.*", "being.*", "care.*"),
        prim!("Mystery", 13, "uncertainty.*", "incompleteness.*", "wonder.*"),
        prim!("Transcendence", 13, "being.*", "emergence.*", "consciousness.*"),
        prim!("Groundlessness", 13, "contingency.*", "nothingness.*", "paradox.*"),
        prim!("Return", 13, "being.*", "continuity.*", "presence.*"),
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
