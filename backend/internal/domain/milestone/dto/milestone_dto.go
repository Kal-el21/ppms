package dto

import "time"

type CreateMilestoneRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type UpdateMilestoneRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Version     int        `json:"version" validate:"required"`
}

type ChangeMilestoneStatusRequest struct {
	Status  string `json:"status" validate:"required,oneof=PLANNED ACTIVE COMPLETED CANCELLED"`
	Version int    `json:"version" validate:"required"`
}

type ReorderMilestoneRequest struct {
	OrderedIDs []uint64 `json:"ordered_ids" validate:"required,min=1"`
}

type MilestoneResponse struct {
	ID          uint64     `json:"id"`
	ProjectID   uint64     `json:"project_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrderIndex  int        `json:"order_index"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
	Version     int        `json:"version"`
}
