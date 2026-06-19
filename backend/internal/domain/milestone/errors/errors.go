package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrMilestoneNotFound = apperrors.New(apperrors.ErrNotFound, "milestone not found")
	ErrVersionMismatch   = apperrors.New(apperrors.ErrConflict, "milestone has been modified by another process, please refresh")
)
