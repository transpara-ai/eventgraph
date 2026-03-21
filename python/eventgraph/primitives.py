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

# Volition
ValuePrimitive           = _def("Value",                1, ["decision.*", "actor.*"])
IntentPrimitive          = _def("Intent",               1, ["value.*", "expectation.*"])
ChoicePrimitive          = _def("Choice",               1, ["intent.*", "value.*", "confidence.*"])
RiskPrimitive            = _def("Risk",                 1, ["intent.*", "uncertainty.*", "value.*"])
# Action
ActPrimitive             = _def("Act",                  1, ["choice.*", "intent.*"])
ConsequencePrimitive     = _def("Consequence",          1, ["act.*", "violation.*"])
CapacityPrimitive        = _def("Capacity",             1, ["actor.*", "resource.*", "trust.*"])
ResourcePrimitive        = _def("Resource",             1, ["act.*", "budget.*"])
# Communication
SignalPrimitive          = _def("Signal",               1, ["act.*", "actor.*"])
ReceptionPrimitive       = _def("Reception",            1, ["*"])
AcknowledgmentPrimitive  = _def("Acknowledgment",       1, ["signal.*"])
CommitmentPrimitive      = _def("Commitment",           1, ["signal.*", "agreement.*", "intent.*"])

# ===================================================================
# Layer 2 — Exchange (12 primitives)
# ===================================================================

# Common Ground
TermPrimitive            = _def("Term",                 2, ["signal.*", "protocol.*"])
ProtocolPrimitive        = _def("Protocol",             2, ["term.*", "agreement.*"])
OfferPrimitive           = _def("Offer",                2, ["term.*", "value.*"])
AcceptancePrimitive      = _def("Acceptance",           2, ["offer.*", "consent.*"])
# Mutual Binding
AgreementPrimitive       = _def("Agreement",            2, ["offer.*", "acceptance.*"])
ObligationPrimitive      = _def("Obligation",           2, ["agreement.*", "commitment.*"])
FulfillmentPrimitive     = _def("Fulfillment",          2, ["obligation.*", "act.*"])
BreachPrimitive          = _def("Breach",               2, ["obligation.*", "violation.*"])
# Value Transfer
ExchangePrimitive        = _def("Exchange",             2, ["offer.*", "fulfillment.*"])
AccountabilityPrimitive  = _def("Accountability",       2, ["obligation.*", "breach.*", "consequence.*"])
DebtPrimitive            = _def("Debt",                 2, ["obligation.*", "exchange.*"])
ReciprocityPrimitive     = _def("Reciprocity",          2, ["exchange.*", "obligation.*", "fulfillment.*"])

# ===================================================================
# Layer 3 — Society (12 primitives)
# ===================================================================

# Collective Identity
GroupPrimitive           = _def("Group",                3, ["actor.*", "agreement.*"])
MembershipPrimitive      = _def("Membership",           3, ["group.*", "consent.*"])
RolePrimitive            = _def("Role",                 3, ["group.*", "capacity.*"])
ConsentPrimitive         = _def("Consent",              3, ["membership.*", "choice.*", "agreement.*"])
# Social Order
NormPrimitive            = _def("Norm",                 3, ["group.*", "pattern.*", "agreement.*"])
ReputationPrimitive      = _def("Reputation",           3, ["trust.*", "consequence.*", "fulfillment.*"])
SanctionPrimitive        = _def("Sanction",             3, ["norm.*", "violation.*", "breach.*"])
AuthorityPrimitive       = _def("Authority",            3, ["role.*", "trust.*", "consent.*"])
# Collective Agency
PropertyPrimitive        = _def("Property",             3, ["exchange.*", "agreement.*", "right.*"])
CommonsPrimitive         = _def("Commons",              3, ["group.*", "property.*", "resource.*"])
GovernancePrimitive      = _def("Governance",           3, ["authority.*", "norm.*", "group.*"])
CollectiveActPrimitive   = _def("CollectiveAct",        3, ["governance.*", "group.*", "choice.*"])

