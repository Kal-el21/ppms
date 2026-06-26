package dto

type CreateLevelRequest struct {
	LevelOrder   int    `json:"level_order" validate:"required,min=1"`
	RoleRequired string `json:"role_required" validate:"required"`
}

type LevelResponse struct {
	ID           uint64 `json:"id"`
	WorkflowID   uint64 `json:"workflow_id"`
	LevelOrder   int    `json:"level_order"`
	RoleRequired string `json:"role_required"`
}
