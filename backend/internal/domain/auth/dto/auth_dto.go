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

// LoginInitResponse dikirim setelah step 1 login (credential valid).
// Berisi session token sementara untuk identify user saat OTP diverifikasi,
// dan flag apakah 2FA aktif (jika tidak aktif, langsung set cookie di step 1).
type LoginInitResponse struct {
	TwoFARequired bool `json:"two_fa_required"`
	// OTPSessionToken: short-lived token (tidak sama dengan JWT auth) yang
	// frontend kirim kembali saat verify OTP, untuk membuktikan sudah selesai
	// step 1 dari flow ini. Disimpan di server juga via sesi pendek.
	OTPSessionToken string `json:"otp_session_token,omitempty"`
	// Jika TwoFARequired=false, langsung isi user + csrf (sama seperti login lama):
	User      *UserSummary `json:"user,omitempty"`
	CSRFToken string       `json:"csrf_token,omitempty"`
}

type VerifyOTPRequest struct {
	OTPSessionToken string `json:"otp_session_token" validate:"required"`
	OTPCode         string `json:"otp_code" validate:"required,len=6"`
}

type ResendOTPRequest struct {
	OTPSessionToken string `json:"otp_session_token" validate:"required"`
}
