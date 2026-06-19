package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrDivisionNotFound  = apperrors.New(apperrors.ErrNotFound, "division not found")
	ErrDivisionNameTaken = apperrors.New(apperrors.ErrConflict, "division name already exists")
)