# ===================================================================
# Layer 4 — Legal (12 primitives)
# ===================================================================

# Codification
LawPrimitive             = _def("Law",                  4, ["norm.*", "governance.*", "authority.*"])
RightPrimitive           = _def("Right",                4, ["law.*", "dignity.*", "consent.*"])
ContractPrimitive        = _def("Contract",             4, ["agreement.*", "obligation.*", "law.*"])
LiabilityPrimitive       = _def("Liability",            4, ["breach.*", "harm.*", "contract.*"])
# Process
DueProcessPrimitive      = _def("DueProcess",           4, ["law.*", "right.*", "adjudication.*"])
AdjudicationPrimitive    = _def("Adjudication",         4, ["dispute.*", "law.*", "evidence.*"])
RemedyPrimitive          = _def("Remedy",               4, ["adjudication.*", "harm.*", "liability.*"])
PrecedentPrimitive       = _def("Precedent",            4, ["adjudication.*", "law.*"])
# Sovereign Structure
JurisdictionPrimitive    = _def("Jurisdiction",         4, ["law.*", "authority.*", "group.*"])
SovereigntyPrimitive     = _def("Sovereignty",          4, ["jurisdiction.*", "authority.*", "governance.*"])
LegitimacyPrimitive      = _def("Legitimacy",           4, ["consent.*", "authority.*", "governance.*"])
TreatyPrimitive          = _def("Treaty",               4, ["sovereignty.*", "agreement.*", "jurisdiction.*"])

# ===================================================================
# Layer 5 — Technology (12 primitives)
# ===================================================================

# Investigation
MethodPrimitive          = _def("Method",               5, ["protocol.*", "knowledge.*"])
MeasurementPrimitive     = _def("Measurement",          5, ["method.*", "data.*", "evidence.*"])
KnowledgePrimitive       = _def("Knowledge",            5, ["measurement.*", "evidence.*", "inference.*"])
ModelPrimitive           = _def("Model",                5, ["knowledge.*", "abstraction.*", "pattern.*"])
# Creation
ToolPrimitive            = _def("Tool",                 5, ["method.*", "capacity.*", "resource.*"])
TechniquePrimitive       = _def("Technique",            5, ["tool.*", "method.*", "knowledge.*"])
InventionPrimitive       = _def("Invention",            5, ["technique.*", "knowledge.*", "creativity.*"])
AbstractionPrimitive     = _def("Abstraction",          5, ["pattern.*", "model.*", "knowledge.*"])
# Systems
InfrastructurePrimitive  = _def("Infrastructure",       5, ["tool.*", "resource.*", "standard.*"])
StandardPrimitive        = _def("Standard",             5, ["protocol.*", "norm.*", "measurement.*"])
EfficiencyPrimitive      = _def("Efficiency",           5, ["measurement.*", "resource.*", "method.*"])
AutomationPrimitive      = _def("Automation",           5, ["tool.*", "technique.*", "efficiency.*"])

# ===================================================================
# Layer 6 — Information (12 primitives)
# ===================================================================

# Representation
SymbolPrimitive          = _def("Symbol",               6, ["signal.*", "pattern.*"])
LanguagePrimitive        = _def("Language",             6, ["symbol.*", "protocol.*", "term.*"])
EncodingPrimitive        = _def("Encoding",             6, ["symbol.*", "language.*", "data.*"])
RecordPrimitive          = _def("Record",               6, ["encoding.*", "event.*", "data.*"])
# Dynamics
ChannelPrimitive         = _def("Channel",              6, ["signal.*", "infrastructure.*"])
CopyPrimitive            = _def("Copy",                 6, ["record.*", "encoding.*"])
NoisePrimitive           = _def("Noise",                6, ["channel.*", "signal.*", "uncertainty.*"])
RedundancyPrimitive      = _def("Redundancy",           6, ["copy.*", "noise.*", "reliability.*"])
# Transformation
DataPrimitive            = _def("Data",                 6, ["record.*", "measurement.*", "encoding.*"])
ComputationPrimitive     = _def("Computation",          6, ["data.*", "algorithm.*", "automation.*"])
AlgorithmPrimitive       = _def("Algorithm",            6, ["method.*", "abstraction.*", "computation.*"])
EntropyPrimitive         = _def("Entropy",              6, ["data.*", "noise.*", "uncertainty.*"])

