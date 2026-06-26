package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrUserNotFound    = apperrors.New(apperrors.ErrNotFound, "user not found")
	ErrEmailTaken      = apperrors.New(apperrors.ErrConflict, "email already registered")
	ErrInvalidDivision = apperrors.New(apperrors.ErrValidation, "division not found")
	ErrInvalidFullName = apperrors.New(apperrors.ErrValidation, "full name cannot be empty")
	ErrUserInactive    = apperrors.New(apperrors.ErrUnauthorized, "user account is inactive")
)
