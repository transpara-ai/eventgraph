package compositions_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/compositions"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestBuildGrammar(t *testing.T) {
	t.Run("BuildAndVersion", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		v1, _ := build.Build(env.ctx, dev.ID(), "eventgraph-cli v1.0.0",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		v2, _ := build.Version(env.ctx, dev.ID(),
			"eventgraph-cli v1.1.0 — added JSON output", v1.ID(), env.convID, signer)
		v3, _ := build.Version(env.ctx, dev.ID(),
			"eventgraph-cli v2.0.0 — breaking: new config format", v2.ID(), env.convID, signer)

		ancestors := env.ancestors(v3.ID(), 10)
		if !containsEvent(ancestors, v1.ID()) {
			t.Error("v3 should trace to v1")
		}
		env.verifyChain()
	})

	t.Run("ShipAndSunset", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		artefact, _ := build.Build(env.ctx, dev.ID(), "auth-lib v1.0",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		shipped, _ := build.Ship(env.ctx, dev.ID(), "auth-lib v1.0 to package registry",
			[]types.EventID{artefact.ID()}, env.convID, signer)
		sunset, _ := build.Sunset(env.ctx, dev.ID(), artefact.ID(),
			"replaced by auth-lib-v2, removal date 2026-09-01", env.convID, signer)

		_ = shipped
		_ = sunset
		env.verifyChain()
	})

	t.Run("TestAndReview", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)
		reviewer := env.actor("Reviewer", 2, event.ActorTypeHuman)

		code, _ := build.Build(env.ctx, dev.ID(), "auth module implementation",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		testResult, _ := build.Test(env.ctx, dev.ID(),
			"45/45 passing, coverage 91%, no regressions",
			[]types.EventID{code.ID()}, env.convID, signer)
		review, _ := build.Review(env.ctx, reviewer.ID(),
			"code quality good, tests comprehensive, approved",
			testResult.ID(), env.convID, signer)

		ancestors := env.ancestors(review.ID(), 10)
		if !containsEvent(ancestors, code.ID()) {
			t.Error("review should trace to code")
		}
		env.verifyChain()
	})

	t.Run("FeedbackAndIterate", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		v1, _ := build.Ship(env.ctx, dev.ID(), "CLI tool v1.0",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		feedback, _ := build.Feedback(env.ctx, dev.ID(),
			"output is hard to read, needs colour coding and table format",
			v1.ID(), env.convID, signer)
		v2, _ := build.Iterate(env.ctx, dev.ID(),
			"CLI tool v1.1 — added colour output and table format",
			feedback.ID(), env.convID, signer)

		ancestors := env.ancestors(v2.ID(), 10)
		if !containsEvent(ancestors, v1.ID()) {
			t.Error("iteration should trace to v1")
		}
		env.verifyChain()
	})

	t.Run("Pipeline", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		ci := env.actor("CI", 1, event.ActorTypeAI)

		result, err := build.Pipeline(env.ctx, ci.ID(),
			"build+test+lint for commit abc123",
			"234/234 passing, coverage 88%",
			"0 lint issues, build time 45s",
			"staging deployment successful",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Pipeline: %v", err)
		}

		ancestors := env.ancestors(result.Deployment.ID(), 10)
		if !containsEvent(ancestors, result.Definition.ID()) {
			t.Error("deployment should trace to definition")
		}
		env.verifyChain()
	})

	t.Run("PostMortem", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		lead := env.actor("Lead", 1, event.ActorTypeHuman)
		eng1 := env.actor("Eng1", 2, event.ActorTypeHuman)
		eng2 := env.actor("Eng2", 3, event.ActorTypeHuman)

		incident, _ := build.Build(env.ctx, lead.ID(),
			"incident: 45-minute production outage",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := build.PostMortem(env.ctx, lead.ID(),
			[]types.ActorID{eng1.ID(), eng2.ID()},
			[]string{
				"connection pool was set to default 10, needs 50+",
				"monitoring didn't alert until connections were fully exhausted",
			},
			"root cause was under-provisioned connection pool + late alerting",
			"1) increase pool to 100 2) add connection utilisation alert at 80%",
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("PostMortem: %v", err)
		}

		if len(result.Feedback) != 2 {
			t.Errorf("expected 2 feedback events, got %d", len(result.Feedback))
		}
		ancestors := env.ancestors(result.Improvements.ID(), 10)
		if !containsEvent(ancestors, incident.ID()) {
			t.Error("improvements should trace to incident")
		}
		env.verifyChain()
	})

	t.Run("DefineAndAutomate", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		workflow, _ := build.Define(env.ctx, dev.ID(), "manual deploy: build, test, ship",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		automated, _ := build.Automate(env.ctx, dev.ID(),
			"replaced manual deploy with CI pipeline",
			workflow.ID(), env.convID, signer)

		ancestors := env.ancestors(automated.ID(), 5)
		if !containsEvent(ancestors, workflow.ID()) {
			t.Error("automation should trace to workflow definition")
		}
		env.verifyChain()
	})

	t.Run("MeasureAndInnovate", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		artefact, _ := build.Build(env.ctx, dev.ID(), "query engine v1",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		measure, _ := build.Measure(env.ctx, dev.ID(), artefact.ID(),
			"latency p50=12ms p99=45ms, throughput 10k qps",
			env.convID, signer)
		innovation, _ := build.Innovate(env.ctx, dev.ID(),
			"vectorized query processing — 5x throughput improvement",
			[]types.EventID{measure.ID()}, env.convID, signer)

		ancestors := env.ancestors(innovation.ID(), 10)
		if !containsEvent(ancestors, artefact.ID()) {
			t.Error("innovation should trace to original artefact")
		}
		env.verifyChain()
	})

	t.Run("Spike", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		result, err := build.Spike(env.ctx, dev.ID(),
			"replace JSON parser with streaming decoder",
			"benchmark: 3x faster for large payloads, memory down 60%",
			"streaming approach viable but needs error recovery work",
			"proceed with streaming decoder, add retry wrapper",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Spike: %v", err)
		}

		ancestors := env.ancestors(result.Decision.ID(), 10)
		if !containsEvent(ancestors, result.Build.ID()) {
			t.Error("decision should trace to build")
		}
		if !containsEvent(ancestors, result.Test.ID()) {
			t.Error("decision should trace to test")
		}
		if !containsEvent(ancestors, result.Feedback.ID()) {
			t.Error("decision should trace to feedback")
		}
		env.verifyChain()
	})

	t.Run("Migration", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		oldLib, _ := build.Build(env.ctx, dev.ID(), "auth_lib v1.0",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := build.Migration(env.ctx, dev.ID(),
			oldLib.ID(),
			"migrate from auth_lib v1 to v2, update all callers",
			"auth_lib v2.0.0 with token rotation support",
			"auth_lib v2.0.0 deployed to staging",
			"integration tests 112/112 passing, no regressions",
			env.convID, signer)
		if err != nil {
			t.Fatalf("Migration: %v", err)
		}

		ancestors := env.ancestors(result.Test.ID(), 10)
		if !containsEvent(ancestors, result.Sunset.ID()) {
			t.Error("test should trace to sunset")
		}
		if !containsEvent(ancestors, oldLib.ID()) {
			t.Error("test should trace to deprecated target")
		}
		env.verifyChain()
	})

	t.Run("TechDebt", func(t *testing.T) {
		env := newTestEnv(t)
		build := compositions.NewBuildGrammar(env.grammar)
		dev := env.actor("Developer", 1, event.ActorTypeHuman)

		codebase, _ := build.Build(env.ctx, dev.ID(), "event_processor module",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := build.TechDebt(env.ctx, dev.ID(),
			codebase.ID(),
			"cyclomatic complexity 42, test coverage 55%, 3 known race conditions",
			"technical debt: tangled event routing logic needs decomposition",
			"split monolithic processor into pipeline stages over 3 sprints",
			env.convID, signer)
		if err != nil {
			t.Fatalf("TechDebt: %v", err)
		}

		ancestors := env.ancestors(result.Iteration.ID(), 10)
		if !containsEvent(ancestors, result.Measure.ID()) {
			t.Error("iteration should trace to measure")
		}
		if !containsEvent(ancestors, result.DebtMark.ID()) {
			t.Error("iteration should trace to debt mark")
		}
		if !containsEvent(ancestors, codebase.ID()) {
			t.Error("iteration should trace to original codebase")
		}
		env.verifyChain()
	})
}
