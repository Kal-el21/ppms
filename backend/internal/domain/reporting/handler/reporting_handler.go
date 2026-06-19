package handler

import (
	"strconv"

	"github.com/Kal-el21/backend/internal/domain/reporting/dto"
	"github.com/Kal-el21/backend/internal/domain/reporting/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type ReportingHandler struct {
	service service.ReportingService
}

func NewReportingHandler(service service.ReportingService) *ReportingHandler {
	return &ReportingHandler{service: service}
}

// Generate menangani laporan TANPA project_id wajib (laporan sistem-wide, ADMIN only).
// Jika project_id tetap dikirim di body, tetap diproses sebagai filter opsional.
func (h *ReportingHandler) Generate(c *gin.Context) {
	var req dto.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	h.generateAndRespond(c, req)
}

// GenerateForProject menangani laporan YANG TERIKAT ke satu project tertentu
// via URL param :id, dan dijaga oleh ProjectContextMiddleware + RequireProjectRole.
// project_id dari URL param selalu override body untuk mencegah inkonsistensi.
func (h *ReportingHandler) GenerateForProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	req.ProjectID = &projectID

	h.generateAndRespond(c, req)
}

func (h *ReportingHandler) generateAndRespond(c *gin.Context, req dto.GenerateReportRequest) {
	fileBytes, fileName, contentType, err := h.service.Generate(req)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(200, contentType, fileBytes)
}
