package dto

import "github.com/Kal-el21/backend/internal/domain/import_export/entity"

// ImportRequest adalah payload JSON yang di-upload untuk restore data. Struktur
// projects-nya identik dengan hasil export sehingga file backup dapat langsung
// di-import kembali.
type ImportRequest struct {
	Version  string                 `json:"version"`
	Projects []entity.ProjectExport `json:"projects"`
}

// ImportResult merangkum hasil proses import untuk ditampilkan ke user.
type ImportResult struct {
	TotalProjects      int      `json:"total_projects"`
	Imported           int      `json:"imported"`
	Skipped            int      `json:"skipped"`
	Errors             []string `json:"errors"`
	ImportedProjectIDs []uint64 `json:"imported_project_ids"`
}
