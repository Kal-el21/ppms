package entity

import "time"

// PortfolioBudgetYear adalah pagu tahunan CAPEX/OPEX untuk seluruh portofolio
// (Entity FR-07.01). Disimpan per tahun dengan plafon (ceiling) masing-masing.
type PortfolioBudgetYear struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	Year         int       `gorm:"unique;not null" json:"year"`
	CapexCeiling float64   `gorm:"not null;default:0" json:"capex_ceiling"`
	OpexCeiling  float64   `gorm:"not null;default:0" json:"opex_ceiling"`
	CreatedBy    *uint64   `json:"created_by,omitempty"`
	UpdatedBy    *uint64   `json:"updated_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int       `gorm:"default:1" json:"version"`
}

func (PortfolioBudgetYear) TableName() string {
	return "portfolio_budget_years"
}
