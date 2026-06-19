package statemachine

import (
	"github.com/Kal-el21/backend/internal/domain/milestone/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

// Sesuai FR-06.04: PLANNED -> ACTIVE -> COMPLETED -> CANCELLED
// CANCELLED ditambahkan sebagai opsi dari ACTIVE & PLANNED juga, mengikuti pola task.
var allowedTransitions = map[entity.MilestoneStatus][]entity.MilestoneStatus{
	entity.MilestonePlanned:   {entity.MilestoneActive, entity.MilestoneCancelled},
	entity.MilestoneActive:    {entity.MilestoneCompleted, entity.MilestoneCancelled},
	entity.MilestoneCompleted: {},
	entity.MilestoneCancelled: {},
}

func CanTransition(from, to entity.MilestoneStatus) bool {
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

func ValidateTransition(from, to entity.MilestoneStatus) error {
	if !CanTransition(from, to) {
		return apperrors.New(
			apperrors.ErrInvalidStateTransition,
			"invalid milestone state transition from "+string(from)+" to "+string(to),
		)
	}
	return nil
}
