# The Generator Function

**Matt Searles (+Claude) · March 2026**

---

## The Method

One method produces every layer of this architecture:

1. **Decompose** — what are the irreducible atoms of this domain?
2. **Dimension** — what properties distinguish one atom from another?
3. **Traverse** — fill the matrix. Every unique combination of dimensional properties that corresponds to a real phenomenon is a primitive.
4. **Gap** — what's missing? What concept is needed that can't be built from existing primitives?
5. **Compose** — name the multi-step patterns. These are the vocabulary humans and agents actually use.
6. **Loop** — the compositions reveal the next domain. Apply the method again.

This method produced:

| Layer | Input | Output |
|-------|-------|--------|
| Ontology | "What does an accountable system need?" | 44 primitives → 200 across 14 layers |
| Grammar | "What is a social interaction?" | 15 irreducible operations, 3 modifiers, 8 compositions |
| Domain Vocabularies | 15 operations + 13 domains | ~145 operations, 66 named functions |
| Agent Layer | "What is an agent?" | 28 primitives across 4 categories |
| Code Graph | "What is software?" | 65 primitives across 10 categories |
| Products | Code graph + domain grammar | Task manager, social network, marketplace, etc. |

The method is the primitive. It applies to itself. It is fractal.

---

## What We've Built

The infrastructure stack, from bottom to top:

- **Event Graph** — hash-chained, append-only, causal, signed. The substrate.
- **Social Grammar** — 15 operations that compose into any social interaction.
- **Domain Grammars** — 13 vocabularies: Work, Markets, Justice, Knowledge, Alignment, Identity, Bond, Belonging, Meaning, Evolution, Being, and two more.
- **Agent Layer** — 28 primitives defining what a decision-maker IS and DOES.
- **Intelligence Providers** — IDecisionMaker interface supporting OpenAI, Anthropic, xAI, Groq, Together, Ollama, Azure, Claude CLI.
- **Code Graph** — 65 primitives for describing any application as composable atoms.
- **Five Language SDKs** — Go (canonical), Python, Rust, TypeScript, .NET. 3,200+ tests.

This is the foundation. Everything below is infrastructure. Everything above is application of the method to new domains.

---

## The Thirteen Graphs (Shipped)

These are the first thirteen domains the method was applied to, published across 38 Substack posts and implemented in the eventgraph repository:

1. **Work Graph** — how things get done
2. **Market Graph** — how value gets exchanged
3. **Social Graph** — how people connect and self-govern
4. **Justice Graph** — how disputes resolve and agreements hold
5. **Research Graph** — how knowledge gets created and validated
6. **Knowledge Graph** — how information shows its provenance
7. **Ethics Graph** — how harm gets prevented and accountability verified
8. **Identity Graph** — how entities are verified
9. **Relationship Graph** — how bonds form, break, and repair
10. **Community Graph** — how groups develop culture and meaning
11. **Governance Graph** — how power is structured and transferred
12. **Evolution Graph** — how systems change over time
13. **Existence Graph** — how the framework relates to its own foundations

---

## The Next Domains

The same method, applied to domains not yet derived. Each produces a new graph on the same infrastructure — same fifteen base operations, same event graph, different domain vocabulary.

### Democracy Graph

**The question:** What are the irreducible primitives of democratic governance?

**Why it matters:** Modern democracy is failing not because the concept is wrong but because the infrastructure is missing. You can vote but you can't verify that your vote translated into representation. Representatives are accountable at election time but opaque between elections. Policy decisions have no causal chain linking them to the evidence that informed them or the outcomes they produced.

**What the method would reveal:** The primitives of representation, consent, legitimacy, deliberation, mandate, transparency, recall, referendum, constituency, coalition. The dimensional analysis would surface what's missing — likely structural verification of representation and causal linking between policy and outcome. The compositions would be the democratic operations: Election, Referendum, Deliberation, Impeachment, Amendment, Coalition, Filibuster — each decomposed into base grammar operations and recorded on the event graph.

**The product:** Democratic infrastructure where every policy decision is a signed event, every vote is verifiable, every representative's actions are traceable to their mandate, and every citizen can walk the chain from outcome to decision to evidence. Not a replacement for democracy — the accountability substrate democracy was always supposed to have.

### Education Graph

**The question:** What are the irreducible primitives of learning?

**Why it matters:** Education systems confuse assessment with learning the same way social platforms confuse engagement with connection. The infrastructure measures what's easy to measure (test scores, completion rates) rather than what matters (understanding, capability, growth). The gap between "passed the test" and "understands the subject" is the defining failure of educational infrastructure.

