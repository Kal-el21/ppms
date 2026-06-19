package dto

type CreateUserRequest struct {
	FullName   string  `json:"full_name" validate:"required,min=2,max=150"`
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=8"`
	SystemRole string  `json:"system_role" validate:"required,oneof=ADMIN USER VIEWER"`
	DivisionID *uint64 `json:"division_id"`
}

type UpdateUserRequest struct {
	FullName   string  `json:"full_name" validate:"required,min=2,max=150"`
	DivisionID *uint64 `json:"division_id"`
}

type AssignRoleRequest struct {
	SystemRole string `json:"system_role" validate:"required,oneof=ADMIN USER VIEWER"`
}

type UserResponse struct {
	ID         uint64  `json:"id"`
	FullName   string  `json:"full_name"`
	Email      string  `json:"email"`
	SystemRole string  `json:"system_role"`
	DivisionID *uint64 `json:"division_id"`
	IsActive   bool    `json:"is_active"`
}
