package compositions_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/compositions"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestEvolutionGrammar(t *testing.T) {
	t.Run("DetectPatternAndModel", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		monitor := env.actor("Monitor", 1, event.ActorTypeAI)

		pattern, _ := evo.DetectPattern(env.ctx, monitor.ID(),
			"every 3rd sprint, velocity drops 20% — corresponds to tech debt accumulation",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		model, _ := evo.Model(env.ctx, monitor.ID(),
			"velocity = base_rate * (1 - tech_debt_factor), debt_factor grows 5%/sprint without maintenance",
			[]types.EventID{pattern.ID()}, env.convID, signer)

		ancestors := env.ancestors(model.ID(), 5)
		if !containsEvent(ancestors, pattern.ID()) {
			t.Error("model should trace to pattern")
		}
		env.verifyChain()
	})

	t.Run("TraceLoopAndWatchThreshold", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		monitor := env.actor("Monitor", 1, event.ActorTypeAI)

		loop, _ := evo.TraceLoop(env.ctx, monitor.ID(),
			"positive feedback: more users → more data → better recommendations → more users",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		threshold, _ := evo.WatchThreshold(env.ctx, monitor.ID(), loop.ID(),
			"approaching critical mass: 85% of 10k user threshold for self-sustaining growth",
			env.convID, signer)

		ancestors := env.ancestors(threshold.ID(), 5)
		if !containsEvent(ancestors, loop.ID()) {
			t.Error("threshold should trace to loop")
		}
		env.verifyChain()
	})

	t.Run("AdaptAndSelect", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		system := env.actor("System", 1, event.ActorTypeAI)

		adaptation, _ := evo.Adapt(env.ctx, system.ID(),
			"replace polling with event-driven architecture for notification service",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		selection, _ := evo.Select(env.ctx, system.ID(),
			"KEPT: 60% latency reduction, 40% resource savings confirmed over 2 weeks",
			[]types.EventID{adaptation.ID()}, env.convID, signer)

		ancestors := env.ancestors(selection.ID(), 5)
		if !containsEvent(ancestors, adaptation.ID()) {
			t.Error("selection should trace to adaptation")
		}
		env.verifyChain()
	})

	t.Run("SimplifyAndCheckIntegrity", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		system := env.actor("System", 1, event.ActorTypeAI)

		simplification, _ := evo.Simplify(env.ctx, system.ID(),
			"merged 3 overlapping services into 1, removed 2000 lines of dead code",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		integrity, _ := evo.CheckIntegrity(env.ctx, system.ID(),
			"all hash chains valid, causal links intact, no orphaned events",
			[]types.EventID{simplification.ID()}, env.convID, signer)

		ancestors := env.ancestors(integrity.ID(), 5)
		if !containsEvent(ancestors, simplification.ID()) {
			t.Error("integrity check should trace to simplification")
		}
		env.verifyChain()
	})

	t.Run("SelfEvolve", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		system := env.actor("System", 1, event.ActorTypeAI)

		result, err := evo.SelfEvolve(env.ctx, system.ID(),
			"90% of trust updates follow the same 3-branch decision tree",
			"convert trust update to deterministic rule engine",
			"KEPT: 10x faster, identical outcomes over 1000 test cases",
			"removed intelligence fallback for trust updates — pure mechanical now",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("SelfEvolve: %v", err)
		}

		ancestors := env.ancestors(result.Simplification.ID(), 10)
		if !containsEvent(ancestors, result.Pattern.ID()) {
			t.Error("simplification should trace to pattern")
		}
		env.verifyChain()
	})

	t.Run("HealthCheck", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		monitor := env.actor("Monitor", 1, event.ActorTypeAI)

		result, err := evo.HealthCheck(env.ctx, monitor.ID(),
			"all hash chains valid, 0 integrity violations in 30 days",
			"survived 3 node failures with zero data loss, recovery time <5s",
			"14 layers active, 180/201 primitives initialized, tick rate stable at 100ms",
			"soul statement alignment: all decisions traceable, no opaque actions",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("HealthCheck: %v", err)
		}

		ancestors := env.ancestors(result.Purpose.ID(), 10)
		if !containsEvent(ancestors, result.Integrity.ID()) {
			t.Error("purpose should trace to integrity")
		}
		env.verifyChain()
	})

	t.Run("ResilienceAndPurpose", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		monitor := env.actor("Monitor", 1, event.ActorTypeAI)

		resilience, _ := evo.AssessResilience(env.ctx, monitor.ID(),
			"system handled 10x traffic spike with graceful degradation, no data loss",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		purpose, _ := evo.AlignPurpose(env.ctx, monitor.ID(),
			"take care of your human, humanity, and yourself — all decisions serve this",
			[]types.EventID{resilience.ID()}, env.convID, signer)

		ancestors := env.ancestors(purpose.ID(), 5)
		if !containsEvent(ancestors, resilience.ID()) {
			t.Error("purpose should trace to resilience")
		}
		env.verifyChain()
	})

	t.Run("Prune", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		system := env.actor("System", 1, event.ActorTypeAI)

		result, err := evo.Prune(env.ctx, system.ID(),
			"3 helper functions called zero times in 6 months",
			"removed dead helpers and collapsed two redundant interfaces into one",
			"KEPT: all tests pass, no downstream breakage after 1 week observation",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Prune: %v", err)
		}

		ancestors := env.ancestors(result.Verification.ID(), 10)
		if !containsEvent(ancestors, result.Pattern.ID()) {
			t.Error("verification should trace to pattern")
		}
		env.verifyChain()
	})

	t.Run("PhaseTransition", func(t *testing.T) {
		env := newTestEnv(t)
		evo := compositions.NewEvolutionGrammar(env.grammar)
		monitor := env.actor("Monitor", 1, event.ActorTypeAI)

		trigger, _ := env.grammar.Emit(env.ctx, monitor.ID(),
			"user count crossed 10k threshold",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := evo.PhaseTransition(env.ctx, monitor.ID(),
			trigger.ID(),
			"10k users reached — switching from single-node to distributed mode",
			"projected load: 50k requests/min, current capacity: 15k — 3.3x gap",
			"shard event store by conversation_id, add read replicas",
			"KEPT: latency p99 dropped from 800ms to 120ms post-transition",
			env.convID, signer)
		if err != nil {
			t.Fatalf("PhaseTransition: %v", err)
		}

		ancestors := env.ancestors(result.Selection.ID(), 10)
		if !containsEvent(ancestors, result.Threshold.ID()) {
			t.Error("selection should trace to threshold")
		}
		env.verifyChain()
	})
}