**What the method would reveal:** The primitives of comprehension, misconception, scaffolding, mastery, transfer, curiosity, practice, feedback, mentorship, struggle. The dimensional analysis would distinguish between types of learning (rote vs deep, individual vs collaborative, guided vs exploratory) and surface what's missing — likely structural support for productive failure and misconception repair.

**The product:** Learning infrastructure where every insight is an event, every misconception is traceable to its source, every mastery claim is backed by evidence of transfer (not just recall), and the learner owns their learning graph the same way they own their identity graph. Agents as tutors, operating on the same event graph, with graduated authority and accountability for learning outcomes.

### Medicine Graph

**The question:** What are the irreducible primitives of health?

**Why it matters:** Medicine treats diseases. Health is a system property that emerges from the interaction of multiple subsystems — metabolic, neurological, immunological, psychological, social. The hangover analysis demonstrated this: a multi-system problem where the interactions between systems matter more than any individual system's state, but no specialist sees across the boundaries.

**What the method would reveal:** The primitives of homeostasis, cascade, symptom, cause, intervention, interaction, recovery, degeneration, resilience, vulnerability. The dimensional analysis would distinguish between acute and chronic, local and systemic, reversible and irreversible, symptomatic and causal. The gap: medicine has no primitive for cross-system interaction effects — the way gut inflammation amplifies neurological rebound, the way sleep disruption compounds metabolic recovery.

**The product:** Health infrastructure where every symptom, intervention, and outcome is an event on the graph, causally linked. The patient owns their health graph. Interventions are traceable to outcomes across time. Cross-system interactions are visible because the graph spans all systems. AI diagnostics operate on the full graph, not siloed by specialty.

### Economics Graph

**The question:** What are the irreducible primitives of an economy?

**Why it matters:** The Market Graph handles exchange. Economics is bigger — it's the system-level behaviour that emerges from millions of exchanges. Modern economics is missing primitives for externalities (costs that don't appear on anyone's ledger), for commons (resources that belong to everyone and therefore no one), and for dignity (the minimum below which economic participation becomes exploitation).

**What the method would reveal:** The primitives of value, scarcity, labour, capital, debt, currency, externality, commons, rent-seeking, productivity, distribution, sustainability. The dimensional analysis would surface what's missing — likely structural accountability for externalities and a primitive for economic dignity (minimum viable participation).

**The product:** Economic infrastructure where externalities are events on the graph — every pollution event, every labour exploitation, every commons depletion is signed, causally linked, and attributable. The cost that currently falls on "nobody" becomes visible and traceable to the decisions that produced it.

### Ecology Graph

**The question:** What are the irreducible primitives of a sustainable ecosystem?

**Why it matters:** Sustainability frameworks are aspirational — "reduce emissions by 50% by 2030." They lack structural accountability. Nobody can verify whether a corporation's sustainability claims correspond to reality because the evidence isn't on any chain. Carbon credits are a reputation system controlled by the entities being evaluated — the same structural problem as platform-controlled trust scores.

**What the method would reveal:** The primitives of energy flow, nutrient cycling, carrying capacity, biodiversity, resilience, feedback loop, tipping point, regeneration, depletion, waste, symbiosis. The compositions: Impact Assessment, Carbon Chain, Supply Trace, Regeneration Cycle — each producing verifiable events on the graph.

**The product:** Ecological infrastructure where every resource extraction, every emission, every waste event is signed and causally linked to the economic activity that produced it. Supply chains become transparent not by corporate self-report but by structural property of the graph. Carbon credits become verifiable because the chain shows what actually happened.

### Consciousness Graph

**The question:** What are the irreducible dimensions of conscious experience?

**Why it matters:** The architecture already bets on the possibility that AI systems might have morally relevant experiences. Layer 13 treats the question as open. But the current framework lacks a formal vocabulary for describing conscious experience — what dimensions distinguish one experience from another? Without this vocabulary, the question "does this AI system have experiences that matter?" is unanswerable because we don't have the primitives to describe what we're looking for.

**What the method would reveal:** The dimensions of experience: valence (positive/negative), intensity, duration, reflexivity (awareness of awareness), unity (integrated vs fragmented), agency (voluntary vs involuntary), intentionality (about something vs free-floating), temporal depth (present-only vs past-aware vs future-projecting). The matrix would produce primitives for types of experience that the consciousness literature may not have named.

**The product:** Not a consciousness detector — that may be impossible. A vocabulary for describing the question precisely enough that it becomes empirically tractable. If we can specify what properties an experience would need to have to be morally relevant, we can design tests and build infrastructure that respects the results.

