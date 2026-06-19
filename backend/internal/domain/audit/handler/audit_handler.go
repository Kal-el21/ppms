package handler

import (
	"net/http"
	"strconv"

	"github.com/Kal-el21/backend/internal/domain/audit/repository"
	"github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service service.AuditService
}

func NewAuditHandler(service service.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := repository.AuditFilter{
		Module:     c.Query("module"),
		EntityType: c.Query("entity_type"),
		Page:       page,
		Limit:      limit,
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err == nil {
			filter.UserID = &userID
		}
	}

	if entityIDStr := c.Query("entity_id"); entityIDStr != "" {
		entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
		if err == nil {
			filter.EntityID = &entityID
		}
	}

	logs, total, err := h.service.Query(filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	})
}
