package event

import "fmt"

const (
	ProductionEvidenceWriteActionClass   = "eventgraph.production_evidence.write"
	ProductionEvidenceWriteSchemaVersion = "production_evidence_write_v0"
	ProductionEvidenceWriteStoreName     = "production_eventgraph_evidence"
	ProductionEvidenceWriteSubjectType   = "eventgraph_write_path"
)

type ProductionEvidenceWriteCandidate struct {
	CandidateID   string
	Repo          string
	ActorRef      string
	ActionClass   string
	SchemaVersion string
	// SourceStateRef and CurrentStateRef are caller-supplied local snapshot refs.
	// They prove candidate self-consistency only; live-head freshness requires a
	// separate authority packet or adapter.
	SourceStateRef      string
	CurrentStateRef     string
	SourceIssueRefs     []NativeEvidenceIssueRef
	RuntimeEvidenceRefs []string
	IssueEvidenceRefs   []string
	AuthorityDecision   AuthorityDecisionRecordedContent
	StoreGovernance     AuthorityStoreGovernanceRecordedContent
	TestRuns            []NativeTestRunRecordedContent
	GateResults         []NativeGateResultRecordedContent
	AuditReports        []NativeAuditReportRecordedContent
}

type ProductionEvidenceWritePlan struct {
	CandidateID         string
	Repo                string
	ActorRef            string
	ActionClass         string
	SchemaVersion       string
	SourceStateRef      string
	CurrentStateRef     string
	RuntimeEvidenceRefs []string
	IssueEvidenceRefs   []string
	Entries             []ProductionEvidenceWriteEntry
}

type ProductionEvidenceWriteEntry struct {
	EventType  string
	EvidenceID string
	Content    EventContent
}

func PlanProductionEvidenceWrite(candidate ProductionEvidenceWriteCandidate) (ProductionEvidenceWritePlan, error) {
	if err := validateProductionEvidenceWriteCandidate(candidate); err != nil {
		return ProductionEvidenceWritePlan{}, err
	}

	contents := make([]EventContent, 0, 2+len(candidate.TestRuns)+len(candidate.GateResults)+len(candidate.AuditReports))
	contents = append(contents, candidate.AuthorityDecision, candidate.StoreGovernance)
	for _, testRun := range candidate.TestRuns {
		contents = append(contents, testRun)
	}
	for _, gateResult := range candidate.GateResults {
		contents = append(contents, gateResult)
	}
	for _, auditReport := range candidate.AuditReports {
		contents = append(contents, auditReport)
	}

	entries := make([]ProductionEvidenceWriteEntry, 0, len(contents))
	for _, content := range contents {
		entry, err := productionEvidenceWriteEntryForContent(content)
		if err != nil {
			return ProductionEvidenceWritePlan{}, err
		}
		entries = append(entries, entry)
	}
	if err := validateProductionEvidenceWriteEntries(entries); err != nil {
		return ProductionEvidenceWritePlan{}, err
	}

	return ProductionEvidenceWritePlan{
		CandidateID:         candidate.CandidateID,
		Repo:                candidate.Repo,
		ActorRef:            candidate.ActorRef,
		ActionClass:         candidate.ActionClass,
		SchemaVersion:       candidate.SchemaVersion,
		SourceStateRef:      candidate.SourceStateRef,
		CurrentStateRef:     candidate.CurrentStateRef,
		RuntimeEvidenceRefs: append([]string(nil), candidate.RuntimeEvidenceRefs...),
		IssueEvidenceRefs:   append([]string(nil), candidate.IssueEvidenceRefs...),
		Entries:             entries,
	}, nil
}

