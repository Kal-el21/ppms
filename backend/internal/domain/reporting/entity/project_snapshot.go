package entity

import (
	"time"
)

// ProjectSnapshot stores a point-in-time snapshot of project data for reporting.
// Updated: added SnapshotData (jsonb) to store full report payload.
type ProjectSnapshot struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	ProjectID      uint64 `gorm:"not null" json:"project_id"`
	SnapshotDate   time.Time `gorm:"type:date;not null" json:"snapshot_date"`

	ProjectStatus  string  `gorm:"not null" json:"project_status"`
	ProjectProgress float64 `json:"project_progress"`

	BudgetAllocated float64 `json:"budget_allocated"`
	BudgetUsed      float64 `json:"budget_used"`

	TotalTasks      int `json:"total_tasks"`
	CompletedTasks  int `json:"completed_tasks"`

	// SnapshotData menyimpan data detail laporan dalam format JSON.
	// Digunakan untuk reporting tanpa query ke tabel live.
	// Contoh: {"tasks": [...], "milestones": [...], "budgets": [...]}
	SnapshotData map[string]interface{} `gorm:"type:jsonb" json:"snapshot_data,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (ProjectSnapshot) TableName() string {
	return "project_snapshots"
}
