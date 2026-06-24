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

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=150"`
}

type Toggle2FARequest struct {
	Enabled bool `json:"enabled"`
}

type ToggleEmailNotificationRequest struct {
	Enabled bool `json:"enabled"`
}

type UserResponse struct {
	ID                       uint64  `json:"id"`
	FullName                 string  `json:"full_name"`
	Email                    string  `json:"email"`
	SystemRole               string  `json:"system_role"`
	DivisionID               *uint64 `json:"division_id"`
	IsActive                 bool    `json:"is_active"`
	ProfilePhotoURL          *string `json:"profile_photo_url,omitempty"`
	TwoFAEnabled             bool    `json:"two_fa_enabled"`
	EmailNotificationEnabled bool    `json:"email_notification_enabled"`
}