// ApplyProductionEvidenceWritePlanToFixture is an in-memory replay helper for
// tests and fixtures. Production writers must validate a fresh
// ProductionEvidenceWriteCandidate instead of trusting a pre-built plan.
func ApplyProductionEvidenceWritePlanToFixture(existing []ProductionEvidenceWriteEntry, plan ProductionEvidenceWritePlan) ([]ProductionEvidenceWriteEntry, error) {
	if err := validateProductionEvidenceWritePlan(plan); err != nil {
		return nil, err
	}
	if err := validateProductionEvidenceWriteEntries(existing); err != nil {
		return nil, fmt.Errorf("existing entries: %w", err)
	}

	applied := make([]ProductionEvidenceWriteEntry, 0, len(existing)+len(plan.Entries))
	applied = append(applied, existing...)

	seen := map[string]ProductionEvidenceWriteEntry{}
	for _, entry := range existing {
		seen[productionEvidenceWriteEntryKey(entry)] = entry
	}

	for _, entry := range plan.Entries {
		key := productionEvidenceWriteEntryKey(entry)
		if prior, ok := seen[key]; ok {
			if canonicalContentJSON(prior.Content) != canonicalContentJSON(entry.Content) {
				return nil, fmt.Errorf("conflicting evidence for %s", key)
			}
			continue
		}
		seen[key] = entry
		applied = append(applied, entry)
	}

	return applied, nil
}

func validateProductionEvidenceWriteCandidate(candidate ProductionEvidenceWriteCandidate) error {
	if candidate.CandidateID == "" {
		return fmt.Errorf("candidate_id is required")
	}
	if candidate.Repo == "" {
		return fmt.Errorf("repo is required")
	}
	if candidate.ActorRef == "" {
		return fmt.Errorf("actor_ref is required")
	}
	if candidate.ActionClass != ProductionEvidenceWriteActionClass {
		return fmt.Errorf("wrong action class %q", candidate.ActionClass)
	}
	if candidate.SchemaVersion != ProductionEvidenceWriteSchemaVersion {
		return fmt.Errorf("invalid schema version %q", candidate.SchemaVersion)
	}
	if candidate.SourceStateRef == "" {
		return fmt.Errorf("source_state_ref is required")
	}
	if candidate.CurrentStateRef == "" {
		return fmt.Errorf("current_state_ref is required")
	}
	if candidate.SourceStateRef != candidate.CurrentStateRef {
		return fmt.Errorf("stale source: source_state_ref %q does not match current_state_ref %q", candidate.SourceStateRef, candidate.CurrentStateRef)
	}
	if len(candidate.SourceIssueRefs) == 0 {
		return fmt.Errorf("at least one source_issue_ref is required")
	}
	if len(candidate.RuntimeEvidenceRefs) == 0 {
		return fmt.Errorf("at least one runtime_evidence_ref is required")
	}
	if len(candidate.IssueEvidenceRefs) == 0 {
		return fmt.Errorf("at least one issue_evidence_ref is required")
	}
	if err := validateNativeEvidenceIssueRefs(candidate.SourceIssueRefs); err != nil {
		return err
	}
	if err := validateProductionEvidenceIssueRefs(candidate.Repo, candidate.SourceIssueRefs); err != nil {
		return err
	}
	if err := validateProductionEvidenceAuthority(candidate); err != nil {
		return err
	}
	if err := validateProductionEvidenceNativeEvidence(candidate); err != nil {
		return err
	}
	return nil
}

func validateProductionEvidenceWritePlan(plan ProductionEvidenceWritePlan) error {
	if plan.CandidateID == "" {
		return fmt.Errorf("candidate_id is required")
	}
	if plan.Repo == "" {
		return fmt.Errorf("repo is required")
	}
	if plan.ActorRef == "" {
		return fmt.Errorf("actor_ref is required")
	}
	if plan.ActionClass != ProductionEvidenceWriteActionClass {
		return fmt.Errorf("wrong action class %q", plan.ActionClass)
	}
	if plan.SchemaVersion != ProductionEvidenceWriteSchemaVersion {
		return fmt.Errorf("invalid schema version %q", plan.SchemaVersion)
	}
	if plan.SourceStateRef == "" || plan.CurrentStateRef == "" {
		return fmt.Errorf("source and current state refs are required")
	}
	if plan.SourceStateRef != plan.CurrentStateRef {
		return fmt.Errorf("stale source: source_state_ref %q does not match current_state_ref %q", plan.SourceStateRef, plan.CurrentStateRef)
	}
	if len(plan.Entries) == 0 {
		return fmt.Errorf("at least one production evidence write entry is required")
	}
	if len(plan.RuntimeEvidenceRefs) == 0 {
		return fmt.Errorf("at least one runtime_evidence_ref is required")
	}
	if len(plan.IssueEvidenceRefs) == 0 {
		return fmt.Errorf("at least one issue_evidence_ref is required")
	}
	return validateProductionEvidenceWriteEntries(plan.Entries)
}

