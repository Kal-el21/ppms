package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	budgetentity "github.com/Kal-el21/backend/internal/domain/budget/entity"
	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type PortfolioBudgetYearHandler struct {
	service  service.PortfolioBudgetYearService
	auditSvc auditservice.AuditService
}

func NewPortfolioBudgetYearHandler(service service.PortfolioBudgetYearService, auditSvc auditservice.AuditService) *PortfolioBudgetYearHandler {
	return &PortfolioBudgetYearHandler{service: service, auditSvc: auditSvc}
}

func (h *PortfolioBudgetYearHandler) Create(c *gin.Context) {
	var req dto.CreatePortfolioBudgetYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	actorID := c.GetUint64("user_id")
	year, err := h.service.Create(actorID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "portfolio_budget_year",
		Action:     "CREATE_BUDGET_YEAR",
		EntityType: "PORTFOLIO_BUDGET_YEAR",
		EntityID:   &year.ID,
		NewData:    year,
	})

	response.Success(c, http.StatusCreated, toBudgetYearResponse(year), "pagu tahunan berhasil dibuat")
}

func (h *PortfolioBudgetYearHandler) GetAll(c *gin.Context) {
	years, err := h.service.GetAll()
	if err != nil {
		response.Error(c, err)
		return
	}
	result := make([]dto.PortfolioBudgetYearResponse, 0, len(years))
	for _, y := range years {
		result = append(result, toBudgetYearResponse(&y))
	}
	response.Success(c, http.StatusOK, result, "")
}

func (h *PortfolioBudgetYearHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid id"))
		return
	}
	year, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, toBudgetYearResponse(year), "")
}

func (h *PortfolioBudgetYearHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid id"))
		return
	}
	var req dto.UpdatePortfolioBudgetYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	actorID := c.GetUint64("user_id")
	year, err := h.service.Update(id, actorID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "portfolio_budget_year",
		Action:     "UPDATE_BUDGET_YEAR",
		EntityType: "PORTFOLIO_BUDGET_YEAR",
		EntityID:   &id,
		NewData:    year,
	})

	response.Success(c, http.StatusOK, toBudgetYearResponse(year), "pagu tahunan berhasil diperbarui")
}

func (h *PortfolioBudgetYearHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid id"))
		return
	}

	actorID := c.GetUint64("user_id")
	if err := h.service.Delete(id); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "portfolio_budget_year",
		Action:     "DELETE_BUDGET_YEAR",
		EntityType: "PORTFOLIO_BUDGET_YEAR",
		EntityID:   &id,
	})

	response.Success(c, http.StatusOK, nil, "pagu tahunan berhasil dihapus")
}

func toBudgetYearResponse(y *budgetentity.PortfolioBudgetYear) dto.PortfolioBudgetYearResponse {
	return dto.PortfolioBudgetYearResponse{
		ID:           y.ID,
		Year:         y.Year,
		CapexCeiling: y.CapexCeiling,
		OpexCeiling:  y.OpexCeiling,
		Version:      y.Version,
		CreatedAt:    y.CreatedAt,
		UpdatedAt:    y.UpdatedAt,
	}
}
