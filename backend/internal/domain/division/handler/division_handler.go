package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/division/dto"
	"github.com/Kal-el21/backend/internal/domain/division/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type DivisionHandler struct {
	service  service.DivisionService
	auditSvc auditservice.AuditService
}

func NewDivisionHandler(service service.DivisionService, auditSvc auditservice.AuditService) *DivisionHandler {
	return &DivisionHandler{service: service, auditSvc: auditSvc}
}

func (h *DivisionHandler) Create(c *gin.Context) {
	var req dto.CreateDivisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	division, err := h.service.Create(req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "division",
		Action:     "CREATE_DIVISION",
		EntityType: "DIVISION",
		EntityID:   &division.ID,
		NewData:    division,
	})

	response.Success(c, http.StatusCreated, division, "division created successfully")
}

func (h *DivisionHandler) GetAll(c *gin.Context) {
	divisions, err := h.service.GetAll()
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, divisions, "")
}

func (h *DivisionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid division id"))
		return
	}

	division, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, division, "")
}

func (h *DivisionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid division id"))
		return
	}

	oldDivision, _ := h.service.GetByID(id)

	var req dto.UpdateDivisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	division, err := h.service.Update(id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "division",
		Action:     "UPDATE_DIVISION",
		EntityType: "DIVISION",
		EntityID:   &id,
		OldData:    oldDivision,
		NewData:    division,
	})

	response.Success(c, http.StatusOK, division, "division updated successfully")
}

func (h *DivisionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid division id"))
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.service.Delete(id, userID); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &userID,
		Module:     "division",
		Action:     "DELETE_DIVISION",
		EntityType: "DIVISION",
		EntityID:   &id,
	})

	response.Success(c, http.StatusOK, nil, "division deleted successfully")
}
