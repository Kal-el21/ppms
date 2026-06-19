package handler

import (
	"net/http"

	"github.com/Kal-el21/backend/internal/domain/dashboard/service"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(service service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	userID := c.GetUint64("user_id")
	role, _ := c.Get("system_role")
	isAdmin := role == "ADMIN"

	summary, err := h.service.GetSummary(userID, isAdmin)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, summary, "")
}
