package handler

import (
	"net/http"
	"strconv"

	"github.com/Kal-el21/backend/internal/domain/notification/dto"
	"github.com/Kal-el21/backend/internal/domain/notification/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service     service.NotificationService
	prefService service.PreferenceService
}

func NewNotificationHandler(service service.NotificationService, prefService service.PreferenceService) *NotificationHandler {
	return &NotificationHandler{service: service, prefService: prefService}
}

func (h *NotificationHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	userID := c.GetUint64("user_id")

	notifications, total, err := h.service.GetByUserID(userID, page, limit, unreadOnly)
	if err != nil {
		response.Error(c, err)
		return
	}

	unreadCount, _ := h.service.CountUnread(userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    notifications,
		"meta": gin.H{
			"page":         page,
			"limit":        limit,
			"total":        total,
			"unread_count": unreadCount,
		},
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid notification id"))
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.service.MarkAsRead(id, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "notification marked as read")
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint64("user_id")

	if err := h.service.MarkAllAsRead(userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "all notifications marked as read")
}

func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	userID := c.GetUint64("user_id")

	prefs, err := h.prefService.GetAll(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, prefs, "")
}

func (h *NotificationHandler) UpdatePreference(c *gin.Context) {
	var req dto.UpdatePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.prefService.Update(userID, req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "preference updated successfully")
}