# ===================================================================
# Layer 7 — Ethics (12 primitives)
# ===================================================================

# Moral Standing
MoralStatusPrimitive     = _def("MoralStatus",          7, ["actor.*", "right.*", "dignity.*"])
DignityPrimitive         = _def("Dignity",              7, ["moralstatus.*", "harm.*", "right.*"])
AutonomyPrimitive        = _def("Autonomy",             7, ["choice.*", "consent.*", "dignity.*"])
FlourishingPrimitive     = _def("Flourishing",          7, ["autonomy.*", "care.*", "growth.*"])
# Moral Obligation
DutyPrimitive            = _def("Duty",                 7, ["obligation.*", "moralstatus.*", "right.*"])
HarmPrimitive            = _def("Harm",                 7, ["violation.*", "consequence.*", "dignity.*"])
CarePrimitive            = _def("Care",                 7, ["harm.*", "vulnerability.*", "duty.*"])
JusticePrimitive         = _def("Justice",              7, ["right.*", "harm.*", "fairness.*"])
# Moral Agency
ConsciencePrimitive      = _def("Conscience",           7, ["duty.*", "harm.*", "value.*"])
VirtuePrimitive          = _def("Virtue",               7, ["conscience.*", "duty.*", "flourishing.*"])
ResponsibilityPrimitive  = _def("Responsibility",       7, ["act.*", "consequence.*", "duty.*"])
MotivePrimitive          = _def("Motive",               7, ["intent.*", "value.*", "conscience.*"])

# ===================================================================
# Layer 8 — Identity (12 primitives)
# ===================================================================

# Self-Knowledge
NarrativePrimitive       = _def("Narrative",            8, ["memory.*", "record.*", "identity.*"])
SelfConceptPrimitive     = _def("SelfConcept",          8, ["narrative.*", "reflection.*", "value.*"])
ReflectionPrimitive      = _def("Reflection",           8, ["act.*", "consequence.*", "conscience.*"])
MemoryPrimitive          = _def("Memory",               8, ["record.*", "experience.*", "narrative.*"])
# Self-Direction
PurposePrimitive         = _def("Purpose",              8, ["value.*", "intent.*", "narrative.*"])
AspirationPrimitive      = _def("Aspiration",           8, ["purpose.*", "value.*", "growth.*"])
AuthenticityPrimitive    = _def("Authenticity",         8, ["selfconcept.*", "value.*", "choice.*"])
ExpressionPrimitive      = _def("Expression",           8, ["signal.*", "authenticity.*", "creativity.*"])
# Self-Becoming
GrowthPrimitive          = _def("Growth",               8, ["reflection.*", "learning.*", "aspiration.*"])
ContinuityPrimitive      = _def("Continuity",           8, ["memory.*", "narrative.*", "identity.*"])
IntegrationPrimitive     = _def("Integration",          8, ["selfconcept.*", "growth.*", "continuity.*"])
CrisisPrimitive          = _def("Crisis",               8, ["continuity.*", "rupture.*", "identity.*"])

# ===================================================================
# Layer 9 — Relationship (12 primitives)
# ===================================================================

