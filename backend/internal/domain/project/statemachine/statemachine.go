package statemachine

import (
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

// allowedTransitions sesuai FR-05.02:
// PLANNED -> ACTIVE -> {ON_HOLD, COMPLETED}, ON_HOLD -> CANCELLED
// Catatan: ON_HOLD juga bisa balik ke ACTIVE (resume), ditambahkan untuk kepraktisan operasional.
var allowedTransitions = map[entity.ProjectStatus][]entity.ProjectStatus{
	entity.ProjectPlanned:   {entity.ProjectActive, entity.ProjectCancelled},
	entity.ProjectActive:    {entity.ProjectOnHold, entity.ProjectCompleted, entity.ProjectCancelled},
	entity.ProjectOnHold:    {entity.ProjectActive, entity.ProjectCancelled},
	entity.ProjectCompleted: {},
	entity.ProjectCancelled: {},
}

func CanTransition(from, to entity.ProjectStatus) bool {
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

func ValidateTransition(from, to entity.ProjectStatus) error {
	if !CanTransition(from, to) {
		return apperrors.New(
			apperrors.ErrInvalidStateTransition,
			"invalid project state transition from "+string(from)+" to "+string(to),
		)
	}
	return nil
}
