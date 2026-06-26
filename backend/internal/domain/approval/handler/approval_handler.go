package handler

import (
	"net/http"
	"strconv"

	"github.com/Kal-el21/backend/internal/domain/approval/dto"
	"github.com/Kal-el21/backend/internal/domain/approval/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	workflowSvc service.ApprovalWorkflowService
	levelSvc    service.ApprovalLevelService
}

func NewApprovalHandler(workflowSvc service.ApprovalWorkflowService, levelSvc service.ApprovalLevelService) *ApprovalHandler {
	return &ApprovalHandler{
		workflowSvc: workflowSvc,
		levelSvc:    levelSvc,
	}
}

func (h *ApprovalHandler) CreateWorkflow(c *gin.Context) {
	var req dto.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	workflow, err := h.workflowSvc.CreateWorkflow(req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, workflow, "workflow created successfully")
}

func (h *ApprovalHandler) GetWorkflow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid workflow id"))
		return
	}

	workflow, err := h.workflowSvc.GetWorkflow(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, workflow, "")
}

func (h *ApprovalHandler) ListWorkflows(c *gin.Context) {
	workflows, err := h.workflowSvc.ListWorkflows()
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, workflows, "")
}

func (h *ApprovalHandler) CreateLevel(c *gin.Context) {
	workflowID, err := strconv.ParseUint(c.Param("workflowId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid workflow id"))
		return
	}

	var req dto.CreateLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	level, err := h.levelSvc.CreateLevel(workflowID, req.LevelOrder, req.RoleRequired)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, level, "level created successfully")
}

func (h *ApprovalHandler) GetLevels(c *gin.Context) {
	workflowID, err := strconv.ParseUint(c.Param("workflowId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid workflow id"))
		return
	}

	levels, err := h.levelSvc.GetLevelsByWorkflow(workflowID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, levels, "")
}
