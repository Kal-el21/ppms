package handler

import (
	"net/http"
	"strconv"

	"github.com/Kal-el21/backend/internal/domain/attachment/service"
	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	service  service.AttachmentService
	auditSvc auditservice.AuditService
}

func NewAttachmentHandler(service service.AttachmentService, auditSvc auditservice.AuditService) *AttachmentHandler {
	return &AttachmentHandler{service: service, auditSvc: auditSvc}
}

func isAdminCtx(c *gin.Context) bool {
	role, _ := c.Get("system_role")
	return role == "ADMIN"
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	entityType := c.PostForm("entity_type")
	entityIDStr := c.PostForm("entity_id")

	entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid entity_id"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "file is required"))
		return
	}

	uploadedBy := c.GetUint64("user_id")

	attachment, err := h.service.Upload(c.Request.Context(), file, entityType, entityID, uploadedBy, isAdminCtx(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &uploadedBy,
		Module:     "attachment",
		Action:     "UPLOAD_FILE",
		EntityType: entityType,
		EntityID:   &entityID,
		NewData: map[string]interface{}{
			"attachment_id": attachment.ID,
			"file_name":     attachment.OriginalName,
			"file_size":     attachment.FileSize,
			"version":       attachment.Version,
		},
	})

	response.Success(c, http.StatusCreated, attachment, "file uploaded successfully")
}

func (h *AttachmentHandler) GetByEntity(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")

	entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid entity_id"))
		return
	}

	userID := c.GetUint64("user_id")

	attachments, err := h.service.GetByEntity(entityType, entityID, userID, isAdminCtx(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, attachments, "")
}

func (h *AttachmentHandler) GetDownloadURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid attachment id"))
		return
	}

	userID := c.GetUint64("user_id")

	url, attachment, err := h.service.GetDownloadURL(c.Request.Context(), id, userID, isAdminCtx(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &userID,
		Module:     "attachment",
		Action:     "DOWNLOAD_FILE",
		EntityType: string(attachment.EntityType),
		EntityID:   &attachment.EntityID,
		NewData:    map[string]uint64{"attachment_id": id},
	})

	response.Success(c, http.StatusOK, gin.H{
		"download_url": url,
		"file_name":    attachment.OriginalName,
		"mime_type":    attachment.MimeType,
	}, "")
}

func (h *AttachmentHandler) GetVersionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid attachment id"))
		return
	}

	userID := c.GetUint64("user_id")

	versions, err := h.service.GetVersionHistory(id, userID, isAdminCtx(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, versions, "")
}

func (h *AttachmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid attachment id"))
		return
	}

	deletedBy := c.GetUint64("user_id")

	if err := h.service.Delete(id, deletedBy, isAdminCtx(c)); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &deletedBy,
		Module:     "attachment",
		Action:     "DELETE_FILE",
		EntityType: "ATTACHMENT",
		EntityID:   &id,
	})

	response.Success(c, http.StatusOK, nil, "file deleted successfully")
}