# Connection
BondPrimitive            = _def("Bond",                 9, ["trust.*", "attachment.*", "commitment.*"])
AttachmentPrimitive      = _def("Attachment",           9, ["bond.*", "care.*", "signal.*"])
RecognitionPrimitive     = _def("Recognition",          9, ["identity.*", "acknowledgment.*", "dignity.*"])
IntimacyPrimitive        = _def("Intimacy",             9, ["bond.*", "vulnerability.*", "trust.*"])
# Relational Dynamics
AttunementPrimitive      = _def("Attunement",           9, ["signal.*", "reception.*", "care.*"])
RupturePrimitive         = _def("Rupture",              9, ["bond.*", "breach.*", "harm.*"])
RepairPrimitive          = _def("Repair",               9, ["rupture.*", "acknowledgment.*", "commitment.*"])
LoyaltyPrimitive         = _def("Loyalty",              9, ["bond.*", "commitment.*", "trust.*"])
# Relational Identity
MutualConstitutionPrimitive = _def("MutualConstitution", 9, ["bond.*", "identity.*", "recognition.*"])
RelationalObligationPrimitive = _def("RelationalObligation", 9, ["bond.*", "obligation.*", "care.*"])
GriefPrimitive           = _def("Grief",                9, ["loss.*", "bond.*", "attachment.*"])
ForgivenessPrimitive     = _def("Forgiveness",          9, ["rupture.*", "repair.*", "trust.*"])

# ===================================================================
# Layer 10 — Community (12 primitives)
# ===================================================================

# Shared Meaning
CulturePrimitive         = _def("Culture",             10, ["norm.*", "narrative.*", "group.*"])
SharedNarrativePrimitive = _def("SharedNarrative",     10, ["narrative.*", "culture.*", "group.*"])
EthosPrimitive           = _def("Ethos",               10, ["value.*", "culture.*", "norm.*"])
SacredPrimitive          = _def("Sacred",              10, ["ethos.*", "culture.*", "meaning.*"])
# Living Practice
TraditionPrimitive       = _def("Tradition",           10, ["culture.*", "practice.*", "continuity.*"])
RitualPrimitive          = _def("Ritual",              10, ["tradition.*", "sacred.*", "community.*"])
PracticePrimitive        = _def("Practice",            10, ["technique.*", "norm.*", "tradition.*"])
PlacePrimitive           = _def("Place",               10, ["community.*", "belonging.*", "infrastructure.*"])
# Communal Experience
BelongingPrimitive       = _def("Belonging",           10, ["membership.*", "bond.*", "culture.*"])
SolidarityPrimitive      = _def("Solidarity",          10, ["belonging.*", "collectiveact.*", "care.*"])
VoicePrimitive           = _def("Voice",               10, ["expression.*", "governance.*", "belonging.*"])
WelcomePrimitive         = _def("Welcome",             10, ["membership.*", "belonging.*", "recognition.*"])

# ===================================================================
# Layer 11 — Culture (12 primitives)
# ===================================================================

# Cultural Awareness
ReflexivityPrimitive     = _def("Reflexivity",         11, ["reflection.*", "culture.*", "selfconcept.*"])
EncounterPrimitive       = _def("Encounter",           11, ["recognition.*", "culture.*", "difference.*"])
TranslationPrimitive     = _def("Translation",         11, ["language.*", "encoding.*", "encounter.*"])
PluralismPrimitive       = _def("Pluralism",           11, ["encounter.*", "culture.*", "governance.*"])
# Cultural Creation
CreativityPrimitive      = _def("Creativity",          11, ["imagination.*", "expression.*", "invention.*"])
AestheticPrimitive       = _def("Aesthetic",           11, ["creativity.*", "culture.*", "value.*"])
InterpretationPrimitive  = _def("Interpretation",      11, ["symbol.*", "narrative.*", "meaning.*"])
DialoguePrimitive        = _def("Dialogue",            11, ["signal.*", "encounter.*", "translation.*"])
# Cultural Dynamics
SyncretismPrimitive      = _def("Syncretism",          11, ["encounter.*", "integration.*", "culture.*"])
CritiquePrimitive        = _def("Critique",            11, ["reflexivity.*", "norm.*", "value.*"])
HegemonyPrimitive        = _def("Hegemony",            11, ["authority.*", "culture.*", "norm.*"])
CulturalEvolutionPrimitive = _def("CulturalEvolution", 11, ["culture.*", "syncretism.*", "creativity.*"])

# ===================================================================
# Layer 12 — Emergence (12 primitives)
# ===================================================================

