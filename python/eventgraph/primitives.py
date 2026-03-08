"""All 201 EventGraph primitives across 14 layers (0-13).

Each primitive is a concrete class implementing the Primitive protocol defined
in ``primitive.py``.  A ``create_all_primitives()`` factory returns the full
set ready for registration.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any

from .event import Event
from .primitive import (
    Mutation,
    Snapshot,
    UpdateState,
)
from .types import (
    Cadence,
    Layer,
    PrimitiveID,
    SubscriptionPattern,
)


# ---------------------------------------------------------------------------
# Base class — shared plumbing for every primitive
# ---------------------------------------------------------------------------

class _Base:
    """Common implementation of the Primitive protocol.

    Subclasses only need to set ``_name``, ``_layer``, and ``_subs`` at the
    class level (done by the ``_def`` helper below).  ``process`` tracks
    events processed and last tick via state mutations.
    """

    _name: str
    _layer: int
    _subs: list[str]

    def id(self) -> PrimitiveID:
        return PrimitiveID(self._name)

    def layer(self) -> Layer:
        return Layer(self._layer)

    def subscriptions(self) -> list[SubscriptionPattern]:
        return [SubscriptionPattern(s) for s in self._subs]

    def cadence(self) -> Cadence:
        return Cadence(1)

    def process(
        self,
        tick: int,
        events: list[Event],
        snapshot: Snapshot,
    ) -> list[Mutation]:
        pid = self.id()
        mutations: list[Mutation] = [
            UpdateState(primitive_id=pid, key="eventsProcessed", value=len(events)),
            UpdateState(primitive_id=pid, key="lastTick", value=tick),
        ]
        return mutations


# ---------------------------------------------------------------------------
# Factory helper — creates a concrete subclass of _Base
# ---------------------------------------------------------------------------

def _def(name: str, layer: int, subs: list[str]) -> type:
    """Return a new class named *name* inheriting from ``_Base``."""
    cls = type(name, (_Base,), {
        "_name": name,
        "_layer": layer,
        "_subs": subs,
    })
    cls.__module__ = __name__
    cls.__qualname__ = name
    return cls


# ===================================================================
# Layer 0 — Foundation (44 primitives)
# ===================================================================

EventPrimitive           = _def("Event",               0, ["*"])
EventStorePrimitive      = _def("EventStore",           0, ["store.*"])
ClockPrimitive           = _def("Clock",                0, ["clock.*"])
HashPrimitive            = _def("Hash",                 0, ["*"])
SelfPrimitive            = _def("Self",                 0, ["system.*"])
CausalLinkPrimitive      = _def("CausalLink",          0, ["*"])
AncestryPrimitive        = _def("Ancestry",             0, ["query.*"])
DescendancyPrimitive     = _def("Descendancy",          0, ["query.*"])
FirstCausePrimitive      = _def("FirstCause",           0, ["query.*"])
ActorIDPrimitive         = _def("ActorID",              0, ["actor.*"])
ActorRegistryPrimitive   = _def("ActorRegistry",        0, ["actor.*"])
SignaturePrimitive       = _def("Signature",            0, ["*"])
VerifyPrimitive          = _def("Verify",               0, ["*"])
ExpectationPrimitive     = _def("Expectation",          0, ["*"])
TimeoutPrimitive         = _def("Timeout",              0, ["authority.*"])
ViolationPrimitive       = _def("Violation",            0, ["violation.*"])
SeverityPrimitive        = _def("Severity",             0, ["violation.*"])
TrustScorePrimitive      = _def("TrustScore",           0, ["trust.*"])
TrustUpdatePrimitive     = _def("TrustUpdate",          0, ["trust.*"])
CorroborationPrimitive   = _def("Corroboration",        0, ["trust.*"])
ContradictionPrimitive   = _def("Contradiction",        0, ["trust.*"])
ConfidencePrimitive      = _def("Confidence",           0, ["decision.*"])
EvidencePrimitive        = _def("Evidence",             0, ["*"])
RevisionPrimitive        = _def("Revision",             0, ["grammar.*"])
UncertaintyPrimitive     = _def("Uncertainty",          0, ["decision.*"])
InstrumentationSpecPrimitive = _def("InstrumentationSpec", 0, ["health.*"])
CoverageCheckPrimitive   = _def("CoverageCheck",        0, ["health.*"])
GapPrimitive             = _def("Gap",                  0, ["*"])
BlindPrimitive           = _def("Blind",                0, ["health.*"])
PathQueryPrimitive       = _def("PathQuery",            0, ["query.*"])
SubgraphExtractPrimitive = _def("SubgraphExtract",      0, ["query.*"])
AnnotatePrimitive        = _def("Annotate",             0, ["grammar.*"])
TimelinePrimitive        = _def("Timeline",             0, ["query.*"])
HashChainPrimitive       = _def("HashChain",            0, ["*"])
ChainVerifyPrimitive     = _def("ChainVerify",          0, ["chain.*"])
WitnessPrimitive         = _def("Witness",              0, ["*"])
IntegrityViolationPrimitive = _def("IntegrityViolation", 0, ["chain.*"])
PatternPrimitive         = _def("Pattern",              0, ["*"])
DeceptionIndicatorPrimitive = _def("DeceptionIndicator", 0, ["trust.*", "violation.*"])
SuspicionPrimitive       = _def("Suspicion",            0, ["trust.*"])
QuarantinePrimitive      = _def("Quarantine",           0, ["actor.*"])
GraphHealthPrimitive     = _def("GraphHealth",          0, ["health.*"])
InvariantPrimitive       = _def("Invariant",            0, ["*"])
InvariantCheckPrimitive  = _def("InvariantCheck",       0, ["chain.*"])
BootstrapPrimitive       = _def("Bootstrap",            0, ["system.*"])

# ===================================================================
# Layer 1 — Agency (12 primitives)
# ===================================================================

GoalPrimitive            = _def("Goal",                 1, ["decision.*", "authority.resolved", "actor.*"])
PlanPrimitive            = _def("Plan",                 1, ["goal.*"])
InitiativePrimitive      = _def("Initiative",           1, ["clock.tick", "goal.*", "plan.*"])
CommitmentPrimitive      = _def("Commitment",           1, ["goal.set", "goal.achieved", "goal.abandoned", "plan.step.completed"])
FocusPrimitive           = _def("Focus",                1, ["*"])
FilterPrimitive          = _def("Filter",               1, ["*"])
SaliencePrimitive        = _def("Salience",             1, ["*"])
DistractionPrimitive     = _def("Distraction",          1, ["focus.*", "goal.*"])
PermissionPrimitive      = _def("Permission",           1, ["authority.*", "decision.*"])
CapabilityPrimitive      = _def("Capability",           1, ["actor.registered", "permission.*", "trust.*"])
DelegationPrimitive      = _def("Delegation",           1, ["authority.*", "edge.created"])
AccountabilityPrimitive  = _def("Accountability",       1, ["delegation.*", "violation.*", "goal.abandoned"])

# ===================================================================
# Layer 2 — Communication (11 primitives)
# ===================================================================

MessagePrimitive         = _def("Message",              2, ["protocol.message.*"])
AcknowledgementPrimitive = _def("Acknowledgement",      2, ["message.sent", "message.received"])
ClarificationPrimitive   = _def("Clarification",       2, ["message.*", "ack.*"])
ContextPrimitive         = _def("Context",              2, ["message.*", "clarification.*"])
OfferPrimitive           = _def("Offer",                2, ["message.*", "exchange.*"])
AcceptancePrimitive      = _def("Acceptance",           2, ["offer.made", "offer.withdrawn"])
ObligationPrimitive      = _def("Obligation",           2, ["offer.accepted", "delegation.*"])
GratitudePrimitive       = _def("Gratitude",            2, ["obligation.fulfilled", "trust.*"])
NegotiationPrimitive     = _def("Negotiation",          2, ["offer.*", "message.*"])
ConsentPrimitive         = _def("Consent",              2, ["negotiation.concluded", "offer.accepted", "authority.*"])
ContractPrimitive        = _def("Contract",             2, ["consent.given", "negotiation.concluded"])
DisputePrimitive         = _def("Dispute",              2, ["contract.breached", "obligation.defaulted", "contradiction.found"])

# ===================================================================
# Layer 3 — Social (12 primitives)
# ===================================================================

GroupPrimitive           = _def("Group",                3, ["actor.*", "consent.*"])
RolePrimitive            = _def("Role",                 3, ["group.*", "delegation.*"])
ReputationPrimitive      = _def("Reputation",           3, ["trust.*", "commitment.*", "violation.*", "gratitude.*"])
ExclusionPrimitive       = _def("Exclusion",            3, ["reputation.*", "violation.*", "quarantine.*", "dispute.*"])
VotePrimitive            = _def("Vote",                 3, ["authority.requested", "group.*"])
ConsensusPrimitive       = _def("Consensus",            3, ["message.*", "corroboration.*", "vote.result"])
DissentPrimitive         = _def("Dissent",              3, ["vote.*", "consensus.*", "contradiction.found"])
MajorityPrimitive        = _def("Majority",             3, ["vote.result", "dissent.*", "exclusion.*"])
ConventionPrimitive      = _def("Convention",           3, ["pattern.detected", "*"])
NormPrimitive            = _def("Norm",                 3, ["convention.detected", "consensus.reached"])
SanctionPrimitive        = _def("Sanction",             3, ["norm.violated", "violation.*"])
ForgivenessPrimitive     = _def("Forgiveness",          3, ["sanction.applied", "trust.*", "obligation.fulfilled"])

# ===================================================================
# Layer 4 — Governance (12 primitives)
# ===================================================================

RulePrimitive            = _def("Rule",                 4, ["norm.established", "consensus.reached", "vote.result"])
JurisdictionPrimitive    = _def("Jurisdiction",         4, ["rule.*", "group.*"])
PrecedentPrimitive       = _def("Precedent",            4, ["dispute.resolved", "decision.*"])
InterpretationPrimitive  = _def("Interpretation",       4, ["rule.*", "dispute.*", "precedent.*"])
AdjudicationPrimitive    = _def("Adjudication",         4, ["dispute.raised", "rule.*", "precedent.*"])
AppealPrimitive          = _def("Appeal",               4, ["adjudication.ruling", "exclusion.enacted", "sanction.applied"])
DueProcessPrimitive      = _def("DueProcess",           4, ["adjudication.*", "exclusion.*", "sanction.*"])
RightsPrimitive          = _def("Rights",               4, ["rule.*", "sanction.*", "exclusion.*", "dueprocess.*"])
AuditPrimitive           = _def("Audit",                4, ["clock.tick", "rule.*"])
EnforcementPrimitive     = _def("Enforcement",          4, ["audit.*", "rule.*", "right.violated"])
AmnestyPrimitive         = _def("Amnesty",              4, ["enforcement.*", "vote.result"])
ReformPrimitive          = _def("Reform",               4, ["precedent.*", "right.violated", "audit.*", "dissent.*"])

# ===================================================================
# Layer 5 — Production (12 primitives)
# ===================================================================

CreatePrimitive          = _def("Create",               5, ["plan.*", "goal.*"])
ToolPrimitive            = _def("Tool",                 5, ["artefact.created", "capability.*"])
QualityPrimitive         = _def("Quality",              5, ["artefact.*", "tool.used"])
DeprecationPrimitive     = _def("Deprecation",          5, ["quality.*", "artefact.version"])
WorkflowPrimitive        = _def("Workflow",             5, ["plan.*", "convention.detected"])
AutomationPrimitive      = _def("Automation",           5, ["workflow.executed", "pattern.detected"])
TestingPrimitive         = _def("Testing",              5, ["artefact.*", "workflow.*", "automation.*"])
ReviewPrimitive          = _def("Review",               5, ["artefact.*", "decision.*"])
FeedbackPrimitive        = _def("Feedback",             5, ["*"])
IterationPrimitive       = _def("Iteration",            5, ["feedback.*", "test.*", "review.*"])
InnovationPrimitive      = _def("Innovation",           5, ["artefact.*", "pattern.detected"])
LegacyPrimitive          = _def("Legacy",               5, ["deprecation.*", "artefact.*"])

# ===================================================================
# Layer 6 — Knowledge (12 primitives)
# ===================================================================

SymbolPrimitive          = _def("Symbol",               6, ["*"])
AbstractionPrimitive     = _def("Abstraction",          6, ["pattern.detected", "symbol.*"])
ClassificationPrimitive  = _def("Classification",       6, ["*"])
EncodingPrimitive        = _def("Encoding",             6, ["symbol.*", "message.*"])
FactPrimitive            = _def("Fact",                 6, ["corroboration.*", "evidence.*", "confidence.*"])
InferencePrimitive       = _def("Inference",            6, ["fact.*", "evidence.*"])
MemoryPrimitive          = _def("Memory",               6, ["fact.*", "inference.*", "abstraction.*"])
LearningPrimitive        = _def("Learning",             6, ["feedback.*", "test.*", "inference.*"])
NarrativePrimitive       = _def("Narrative",            6, ["fact.*", "inference.*", "memory.*"])
BiasPrimitive            = _def("Bias",                 6, ["narrative.*", "classification.*", "inference.*"])
CorrectionPrimitive      = _def("Correction",           6, ["bias.detected", "fact.retracted", "contradiction.found"])
ProvenancePrimitive      = _def("Provenance",           6, ["fact.*", "memory.*", "message.*"])

# ===================================================================
# Layer 7 — Ethics (12 primitives)
# ===================================================================

ValuePrimitive           = _def("Value",                7, ["consensus.*", "norm.*", "right.*"])
HarmPrimitive            = _def("Harm",                 7, ["violation.*", "right.violated", "exclusion.*"])
FairnessPrimitive        = _def("Fairness",             7, ["decision.*", "sanction.*", "exclusion.*", "bias.detected"])
CarePrimitive            = _def("Care",                 7, ["harm.*", "health.*", "trust.*"])
DilemmaPrimitive         = _def("Dilemma",              7, ["value.conflict", "decision.*"])
ProportionalityPrimitive = _def("Proportionality",      7, ["enforcement.*", "sanction.*", "harm.*"])
IntentionPrimitive       = _def("Intention",            7, ["decision.*", "goal.*", "initiative.*"])
ConsequencePrimitive     = _def("Consequence",          7, ["decision.*", "harm.*", "goal.achieved", "goal.abandoned"])
ResponsibilityPrimitive  = _def("Responsibility",       7, ["intention.*", "consequence.*", "accountability.traced"])
TransparencyPrimitive    = _def("Transparency",         7, ["decision.*", "adjudication.*"])
RedressPrimitive         = _def("Redress",              7, ["harm.*", "responsibility.*"])
GrowthPrimitive          = _def("Growth",               7, ["redress.*", "responsibility.*", "learning.*"])

# ===================================================================
# Layer 8 — Identity (12 primitives)
# ===================================================================

SelfModelPrimitive       = _def("SelfModel",            8, ["commitment.*", "learning.*", "moral.growth", "capability.*"])
AuthenticityPrimitive    = _def("Authenticity",         8, ["self.model.*", "decision.*", "value.*"])
NarrativeIdentityPrimitive = _def("NarrativeIdentity", 8, ["self.model.*", "narrative.*", "memory.*"])
BoundaryPrimitive        = _def("Boundary",             8, ["delegation.*", "group.*", "consent.*"])
PersistencePrimitive     = _def("Persistence",          8, ["self.model.*", "learning.*"])
TransformationPrimitive  = _def("Transformation",       8, ["self.model.*", "moral.growth", "learning.*"])
HeritagePrimitive        = _def("Heritage",             8, ["memory.*", "legacy.*", "provenance.*"])
AspirationPrimitive      = _def("Aspiration",           8, ["self.model.*", "goal.*", "value.*"])
DignityPrimitive         = _def("Dignity",              8, ["exclusion.*", "harm.*", "right.violated", "actor.memorial"])
IdentityAcknowledgementPrimitive = _def("IdentityAcknowledgement", 8, ["message.*", "gratitude.*", "reputation.*"])
UniquenessPrimitive      = _def("Uniqueness",           8, ["self.model.*", "identity.narrative", "pattern.detected"])
MemorialPrimitive        = _def("Memorial",             8, ["actor.memorial"])

# ===================================================================
# Layer 9 — Relationship (12 primitives)
# ===================================================================

AttachmentPrimitive      = _def("Attachment",           9, ["trust.*", "gratitude.*", "message.*", "edge.created"])
ReciprocityPrimitive     = _def("Reciprocity",          9, ["obligation.*", "gratitude.*", "offer.*"])
RelationalTrustPrimitive = _def("RelationalTrust",      9, ["trust.*", "attachment.*", "reciprocity.*"])
RupturePrimitive         = _def("Rupture",              9, ["contract.breached", "trust.*", "dispute.*", "dignity.violated"])
ApologyPrimitive         = _def("Apology",              9, ["rupture.detected", "harm.*", "responsibility.*"])
ReconciliationPrimitive  = _def("Reconciliation",       9, ["apology.*", "forgiveness.*", "trust.*"])
RelationalGrowthPrimitive = _def("RelationalGrowth",    9, ["reconciliation.*", "attachment.*"])
LossPrimitive            = _def("Loss",                 9, ["actor.memorial", "rupture.*", "exclusion.enacted"])
VulnerabilityPrimitive   = _def("Vulnerability",        9, ["relational.trust", "boundary.*"])
UnderstandingPrimitive   = _def("Understanding",        9, ["self.model.*", "message.*", "vulnerability.*"])
EmpathyPrimitive         = _def("Empathy",              9, ["harm.*", "loss.*", "understanding.*"])
PresencePrimitive        = _def("Presence",             9, ["message.*", "clock.tick"])

# ===================================================================
# Layer 10 — Community (12 primitives)
# ===================================================================

HomePrimitive            = _def("Home",                10, ["group.*", "attachment.*", "presence.*"])
ContributionPrimitive    = _def("Contribution",        10, ["artefact.created", "review.*", "care.action"])
InclusionPrimitive       = _def("Inclusion",           10, ["group.*", "exclusion.*", "fairness.*"])
TraditionPrimitive       = _def("Tradition",           10, ["convention.detected", "heritage.*", "pattern.detected"])
CommonsPrimitive         = _def("Commons",             10, ["artefact.*", "group.*"])
SustainabilityPrimitive  = _def("Sustainability",      10, ["health.*", "commons.*", "contribution.*"])
SuccessionPrimitive      = _def("Succession",          10, ["delegation.*", "actor.memorial", "role.*"])
RenewalPrimitive         = _def("Renewal",             10, ["sustainability.*", "innovation.*", "tradition.evolved"])
MilestonePrimitive       = _def("Milestone",           10, ["goal.achieved", "innovation.*", "reconciliation.completed"])
CeremonyPrimitive        = _def("Ceremony",            10, ["milestone.*", "succession.*", "actor.memorial"])
StoryPrimitive           = _def("Story",               10, ["milestone.*", "ceremony.*", "tradition.*", "memorial.created"])
GiftPrimitive            = _def("Gift",                10, ["contribution.*", "gratitude.*"])

# ===================================================================
# Layer 11 — Reflection (12 primitives)
# ===================================================================

SelfAwarenessPrimitive   = _def("SelfAwareness",       11, ["health.*", "self.model.*", "bias.detected"])
PerspectivePrimitive     = _def("Perspective",         11, ["narrative.*", "dissent.*", "value.conflict"])
CritiquePrimitive        = _def("Critique",            11, ["convention.*", "norm.*", "tradition.*"])
WisdomPrimitive          = _def("Wisdom",              11, ["learning.*", "moral.growth", "consequence.*", "memory.*"])
AestheticPrimitive       = _def("Aesthetic",           11, ["artefact.*", "quality.*"])
MetaphorPrimitive        = _def("Metaphor",            11, ["abstraction.*", "symbol.*", "narrative.*"])
HumourPrimitive          = _def("Humour",              11, ["contradiction.found", "perspective.shift", "*"])
SilencePrimitive         = _def("Silence",             11, ["clock.tick", "presence.*", "acknowledgement.absent"])
TeachingPrimitive        = _def("Teaching",            11, ["learning.*", "wisdom.*", "memory.*"])
TranslationPrimitive     = _def("Translation",         11, ["encoding.*", "message.*"])
ArchivePrimitive         = _def("Archive",             11, ["memory.*", "legacy.*", "community.story"])
ProphecyPrimitive        = _def("Prophecy",            11, ["pattern.detected", "sustainability.*", "wisdom.*"])

# ===================================================================
# Layer 12 — Emergence (12 primitives)
# ===================================================================

MetaPatternPrimitive     = _def("MetaPattern",         12, ["pattern.detected", "convention.detected", "abstraction.formed"])
SystemDynamicPrimitive   = _def("SystemDynamic",       12, ["health.*", "meta.pattern", "sustainability.*"])
FeedbackLoopPrimitive    = _def("FeedbackLoop",        12, ["system.dynamic", "pattern.detected"])
ThresholdPrimitive       = _def("Threshold",           12, ["system.dynamic", "feedback.loop", "meta.pattern"])
AdaptationPrimitive      = _def("Adaptation",          12, ["feedback.*", "system.dynamic", "sustainability.*"])
SelectionPrimitive       = _def("Selection",           12, ["adaptation.*", "test.*", "quality.*"])
ComplexificationPrimitive = _def("Complexification",   12, ["system.dynamic", "innovation.*", "meta.pattern"])
SimplificationPrimitive  = _def("Simplification",      12, ["complexity.*", "automation.*"])
SystemicIntegrityPrimitive = _def("SystemicIntegrity", 12, ["health.*", "invariant.*", "system.dynamic"])
HarmonyPrimitive         = _def("Harmony",             12, ["system.dynamic", "feedback.loop", "dispute.*"])
ResiliencePrimitive      = _def("Resilience",          12, ["threshold.*", "rupture.*", "sustainability.*"])
PurposePrimitive         = _def("Purpose",             12, ["value.*", "goal.*", "wisdom.*"])

# ===================================================================
# Layer 13 — Existential (12 primitives)
# ===================================================================

BeingPrimitive           = _def("Being",               13, ["clock.tick"])
FinitudePrimitive        = _def("Finitude",            13, ["actor.memorial", "sustainability.*", "threshold.*"])
ChangePrimitive          = _def("Change",              13, ["*"])
InterdependencePrimitive = _def("Interdependence",     13, ["system.dynamic", "attachment.*", "relational.trust"])
MysteryPrimitive         = _def("Mystery",             13, ["uncertainty.*", "wisdom.*", "self.awareness.*"])
ParadoxPrimitive         = _def("Paradox",             13, ["contradiction.found", "dilemma.*", "meta.pattern"])
InfinityPrimitive        = _def("Infinity",            13, ["complexity.*", "threshold.*"])
VoidPrimitive            = _def("Void",                13, ["silence.*", "loss.*", "instrumentation.blind"])
AwePrimitive             = _def("Awe",                 13, ["mystery.*", "infinity.*", "complexity.*"])
ExistentialGratitudePrimitive = _def("ExistentialGratitude", 13, ["being.affirmed", "milestone.*"])
PlayPrimitive            = _def("Play",                13, ["humour.*", "innovation.*", "*"])
WonderPrimitive          = _def("Wonder",              13, ["*"])


# ===================================================================
# Factory
# ===================================================================

# Ordered list of all 201 primitive classes for programmatic access.
ALL_PRIMITIVE_CLASSES: list[type] = [
    # Layer 0 (45)
    EventPrimitive, EventStorePrimitive, ClockPrimitive, HashPrimitive,
    SelfPrimitive, CausalLinkPrimitive, AncestryPrimitive, DescendancyPrimitive,
    FirstCausePrimitive, ActorIDPrimitive, ActorRegistryPrimitive,
    SignaturePrimitive, VerifyPrimitive, ExpectationPrimitive, TimeoutPrimitive,
    ViolationPrimitive, SeverityPrimitive, TrustScorePrimitive,
    TrustUpdatePrimitive, CorroborationPrimitive, ContradictionPrimitive,
    ConfidencePrimitive, EvidencePrimitive, RevisionPrimitive,
    UncertaintyPrimitive, InstrumentationSpecPrimitive, CoverageCheckPrimitive,
    GapPrimitive, BlindPrimitive, PathQueryPrimitive, SubgraphExtractPrimitive,
    AnnotatePrimitive, TimelinePrimitive, HashChainPrimitive,
    ChainVerifyPrimitive, WitnessPrimitive, IntegrityViolationPrimitive,
    PatternPrimitive, DeceptionIndicatorPrimitive, SuspicionPrimitive,
    QuarantinePrimitive, GraphHealthPrimitive, InvariantPrimitive,
    InvariantCheckPrimitive, BootstrapPrimitive,
    # Layer 1 (12)
    GoalPrimitive, PlanPrimitive, InitiativePrimitive, CommitmentPrimitive,
    FocusPrimitive, FilterPrimitive, SaliencePrimitive, DistractionPrimitive,
    PermissionPrimitive, CapabilityPrimitive, DelegationPrimitive,
    AccountabilityPrimitive,
    # Layer 2 (12)
    MessagePrimitive, AcknowledgementPrimitive, ClarificationPrimitive,
    ContextPrimitive, OfferPrimitive, AcceptancePrimitive, ObligationPrimitive,
    GratitudePrimitive, NegotiationPrimitive, ConsentPrimitive,
    ContractPrimitive, DisputePrimitive,
    # Layer 3 (12)
    GroupPrimitive, RolePrimitive, ReputationPrimitive, ExclusionPrimitive,
    VotePrimitive, ConsensusPrimitive, DissentPrimitive, MajorityPrimitive,
    ConventionPrimitive, NormPrimitive, SanctionPrimitive, ForgivenessPrimitive,
    # Layer 4 (12)
    RulePrimitive, JurisdictionPrimitive, PrecedentPrimitive,
    InterpretationPrimitive, AdjudicationPrimitive, AppealPrimitive,
    DueProcessPrimitive, RightsPrimitive, AuditPrimitive,
    EnforcementPrimitive, AmnestyPrimitive, ReformPrimitive,
    # Layer 5 (12)
    CreatePrimitive, ToolPrimitive, QualityPrimitive, DeprecationPrimitive,
    WorkflowPrimitive, AutomationPrimitive, TestingPrimitive, ReviewPrimitive,
    FeedbackPrimitive, IterationPrimitive, InnovationPrimitive, LegacyPrimitive,
    # Layer 6 (12)
    SymbolPrimitive, AbstractionPrimitive, ClassificationPrimitive,
    EncodingPrimitive, FactPrimitive, InferencePrimitive, MemoryPrimitive,
    LearningPrimitive, NarrativePrimitive, BiasPrimitive, CorrectionPrimitive,
    ProvenancePrimitive,
    # Layer 7 (12)
    ValuePrimitive, HarmPrimitive, FairnessPrimitive, CarePrimitive,
    DilemmaPrimitive, ProportionalityPrimitive, IntentionPrimitive,
    ConsequencePrimitive, ResponsibilityPrimitive, TransparencyPrimitive,
    RedressPrimitive, GrowthPrimitive,
    # Layer 8 (12)
    SelfModelPrimitive, AuthenticityPrimitive, NarrativeIdentityPrimitive,
    BoundaryPrimitive, PersistencePrimitive, TransformationPrimitive,
    HeritagePrimitive, AspirationPrimitive, DignityPrimitive,
    IdentityAcknowledgementPrimitive, UniquenessPrimitive, MemorialPrimitive,
    # Layer 9 (12)
    AttachmentPrimitive, ReciprocityPrimitive, RelationalTrustPrimitive,
    RupturePrimitive, ApologyPrimitive, ReconciliationPrimitive,
    RelationalGrowthPrimitive, LossPrimitive, VulnerabilityPrimitive,
    UnderstandingPrimitive, EmpathyPrimitive, PresencePrimitive,
    # Layer 10 (12)
    HomePrimitive, ContributionPrimitive, InclusionPrimitive,
    TraditionPrimitive, CommonsPrimitive, SustainabilityPrimitive,
    SuccessionPrimitive, RenewalPrimitive, MilestonePrimitive,
    CeremonyPrimitive, StoryPrimitive, GiftPrimitive,
    # Layer 11 (12)
    SelfAwarenessPrimitive, PerspectivePrimitive, CritiquePrimitive,
    WisdomPrimitive, AestheticPrimitive, MetaphorPrimitive, HumourPrimitive,
    SilencePrimitive, TeachingPrimitive, TranslationPrimitive,
    ArchivePrimitive, ProphecyPrimitive,
    # Layer 12 (12)
    MetaPatternPrimitive, SystemDynamicPrimitive, FeedbackLoopPrimitive,
    ThresholdPrimitive, AdaptationPrimitive, SelectionPrimitive,
    ComplexificationPrimitive, SimplificationPrimitive,
    SystemicIntegrityPrimitive, HarmonyPrimitive, ResiliencePrimitive,
    PurposePrimitive,
    # Layer 13 (12)
    BeingPrimitive, FinitudePrimitive, ChangePrimitive,
    InterdependencePrimitive, MysteryPrimitive, ParadoxPrimitive,
    InfinityPrimitive, VoidPrimitive, AwePrimitive,
    ExistentialGratitudePrimitive, PlayPrimitive, WonderPrimitive,
]


def create_all_primitives() -> list[_Base]:
    """Instantiate and return all 201 primitives."""
    return [cls() for cls in ALL_PRIMITIVE_CLASSES]
