package entity

import "time"

type MilestoneStatus string

const (
	MilestonePlanned   MilestoneStatus = "PLANNED"
	MilestoneActive    MilestoneStatus = "ACTIVE"
	MilestoneCompleted MilestoneStatus = "COMPLETED"
	MilestoneCancelled MilestoneStatus = "CANCELLED"
)

type Milestone struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	ProjectID uint64 `gorm:"not null" json:"project_id"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	OrderIndex int `gorm:"default:0" json:"order_index"`

	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	Status MilestoneStatus `gorm:"not null;default:PLANNED" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`
}

func (Milestone) TableName() string {
	return "milestones"
}
