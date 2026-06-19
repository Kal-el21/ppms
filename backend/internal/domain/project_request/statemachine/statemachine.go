package statemachine

import (
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

// allowedTransitions mendefinisikan transisi valid sesuai FR-04.06.
// DRAFT -> SUBMITTED -> UNDER_REVIEW -> APPROVED | REJECTED -> REVISED -> SUBMITTED
var allowedTransitions = map[entity.RequestStatus][]entity.RequestStatus{
	entity.StatusDraft:       {entity.StatusSubmitted},
	entity.StatusSubmitted:   {entity.StatusUnderReview},
	entity.StatusUnderReview: {entity.StatusApproved, entity.StatusRejected},
	entity.StatusRejected:    {entity.StatusRevised},
	entity.StatusRevised:     {entity.StatusSubmitted},
	entity.StatusApproved:    {}, // terminal state
}

func CanTransition(from, to entity.RequestStatus) bool {
	allowed, exists := allowedTransitions[from]
	if !exists {
		return false
	}
	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}

func ValidateTransition(from, to entity.RequestStatus) error {
	if !CanTransition(from, to) {
		return apperrors.New(
			apperrors.ErrInvalidStateTransition,
			"invalid state transition from "+string(from)+" to "+string(to),
		)
	}
	return nil
}
