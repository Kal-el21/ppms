package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrProjectNotFound     = apperrors.New(apperrors.ErrNotFound, "project not found")
	ErrVersionMismatch     = apperrors.New(apperrors.ErrConflict, "project has been modified by another process, please refresh")
	ErrMemberNotFound      = apperrors.New(apperrors.ErrNotFound, "project member not found")
	ErrMemberAlreadyExists = apperrors.New(apperrors.ErrConflict, "user is already a member of this project")
	ErrLastPMProtection    = apperrors.New(apperrors.ErrLastPMProtection, "cannot remove or demote the last active project manager")
	ErrCannotModifyPM      = apperrors.New(apperrors.ErrInsufficientProjectRole, "project manager cannot modify another project manager's role")
	ErrCannotRemoveAdmin   = apperrors.New(apperrors.ErrInsufficientProjectRole, "cannot remove a user with system role ADMIN from project via this action")
)
