package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/project/dto"
	projectentity "github.com/Kal-el21/backend/internal/domain/project/entity"
	"github.com/Kal-el21/backend/internal/domain/project/repository"
	"github.com/Kal-el21/backend/internal/domain/project/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service       service.ProjectService
	provisioning service.ProvisioningService
	auditSvc      auditservice.AuditService
}

func NewProjectHandler(service service.ProjectService, provisioning service.ProvisioningService, auditSvc auditservice.AuditService) *ProjectHandler {
	return &ProjectHandler{service: service, provisioning: provisioning, auditSvc: auditSvc}
}

func (h *ProjectHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := repository.ProjectListFilter{
		Search:         c.Query("search"),
		Status:         c.Query("status"),
		BudgetType:     c.Query("budget_type"),
		InitiationType: c.Query("initiation_type"),
		Priority:       c.Query("priority"),
		Sort:           c.Query("sort"),
		Page:           page,
		Limit:          limit,
	}

	role, _ := c.Get("system_role")
	userID := c.GetUint64("user_id")

	// USER/VIEWER hanya boleh lihat project yang mereka ikuti.
	if role != "ADMIN" {
		uid := userID
		filter.MemberUserID = &uid
	}

	rows, total, err := h.service.List(filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	projectResponses := make([]dto.ProjectResponse, 0, len(rows))
	for _, row := range rows {
		projectResponses = append(projectResponses, h.buildListResponse(row))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projectResponses,
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *ProjectHandler) GetDeadlineProjects(c *gin.Context) {
	window := c.Query("window")
	if window == "" {
		window = "90"
	}

	projects, err := h.service.GetDeadlineProjects(window)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.buildProjectResponse(*project),
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

	response.Success(c, http.StatusOK, h.buildProjectResponse(*project), "project updated successfully")
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

	response.Success(c, http.StatusOK, h.buildProjectResponse(*project), "project status updated successfully")
}

func (h *ProjectHandler) CreateDirect(c *gin.Context) {
	var req dto.CreateProjectDirectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	actorID := c.GetUint64("user_id")
	project, err := h.provisioning.CreateDirect(actorID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "project",
		Action:     "CREATE_PROJECT_DIRECT",
		EntityType: "PROJECT",
		EntityID:   &project.ID,
		NewData:    project,
	})

	response.Success(c, http.StatusCreated, h.buildProjectResponse(*project), "project created successfully")
}

func (h *ProjectHandler) resolveHealth(project projectentity.Project, progress float64) string {
	if project.Health != nil {
		return *project.Health
	}
	return h.service.CalculateHealth(&project, progress)
}

func (h *ProjectHandler) buildProjectResponse(project projectentity.Project) dto.ProjectResponse {
	progress, _ := h.service.CalculateProgress(project.ID)

	return dto.ProjectResponse{
		ID:               project.ID,
		ProjectRequestID: project.ProjectRequestID,
		ProjectCode:      project.ProjectCode,
		Name:             project.Name,
		Description:      project.Description,
		Category:         project.Category,
		InitiationType:   project.InitiationType,
		Priority:         project.Priority,
		Notes:            project.Notes,
		StartDate:        project.StartDate,
		EndDate:          project.EndDate,
		CompletedAt:      project.CompletedAt,
		Status:           string(project.Status),
		Progress:         progress,
		Health:           h.resolveHealth(project, progress),
		Version:          project.Version,
		CreatedAt:        project.CreatedAt,
		UpdatedAt:        project.UpdatedAt,
	}
}

func (h *ProjectHandler) buildListResponse(row repository.ProjectListRow) dto.ProjectResponse {
	project := row.Project
	progress, _ := h.service.CalculateProgress(project.ID)

	return dto.ProjectResponse{
		ID:               project.ID,
		ProjectRequestID: project.ProjectRequestID,
		ProjectCode:      project.ProjectCode,
		Name:             project.Name,
		Description:      project.Description,
		Category:         project.Category,
		InitiationType:   project.InitiationType,
		Priority:         project.Priority,
		Notes:            project.Notes,
		StartDate:        project.StartDate,
		EndDate:          project.EndDate,
		CompletedAt:      project.CompletedAt,
		Status:           string(project.Status),
		Progress:         progress,
		Health:           h.resolveHealth(project, progress),
		BudgetType:       row.BudgetType,
		BudgetAllocated:  row.BudgetAllocated,
		BudgetUsed:       row.BudgetUsed,
		Version:          project.Version,
		CreatedAt:        project.CreatedAt,
		UpdatedAt:        project.UpdatedAt,
	}
}
