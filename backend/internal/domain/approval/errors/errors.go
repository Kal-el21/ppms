package errors

import "errors"

var (
	ErrWorkflowNotFound       = errors.New("approval workflow not found")
	ErrLevelNotFound          = errors.New("approval level not found")
	ErrWorkflowAlreadyExists  = errors.New("approval workflow already exists")
)