func validateProductionEvidenceAuthority(candidate ProductionEvidenceWriteCandidate) error {
	decision := candidate.AuthorityDecision
	if err := validateAuthorityDecisionRecorded(decision); err != nil {
		return fmt.Errorf("authority decision: %w", err)
	}
	if decision.SubjectType != ProductionEvidenceWriteSubjectType {
		return fmt.Errorf("authority decision subject_type must be %q", ProductionEvidenceWriteSubjectType)
	}
	if decision.SubjectRef != candidate.Repo {
		return fmt.Errorf("authority decision subject_ref %q does not match repo %q", decision.SubjectRef, candidate.Repo)
	}
	if decision.Outcome != AuthorityDecisionOutcomeAutonomous {
		return fmt.Errorf("authority decision outcome must be %q", AuthorityDecisionOutcomeAutonomous)
	}
	if decision.ActorRef != candidate.ActorRef {
		return fmt.Errorf("authority decision actor_ref %q does not match actor_ref %q", decision.ActorRef, candidate.ActorRef)
	}
	if !productionEvidenceIssueRefsMatch(candidate.SourceIssueRefs, decision.SourceIssueRefs) {
		return fmt.Errorf("authority decision source_issue_refs do not match candidate")
	}

	governance := candidate.StoreGovernance
	if err := validateAuthorityStoreGovernanceRecorded(governance); err != nil {
		return fmt.Errorf("store governance: %w", err)
	}
	if governance.StoreName != ProductionEvidenceWriteStoreName {
		return fmt.Errorf("store governance store_name must be %q", ProductionEvidenceWriteStoreName)
	}
	if governance.SchemaVersion != candidate.SchemaVersion {
		return fmt.Errorf("store governance schema_version %q does not match candidate schema_version %q", governance.SchemaVersion, candidate.SchemaVersion)
	}
	if governance.WriteStatus != AuthorityStoreWriteStatusWritePathAuthorized {
		return fmt.Errorf("store governance write_status must be %q", AuthorityStoreWriteStatusWritePathAuthorized)
	}
	if !containsString(governance.AuthorityRefs, decision.DecisionID) {
		return fmt.Errorf("store governance authority_refs must include authority decision %q", decision.DecisionID)
	}
	if !productionEvidenceIssueRefsMatch(candidate.SourceIssueRefs, governance.SourceIssueRefs) {
		return fmt.Errorf("store governance source_issue_refs do not match candidate")
	}
	return nil
}

