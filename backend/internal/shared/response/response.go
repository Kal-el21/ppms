package response

import (
	"net/http"

	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		statusCode := mapCodeToStatus(appErr.Code)
		c.JSON(statusCode, ErrorResponse{
			Success: false,
			Code:    string(appErr.Code),
			Message: appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Code:    string(apperrors.ErrInternal),
		Message: "Internal server error",
	})
}

func mapCodeToStatus(code apperrors.ErrorCode) int {
	switch code {
	case apperrors.ErrInsufficientSystemRole, apperrors.ErrInsufficientProjectRole,
		apperrors.ErrNotProjectMember, apperrors.ErrResourceNotOwned, apperrors.ErrUnauthorized:
		return http.StatusForbidden
	case apperrors.ErrNotFound:
		return http.StatusNotFound
	case apperrors.ErrValidation, apperrors.ErrInvalidStateTransition,
		apperrors.ErrFileTooLarge, apperrors.ErrUnsupportedFile:
		return http.StatusBadRequest
	case apperrors.ErrConflict, apperrors.ErrProjectLocked,
		apperrors.ErrLastPMProtection, apperrors.ErrDuplicateEntry:
		return http.StatusConflict
	case apperrors.ErrRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
