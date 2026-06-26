package entity

import "time"

type ProjectStatus string

const (
	ProjectPlanned   ProjectStatus = "PLANNED"
	ProjectActive    ProjectStatus = "ACTIVE"
	ProjectOnHold    ProjectStatus = "ON_HOLD"
	ProjectCompleted ProjectStatus = "COMPLETED"
	ProjectCancelled ProjectStatus = "CANCELLED"
)

type Project struct {
	ID               uint64  `gorm:"primaryKey" json:"id"`
	ProjectRequestID *uint64 `json:"project_request_id"`
	ProjectCode      string  `gorm:"unique" json:"project_code,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	CompletedAt *time.Time `json:"completed_at,omitempty"`

	Status ProjectStatus `gorm:"not null;default:PLANNED" json:"status"`

	CreatedBy uint64 `gorm:"not null" json:"created_by"`
	UpdatedBy *uint64 `json:"updated_by,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`
}

func (Project) TableName() string {
	return "projects"
}
