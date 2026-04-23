package compositions_test

import (
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/compositions"
	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/types"
)

func TestBelongingGrammar(t *testing.T) {
	t.Run("SettleAndContribute", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		community := env.actor("Community", 1, event.ActorTypeCommittee)
		newcomer := env.actor("Newcomer", 2, event.ActorTypeHuman)

		settle, _ := belonging.Settle(env.ctx, newcomer.ID(), community.ID(),
			types.Some(types.MustDomainScope("community_alpha")),
			env.boot.ID(), env.convID, signer)
		contrib, _ := belonging.Contribute(env.ctx, newcomer.ID(),
			"first code review submitted",
			[]types.EventID{settle.ID()}, env.convID, signer)

		ancestors := env.ancestors(contrib.ID(), 5)
		if !containsEvent(ancestors, settle.ID()) {
			t.Error("contribution should trace to settlement")
		}
		env.verifyChain()
	})

	t.Run("IncludeAndPractice", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		sponsor := env.actor("Sponsor", 1, event.ActorTypeHuman)
		newcomer := env.actor("Newcomer", 2, event.ActorTypeHuman)

		inclusion, _ := belonging.Include(env.ctx, sponsor.ID(),
			"mentorship pairing for newcomer onboarding",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		practice, _ := belonging.Practice(env.ctx, newcomer.ID(),
			"attended Friday demo day",
			[]types.EventID{inclusion.ID()}, env.convID, signer)

		ancestors := env.ancestors(practice.ID(), 5)
		if !containsEvent(ancestors, inclusion.ID()) {
			t.Error("practice should trace to inclusion")
		}
		env.verifyChain()
	})

	t.Run("StewardAndSustain", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)
		lead := env.actor("Lead", 2, event.ActorTypeHuman)

		stewardship, _ := belonging.Steward(env.ctx, admin.ID(), lead.ID(),
			types.MustDomainScope("backend_team"), types.MustWeight(0.7),
			env.boot.ID(), env.convID, signer)
		sustain, _ := belonging.Sustain(env.ctx, lead.ID(),
			"team healthy: 12 active members, 3 new joins, low turnover",
			[]types.EventID{stewardship.ID()}, env.convID, signer)

		ancestors := env.ancestors(sustain.ID(), 5)
		if !containsEvent(ancestors, stewardship.ID()) {
			t.Error("sustainability should trace to stewardship")
		}
		env.verifyChain()
	})

	t.Run("CelebrateAndTell", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)

		milestone, _ := belonging.Contribute(env.ctx, alice.ID(),
			"community reached 1000 contributions",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		celebration, _ := belonging.Celebrate(env.ctx, alice.ID(),
			"1000 contributions milestone — thank you all",
			[]types.EventID{milestone.ID()}, env.convID, signer)
		story, _ := belonging.Tell(env.ctx, alice.ID(),
			"started with 3 people, now 50 active contributors in 6 months",
			[]types.EventID{celebration.ID()}, env.convID, signer)

		ancestors := env.ancestors(story.ID(), 10)
		if !containsEvent(ancestors, milestone.ID()) {
			t.Error("story should trace to milestone")
		}
		env.verifyChain()
	})

	t.Run("GiftAndPassOn", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		elder := env.actor("Elder", 1, event.ActorTypeHuman)
		newcomer := env.actor("Newcomer", 2, event.ActorTypeHuman)

		gift, _ := belonging.Gift(env.ctx, elder.ID(),
			"shared private mentorship notes from 10 years of experience",
			[]types.EventID{env.boot.ID()}, env.convID, signer)
		passOn, _ := belonging.PassOn(env.ctx, elder.ID(), newcomer.ID(),
			types.MustDomainScope("mentorship"), "stewardship of the mentorship program",
			gift.ID(), env.convID, signer)

		ancestors := env.ancestors(passOn.ID(), 5)
		if !containsEvent(ancestors, gift.ID()) {
			t.Error("pass-on should trace to gift")
		}
		env.verifyChain()
	})

	t.Run("Festival", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)

		milestone, _ := belonging.Contribute(env.ctx, alice.ID(),
			"community reached 500 active members",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := belonging.Festival(env.ctx, alice.ID(),
			"celebrating 500 members milestone",
			"annual community gathering tradition",
			"from a small group to a thriving community in one year",
			"open source toolkit for new contributors",
			[]types.EventID{milestone.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Festival: %v", err)
		}

		ancestors := env.ancestors(result.Gift.ID(), 10)
		if !containsEvent(ancestors, result.Celebration.ID()) {
			t.Error("gift should trace to celebration")
		}
		if !containsEvent(ancestors, result.Practice.ID()) {
			t.Error("gift should trace to practice")
		}
		if !containsEvent(ancestors, result.Story.ID()) {
			t.Error("gift should trace to story")
		}
		if !containsEvent(ancestors, milestone.ID()) {
			t.Error("gift should trace to milestone")
		}
		env.verifyChain()
	})

	t.Run("CommonsGovernance", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		admin := env.actor("Admin", 1, event.ActorTypeHuman)
		steward := env.actor("Steward", 2, event.ActorTypeHuman)

		result, err := belonging.CommonsGovernance(env.ctx, admin.ID(), steward.ID(),
			types.MustDomainScope("infrastructure"),
			types.MustWeight(0.8),
			"infrastructure healthy, 99.9 percent uptime last quarter",
			"all deployments require two approvals",
			"audit complete: no policy violations found",
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("CommonsGovernance: %v", err)
		}

		ancestors := env.ancestors(result.Audit.ID(), 10)
		if !containsEvent(ancestors, result.Stewardship.ID()) {
			t.Error("audit should trace to stewardship")
		}
		if !containsEvent(ancestors, result.Assessment.ID()) {
			t.Error("audit should trace to assessment")
		}
		if !containsEvent(ancestors, result.Legislation.ID()) {
			t.Error("audit should trace to legislation")
		}
		env.verifyChain()
	})

	t.Run("Renewal", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		alice := env.actor("Alice", 1, event.ActorTypeHuman)

		practice, _ := belonging.Practice(env.ctx, alice.ID(),
			"weekly retrospectives every Friday",
			[]types.EventID{env.boot.ID()}, env.convID, signer)

		result, err := belonging.Renewal(env.ctx, alice.ID(),
			"community still vibrant, practices need refreshing",
			"moved retrospectives to async format for global members",
			"evolved from co-located team to distributed community",
			[]types.EventID{practice.ID()}, env.convID, signer)
		if err != nil {
			t.Fatalf("Renewal: %v", err)
		}

		ancestors := env.ancestors(result.Story.ID(), 10)
		if !containsEvent(ancestors, result.Assessment.ID()) {
			t.Error("story should trace to assessment")
		}
		if !containsEvent(ancestors, result.Practice.ID()) {
			t.Error("story should trace to practice")
		}
		if !containsEvent(ancestors, practice.ID()) {
			t.Error("story should trace to original practice event")
		}
		env.verifyChain()
	})

	t.Run("Onboard", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		sponsor := env.actor("Sponsor", 1, event.ActorTypeHuman)
		newcomer := env.actor("Newcomer", 2, event.ActorTypeHuman)
		community := env.actor("Community", 3, event.ActorTypeCommittee)

		result, err := belonging.Onboard(env.ctx, sponsor.ID(), newcomer.ID(), community.ID(),
			types.Some(types.MustDomainScope("community_alpha")),
			"mentorship pairing arranged",
			"attended Friday demo day",
			"first code review submitted",
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Onboard: %v", err)
		}

		ancestors := env.ancestors(result.Contribution.ID(), 10)
		if !containsEvent(ancestors, result.Inclusion.ID()) {
			t.Error("contribution should trace to inclusion")
		}
		env.verifyChain()
	})

	t.Run("Succession", func(t *testing.T) {
		env := newTestEnv(t)
		belonging := compositions.NewBelongingGrammar(env.grammar)
		outgoing := env.actor("Outgoing", 1, event.ActorTypeHuman)
		incoming := env.actor("Incoming", 2, event.ActorTypeHuman)

		result, err := belonging.Succession(env.ctx, outgoing.ID(), incoming.ID(),
			"community healthy, ready for transition",
			types.MustDomainScope("governance"),
			"celebrating 3 years of stewardship",
			"started with 5 members, grew to 200, survived 2 governance crises",
			env.boot.ID(), env.convID, signer)
		if err != nil {
			t.Fatalf("Succession: %v", err)
		}

		ancestors := env.ancestors(result.Story.ID(), 10)
		if !containsEvent(ancestors, result.Assessment.ID()) {
			t.Error("story should trace to assessment")
		}
		env.verifyChain()
	})
}
