package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/task/dto"
	"github.com/Kal-el21/backend/internal/domain/task/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service  service.TaskService
	auditSvc auditservice.AuditService
}

func NewTaskHandler(service service.TaskService, auditSvc auditservice.AuditService) *TaskHandler {
	return &TaskHandler{service: service, auditSvc: auditSvc}
}

func isMemberOnly(c *gin.Context) bool {
	role, _ := c.Get("project_role")
	return role == "MEMBER"
}

func (h *TaskHandler) Create(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	task, err := h.service.Create(projectID, userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &userID,
		Module:     "task",
		Action:     "CREATE_TASK",
		EntityType: "TASK",
		EntityID:   &task.ID,
		NewData:    task,
	})

	response.Success(c, http.StatusCreated, task, "task created successfully")
}

func (h *TaskHandler) GetList(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	tasks, err := h.service.GetByProjectID(projectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, tasks, "")
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	oldTask, _ := h.service.GetByID(taskID)

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	task, err := h.service.Update(taskID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "task",
		Action:     "UPDATE_TASK",
		EntityType: "TASK",
		EntityID:   &taskID,
		OldData:    oldTask,
		NewData:    task,
	})

	response.Success(c, http.StatusOK, task, "task updated successfully")
}

func (h *TaskHandler) ChangeStatus(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	oldTask, _ := h.service.GetByID(taskID)

	var req dto.ChangeTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	task, err := h.service.ChangeStatus(taskID, userID, isMemberOnly(c), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	oldStatus := ""
	if oldTask != nil {
		oldStatus = string(oldTask.Status)
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &userID,
		Module:     "task",
		Action:     "CHANGE_TASK_STATUS",
		EntityType: "TASK",
		EntityID:   &taskID,
		OldData:    map[string]string{"status": oldStatus},
		NewData:    map[string]string{"status": req.Status},
	})

	response.Success(c, http.StatusOK, task, "task status updated successfully")
}

func (h *TaskHandler) UpdateProgress(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	oldTask, _ := h.service.GetByID(taskID)

	var req dto.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	task, err := h.service.UpdateProgress(taskID, userID, isMemberOnly(c), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	oldProgress := 0
	if oldTask != nil {
		oldProgress = oldTask.Progress
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &userID,
		Module:     "task",
		Action:     "UPDATE_TASK_PROGRESS",
		EntityType: "TASK",
		EntityID:   &taskID,
		OldData:    map[string]int{"progress": oldProgress},
		NewData:    map[string]int{"progress": req.Progress},
	})

	response.Success(c, http.StatusOK, task, "task progress updated successfully")
}

func (h *TaskHandler) AssignUsers(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	var req dto.AssignUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	assignedBy := c.GetUint64("user_id")

	if err := h.service.AssignUsers(taskID, assignedBy, req); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &assignedBy,
		Module:     "task",
		Action:     "ASSIGN_TASK_USERS",
		EntityType: "TASK",
		EntityID:   &taskID,
		NewData:    map[string]interface{}{"user_ids": req.UserIDs},
	})

	response.Success(c, http.StatusOK, nil, "users assigned successfully")
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.service.AddComment(taskID, userID, req.Comment); err != nil {
		response.Error(c, err)
		return
	}

	// Catatan: komentar TIDAK di-audit-log secara detail (isi komentar tidak disimpan
	// di audit_logs) untuk menghindari duplikasi data dan menjaga audit_logs tetap ringkas.
	// task_comments sendiri sudah menjadi log permanen (tidak ada delete endpoint).

	response.Success(c, http.StatusCreated, nil, "comment added successfully")
}

func (h *TaskHandler) GetComments(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	comments, total, err := h.service.GetComments(taskID, page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comments,
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid task id"))
		return
	}

	actorID := c.GetUint64("user_id")

	if err := h.service.Delete(taskID, actorID); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "task",
		Action:     "DELETE_TASK",
		EntityType: "TASK",
		EntityID:   &taskID,
	})

	response.Success(c, http.StatusOK, nil, "task deleted successfully")
}