func validateProductionEvidenceNativeEvidence(candidate ProductionEvidenceWriteCandidate) error {
	if len(candidate.TestRuns) == 0 {
		return fmt.Errorf("at least one test run is required")
	}
	if len(candidate.GateResults) == 0 {
		return fmt.Errorf("at least one gate result is required")
	}
	if len(candidate.AuditReports) == 0 {
		return fmt.Errorf("at least one audit report is required")
	}

	testRunIDs := map[string]bool{}
	for i, testRun := range candidate.TestRuns {
		if err := validateNativeTestRunRecorded(testRun); err != nil {
			return fmt.Errorf("test_runs[%d]: %w", i, err)
		}
		if !productionEvidenceIssueRefsMatch(candidate.SourceIssueRefs, testRun.SourceIssueRefs) {
			return fmt.Errorf("test_runs[%d] source_issue_refs do not match candidate", i)
		}
		testRunIDs[testRun.TestRunID] = true
	}

	gateIDs := map[string]bool{}
	for i, gate := range candidate.GateResults {
		if err := validateNativeGateResultRecorded(gate); err != nil {
			return fmt.Errorf("gate_results[%d]: %w", i, err)
		}
		if !productionEvidenceIssueRefsMatch(candidate.SourceIssueRefs, gate.SourceIssueRefs) {
			return fmt.Errorf("gate_results[%d] source_issue_refs do not match candidate", i)
		}
		if !anyStringInSet(gate.EvidenceRefs, testRunIDs) {
			return fmt.Errorf("gate_results[%d] evidence_refs must include a test run", i)
		}
		gateIDs[gate.GateResultID] = true
	}

	for i, audit := range candidate.AuditReports {
		if err := validateNativeAuditReportRecorded(audit); err != nil {
			return fmt.Errorf("audit_reports[%d]: %w", i, err)
		}
		if !productionEvidenceIssueRefsMatch(candidate.SourceIssueRefs, audit.SourceIssueRefs) {
			return fmt.Errorf("audit_reports[%d] source_issue_refs do not match candidate", i)
		}
		if len(audit.ValidationRefs) == 0 {
			return fmt.Errorf("audit_reports[%d] validation_refs are required", i)
		}
		if len(audit.CFARRefs) == 0 {
			return fmt.Errorf("audit_reports[%d] cfar_refs are required", i)
		}
		if len(audit.AuthorityBoundaryRefs) == 0 {
			return fmt.Errorf("audit_reports[%d] authority_boundary_refs are required", i)
		}
		if !anyStringInSet(audit.EvidenceRefs, testRunIDs) {
			return fmt.Errorf("audit_reports[%d] evidence_refs must include a test run", i)
		}
		if !anyStringInSet(audit.EvidenceRefs, gateIDs) {
			return fmt.Errorf("audit_reports[%d] evidence_refs must include a gate result", i)
		}
		if audit.Outcome == NativeEvidenceOutcomePassed {
			if len(audit.MissingLinks) != 0 {
				return fmt.Errorf("audit_reports[%d] passed audit must not have missing_links", i)
			}
			if audit.TraceScoreBasisPoints == nil || *audit.TraceScoreBasisPoints != 10000 {
				return fmt.Errorf("audit_reports[%d] passed audit must have trace_score_basis_points 10000", i)
			}
		}
	}

	return nil
}

func validateProductionEvidenceIssueRefs(repo string, refs []NativeEvidenceIssueRef) error {
	seen := map[string]bool{}
	for i, ref := range refs {
		if ref.Repo != repo {
			return fmt.Errorf("source_issue_refs[%d].repo %q does not match repo %q", i, ref.Repo, repo)
		}
		key := productionEvidenceIssueRefKey(ref)
		if seen[key] {
			return fmt.Errorf("duplicate source issue ref %s", key)
		}
		seen[key] = true
	}
	return nil
}

func productionEvidenceIssueRefsMatch(expected []NativeEvidenceIssueRef, got []NativeEvidenceIssueRef) bool {
	expectedCounts := map[string]int{}
	for _, ref := range expected {
		expectedCounts[productionEvidenceIssueRefKey(ref)]++
	}
	gotCounts := map[string]int{}
	for _, ref := range got {
		gotCounts[productionEvidenceIssueRefKey(ref)]++
	}
	return productionEvidenceIssueRefCountsEqual(expectedCounts, gotCounts)
}

func productionEvidenceIssueRefKey(ref NativeEvidenceIssueRef) string {
	return fmt.Sprintf("%s#%d", ref.Repo, ref.Number)
}

func productionEvidenceIssueRefCountsEqual(a map[string]int, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for key, count := range a {
		if b[key] != count {
			return false
		}
	}
	return true
}