# Principles of Complexity
EmergencePrimitive       = _def("Emergence",           12, ["pattern.*", "complexity.*", "selforganization.*"])
SelfOrganizationPrimitive = _def("SelfOrganization",   12, ["feedback.*", "pattern.*", "autonomy.*"])
FeedbackPrimitive        = _def("Feedback",            12, ["consequence.*", "signal.*", "measurement.*"])
ComplexityPrimitive      = _def("Complexity",          12, ["pattern.*", "emergence.*", "entropy.*"])
# Limits and Self-Reference
ConsciousnessPrimitive   = _def("Consciousness",      12, ["selfconcept.*", "reflection.*", "recursion.*"])
RecursionPrimitive       = _def("Recursion",           12, ["selforganization.*", "feedback.*", "pattern.*"])
ParadoxPrimitive         = _def("Paradox",             12, ["recursion.*", "contradiction.*", "incompleteness.*"])
IncompletenesPrimitive   = _def("Incompleteness",      12, ["knowledge.*", "paradox.*", "boundary.*"])
# Dynamic Architecture
PhaseTransitionPrimitive = _def("PhaseTransition",     12, ["threshold.*", "emergence.*", "complexity.*"])
DownwardCausationPrimitive = _def("DownwardCausation", 12, ["emergence.*", "constraint.*", "pattern.*"])
AutopoiesisPrimitive     = _def("Autopoiesis",        12, ["selforganization.*", "boundary.*", "recursion.*"])
CoEvolutionPrimitive     = _def("CoEvolution",         12, ["feedback.*", "adaptation.*", "mutualconstitution.*"])

# ===================================================================
# Layer 13 — Existence (12 primitives)
# ===================================================================

