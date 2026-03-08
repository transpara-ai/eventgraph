package compositions_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/compositions"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

func TestKnowledgeGrammar(t *testing.T) {
	t.Run("ClaimAndCategorize", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		claim, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"Go 1.24 supports generic type aliases, confidence 0.95",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		category, _ := knowledge.Categorize(env.ctx, analyst.ID(),
			claim.ID(), "programming_languages/go/features", env.convID, signer)

		ancestors := env.ancestors(category.ID(), 5)
		if !containsEvent(ancestors, claim.ID()) {
			t.Error("category should reference claim")
		}
		env.verifyChain()
	})

	t.Run("AbstractAndInfer", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		fact1, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"Service A handles 10k RPS on Go",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		fact2, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"Service B handles 12k RPS on Go",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		abstraction, _ := knowledge.Abstract(env.ctx, analyst.ID(),
			"Go services typically handle 10k+ RPS",
			[]types.EventID{fact1.ID(), fact2.ID()}, env.convID, signer)
		inference, _ := knowledge.Infer(env.ctx, analyst.ID(),
			"new Go service C should handle 10k+ RPS, confidence 0.7",
			abstraction.ID(), env.convID, signer)

		ancestors := env.ancestors(inference.ID(), 10)
		if !containsEvent(ancestors, fact1.ID()) {
			t.Error("inference should trace to fact1")
		}
		env.verifyChain()
	})

	t.Run("ChallengeAndCorrect", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)
		reviewer := env.actor("Reviewer", 2, event.ActorTypeAI)

		claim, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"Python is faster than Go for web servers",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		challenge, _ := knowledge.Challenge(env.ctx, reviewer.ID(),
			"benchmark shows Go 3x faster than Python for HTTP serving",
			claim.ID(), env.convID, signer)
		correction, _ := knowledge.Correct(env.ctx, analyst.ID(),
			"Go is significantly faster than Python for web servers",
			challenge.ID(), env.convID, signer)

		ancestors := env.ancestors(correction.ID(), 10)
		if !containsEvent(ancestors, claim.ID()) {
			t.Error("correction should trace to original claim")
		}
		env.verifyChain()
	})

	t.Run("DetectBias", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)
		reviewer := env.actor("Reviewer", 2, event.ActorTypeAI)

		claim, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"framework X is the best for microservices",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		bias, _ := knowledge.DetectBias(env.ctx, reviewer.ID(), claim.ID(),
			"vendor bias: all cited sources are from framework X's company",
			env.convID, signer)

		_ = bias
		env.verifyChain()
	})

	t.Run("Learn", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		mistake, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"predicted Service X handles 10k RPS, actual was 6k",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		learning, _ := knowledge.Learn(env.ctx, analyst.ID(),
			"always verify benchmarks include production conditions",
			[]types.EventID{mistake.ID()}, env.convID, signer)

		ancestors := env.ancestors(learning.ID(), 5)
		if !containsEvent(ancestors, mistake.ID()) {
			t.Error("learning should trace to mistake")
		}
		env.verifyChain()
	})

	t.Run("FactCheck", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)
		checker := env.actor("FactChecker", 2, event.ActorTypeAI)

		claim, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"event sourcing always improves performance",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := knowledge.FactCheck(env.ctx, checker.ID(), claim.ID(),
			"source: blog post, no benchmarks cited",
			"absolute claim without qualification, no counter-evidence considered",
			"MISLEADING — event sourcing improves auditability but can decrease read performance",
			env.convID, signer)
		if err != nil {
			t.Fatalf("FactCheck: %v", err)
		}

		ancestors := env.ancestors(result.Verdict.ID(), 5)
		if !containsEvent(ancestors, claim.ID()) {
			t.Error("verdict should trace to claim")
		}
		env.verifyChain()
	})

	t.Run("EncodeAndRecall", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		claim, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"Go generics reduce boilerplate by ~40%",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		encoded, _ := knowledge.Encode(env.ctx, analyst.ID(),
			"JSON: {\"language\":\"Go\",\"feature\":\"generics\",\"reduction\":0.4}",
			claim.ID(), env.convID, signer)

		ancestors := env.ancestors(encoded.ID(), 5)
		if !containsEvent(ancestors, claim.ID()) {
			t.Error("encoding should trace to original claim")
		}
		env.verifyChain()
	})

	t.Run("RememberAndRecall", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		memory, _ := knowledge.Remember(env.ctx, analyst.ID(),
			"connection pool defaults: Postgres=100, MySQL=151",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		recall, _ := knowledge.Recall(env.ctx, analyst.ID(),
			"what are the default connection pool sizes?",
			[]types.EventID{memory.ID()}, env.convID, signer)

		ancestors := env.ancestors(recall.ID(), 5)
		if !containsEvent(ancestors, memory.ID()) {
			t.Error("recall should trace to memory")
		}
		env.verifyChain()
	})

	t.Run("AbstractRequiresTwo", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		fact, _ := knowledge.Claim(env.ctx, analyst.ID(), "single fact",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		_, err := knowledge.Abstract(env.ctx, analyst.ID(),
			"invalid generalization from one instance",
			[]types.EventID{fact.ID()}, env.convID, signer)
		if err == nil {
			t.Error("Abstract with < 2 instances should fail")
		}
	})

	t.Run("Retract", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		claim, _ := knowledge.Claim(env.ctx, analyst.ID(),
			"library X has no known vulnerabilities",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		retraction, err := knowledge.Retract(env.ctx, analyst.ID(),
			claim.ID(), "CVE-2026-1234 discovered after publication",
			env.convID, signer)
		if err != nil {
			t.Fatalf("Retract: %v", err)
		}

		ancestors := env.ancestors(retraction.ID(), 5)
		if !containsEvent(ancestors, claim.ID()) {
			t.Error("retraction should trace to claim")
		}
		env.verifyChain()
	})

	t.Run("Verify", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		result, err := knowledge.Verify(env.ctx, analyst.ID(),
			"Go 1.24 supports range over int, confirmed in release notes",
			"source: go.dev/doc/go1.24, section 'Language Changes'",
			"independently verified via playground test and spec diff",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}

		ancestors := env.ancestors(result.Corroboration.ID(), 10)
		if !containsEvent(ancestors, result.Claim.ID()) {
			t.Error("corroboration should trace to claim")
		}
		if !containsEvent(ancestors, result.Provenance.ID()) {
			t.Error("corroboration should trace to provenance")
		}
		env.verifyChain()
	})

	t.Run("Survey", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		result, err := knowledge.Survey(env.ctx, analyst.ID(),
			[]string{
				"what is the p99 latency of service_alpha?",
				"what is the p99 latency of service_beta?",
				"what is the p99 latency of service_gamma?",
			},
			"Go services in this cluster have p99 latency between 10ms and 30ms",
			"cluster wide latency is well within SLA, no action needed",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Survey: %v", err)
		}

		if len(result.Recalls) != 3 {
			t.Errorf("expected 3 recall events, got %d", len(result.Recalls))
		}
		ancestors := env.ancestors(result.Synthesis.ID(), 10)
		if !containsEvent(ancestors, result.Abstraction.ID()) {
			t.Error("synthesis should trace to abstraction")
		}
		if !containsEvent(ancestors, result.Recalls[0].ID()) {
			t.Error("synthesis should trace to first recall")
		}
		env.verifyChain()
	})

	t.Run("SurveyRequiresTwoQueries", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		_, err := knowledge.Survey(env.ctx, analyst.ID(),
			[]string{"only one query"},
			"invalid generalization",
			"invalid synthesis",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err == nil {
			t.Error("Survey with < 2 queries should fail")
		}
	})

	t.Run("KnowledgeBase", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		result, err := knowledge.KnowledgeBase(env.ctx, analyst.ID(),
			[]string{
				"Go uses goroutines for concurrency",
				"Rust uses async/await for concurrency",
			},
			[]string{
				"programming_languages/go/concurrency",
				"programming_languages/rust/concurrency",
			},
			"concurrency_models_comparison",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("KnowledgeBase: %v", err)
		}

		if len(result.Claims) != 2 {
			t.Errorf("expected 2 claim events, got %d", len(result.Claims))
		}
		if len(result.Categories) != 2 {
			t.Errorf("expected 2 category events, got %d", len(result.Categories))
		}
		ancestors := env.ancestors(result.Memory.ID(), 10)
		if !containsEvent(ancestors, result.Claims[0].ID()) {
			t.Error("memory should trace to first claim")
		}
		env.verifyChain()
	})

	t.Run("Transfer", func(t *testing.T) {
		env := newTestEnv(t)
		knowledge := compositions.NewKnowledgeGrammar(env.grammar)
		analyst := env.actor("Analyst", 1, event.ActorTypeAI)

		result, err := knowledge.Transfer(env.ctx, analyst.ID(),
			"what are best practices for connection pool sizing?",
			"JSON: {\"min\":10,\"max\":100,\"per_core\":5,\"idle_timeout_s\":300}",
			"connection pools should scale with core count, not request volume",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Transfer: %v", err)
		}

		ancestors := env.ancestors(result.Learn.ID(), 10)
		if !containsEvent(ancestors, result.Recall.ID()) {
			t.Error("learn should trace to recall")
		}
		if !containsEvent(ancestors, result.Encode.ID()) {
			t.Error("learn should trace to encode")
		}
		env.verifyChain()
	})
}
