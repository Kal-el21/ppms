package errors

type ErrorCode string

const (
	// Permission Denial Codes (Permission Matrix FR-15.02)
	ErrInsufficientSystemRole  ErrorCode = "INSUFFICIENT_SYSTEM_ROLE"
	ErrInsufficientProjectRole ErrorCode = "INSUFFICIENT_PROJECT_ROLE"
	ErrNotProjectMember        ErrorCode = "NOT_PROJECT_MEMBER"
	ErrResourceNotOwned        ErrorCode = "RESOURCE_NOT_OWNED"
	ErrProjectLocked           ErrorCode = "PROJECT_LOCKED"
	ErrLastPMProtection        ErrorCode = "LAST_PM_PROTECTION"
	ErrInvalidStateTransition  ErrorCode = "INVALID_STATE_TRANSITION"

	// General Errors
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrConflict     ErrorCode = "CONFLICT"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"

	// Tambahan Phase 7 — sebelumnya beberapa handler memakai ErrValidation
	// secara generik untuk kasus ini, sekarang dipisah agar frontend bisa
	// menampilkan pesan yang lebih spesifik per kode.
	ErrRateLimited     ErrorCode = "RATE_LIMITED"
	ErrFileTooLarge    ErrorCode = "FILE_TOO_LARGE"
	ErrUnsupportedFile ErrorCode = "UNSUPPORTED_FILE_TYPE"
	ErrDuplicateEntry  ErrorCode = "DUPLICATE_ENTRY"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
