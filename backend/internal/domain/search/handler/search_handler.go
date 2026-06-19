package handler

import (
	"net/http"

	"github.com/Kal-el21/backend/internal/domain/search/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service service.SearchService
}

func NewSearchHandler(service service.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "search query must be at least 2 characters"))
		return
	}

	userID := c.GetUint64("user_id")
	role, _ := c.Get("system_role")
	isAdmin := role == "ADMIN"

	results, err := h.service.Search(query, userID, isAdmin)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, results, "")
}
