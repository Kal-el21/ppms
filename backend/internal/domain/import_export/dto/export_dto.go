package dto

import (
	"time"

	"github.com/Kal-el21/backend/internal/domain/import_export/entity"
)

// ExportResponse membungkus seluruh project yang dapat diakses beserta relasinya
// dalam satu payload JSON untuk keperluan backup penuh.
type ExportResponse struct {
	Version    string                 `json:"version"`
	ExportedAt time.Time              `json:"exported_at"`
	ExportedBy uint64                 `json:"exported_by"`
	Projects   []entity.ProjectExport `json:"projects"`
}
