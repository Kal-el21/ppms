package errors

import apperrors "github.com/Kal-el21/backend/internal/shared/errors"

var (
	ErrBudgetNotFound         = apperrors.New(apperrors.ErrNotFound, "budget not found")
	ErrBudgetAlreadyExists    = apperrors.New(apperrors.ErrConflict, "budget already exists for this project")
	ErrVersionMismatch        = apperrors.New(apperrors.ErrConflict, "budget has been modified by another process, please refresh")
	ErrAdjustmentRequiresType = apperrors.New(apperrors.ErrValidation, "adjustment_type and reason are required when transaction type is ADJUSTMENT")
	ErrDuplicateTransaction   = apperrors.New(apperrors.ErrDuplicateEntry, "duplicate transaction detected, idempotency key already used")
	ErrTransactionImmutable   = apperrors.New(apperrors.ErrValidation, "budget transactions cannot be modified or deleted, use ADJUSTMENT instead")
)
