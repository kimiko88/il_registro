package documents

import "errors"

type WorkflowEngine struct{}

func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{}
}

func (w *WorkflowEngine) CanTransition(current, next DocStatus) error {
	// Valid Paths:
	// Draft -> Submitted
	// Submitted -> Review | Draft (Rejected)
	// Review -> Approved | Submitted (Rejected)
	// Approved -> Signed

	switch current {
	case StatusDraft:
		if next == StatusSubmitted {
			return nil
		}
	case StatusSubmitted:
		if next == StatusReview || next == StatusDraft || next == StatusRejected {
			return nil
		}
	case StatusReview:
		if next == StatusApproved || next == StatusSubmitted || next == StatusRejected {
			return nil
		}
	case StatusApproved:
		if next == StatusSigned {
			return nil
		}
	}
	return errors.New("invalid status transition")
}

func (w *WorkflowEngine) GetNextStatus(action string, currentRole string) (DocStatus, error) {
	// Helpers to guess next status based on action
	switch action {
	case "submit":
		return StatusSubmitted, nil
	case "review_pass":
		return StatusReview, nil
	case "approve":
		return StatusApproved, nil
	case "reject":
		return StatusRejected, nil // Or Draft
	}
	return "", errors.New("unknown action")
}