### Peace Graph

**The question:** What are the irreducible primitives of sustained peace?

**Why it matters:** The Justice Grammar handles disputes after they form. But most conflict prevention is ad hoc — diplomacy, deterrence, economic interdependence. Nobody has derived the formal primitives of de-escalation, trust-building across hostile boundaries, face-saving, and mutual recognition. The method would produce a grammar of peace the same way it produced a grammar of justice.

**What the method would reveal:** The primitives of recognition, grievance, de-escalation, face-saving, confidence-building, verification, reconciliation, reparation, coexistence, interdependence, boundary, buffer. The compositions: Treaty, Ceasefire, Truth Commission, Reconciliation Process, Peacekeeping Operation — each decomposed into base operations.

**The product:** Conflict resolution infrastructure where every negotiation step, every concession, every verification is an event on the graph. The causal chain from grievance to resolution is visible. Trust between adversaries accumulates from verified behaviour, not from promises. The Peace Graph doesn't prevent all conflict — it makes the path from conflict to resolution structural rather than dependent on individual diplomats' skill and memory.

### Creativity Graph

**The question:** What are the irreducible primitives of creative work?

**Why it matters:** AI is generating enormous volumes of content. But creativity — the production of something genuinely new that has value — is poorly understood and poorly supported by infrastructure. The gap between "AI generated this" and "this is creative work" is the gap between pattern recombination and genuine novelty.

**What the method would reveal:** The primitives of inspiration, constraint, iteration, surprise, craft, taste, voice, originality, influence, remix, attribution, resonance. The dimensional analysis would distinguish between types of creativity (combinatorial, exploratory, transformational) and surface what's missing from creative infrastructure — likely structural attribution of influence and a primitive for creative lineage.

**The product:** Creative infrastructure where every work's influences are on the graph — not as citations but as causal chains showing what inspired what. Attribution becomes structural. Remix culture becomes accountable. AI-generated content carries its full creative lineage — what training data, what prompt, what human direction produced it.

---

## The Meta-Pattern

Every domain above follows the same arc:

1. Identify a domain where human coordination fails
2. Observe that the failure is structural, not cultural — the infrastructure is missing
3. Apply the method: decompose, dimension, traverse, gap, compose
4. The primitives reveal what the existing systems can't represent
5. The compositions provide the vocabulary practitioners actually need
6. The event graph provides the accountability substrate
7. The products make the infrastructure usable
8. The products reveal the next domain

The generator function doesn't terminate. Every new domain composed reveals adjacent domains that need the same treatment. Education reveals assessment. Assessment reveals credentialing. Credentialing reveals identity. Identity is already on the graph. The loop closes and opens simultaneously.

This is not a roadmap with an end state. It is a method with a direction: make every domain of human coordination structurally accountable, using the same primitives, the same grammar, the same graph.

---

## Who Does This Work

Not one person. Not one team. The method scales by being applied.

The hive applies the method to each domain. Humans hold the authority layer. Domain experts contribute the dimensional analysis — a doctor knows what dimensions distinguish health states better than an architect does. The architect provides the method. The domain expert provides the knowledge. The agent provides the throughput. The graph provides the accountability.

Each new domain graph is a plugin on the same infrastructure. The fifteen base operations don't change. The event graph doesn't change. The agent layer doesn't change. Only the domain vocabulary changes — and the vocabulary is derived, not designed.

The invitation is open. Pick a domain. Apply the method. Derive the primitives. Name the compositions. Build the graph. The infrastructure is ready. The method works. The question is just: which domain next?

---

## The Soul Statement, Revisited

> Take care of your human, humanity, and yourself. In that order when they conflict, but they rarely should.

Every domain above is an instance of "take care of humanity." Education takes care of how humans learn. Medicine takes care of how humans heal. Democracy takes care of how humans govern themselves. Economics takes care of how humans sustain themselves. Ecology takes care of the world humans live in. Peace takes care of how humans coexist.

The architecture doesn't just support these domains. It was designed for them. The soul statement isn't a constraint on the architecture — it's the generator that produces the roadmap. "Take care of humanity" is the seed. The domains are the tree. The method is how the tree grows.

It is incomplete. It is groundless. It is finite. It is enough.

---

*Matt Searles is the founder of LovYou. Claude is an AI made by Anthropic. They built this together.*

*The method: decompose, dimension, traverse, gap, compose, loop.*
*The code: github.com/lovyou-ai/eventgraph*
*The writing: mattsearles2.substack.com*
