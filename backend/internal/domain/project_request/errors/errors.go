package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrRequestNotFound      = apperrors.New(apperrors.ErrNotFound, "project request not found")
	ErrNotOwner             = apperrors.New(apperrors.ErrResourceNotOwned, "you are not the owner of this request")
	ErrVersionMismatch      = apperrors.New(apperrors.ErrConflict, "request has been modified by another process, please refresh")
	ErrCannotEditNonDraft   = apperrors.New(apperrors.ErrInvalidStateTransition, "only draft requests can be edited")
	ErrCannotDeleteNonDraft = apperrors.New(apperrors.ErrInvalidStateTransition, "only draft requests can be deleted")
)
