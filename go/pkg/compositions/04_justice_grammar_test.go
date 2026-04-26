package compositions_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/compositions"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestJusticeGrammar(t *testing.T) {
	t.Run("Legislate", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)

		rule, _ := justice.Legislate(env.ctx, admin.ID(),
			"all deployments require passing CI",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		if rule.Source() != admin.ID() {
			t.Error("rule source should be admin")
		}
		env.verifyChain()
	})

	t.Run("AmendAndRepeal", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)

		rule, _ := justice.Legislate(env.ctx, admin.ID(), "no deploys on Friday",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		amended, _ := justice.Amend(env.ctx, admin.ID(),
			"no deploys on Friday after 2pm", rule.ID(), env.convID, signer)
		ancestors := env.ancestors(amended.ID(), 5)
		if !containsEvent(ancestors, rule.ID()) {
			t.Error("amendment should trace to original rule")
		}
		env.verifyChain()
	})

	t.Run("FileAndJudge", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		plaintiff := env.actor("Plaintiff", 1, event.ActorTypeHuman)
		judge := env.actor("Judge", 2, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, plaintiff.ID(), "incident occurred",
			env.convID, []types.EventID{env.boot.ID()}, signer)
		filing, _ := justice.File(env.ctx, plaintiff.ID(),
			"violated code of conduct section 3", incident.ID(), env.convID, signer)
		ruling, _ := justice.Judge(env.ctx, judge.ID(), "violation confirmed, warning issued",
			[]types.EventID{filing.ID()}, env.convID, signer)

		ancestors := env.ancestors(ruling.ID(), 10)
		if !containsEvent(ancestors, incident.ID()) {
			t.Error("ruling should trace to incident")
		}
		env.verifyChain()
	})

	t.Run("Appeal", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		defendant := env.actor("Defendant", 1, event.ActorTypeHuman)
		judge := env.actor("Judge", 2, event.ActorTypeHuman)

		ruling, _ := justice.Judge(env.ctx, judge.ID(), "suspension for 7 days",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		appeal, _ := justice.Appeal(env.ctx, defendant.ID(),
			"no warning was given, due process violated",
			ruling.ID(), env.convID, signer)

		ancestors := env.ancestors(appeal.ID(), 5)
		if !containsEvent(ancestors, ruling.ID()) {
			t.Error("appeal should trace to ruling")
		}
		env.verifyChain()
	})

	t.Run("Trial", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		plaintiff := env.actor("Plaintiff", 1, event.ActorTypeHuman)
		defendant := env.actor("Defendant", 2, event.ActorTypeHuman)
		judge := env.actor("Judge", 3, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, plaintiff.ID(), "contract breach",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := justice.Trial(env.ctx,
			plaintiff.ID(), defendant.ID(), judge.ID(),
			"failed to deliver on time",
			"delivery was 2 weeks late, contract specified penalty",
			"force majeure: supply chain disruption",
			"contract is clear on deadline penalties",
			"supply chain issues were foreseeable",
			"partial penalty: 50% reduction due to mitigating circumstances",
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Trial: %v", err)
		}

		if len(result.Submissions) != 2 {
			t.Errorf("expected 2 submissions, got %d", len(result.Submissions))
		}
		ancestors := env.ancestors(result.Ruling.ID(), 15)
		if !containsEvent(ancestors, result.Filing.ID()) {
			t.Error("ruling should trace to filing")
		}
		env.verifyChain()
	})

	t.Run("Pardon", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		authority := env.actor("Authority", 1, event.ActorTypeHuman)
		offender := env.actor("Offender", 2, event.ActorTypeHuman)

		ruling, _ := justice.Judge(env.ctx, authority.ID(), "suspension",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		pardon, _ := justice.Pardon(env.ctx, authority.ID(), offender.ID(),
			"time served, good behaviour", types.MustDomainScope("community"),
			ruling.ID(), env.convID, signer)

		ancestors := env.ancestors(pardon.ID(), 5)
		if !containsEvent(ancestors, ruling.ID()) {
			t.Error("pardon should trace to ruling")
		}
		env.verifyChain()
	})

	t.Run("Repeal", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)

		rule, _ := justice.Legislate(env.ctx, admin.ID(), "mandatory Friday demos",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		repeal, _ := justice.Repeal(env.ctx, admin.ID(), rule.ID(),
			"demos now async via recorded video", env.convID, signer)

		ancestors := env.ancestors(repeal.ID(), 5)
		if !containsEvent(ancestors, rule.ID()) {
			t.Error("repeal should trace to original rule")
		}
		env.verifyChain()
	})

	t.Run("Enforce", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		judge := env.actor("Judge", 1, event.ActorTypeHuman)
		executor := env.actor("Executor", 2, event.ActorTypeHuman)

		ruling, _ := justice.Judge(env.ctx, judge.ID(), "access revoked for 30 days",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		enforcement, _ := justice.Enforce(env.ctx, judge.ID(), executor.ID(),
			types.MustDomainScope("access_control"), types.MustWeight(0.8),
			ruling.ID(), env.convID, signer)

		ancestors := env.ancestors(enforcement.ID(), 5)
		if !containsEvent(ancestors, ruling.ID()) {
			t.Error("enforcement should trace to ruling")
		}
		env.verifyChain()
	})

	t.Run("Audit", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		auditor := env.actor("Auditor", 1, event.ActorTypeHuman)

		action, _ := env.grammar.Emit(env.ctx, auditor.ID(), "deployment without review",
			env.convID, []types.EventID{env.boot.ID()}, signer)
		audit, _ := justice.Audit(env.ctx, auditor.ID(), action.ID(),
			"violation of review policy, no approval found", env.convID, signer)

		ancestors := env.ancestors(audit.ID(), 5)
		if !containsEvent(ancestors, action.ID()) {
			t.Error("audit should trace to audited action")
		}
		env.verifyChain()
	})

	t.Run("Reform", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)
		judge := env.actor("Judge", 2, event.ActorTypeHuman)

		ruling, _ := justice.Judge(env.ctx, judge.ID(),
			"insufficient: current rules don't address async collaboration",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		reform, _ := justice.Reform(env.ctx, admin.ID(),
			"add async collaboration guidelines to code of conduct",
			ruling.ID(), env.convID, signer)

		ancestors := env.ancestors(reform.ID(), 5)
		if !containsEvent(ancestors, ruling.ID()) {
			t.Error("reform should trace to precedent ruling")
		}
		env.verifyChain()
	})

	t.Run("ConstitutionalAmendment", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		proposer := env.actor("Proposer", 1, event.ActorTypeHuman)

		precedent, _ := justice.Judge(env.ctx, proposer.ID(),
			"current governance structure insufficient",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := justice.ConstitutionalAmendment(env.ctx, proposer.ID(),
			"restructure voting rights to include all contributors",
			"all contributors with 3+ merged PRs gain voting rights",
			"assessment: expands representation, no rights diminished",
			precedent.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("ConstitutionalAmendment: %v", err)
		}

		ancestors := env.ancestors(result.RightsCheck.ID(), 10)
		if !containsEvent(ancestors, result.Reform.ID()) {
			t.Error("rights check should trace to reform")
		}
		if !containsEvent(ancestors, result.Legislation.ID()) {
			t.Error("rights check should trace to legislation")
		}
		env.verifyChain()
	})

	t.Run("Injunction", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		petitioner := env.actor("Petitioner", 1, event.ActorTypeHuman)
		judge := env.actor("Judge", 2, event.ActorTypeHuman)
		executor := env.actor("Executor", 3, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, petitioner.ID(),
			"ongoing data breach detected",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := justice.Injunction(env.ctx, petitioner.ID(), judge.ID(),
			executor.ID(), "emergency: active data exfiltration",
			"cease all external API access immediately",
			types.MustDomainScope("security"), types.MustWeight(0.9),
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Injunction: %v", err)
		}

		ancestors := env.ancestors(result.Enforcement.ID(), 10)
		if !containsEvent(ancestors, result.Filing.ID()) {
			t.Error("enforcement should trace to filing")
		}
		if !containsEvent(ancestors, result.Ruling.ID()) {
			t.Error("enforcement should trace to ruling")
		}
		env.verifyChain()
	})

	t.Run("Plea", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		defendant := env.actor("Defendant", 1, event.ActorTypeHuman)
		prosecutor := env.actor("Prosecutor", 2, event.ActorTypeHuman)
		executor := env.actor("Executor", 3, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, prosecutor.ID(),
			"unauthorized access to production database",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := justice.Plea(env.ctx, defendant.ID(), prosecutor.ID(),
			executor.ID(), "accessed production DB without authorization",
			"accept 7-day suspension instead of full tribunal",
			types.MustDomainScope("access_control"), types.MustWeight(0.6),
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Plea: %v", err)
		}

		ancestors := env.ancestors(result.Enforcement.ID(), 10)
		if !containsEvent(ancestors, result.Filing.ID()) {
			t.Error("enforcement should trace to filing")
		}
		if !containsEvent(ancestors, result.Acceptance.ID()) {
			t.Error("enforcement should trace to acceptance")
		}
		env.verifyChain()
	})

	t.Run("ClassAction", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		p1 := env.actor("Plaintiff1", 1, event.ActorTypeHuman)
		p2 := env.actor("Plaintiff2", 2, event.ActorTypeHuman)
		defendant := env.actor("Defendant", 3, event.ActorTypeHuman)
		judge := env.actor("Judge", 4, event.ActorTypeHuman)

		incident, _ := env.grammar.Emit(env.ctx, p1.ID(),
			"systematic policy violations",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := justice.ClassAction(env.ctx,
			[]types.ActorID{p1.ID(), p2.ID()},
			defendant.ID(), judge.ID(),
			[]string{"denied access to shared resources", "excluded from decision process"},
			"logs showing systematic exclusion over 3 months",
			"pattern of deliberate exclusion violates community charter",
			"access was restricted due to security audit",
			"security audit was completed, restrictions were standard procedure",
			"violation confirmed: restrictions exceeded audit scope, remediation ordered",
			incident.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("ClassAction: %v", err)
		}

		if len(result.Filings) != 2 {
			t.Errorf("expected 2 filings, got %d", len(result.Filings))
		}
		ancestors := env.ancestors(result.Trial.Ruling.ID(), 20)
		if !containsEvent(ancestors, result.Merged.ID()) {
			t.Error("ruling should trace to merged filing")
		}
		env.verifyChain()
	})

	t.Run("Recall", func(t *testing.T) {
		env := newTestEnv(t)
		justice := compositions.NewJusticeGrammar(env.grammar)
		auditor := env.actor("Auditor", 1, event.ActorTypeHuman)
		community := env.actor("Community", 2, event.ActorTypeCommittee)
		official := env.actor("Official", 3, event.ActorTypeHuman)

		action, _ := env.grammar.Emit(env.ctx, official.ID(),
			"unilateral policy change without consultation",
			env.convID, []types.EventID{env.boot.ID()}, signer)

		result, err := justice.Recall(env.ctx, auditor.ID(), community.ID(),
			official.ID(),
			"audit: official bypassed required approval process 4 times",
			"motion to recall official from governance role",
			types.MustDomainScope("governance"),
			action.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Recall: %v", err)
		}

		ancestors := env.ancestors(result.Revocation.ID(), 15)
		if !containsEvent(ancestors, result.Audit.ID()) {
			t.Error("revocation should trace to audit")
		}
		if !containsEvent(ancestors, result.Filing.ID()) {
			t.Error("revocation should trace to filing")
		}
		if !containsEvent(ancestors, result.Consent.ID()) {
			t.Error("revocation should trace to consent")
		}
		env.verifyChain()
	})
}
