"""Trust model — asymmetric, non-transitive, decaying trust with domain scoping."""

from __future__ import annotations

import threading
from dataclasses import dataclass, field
from typing import Any, Protocol, runtime_checkable

from .event import Event
from .types import ActorID, DomainScope, EventID, Score, Weight


# ── TrustMetrics ─────────────────────────────────────────────────────────


@dataclass(frozen=True, slots=True)
class TrustMetrics:
    """Immutable snapshot of an actor's trust state."""

    _actor: ActorID
    _overall: Score
    _by_domain: dict[str, Score]
    _confidence: Score
    _trend: Weight
    _evidence: list[EventID]
    _last_updated: int
    _decay_rate: Score

    def __init__(
        self,
        actor: ActorID,
        overall: Score,
        by_domain: dict[str, Score],
        confidence: Score,
        trend: Weight,
        evidence: list[EventID],
        last_updated: int,
        decay_rate: Score,
    ) -> None:
        object.__setattr__(self, "_actor", actor)
        object.__setattr__(self, "_overall", overall)
        object.__setattr__(self, "_by_domain", dict(by_domain))
        object.__setattr__(self, "_confidence", confidence)
        object.__setattr__(self, "_trend", trend)
        object.__setattr__(self, "_evidence", list(evidence))
        object.__setattr__(self, "_last_updated", last_updated)
        object.__setattr__(self, "_decay_rate", decay_rate)

    @property
    def actor(self) -> ActorID:
        return self._actor

    @property
    def overall(self) -> Score:
        return self._overall

    @property
    def by_domain(self) -> dict[str, Score]:
        return dict(self._by_domain)

    @property
    def confidence(self) -> Score:
        return self._confidence

    @property
    def trend(self) -> Weight:
        return self._trend

    @property
    def evidence(self) -> list[EventID]:
        return list(self._evidence)

    @property
    def last_updated(self) -> int:
        return self._last_updated

    @property
    def decay_rate(self) -> Score:
        return self._decay_rate


# ── TrustConfig ──────────────────────────────────────────────────────────


@dataclass(frozen=True, slots=True)
class TrustConfig:
    """Configuration for the default trust model."""

    initial_trust: Score = field(default_factory=lambda: Score(0.0))
    decay_rate: Score = field(default_factory=lambda: Score(0.01))
    max_adjustment: Weight = field(default_factory=lambda: Weight(0.1))
    observed_event_delta: float = 0.01
    trend_decay_rate: float = 0.01


# ── TrustModel Protocol ─────────────────────────────────────────────────


@runtime_checkable
class TrustModel(Protocol):
    """Protocol for trust calculation, update, and decay."""

    def score(self, actor: ActorID) -> TrustMetrics: ...

    def score_in_domain(self, actor: ActorID, domain: DomainScope) -> TrustMetrics: ...

    def update(self, actor: ActorID, evidence: Event) -> TrustMetrics: ...

    def update_between(
        self, from_actor: ActorID, to_actor: ActorID, evidence: Event
    ) -> TrustMetrics: ...

    def decay(self, actor: ActorID, elapsed_seconds: float) -> TrustMetrics: ...

    def between(self, from_actor: ActorID, to_actor: ActorID) -> TrustMetrics: ...


# ── Internal state ───────────────────────────────────────────────────────


class _TrustState:
    """Mutable internal trust state for a single actor or directed pair."""

    __slots__ = ("score", "by_domain", "evidence", "last_updated", "trend")

    def __init__(self, initial_score: float, last_updated: int) -> None:
        self.score: float = initial_score
        self.by_domain: dict[str, float] = {}
        self.evidence: list[EventID] = []
        self.last_updated: int = last_updated
        self.trend: float = 0.0


# ── DefaultTrustModel ────────────────────────────────────────────────────


def _now_nanos() -> int:
    import time

    return int(time.time() * 1_000_000_000)


