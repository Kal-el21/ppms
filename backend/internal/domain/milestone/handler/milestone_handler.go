package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/milestone/dto"
	"github.com/Kal-el21/backend/internal/domain/milestone/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type MilestoneHandler struct {
	service  service.MilestoneService
	auditSvc auditservice.AuditService
}

func NewMilestoneHandler(service service.MilestoneService, auditSvc auditservice.AuditService) *MilestoneHandler {
	return &MilestoneHandler{service: service, auditSvc: auditSvc}
}

func (h *MilestoneHandler) Create(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	milestone, err := h.service.Create(projectID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "milestone",
		Action:     "CREATE_MILESTONE",
		EntityType: "MILESTONE",
		EntityID:   &milestone.ID,
		NewData:    milestone,
	})

	response.Success(c, http.StatusCreated, milestone, "milestone created successfully")
}

func (h *MilestoneHandler) GetList(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	milestones, err := h.service.GetByProjectID(projectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, milestones, "")
}

func (h *MilestoneHandler) Update(c *gin.Context) {
	milestoneID, err := strconv.ParseUint(c.Param("milestoneId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid milestone id"))
		return
	}

	oldMilestone, _ := h.service.GetByID(milestoneID)

	var req dto.UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	milestone, err := h.service.Update(milestoneID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "milestone",
		Action:     "UPDATE_MILESTONE",
		EntityType: "MILESTONE",
		EntityID:   &milestoneID,
		OldData:    oldMilestone,
		NewData:    milestone,
	})

	response.Success(c, http.StatusOK, milestone, "milestone updated successfully")
}

func (h *MilestoneHandler) ChangeStatus(c *gin.Context) {
	milestoneID, err := strconv.ParseUint(c.Param("milestoneId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid milestone id"))
		return
	}

	oldMilestone, _ := h.service.GetByID(milestoneID)

	var req dto.ChangeMilestoneStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	milestone, err := h.service.ChangeStatus(milestoneID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	oldStatus := ""
	if oldMilestone != nil {
		oldStatus = string(oldMilestone.Status)
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "milestone",
		Action:     "CHANGE_MILESTONE_STATUS",
		EntityType: "MILESTONE",
		EntityID:   &milestoneID,
		OldData:    map[string]string{"status": oldStatus},
		NewData:    map[string]string{"status": req.Status},
	})

	response.Success(c, http.StatusOK, milestone, "milestone status updated successfully")
}

func (h *MilestoneHandler) Reorder(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.ReorderMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	if err := h.service.Reorder(projectID, req); err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "milestone",
		Action:     "REORDER_MILESTONES",
		EntityType: "PROJECT",
		EntityID:   &projectID,
		NewData:    map[string]interface{}{"ordered_ids": req.OrderedIDs},
	})

	response.Success(c, http.StatusOK, nil, "milestones reordered successfully")
}

func (h *MilestoneHandler) Delete(c *gin.Context) {
	milestoneID, err := strconv.ParseUint(c.Param("milestoneId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid milestone id"))
		return
	}

	actorID := c.GetUint64("user_id")

	if err := h.service.Delete(milestoneID, actorID); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "milestone",
		Action:     "DELETE_MILESTONE",
		EntityType: "MILESTONE",
		EntityID:   &milestoneID,
	})

	response.Success(c, http.StatusOK, nil, "milestone deleted successfully")
}
