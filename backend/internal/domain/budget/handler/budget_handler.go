package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	service  service.BudgetService
	auditSvc auditservice.AuditService
}

func NewBudgetHandler(service service.BudgetService, auditSvc auditservice.AuditService) *BudgetHandler {
	return &BudgetHandler{service: service, auditSvc: auditSvc}
}

func (h *BudgetHandler) Create(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	budget, err := h.service.Create(projectID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "budget",
		Action:     "CREATE_BUDGET",
		EntityType: "PROJECT",
		EntityID:   &projectID,
		NewData:    budget,
	})

	response.Success(c, http.StatusCreated, budget, "budget created successfully")
}

func (h *BudgetHandler) GetByProjectID(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	budget, err := h.service.GetByProjectID(projectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, budget, "")
}

func (h *BudgetHandler) Update(c *gin.Context) {
	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid budget id"))
		return
	}

	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	budget, err := h.service.Update(budgetID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "budget",
		Action:     "UPDATE_BUDGET_ALLOCATION",
		EntityType: "BUDGET",
		EntityID:   &budgetID,
		NewData:    map[string]float64{"allocated_budget": req.AllocatedBudget},
	})

	response.Success(c, http.StatusOK, budget, "budget updated successfully")
}
