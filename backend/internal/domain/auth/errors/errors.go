package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrInvalidCredentials = apperrors.New(apperrors.ErrUnauthorized, "invalid email or password")
	ErrUserInactive       = apperrors.New(apperrors.ErrUnauthorized, "user account is inactive")
	ErrInvalidToken       = apperrors.New(apperrors.ErrUnauthorized, "invalid or expired token")
	ErrSessionRevoked     = apperrors.New(apperrors.ErrUnauthorized, "session has been revoked")
	ErrWrongOldPassword   = apperrors.New(apperrors.ErrValidation, "old password is incorrect")
	ErrTooManyAttempts    = apperrors.New(apperrors.ErrUnauthorized, "too many login attempts, please try again later")
)