class DefaultTrustModel:
    """Thread-safe trust model with linear decay and equal weighting.

    Implements the TrustModel protocol.
    """

    def __init__(self, config: TrustConfig | None = None) -> None:
        self._config = config or TrustConfig()
        self._lock = threading.Lock()
        self._scores: dict[str, _TrustState] = {}  # actor_id.value -> state
        self._directed: dict[tuple[str, str], _TrustState] = {}  # (from, to) -> state

    # ── Internal helpers ─────────────────────────────────────────────────

    def _get_or_create(self, actor_id: ActorID) -> _TrustState:
        """Get or create trust state. Caller must hold self._lock."""
        key = actor_id.value
        if key not in self._scores:
            self._scores[key] = _TrustState(
                self._config.initial_trust.value, _now_nanos()
            )
        return self._scores[key]

    def _get_or_default(self, actor_id: ActorID) -> _TrustState:
        """Get trust state or return a default without mutating storage.
        Caller must hold self._lock.
        """
        key = actor_id.value
        if key in self._scores:
            return self._scores[key]
        return _TrustState(self._config.initial_trust.value, _now_nanos())

    def _build_metrics(self, actor_id: ActorID, state: _TrustState) -> TrustMetrics:
        evidence_count = len(state.evidence)
        confidence = min(1.0, evidence_count / 50.0)
        by_domain = {k: Score(v) for k, v in state.by_domain.items()}
        return TrustMetrics(
            actor=actor_id,
            overall=Score(state.score),
            by_domain=by_domain,
            confidence=Score(confidence),
            trend=Weight(state.trend),
            evidence=list(state.evidence),
            last_updated=state.last_updated,
            decay_rate=self._config.decay_rate,
        )

    def _build_domain_metrics(
        self, actor_id: ActorID, state: _TrustState, domain_score: float
    ) -> TrustMetrics:
        evidence_count = len(state.evidence)
        confidence = min(1.0, evidence_count / 50.0)
        by_domain = {k: Score(v) for k, v in state.by_domain.items()}
        return TrustMetrics(
            actor=actor_id,
            overall=Score(domain_score),
            by_domain=by_domain,
            confidence=Score(confidence),
            trend=Weight(state.trend),
            evidence=list(state.evidence),
            last_updated=state.last_updated,
            decay_rate=self._config.decay_rate,
        )

    def _extract_delta(self, evidence: Event, current_score: float) -> float:
        """Extract trust delta from evidence event content."""
        content = evidence.content
        if isinstance(content, dict) and "current" in content:
            return content["current"] - current_score
        return self._config.observed_event_delta

    def _apply_update(self, state: _TrustState, evidence: Event) -> None:
        """Apply an evidence event to a trust state. Caller must hold self._lock."""
        # Deduplicate
        for eid in state.evidence:
            if eid.value == evidence.id.value:
                return

        delta = self._extract_delta(evidence, state.score)

        # Clamp to max_adjustment
        max_adj = self._config.max_adjustment.value
        if delta > max_adj:
            delta = max_adj
        if delta < -max_adj:
            delta = -max_adj

        # Apply delta, clamp to [0, 1]
        new_score = max(0.0, min(1.0, state.score + delta))
        state.score = new_score

        # Update domain-specific score if evidence carries domain info
        content = evidence.content
        if isinstance(content, dict) and "domain" in content:
            domain_key = content["domain"]
            if domain_key in state.by_domain:
                domain_score = max(0.0, min(1.0, state.by_domain[domain_key] + delta))
            else:
                domain_score = max(
                    0.0, min(1.0, self._config.initial_trust.value + delta)
                )
            state.by_domain[domain_key] = domain_score

        # Update trend
        if delta > 0:
            state.trend = min(1.0, state.trend + 0.1)
        elif delta < 0:
            state.trend = max(-1.0, state.trend - 0.1)

        # Track evidence (cap at 100)
        state.evidence.append(evidence.id)
        if len(state.evidence) > 100:
            state.evidence = state.evidence[len(state.evidence) - 100 :]

        state.last_updated = _now_nanos()

    def _decay_state(
        self, state: _TrustState, decay_amount: float, trend_decay: float
    ) -> None:
        """Apply linear decay to a trust state. Caller must hold self._lock."""
        state.score = max(0.0, state.score - decay_amount)

        for domain_key in state.by_domain:
            state.by_domain[domain_key] = max(
                0.0, state.by_domain[domain_key] - decay_amount
            )

        if state.trend > 0:
            state.trend = max(0.0, state.trend - trend_decay)
        elif state.trend < 0:
            state.trend = min(0.0, state.trend + trend_decay)

        state.last_updated = _now_nanos()

    # ── Public API (TrustModel protocol) ─────────────────────────────────

    def score(self, actor: ActorID) -> TrustMetrics:
        """Return current trust metrics for an actor."""
        with self._lock:
            state = self._get_or_default(actor)
            return self._build_metrics(actor, state)

    def score_in_domain(self, actor: ActorID, domain: DomainScope) -> TrustMetrics:
        """Return domain-specific trust, falling back to global with halved confidence."""
        with self._lock:
            state = self._get_or_default(actor)

            if domain.value in state.by_domain:
                return self._build_domain_metrics(
                    actor, state, state.by_domain[domain.value]
                )

            # Fall back to global score with halved confidence
            evidence_count = len(state.evidence)
            global_confidence = min(1.0, evidence_count / 50.0)
            by_domain = {k: Score(v) for k, v in state.by_domain.items()}
            return TrustMetrics(
                actor=actor,
                overall=Score(state.score),
                by_domain=by_domain,
                confidence=Score(global_confidence * 0.5),
                trend=Weight(state.trend),
                evidence=list(state.evidence),
                last_updated=state.last_updated,
                decay_rate=self._config.decay_rate,
            )

    def update(self, actor: ActorID, evidence: Event) -> TrustMetrics:
        """Update trust based on evidence event."""
        with self._lock:
            state = self._get_or_create(actor)

            # Check deduplication before applying
            for eid in state.evidence:
                if eid.value == evidence.id.value:
                    return self._build_metrics(actor, state)

            self._apply_update(state, evidence)
            return self._build_metrics(actor, state)

    def update_between(
        self, from_actor: ActorID, to_actor: ActorID, evidence: Event
    ) -> TrustMetrics:
        """Update directional trust from one actor toward another."""
        with self._lock:
            key = (from_actor.value, to_actor.value)
            if key not in self._directed:
                self._directed[key] = _TrustState(
                    self._config.initial_trust.value, _now_nanos()
                )
            state = self._directed[key]

            # Check deduplication before applying
            for eid in state.evidence:
                if eid.value == evidence.id.value:
                    return self._build_metrics(to_actor, state)

            self._apply_update(state, evidence)
            return self._build_metrics(to_actor, state)

    def decay(self, actor: ActorID, elapsed_seconds: float) -> TrustMetrics:
        """Decay trust over elapsed time. Also decays directed trust where actor is from."""
        with self._lock:
            if elapsed_seconds <= 0:
                state = self._get_or_default(actor)
                return self._build_metrics(actor, state)

            days = elapsed_seconds / 86400.0
            decay_amount = self._config.decay_rate.value * days
            trend_decay = self._config.trend_decay_rate * days

            # Decay undirected trust
            key = actor.value
            has_state = key in self._scores
            if has_state:
                self._decay_state(self._scores[key], decay_amount, trend_decay)

            # Decay directed trust where actor is the from side
            for dkey, dstate in self._directed.items():
                if dkey[0] == actor.value:
                    self._decay_state(dstate, decay_amount, trend_decay)

            if not has_state:
                defaults = _TrustState(self._config.initial_trust.value, _now_nanos())
                return self._build_metrics(actor, defaults)

            return self._build_metrics(actor, self._scores[key])

    def between(self, from_actor: ActorID, to_actor: ActorID) -> TrustMetrics:
        """Return directed trust from one actor to another."""
        with self._lock:
            key = (from_actor.value, to_actor.value)
            state = self._directed.get(key)
            if state is None:
                return TrustMetrics(
                    actor=to_actor,
                    overall=self._config.initial_trust,
                    by_domain={},
                    confidence=Score(0.0),
                    trend=Weight(0.0),
                    evidence=[],
                    last_updated=_now_nanos(),
                    decay_rate=self._config.decay_rate,
                )
            return self._build_metrics(to_actor, state)
