package v39

import (
	"errors"
	"fmt"
)

const (
	EdgeDerivedFrom         = "DERIVED_FROM"
	EdgeRequires            = "REQUIRES"
	EdgeDecomposedInto      = "DECOMPOSED_INTO"
	EdgeAssignedTo          = "ASSIGNED_TO"
	EdgeInvoked             = "INVOKED"
	EdgeUsedEnvelope        = "USED_ENVELOPE"
	EdgeProduced            = "PRODUCED"
	EdgeModified            = "MODIFIED"
	EdgeImplements          = "IMPLEMENTS"
	EdgeVerifies            = "VERIFIES"
	EdgeFailedBy            = "FAILED_BY"
	EdgeRepairedBy          = "REPAIRED_BY"
	EdgeWaivedBy            = "WAIVED_BY"
	EdgeCertifiedBy         = "CERTIFIED_BY"
	EdgePackagedAs          = "PACKAGED_AS"
	EdgeAuditedBy           = "AUDITED_BY"
	EdgeRequestedAuthority  = "REQUESTED_AUTHORITY"
	EdgeDecidedBy           = "DECIDED_BY"
	EdgeReceiptedBy         = "RECEIPTED_BY"
	EdgeTransitionedBy      = "TRANSITIONED_BY"
	EdgeObservedFailure     = "OBSERVED_FAILURE"
	EdgeSupersedes          = "SUPERSEDES"
	EdgeUsedCapability      = "USED_CAPABILITY"
	EdgeReferencedMemory    = "REFERENCED_MEMORY"
	EdgeReferencedKnowledge = "REFERENCED_KNOWLEDGE"
	EdgeApprovedBy          = "APPROVED_BY"
)

var ErrRequiredPathMissing = errors.New("dark factory v3.9 required path missing")

type RequiredPath struct {
	Name      string   `json:"name"`
	NodeIDs   []string `json:"node_ids"`
	EdgeIDs   []string `json:"edge_ids,omitempty"`
	Missing   []string `json:"missing,omitempty"`
	Completed bool     `json:"completed"`
}

func (p RequiredPath) Err() error {
	if p.Completed {
		return nil
	}
	return fmt.Errorf("%w: %s: %v", ErrRequiredPathMissing, p.Name, p.Missing)
}

func (s *InMemoryStore) QueryRequiredPath(startID string, edgeTypes ...string) (RequiredPath, error) {
	path := RequiredPath{Name: "query_required_path", NodeIDs: []string{startID}}
	current := startID
	for _, edgeType := range edgeTypes {
		edge, ok := s.firstOutgoingEdge(current, edgeType)
		if !ok {
			path.Missing = append(path.Missing, fmt.Sprintf("%s from %s", edgeType, current))
			return path, path.Err()
		}
		path.EdgeIDs = append(path.EdgeIDs, edge.ID)
		path.NodeIDs = append(path.NodeIDs, edge.ToID)
		current = edge.ToID
	}
	path.Completed = true
	return path, nil
}

