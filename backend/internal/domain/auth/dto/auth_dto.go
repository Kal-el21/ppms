package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User      UserSummary `json:"user"`
	CSRFToken string      `json:"csrf_token"`
	// AccessToken & RefreshToken TIDAK LAGI dikirim di body sejak migrasi cookie.
	// Token dikirim via Set-Cookie header, ditangani di handler, bukan di sini.
}

type UserSummary struct {
	ID         uint64  `json:"id"`
	FullName   string  `json:"full_name"`
	Email      string  `json:"email"`
	SystemRole string  `json:"system_role"`
	DivisionID *uint64 `json:"division_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	CSRFToken string `json:"csrf_token"`
	// AccessToken & RefreshToken juga via Set-Cookie, bukan body
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
