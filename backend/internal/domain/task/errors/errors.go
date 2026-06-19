package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrTaskNotFound    = apperrors.New(apperrors.ErrNotFound, "task not found")
	ErrVersionMismatch = apperrors.New(apperrors.ErrConflict, "task has been modified by another process, please refresh")
	ErrNotAssigned     = apperrors.New(apperrors.ErrInsufficientProjectRole, "only assigned members can update this task")
	ErrInvalidProgress = apperrors.New(apperrors.ErrValidation, "progress must be between 0 and 100")
)