# The Given
BeingPrimitive           = _def("Being",               13, ["clock.tick", "existence.*"])
NothingnessPrimitive     = _def("Nothingness",         13, ["void.*", "absence.*", "being.*"])
FinitudePrimitive        = _def("Finitude",            13, ["being.*", "boundary.*", "mortality.*"])
ContingencyPrimitive     = _def("Contingency",         13, ["being.*", "uncertainty.*", "causality.*"])
# The Response
WonderPrimitive          = _def("Wonder",              13, ["mystery.*", "emergence.*", "being.*"])
ExistentialAcceptancePrimitive = _def("ExistentialAcceptance", 13, ["finitude.*", "contingency.*", "grief.*"])
PresencePrimitive        = _def("Presence",            13, ["being.*", "attunement.*", "clock.tick"])
GratitudePrimitive       = _def("Gratitude",           13, ["presence.*", "being.*", "gift.*"])
# The Horizon
MysteryPrimitive         = _def("Mystery",             13, ["wonder.*", "incompleteness.*", "being.*"])
TranscendencePrimitive   = _def("Transcendence",       13, ["mystery.*", "being.*", "meaning.*"])
GroundlessnessPrimitive  = _def("Groundlessness",      13, ["contingency.*", "nothingness.*", "paradox.*"])
ReturnPrimitive          = _def("Return",              13, ["presence.*", "being.*", "recursion.*"])


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
    # Layer 1 — Agency (12)
    ValuePrimitive, IntentPrimitive, ChoicePrimitive, RiskPrimitive,
    ActPrimitive, ConsequencePrimitive, CapacityPrimitive, ResourcePrimitive,
    SignalPrimitive, ReceptionPrimitive, AcknowledgmentPrimitive,
    CommitmentPrimitive,
    # Layer 2 — Exchange (12)
    TermPrimitive, ProtocolPrimitive, OfferPrimitive, AcceptancePrimitive,
    AgreementPrimitive, ObligationPrimitive, FulfillmentPrimitive,
    BreachPrimitive, ExchangePrimitive, AccountabilityPrimitive,
    DebtPrimitive, ReciprocityPrimitive,
    # Layer 3 — Society (12)
    GroupPrimitive, MembershipPrimitive, RolePrimitive, ConsentPrimitive,
    NormPrimitive, ReputationPrimitive, SanctionPrimitive, AuthorityPrimitive,
    PropertyPrimitive, CommonsPrimitive, GovernancePrimitive,
    CollectiveActPrimitive,
    # Layer 4 — Legal (12)
    LawPrimitive, RightPrimitive, ContractPrimitive, LiabilityPrimitive,
    DueProcessPrimitive, AdjudicationPrimitive, RemedyPrimitive,
    PrecedentPrimitive, JurisdictionPrimitive, SovereigntyPrimitive,
    LegitimacyPrimitive, TreatyPrimitive,
    # Layer 5 — Technology (12)
    MethodPrimitive, MeasurementPrimitive, KnowledgePrimitive, ModelPrimitive,
    ToolPrimitive, TechniquePrimitive, InventionPrimitive, AbstractionPrimitive,
    InfrastructurePrimitive, StandardPrimitive, EfficiencyPrimitive,
    AutomationPrimitive,
    # Layer 6 — Information (12)
    SymbolPrimitive, LanguagePrimitive, EncodingPrimitive, RecordPrimitive,
    ChannelPrimitive, CopyPrimitive, NoisePrimitive, RedundancyPrimitive,
    DataPrimitive, ComputationPrimitive, AlgorithmPrimitive, EntropyPrimitive,
    # Layer 7 — Ethics (12)
    MoralStatusPrimitive, DignityPrimitive, AutonomyPrimitive,
    FlourishingPrimitive, DutyPrimitive, HarmPrimitive, CarePrimitive,
    JusticePrimitive, ConsciencePrimitive, VirtuePrimitive,
    ResponsibilityPrimitive, MotivePrimitive,
    # Layer 8 — Identity (12)
    NarrativePrimitive, SelfConceptPrimitive, ReflectionPrimitive,
    MemoryPrimitive, PurposePrimitive, AspirationPrimitive,
    AuthenticityPrimitive, ExpressionPrimitive, GrowthPrimitive,
    ContinuityPrimitive, IntegrationPrimitive, CrisisPrimitive,
    # Layer 9 — Relationship (12)
    BondPrimitive, AttachmentPrimitive, RecognitionPrimitive,
    IntimacyPrimitive, AttunementPrimitive, RupturePrimitive,
    RepairPrimitive, LoyaltyPrimitive, MutualConstitutionPrimitive,
    RelationalObligationPrimitive, GriefPrimitive, ForgivenessPrimitive,
    # Layer 10 — Community (12)
    CulturePrimitive, SharedNarrativePrimitive, EthosPrimitive,
    SacredPrimitive, TraditionPrimitive, RitualPrimitive, PracticePrimitive,
    PlacePrimitive, BelongingPrimitive, SolidarityPrimitive, VoicePrimitive,
    WelcomePrimitive,
    # Layer 11 — Culture (12)
    ReflexivityPrimitive, EncounterPrimitive, TranslationPrimitive,
    PluralismPrimitive, CreativityPrimitive, AestheticPrimitive,
    InterpretationPrimitive, DialoguePrimitive, SyncretismPrimitive,
    CritiquePrimitive, HegemonyPrimitive, CulturalEvolutionPrimitive,
    # Layer 12 — Emergence (12)
    EmergencePrimitive, SelfOrganizationPrimitive, FeedbackPrimitive,
    ComplexityPrimitive, ConsciousnessPrimitive, RecursionPrimitive,
    ParadoxPrimitive, IncompletenesPrimitive, PhaseTransitionPrimitive,
    DownwardCausationPrimitive, AutopoiesisPrimitive, CoEvolutionPrimitive,
    # Layer 13 — Existence (12)
    BeingPrimitive, NothingnessPrimitive, FinitudePrimitive,
    ContingencyPrimitive, WonderPrimitive, ExistentialAcceptancePrimitive,
    PresencePrimitive, GratitudePrimitive, MysteryPrimitive,
    TranscendencePrimitive, GroundlessnessPrimitive, ReturnPrimitive,
]


def create_all_primitives() -> list[_Base]:
    """Instantiate and return all 201 primitives."""
    return [cls() for cls in ALL_PRIMITIVE_CLASSES]
