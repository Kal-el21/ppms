package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/project_request/dto"
	"github.com/Kal-el21/backend/internal/domain/project_request/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	service  service.RequestService
	auditSvc auditservice.AuditService
}

func NewRequestHandler(service service.RequestService, auditSvc auditservice.AuditService) *RequestHandler {
	return &RequestHandler{service: service, auditSvc: auditSvc}
}

func isAdmin(c *gin.Context) bool {
	role, _ := c.Get("system_role")
	return role == "ADMIN"
}

func (h *RequestHandler) CreateDraft(c *gin.Context) {
	var req dto.CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	requesterID := c.GetUint64("user_id")

	request, err := h.service.CreateDraft(requesterID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, request, "draft created successfully")
}

func (h *RequestHandler) UpdateDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	var req dto.UpdateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	requesterID := c.GetUint64("user_id")

	request, err := h.service.UpdateDraft(id, requesterID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, request, "draft updated successfully")
}

func (h *RequestHandler) Submit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	requesterID := c.GetUint64("user_id")

	request, err := h.service.Submit(id, requesterID)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &requesterID,
		Module:     "project_request",
		Action:     "SUBMIT_REQUEST",
		EntityType: "PROJECT_REQUEST",
		EntityID:   &id,
		NewData:    map[string]string{"status": string(request.Status)},
	})

	response.Success(c, http.StatusOK, request, "request submitted for review")
}

func (h *RequestHandler) Review(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	reviewerID := c.GetUint64("user_id")

	request, err := h.service.Review(id, reviewerID, req, isAdmin(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &reviewerID,
		Module:     "project_request",
		Action:     "REVIEW_REQUEST_" + req.Action,
		EntityType: "PROJECT_REQUEST",
		EntityID:   &id,
		NewData:    map[string]string{"action": req.Action, "comment": req.Comment, "resulting_status": string(request.Status)},
	})

	response.Success(c, http.StatusOK, request, "review submitted successfully")
}

func (h *RequestHandler) Revise(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	var req dto.ReviseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	requesterID := c.GetUint64("user_id")

	request, err := h.service.Revise(id, requesterID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &requesterID,
		Module:     "project_request",
		Action:     "REVISE_REQUEST",
		EntityType: "PROJECT_REQUEST",
		EntityID:   &id,
		NewData:    map[string]string{"revision_reason": req.RevisionReason},
	})

	response.Success(c, http.StatusOK, request, "request revised successfully")
}

func (h *RequestHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	requesterID := c.GetUint64("user_id")

	request, err := h.service.GetByID(id, requesterID, isAdmin(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, request, "")
}

func (h *RequestHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusFilter := c.Query("status")

	var requests interface{}
	var total int64
	var err error

	if isAdmin(c) {
		requests, total, err = h.service.GetAllRequests(page, limit, statusFilter)
	} else {
		requesterID := c.GetUint64("user_id")
		requests, total, err = h.service.GetOwnRequests(requesterID, page, limit)
	}

	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    requests,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *RequestHandler) GetRevisionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	revisions, err := h.service.GetRevisionHistory(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, revisions, "")
}

func (h *RequestHandler) GetApprovalHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	approvals, err := h.service.GetApprovalHistory(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, approvals, "")
}

func (h *RequestHandler) DeleteDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid request id"))
		return
	}

	requesterID := c.GetUint64("user_id")

	if err := h.service.DeleteDraft(id, requesterID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "draft deleted successfully")
}
