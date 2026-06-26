package statemachine

import (
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

// allowedTransitions mendefinisikan transisi valid sesuai FR-04.06.
// DRAFT -> SUBMITTED -> UNDER_REVIEW -> APPROVED | REJECTED | REVISION_REQUESTED -> REVISED -> APPROVED/REJECTED/REVISION_REQUESTED
// REVISED -> SUBMITTED tetap didukung untuk kompatibilitas alur lama.
var allowedTransitions = map[entity.RequestStatus][]entity.RequestStatus{
	entity.StatusDraft:             {entity.StatusSubmitted},
	entity.StatusSubmitted:         {entity.StatusUnderReview},
	entity.StatusUnderReview:       {entity.StatusApproved, entity.StatusRejected, entity.StatusRevisionRequested},
	entity.StatusRevisionRequested: {entity.StatusRevised},
	entity.StatusRejected:          {}, // terminal state
	entity.StatusRevised:           {entity.StatusSubmitted, entity.StatusApproved, entity.StatusRejected, entity.StatusRevisionRequested},
	entity.StatusApproved:          {}, // terminal state
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
