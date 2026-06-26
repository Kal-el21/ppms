package dto

type CreateWorkflowRequest struct {
	Name string `json:"name" validate:"required,min=3,max=200"`
}

type WorkflowResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
}
