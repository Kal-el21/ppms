package entity

import "time"

// Export schema version. Bump ini jika struktur JSON export/import berubah
// secara tidak kompatibel agar importer bisa menolak versi yang tidak dikenal.
const SchemaVersion = "1.0"

// ProjectExport adalah representasi portable satu project beserta seluruh
// relasinya (members, milestones, tasks, budget + transactions). Struktur ini
// dipakai bersama oleh flow export maupun import (backup/restore JSON).
type ProjectExport struct {
	ProjectCode    string     `json:"project_code"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Category       string     `json:"category"`
	InitiationType *string    `json:"initiation_type,omitempty"`
	Priority       string     `json:"priority"`
	Notes          string     `json:"notes"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	Status         string     `json:"status"`
	Health         *string    `json:"health,omitempty"`

	Members    []MemberExport    `json:"members"`
	Milestones []MilestoneExport `json:"milestones"`
	Tasks      []TaskExport      `json:"tasks"`
	Budget     *BudgetExport     `json:"budget,omitempty"`
}

type MemberExport struct {
	UserID      uint64 `json:"user_id"`
	ProjectRole string `json:"project_role"`
	Status      string `json:"status"`
}

// MilestoneExport menyertakan RefID (id milestone asli saat export) agar task
// dapat direlasikan kembali ke milestone yang benar saat import lewat
// MilestoneRefID pada TaskExport.
type MilestoneExport struct {
	RefID       uint64     `json:"ref_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrderIndex  int        `json:"order_index"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `json:"status"`
}

type TaskExport struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`
	Status         string     `json:"status"`
	Progress       int        `json:"progress"`
	OrderIndex     int        `json:"order_index"`
	MilestoneRefID *uint64    `json:"milestone_ref_id,omitempty"`
	StartDate      *time.Time `json:"start_date"`
	DueDate        *time.Time `json:"due_date"`
}

type BudgetExport struct {
	BudgetType      *string             `json:"budget_type,omitempty"`
	BudgetName      string              `json:"budget_name"`
	AllocatedBudget float64             `json:"allocated_budget"`
	Transactions    []TransactionExport `json:"transactions"`
}

type TransactionExport struct {
	Type            string    `json:"type"`
	AdjustmentType  *string   `json:"adjustment_type,omitempty"`
	Amount          float64   `json:"amount"`
	Reason          string    `json:"reason"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}
