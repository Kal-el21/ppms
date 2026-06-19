package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrAttachmentNotFound = apperrors.New(apperrors.ErrNotFound, "attachment not found")
	ErrFileTooLarge       = apperrors.New(apperrors.ErrFileTooLarge, "file size exceeds maximum allowed limit of 25MB")
	ErrUnsupportedMime    = apperrors.New(apperrors.ErrUnsupportedFile, "file type is not supported")
	ErrInvalidEntityType  = apperrors.New(apperrors.ErrValidation, "invalid entity type")
)
