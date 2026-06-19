package statemachine

import (
	"github.com/Kal-el21/backend/internal/domain/task/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

// Sesuai FR-07.05: TODO -> IN_PROGRESS -> DONE, CANCELLED dari state apapun.
var allowedTransitions = map[entity.TaskStatus][]entity.TaskStatus{
	entity.StatusTodo:       {entity.StatusInProgress, entity.StatusCancelled},
	entity.StatusInProgress: {entity.StatusDone, entity.StatusCancelled},
	entity.StatusDone:       {entity.StatusCancelled},
	entity.StatusCancelled:  {},
}

func CanTransition(from, to entity.TaskStatus) bool {
	if to == entity.StatusCancelled && from != entity.StatusCancelled {
		return true // CANCELLED from any state per FR-07.05
	}
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

func ValidateTransition(from, to entity.TaskStatus) error {
	if !CanTransition(from, to) {
		return apperrors.New(
			apperrors.ErrInvalidStateTransition,
			"invalid task state transition from "+string(from)+" to "+string(to),
		)
	}
	return nil
}
