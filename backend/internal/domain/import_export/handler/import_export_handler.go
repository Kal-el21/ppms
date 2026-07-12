package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/import_export/dto"
	"github.com/Kal-el21/backend/internal/domain/import_export/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// maxImportSize membatasi ukuran file import agar tidak membebani memori server.
const maxImportSize = 20 << 20 // 20 MB

type ImportExportHandler struct {
	service  service.ImportExportService
	auditSvc auditservice.AuditService
}

func NewImportExportHandler(svc service.ImportExportService, auditSvc auditservice.AuditService) *ImportExportHandler {
	return &ImportExportHandler{service: svc, auditSvc: auditSvc}
}

func isAdminCtx(c *gin.Context) bool {
	role, _ := c.Get("system_role")
	return role == "ADMIN"
}

// Export mengembalikan seluruh data project yang dapat diakses sebagai file JSON
// yang otomatis ter-download di browser.
func (h *ImportExportHandler) Export(c *gin.Context) {
	userID := c.GetUint64("user_id")

	data, err := h.service.ExportAll(userID, isAdminCtx(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to serialize export"))
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &userID,
		Module:     "import_export",
		Action:     "EXPORT_DATA",
		EntityType: "PROJECT",
		NewData: map[string]interface{}{
			"project_count": len(data.Projects),
		},
	})

	filename := fmt.Sprintf("ppms-backup-%s.json", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/json", payload)
}

// Import menerima file JSON hasil export dan membuat ulang project di database.
func (h *ImportExportHandler) Import(c *gin.Context) {
	actorID := c.GetUint64("user_id")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "file is required"))
		return
	}

	if fileHeader.Size > maxImportSize {
		response.Error(c, apperrors.New(apperrors.ErrFileTooLarge, "import file exceeds 20MB limit"))
		return
	}

	opened, err := fileHeader.Open()
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "failed to open uploaded file"))
		return
	}
	defer opened.Close()

	raw, err := io.ReadAll(io.LimitReader(opened, maxImportSize))
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "failed to read uploaded file"))
		return
	}

	var payload dto.ImportRequest
	if err := json.Unmarshal(raw, &payload); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid JSON format: "+err.Error()))
		return
	}

	result, err := h.service.ImportData(payload, actorID)
	if err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "import_export",
		Action:     "IMPORT_DATA",
		EntityType: "PROJECT",
		NewData: map[string]interface{}{
			"total_projects": result.TotalProjects,
			"imported":       result.Imported,
			"skipped":        result.Skipped,
		},
	})

	response.Success(c, http.StatusOK, result, "import completed")
}