func (s *InMemoryStore) FactoryOrderRequirementAcceptanceTask(factoryOrderID string) (RequiredPath, error) {
	path := RequiredPath{Name: "FactoryOrder -> Requirement -> AcceptanceCriterion -> Task", NodeIDs: []string{factoryOrderID}}
	req, ok := s.firstRequirementForFactoryOrder(factoryOrderID)
	if !ok {
		path.Missing = append(path.Missing, "Requirement for FactoryOrder "+factoryOrderID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, req.CommonNode.ID)

	ac, ok := s.firstAcceptanceCriterionForRequirement(req.CommonNode.ID)
	if !ok {
		path.Missing = append(path.Missing, "AcceptanceCriterion for Requirement "+req.CommonNode.ID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, ac.CommonNode.ID)

	edge, ok := s.firstOutgoingEdge(ac.CommonNode.ID, EdgeDecomposedInto)
	if !ok {
		edge, ok = s.firstOutgoingEdge(ac.CommonNode.ID, EdgeRequires)
	}
	if !ok {
		path.Missing = append(path.Missing, "Task edge from AcceptanceCriterion "+ac.CommonNode.ID)
		return path, path.Err()
	}
	if task, ok := s.mustGetTask(edge.ToID); !ok {
		path.Missing = append(path.Missing, "Task "+edge.ToID)
		return path, path.Err()
	} else if task.FactoryOrderID == nil || *task.FactoryOrderID != factoryOrderID {
		path.Missing = append(path.Missing, "Task "+task.CommonNode.ID+" linked to FactoryOrder "+factoryOrderID)
		return path, path.Err()
	}
	path.EdgeIDs = append(path.EdgeIDs, edge.ID)
	path.NodeIDs = append(path.NodeIDs, edge.ToID)
	path.Completed = true
	return path, nil
}

func (s *InMemoryStore) TaskRuntimeEnvelopeResult(taskID string) (RequiredPath, error) {
	path, err := s.QueryRequiredPath(taskID, EdgeUsedEnvelope, EdgeProduced)
	path.Name = "Task -> RuntimeEnvelope -> RuntimeResult"
	if err != nil {
		return path, err
	}
	if _, ok := s.mustGetRuntimeEnvelope(path.NodeIDs[1]); !ok {
		path.Completed = false
		path.Missing = append(path.Missing, "RuntimeEnvelope "+path.NodeIDs[1])
		return path, path.Err()
	}
	if _, ok := s.mustGetRuntimeResult(path.NodeIDs[2]); !ok {
		path.Completed = false
		path.Missing = append(path.Missing, "RuntimeResult "+path.NodeIDs[2])
		return path, path.Err()
	}
	return path, nil
}

func (s *InMemoryStore) TaskArtifact(taskID string) (RequiredPath, error) {
	path := RequiredPath{Name: "Task -> Artifact", NodeIDs: []string{taskID}}
	for _, r := range s.ByType(TypeArtifact) {
		artifact := r.(*Artifact)
		if artifact.TaskID != nil && *artifact.TaskID == taskID {
			path.NodeIDs = append(path.NodeIDs, artifact.CommonNode.ID)
			path.Completed = true
			return path, nil
		}
	}
	if edge, ok := s.firstOutgoingEdge(taskID, EdgeProduced); ok {
		if _, ok := s.mustGetArtifact(edge.ToID); ok {
			path.EdgeIDs = append(path.EdgeIDs, edge.ID)
			path.NodeIDs = append(path.NodeIDs, edge.ToID)
			path.Completed = true
			return path, nil
		}
	}
	path.Missing = append(path.Missing, "Artifact for Task "+taskID)
	return path, path.Err()
}

func (s *InMemoryStore) TaskTestCaseRunGateResult(taskID string) (RequiredPath, error) {
	path := RequiredPath{Name: "Task -> TestCase -> TestRun -> GateResult", NodeIDs: []string{taskID}}
	edge, ok := s.firstOutgoingEdge(taskID, EdgeVerifies)
	if !ok {
		path.Missing = append(path.Missing, "TestCase edge from Task "+taskID)
		return path, path.Err()
	}
	tc, ok := s.mustGetTestCase(edge.ToID)
	if !ok {
		path.Missing = append(path.Missing, "TestCase "+edge.ToID)
		return path, path.Err()
	}
	path.EdgeIDs = append(path.EdgeIDs, edge.ID)
	path.NodeIDs = append(path.NodeIDs, tc.CommonNode.ID)

	tr, ok := s.firstTestRunForTestCase(tc.CommonNode.ID)
	if !ok {
		path.Missing = append(path.Missing, "TestRun for TestCase "+tc.CommonNode.ID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, tr.CommonNode.ID)

	gr, ok := s.firstGateResultForEvidence(tr.CommonNode.ID)
	if !ok {
		path.Missing = append(path.Missing, "GateResult evidence for TestRun "+tr.CommonNode.ID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, gr.CommonNode.ID)
	path.Completed = true
	return path, nil
}

func (s *InMemoryStore) GateResultFailureRepairWaiver(gateResultID string) (RequiredPath, error) {
	path := RequiredPath{Name: "GateResult -> Failure / RepairAttempt / Waiver", NodeIDs: []string{gateResultID}}
	gr, ok := s.mustGetGateResult(gateResultID)
	if !ok {
		path.Missing = append(path.Missing, "GateResult "+gateResultID)
		return path, path.Err()
	}
	if gr.WaiverRef != nil && *gr.WaiverRef != "" {
		if _, ok := s.mustGetWaiver(*gr.WaiverRef); !ok {
			path.Missing = append(path.Missing, "Waiver "+*gr.WaiverRef)
			return path, path.Err()
		}
		path.NodeIDs = append(path.NodeIDs, *gr.WaiverRef)
		path.Completed = true
		return path, nil
	}
	failure, ok := s.firstFailureForGateResult(gateResultID)
	if !ok {
		path.Missing = append(path.Missing, "Failure or Waiver for GateResult "+gateResultID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, failure.CommonNode.ID)
	repair, ok := s.firstRepairAttemptForFailure(failure.CommonNode.ID)
	if !ok {
		path.Missing = append(path.Missing, "RepairAttempt for Failure "+failure.CommonNode.ID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, repair.CommonNode.ID)
	path.Completed = true
	return path, nil
}

func (s *InMemoryStore) FactoryRuntimeVersionPath(factoryOrderOrReleaseCandidateID string) (RequiredPath, error) {
	path := RequiredPath{Name: "FactoryOrder or ReleaseCandidate -> FactoryRuntimeVersion", NodeIDs: []string{factoryOrderOrReleaseCandidateID}}
	if rc, ok := s.mustGetReleaseCandidate(factoryOrderOrReleaseCandidateID); ok {
		if rc.FactoryRuntimeVersionID == nil || *rc.FactoryRuntimeVersionID == "" {
			path.Missing = append(path.Missing, "FactoryRuntimeVersion for ReleaseCandidate "+rc.CommonNode.ID)
			return path, path.Err()
		}
		if _, ok := s.mustGetFactoryRuntimeVersion(*rc.FactoryRuntimeVersionID); !ok {
			path.Missing = append(path.Missing, "FactoryRuntimeVersion "+*rc.FactoryRuntimeVersionID)
			return path, path.Err()
		}
		path.NodeIDs = append(path.NodeIDs, *rc.FactoryRuntimeVersionID)
		path.Completed = true
		return path, nil
	}
	for _, r := range s.ByType(TypeReleaseCandidate) {
		rc := r.(*ReleaseCandidate)
		if rc.FactoryOrderID == factoryOrderOrReleaseCandidateID && rc.FactoryRuntimeVersionID != nil {
			if _, ok := s.mustGetFactoryRuntimeVersion(*rc.FactoryRuntimeVersionID); ok {
				path.NodeIDs = append(path.NodeIDs, rc.CommonNode.ID, *rc.FactoryRuntimeVersionID)
				path.Completed = true
				return path, nil
			}
		}
	}
	path.Missing = append(path.Missing, "FactoryRuntimeVersion for "+factoryOrderOrReleaseCandidateID)
	return path, path.Err()
}

func (s *InMemoryStore) ReleaseCandidateCertificationOrRejection(releaseCandidateID string) (RequiredPath, error) {
	path := RequiredPath{Name: "ReleaseCandidate -> Certification or Rejection", NodeIDs: []string{releaseCandidateID}}
	for _, r := range s.ByType(TypeCertification) {
		cert := r.(*Certification)
		if cert.ReleaseCandidateID == releaseCandidateID {
			path.NodeIDs = append(path.NodeIDs, cert.CommonNode.ID)
			path.Completed = true
			return path, nil
		}
	}
	for _, r := range s.ByType(TypeRejection) {
		rejection := r.(*Rejection)
		if rejection.ReleaseCandidateID == releaseCandidateID {
			path.NodeIDs = append(path.NodeIDs, rejection.CommonNode.ID)
			path.Completed = true
			return path, nil
		}
	}
	path.Missing = append(path.Missing, "Certification or Rejection for ReleaseCandidate "+releaseCandidateID)
	return path, path.Err()
}

func (s *InMemoryStore) DecisionAuditReport(decisionID string) (RequiredPath, error) {
	path, err := s.QueryRequiredPath(decisionID, EdgeAuditedBy)
	path.Name = "Certification/Rejection -> AuditReport"
	if err != nil {
		return path, err
	}
	if _, ok := s.mustGetAuditReport(path.NodeIDs[1]); !ok {
		path.Completed = false
		path.Missing = append(path.Missing, "AuditReport "+path.NodeIDs[1])
		return path, path.Err()
	}
	return path, nil
}

func (s *InMemoryStore) AuthorityRequestDecisionReceipt(authorityRequestID string) (RequiredPath, error) {
	path := RequiredPath{Name: "AuthorityRequest -> AuthorityDecision -> ExecutionReceipt", NodeIDs: []string{authorityRequestID}}
	decision, ok := s.firstAuthorityDecisionForRequest(authorityRequestID)
	if !ok {
		path.Missing = append(path.Missing, "AuthorityDecision for AuthorityRequest "+authorityRequestID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, decision.CommonNode.ID)
	receipt, ok := s.firstExecutionReceiptForDecision(decision.CommonNode.ID)
	if !ok {
		path.Missing = append(path.Missing, "ExecutionReceipt for AuthorityDecision "+decision.CommonNode.ID)
		return path, path.Err()
	}
	path.NodeIDs = append(path.NodeIDs, receipt.CommonNode.ID)
	path.Completed = true
	return path, nil
}

func (s *InMemoryStore) ActorAuthorityRequestDecisionReceipt(authorityRequestID string) (RequiredPath, error) {
	path := RequiredPath{Name: "ActorIdentity / AuthorityRequest / AuthorityDecision / ExecutionReceipt", NodeIDs: []string{authorityRequestID}}
	requestRecord, err := s.Get(authorityRequestID)
	if err != nil {
		path.Missing = append(path.Missing, "AuthorityRequest "+authorityRequestID)
		return path, path.Err()
	}
	request, ok := requestRecord.(*AuthorityRequest)
	if !ok {
		path.Missing = append(path.Missing, "AuthorityRequest "+authorityRequestID)
		return path, path.Err()
	}
	if !s.hasActorIdentity(request.ActorID) {
		path.Missing = append(path.Missing, "ActorIdentity for actor "+request.ActorID)
		return path, path.Err()
	}
	authorityPath, err := s.AuthorityRequestDecisionReceipt(authorityRequestID)
	path.NodeIDs = append(path.NodeIDs, authorityPath.NodeIDs[1:]...)
	path.EdgeIDs = append(path.EdgeIDs, authorityPath.EdgeIDs...)
	path.Missing = append(path.Missing, authorityPath.Missing...)
	path.Completed = authorityPath.Completed
	if err != nil {
		return path, err
	}
	return path, nil
}

func (s *InMemoryStore) firstOutgoingEdge(fromID, edgeType string) (CommonEdge, bool) {
	for _, edge := range s.EdgesFrom(fromID) {
		if edge.Type == edgeType {
			return edge, true
		}
	}
	return CommonEdge{}, false
}

func (s *InMemoryStore) firstRequirementForFactoryOrder(factoryOrderID string) (*Requirement, bool) {
	for _, r := range s.ByType(TypeRequirement) {
		req := r.(*Requirement)
		if req.FactoryOrderID == factoryOrderID {
			return req, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstAcceptanceCriterionForRequirement(requirementID string) (*AcceptanceCriterion, bool) {
	for _, r := range s.ByType(TypeAcceptanceCriterion) {
		ac := r.(*AcceptanceCriterion)
		if ac.RequirementID == requirementID {
			return ac, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstTestRunForTestCase(testCaseID string) (*TestRun, bool) {
	for _, r := range s.ByType(TypeTestRun) {
		tr := r.(*TestRun)
		if tr.TestCaseID != nil && *tr.TestCaseID == testCaseID {
			return tr, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstGateResultForEvidence(evidenceID string) (*GateResult, bool) {
	for _, r := range s.ByType(TypeGateResult) {
		gr := r.(*GateResult)
		for _, ref := range gr.EvidenceRefs {
			if ref == evidenceID {
				return gr, true
			}
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstFailureForGateResult(gateResultID string) (*Failure, bool) {
	for _, r := range s.ByType(TypeFailure) {
		failure := r.(*Failure)
		if failure.GateResultID != nil && *failure.GateResultID == gateResultID {
			return failure, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstRepairAttemptForFailure(failureID string) (*RepairAttempt, bool) {
	for _, r := range s.ByType(TypeRepairAttempt) {
		repair := r.(*RepairAttempt)
		if repair.FailureID == failureID {
			return repair, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstAuthorityDecisionForRequest(authorityRequestID string) (*AuthorityDecision, bool) {
	for _, r := range s.ByType(TypeAuthorityDecision) {
		decision := r.(*AuthorityDecision)
		if decision.AuthorityRequestID == authorityRequestID {
			return decision, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) firstExecutionReceiptForDecision(authorityDecisionID string) (*ExecutionReceipt, bool) {
	for _, r := range s.ByType(TypeExecutionReceipt) {
		receipt := r.(*ExecutionReceipt)
		if receipt.AuthorityDecisionID == authorityDecisionID {
			return receipt, true
		}
	}
	return nil, false
}

func (s *InMemoryStore) hasActorIdentity(actorID string) bool {
	for _, r := range s.ByType(TypeActorIdentity) {
		identity := r.(*ActorIdentity)
		if identity.ActorID == actorID {
			return true
		}
	}
	return false
}

func (s *InMemoryStore) mustGetTask(id string) (*Task, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	task, ok := r.(*Task)
	return task, ok
}

func (s *InMemoryStore) mustGetRuntimeEnvelope(id string) (*RuntimeEnvelope, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	env, ok := r.(*RuntimeEnvelope)
	return env, ok
}

func (s *InMemoryStore) mustGetRuntimeResult(id string) (*RuntimeResult, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	result, ok := r.(*RuntimeResult)
	return result, ok
}

func (s *InMemoryStore) mustGetArtifact(id string) (*Artifact, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	artifact, ok := r.(*Artifact)
	return artifact, ok
}

func (s *InMemoryStore) mustGetTestCase(id string) (*TestCase, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	tc, ok := r.(*TestCase)
	return tc, ok
}

func (s *InMemoryStore) mustGetGateResult(id string) (*GateResult, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	gr, ok := r.(*GateResult)
	return gr, ok
}

func (s *InMemoryStore) mustGetWaiver(id string) (*Waiver, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	waiver, ok := r.(*Waiver)
	return waiver, ok
}

func (s *InMemoryStore) mustGetReleaseCandidate(id string) (*ReleaseCandidate, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	rc, ok := r.(*ReleaseCandidate)
	return rc, ok
}

func (s *InMemoryStore) mustGetFactoryRuntimeVersion(id string) (*FactoryRuntimeVersion, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	frv, ok := r.(*FactoryRuntimeVersion)
	return frv, ok
}

func (s *InMemoryStore) mustGetAuditReport(id string) (*AuditReport, bool) {
	r, err := s.Get(id)
	if err != nil {
		return nil, false
	}
	audit, ok := r.(*AuditReport)
	return audit, ok
}
