package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/handover/dto"
	"github.com/Kal-el21/backend/internal/domain/handover/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type HandoverHandler struct {
	service  service.HandoverService
	auditSvc auditservice.AuditService
}

func NewHandoverHandler(service service.HandoverService, auditSvc auditservice.AuditService) *HandoverHandler {
	return &HandoverHandler{service: service, auditSvc: auditSvc}
}

func (h *HandoverHandler) Create(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.CreateHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	senderID := c.GetUint64("user_id")

	handover, err := h.service.Create(projectID, senderID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &senderID,
		Module:     "handover",
		Action:     "CREATE_HANDOVER",
		EntityType: "HANDOVER",
		EntityID:   &handover.ID,
		NewData:    handover,
	})

	response.Success(c, http.StatusCreated, handover, "handover recorded successfully")
}

func (h *HandoverHandler) GetList(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	statusFilter := c.Query("status")

	handovers, err := h.service.GetByProjectID(projectID, statusFilter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, handovers, "")
}

func (h *HandoverHandler) MarkReceived(c *gin.Context) {
	handoverID, err := strconv.ParseUint(c.Param("handoverId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid handover id"))
		return
	}

	var req dto.MarkReceivedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	receiverID := c.GetUint64("user_id")

	handover, err := h.service.MarkReceived(handoverID, receiverID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &receiverID,
		Module:     "handover",
		Action:     "MARK_HANDOVER_RECEIVED",
		EntityType: "HANDOVER",
		EntityID:   &handoverID,
		NewData:    map[string]string{"status": string(handover.Status)},
	})

	response.Success(c, http.StatusOK, handover, "handover marked as received")
}

func (h *HandoverHandler) MarkReturned(c *gin.Context) {
	handoverID, err := strconv.ParseUint(c.Param("handoverId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid handover id"))
		return
	}

	var req dto.MarkReturnedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	handover, err := h.service.MarkReturned(handoverID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "handover",
		Action:     "MARK_HANDOVER_RETURNED",
		EntityType: "HANDOVER",
		EntityID:   &handoverID,
		NewData:    map[string]string{"status": string(handover.Status), "reason": req.Reason},
	})

	response.Success(c, http.StatusOK, handover, "handover marked as returned")
}