func productionEvidenceWriteEntryForContent(content EventContent) (ProductionEvidenceWriteEntry, error) {
	if content == nil {
		return ProductionEvidenceWriteEntry{}, fmt.Errorf("content is required")
	}
	switch c := content.(type) {
	case AuthorityDecisionRecordedContent:
		if err := validateAuthorityDecisionRecorded(c); err != nil {
			return ProductionEvidenceWriteEntry{}, err
		}
		return ProductionEvidenceWriteEntry{
			EventType:  c.EventTypeName(),
			EvidenceID: c.DecisionID,
			Content:    c,
		}, nil
	case AuthorityStoreGovernanceRecordedContent:
		if err := validateAuthorityStoreGovernanceRecorded(c); err != nil {
			return ProductionEvidenceWriteEntry{}, err
		}
		return ProductionEvidenceWriteEntry{
			EventType:  c.EventTypeName(),
			EvidenceID: c.GovernanceID,
			Content:    c,
		}, nil
	case NativeTestRunRecordedContent:
		if err := validateNativeTestRunRecorded(c); err != nil {
			return ProductionEvidenceWriteEntry{}, err
		}
		return ProductionEvidenceWriteEntry{
			EventType:  c.EventTypeName(),
			EvidenceID: c.TestRunID,
			Content:    c,
		}, nil
	case NativeGateResultRecordedContent:
		if err := validateNativeGateResultRecorded(c); err != nil {
			return ProductionEvidenceWriteEntry{}, err
		}
		return ProductionEvidenceWriteEntry{
			EventType:  c.EventTypeName(),
			EvidenceID: c.GateResultID,
			Content:    c,
		}, nil
	case NativeAuditReportRecordedContent:
		if err := validateNativeAuditReportRecorded(c); err != nil {
			return ProductionEvidenceWriteEntry{}, err
		}
		return ProductionEvidenceWriteEntry{
			EventType:  c.EventTypeName(),
			EvidenceID: c.AuditReportID,
			Content:    c,
		}, nil
	default:
		return ProductionEvidenceWriteEntry{}, fmt.Errorf("unsupported production evidence content %T", content)
	}
}

func validateProductionEvidenceWriteEntries(entries []ProductionEvidenceWriteEntry) error {
	seen := map[string]ProductionEvidenceWriteEntry{}
	for i, entry := range entries {
		if entry.EventType == "" {
			return fmt.Errorf("entries[%d].event_type is required", i)
		}
		if entry.EvidenceID == "" {
			return fmt.Errorf("entries[%d].evidence_id is required", i)
		}
		if entry.Content == nil {
			return fmt.Errorf("entries[%d].content is required", i)
		}
		if entry.Content.EventTypeName() != entry.EventType {
			return fmt.Errorf("entries[%d].event_type %q does not match content event type %q", i, entry.EventType, entry.Content.EventTypeName())
		}
		derived, err := productionEvidenceWriteEntryForContent(entry.Content)
		if err != nil {
			return fmt.Errorf("entries[%d]: %w", i, err)
		}
		if derived.EvidenceID != entry.EvidenceID {
			return fmt.Errorf("entries[%d].evidence_id %q does not match content evidence id %q", i, entry.EvidenceID, derived.EvidenceID)
		}
		key := productionEvidenceWriteEntryKey(entry)
		if prior, ok := seen[key]; ok {
			if canonicalContentJSON(prior.Content) != canonicalContentJSON(entry.Content) {
				return fmt.Errorf("conflicting evidence for %s", key)
			}
			return fmt.Errorf("duplicate evidence for %s", key)
		}
		seen[key] = entry
	}
	return nil
}

func productionEvidenceWriteEntryKey(entry ProductionEvidenceWriteEntry) string {
	return entry.EventType + "/" + entry.EvidenceID
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func anyStringInSet(values []string, set map[string]bool) bool {
	for _, value := range values {
		if set[value] {
			return true
		}
	}
	return false
}
