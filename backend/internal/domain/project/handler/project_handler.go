package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/project/dto"
	"github.com/Kal-el21/backend/internal/domain/project/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service  service.ProjectService
	auditSvc auditservice.AuditService
}

func NewProjectHandler(service service.ProjectService, auditSvc auditservice.AuditService) *ProjectHandler {
	return &ProjectHandler{service: service, auditSvc: auditSvc}
}

func (h *ProjectHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	role, _ := c.Get("system_role")
	userID := c.GetUint64("user_id")

	var projects interface{}
	var total int64
	var err error

	// ADMIN melihat semua project; USER/VIEWER hanya project yang mereka ikuti
	if role == "ADMIN" {
		projects, total, err = h.service.GetAll(page, limit, status)
	} else {
		projects, total, err = h.service.GetMyProjects(userID, page, limit)
	}

	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	project, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	progress, _ := h.service.CalculateProgress(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":          project.ID,
			"name":        project.Name,
			"description": project.Description,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
			"status":      project.Status,
			"progress":    progress,
			"version":     project.Version,
			"created_at":  project.CreatedAt,
		},
	})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	project, err := h.service.Update(id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, project, "project updated successfully")
}

func (h *ProjectHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	project, err := h.service.ChangeStatus(id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "project",
		Action:     "CHANGE_PROJECT_STATUS",
		EntityType: "PROJECT",
		EntityID:   &id,
		NewData:    map[string]string{"status": req.Status},
	})

	response.Success(c, http.StatusOK, project, "project status updated successfully")
}
