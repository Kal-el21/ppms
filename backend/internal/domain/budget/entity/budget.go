package entity

import "time"

type Budget struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	ProjectID uint64 `gorm:"unique;not null" json:"project_id"`

	BudgetType *string `json:"budget_type,omitempty"`
	BudgetName string  `json:"budget_name"`

	AllocatedBudget float64 `gorm:"not null;default:0" json:"allocated_budget"`

	CreatedBy uint64  `gorm:"not null" json:"created_by"`
	UpdatedBy *uint64 `json:"updated_by,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`
}

func (Budget) TableName() string {
	return "budgets"
}
