package handler

import (
	"net/http"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/auth/dto"
	"github.com/Kal-el21/backend/internal/domain/auth/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service  service.AuthService
	auditSvc auditservice.AuditService
}

func NewAuthHandler(service service.AuthService, auditSvc auditservice.AuditService) *AuthHandler {
	return &AuthHandler{service: service, auditSvc: auditSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	result, err := h.service.Login(req, ipAddress, deviceInfo)
	if err != nil {
		// Log percobaan login gagal juga, untuk audit security (tanpa user_id karena belum teridentifikasi)
		h.auditSvc.Log(auditservice.LogParams{
			Module:    "auth",
			Action:    "LOGIN_FAILED",
			IPAddress: ipAddress,
			UserAgent: deviceInfo,
			NewData:   map[string]string{"email": req.Email},
		})
		response.Error(c, err)
		return
	}

	userID := result.User.ID
	h.auditSvc.Log(auditservice.LogParams{
		UserID:    &userID,
		Module:    "auth",
		Action:    "LOGIN_SUCCESS",
		IPAddress: ipAddress,
		UserAgent: deviceInfo,
	})

	response.Success(c, http.StatusOK, result, "login successful")
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	result, err := h.service.RefreshToken(req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, "token refreshed")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	if err := h.service.Logout(req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "logout successful")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.service.ChangePassword(userID, req); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:    &userID,
		Module:    "auth",
		Action:    "PASSWORD_CHANGED",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	response.Success(c, http.StatusOK, nil, "password changed successfully, all other sessions revoked")
}

func (h *AuthHandler) RevokeAllSessions(c *gin.Context) {
	userID := c.GetUint64("user_id")

	if err := h.service.RevokeAllSessions(userID, "manual_revoke"); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID: &userID,
		Module: "auth",
		Action: "SESSIONS_REVOKED",
	})

	response.Success(c, http.StatusOK, nil, "all sessions revoked")
}
