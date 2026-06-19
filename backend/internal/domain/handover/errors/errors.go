package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrHandoverNotFound    = apperrors.New(apperrors.ErrNotFound, "handover not found")
	ErrVersionMismatch     = apperrors.New(apperrors.ErrConflict, "handover has been modified by another process, please refresh")
	ErrInvalidStatusChange = apperrors.New(apperrors.ErrInvalidStateTransition, "handover status cannot be changed from its current state")
)
