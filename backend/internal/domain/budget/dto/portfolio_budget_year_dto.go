package dto

import "time"

type CreatePortfolioBudgetYearRequest struct {
	Year         int     `json:"year" validate:"required"`
	CapexCeiling float64 `json:"capex_ceiling" validate:"gte=0"`
	OpexCeiling  float64 `json:"opex_ceiling" validate:"gte=0"`
}

type UpdatePortfolioBudgetYearRequest struct {
	CapexCeiling float64 `json:"capex_ceiling" validate:"gte=0"`
	OpexCeiling  float64 `json:"opex_ceiling" validate:"gte=0"`
	Version      int     `json:"version" validate:"required"`
}

type PortfolioBudgetYearResponse struct {
	ID           uint64    `json:"id"`
	Year         int       `json:"year"`
	CapexCeiling float64   `json:"capex_ceiling"`
	OpexCeiling  float64   `json:"opex_ceiling"`
	Version      int       `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
