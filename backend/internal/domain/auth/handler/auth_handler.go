package handler

import (
	"net/http"

	"github.com/Kal-el21/backend/configs"
	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/auth/dto"
	"github.com/Kal-el21/backend/internal/domain/auth/service"
	"github.com/Kal-el21/backend/internal/middleware"
	"github.com/Kal-el21/backend/internal/shared/cookie"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service  service.AuthService
	auditSvc auditservice.AuditService
	cfg      *configs.Config
}

func NewAuthHandler(service service.AuthService, auditSvc auditservice.AuditService, cfg *configs.Config) *AuthHandler {
	return &AuthHandler{service: service, auditSvc: auditSvc, cfg: cfg}
}

func (h *AuthHandler) cookieConfig() cookie.Config {
	sameSite := http.SameSiteLaxMode
	// Strict lebih aman tapi memutus flow OAuth-style redirect dari domain lain;
	// untuk PPMS internal (single SPA origin, tidak ada redirect cross-site),
	// Lax sudah cukup ketat sembari tetap mengizinkan navigasi normal antar halaman.
	return cookie.Config{
		Domain:   h.cfg.CookieDomain,
		Secure:   h.cfg.CookieSecure,
		SameSite: sameSite,
	}
}

// Tambahkan ke AuthHandler:

// Login sekarang menjadi "InitLogin" — step 1 dari 2FA flow.
// Endpoint yang dipakai frontend tetap POST /auth/login, tapi response-nya
// sekarang bisa return {two_fa_required: true} alih-alih langsung set cookie.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	initResp, err := h.service.InitLogin(req, ipAddress, deviceInfo)
	if err != nil {
		h.auditSvc.Log(auditservice.LogParams{
			Module:    "auth",
			Action:    "LOGIN_FAILED",
			IPAddress: ipAddress,
			NewData:   map[string]string{"email": req.Email},
		})
		response.Error(c, err)
		return
	}

	if initResp.TwoFARequired {
		// Step 1 selesai, tunggu OTP dari email
		response.Success(c, http.StatusOK, dto.LoginInitResponse{
			TwoFARequired:   true,
			OTPSessionToken: initResp.OTPSessionToken,
		}, "otp sent to your registered email")
		return
	}

	// 2FA tidak aktif: langsung set cookie
	csrfToken, err := middleware.GenerateCSRFToken()
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to generate csrf token"))
		return
	}

	// Ambil token dari LoginResult yang sudah dieksekusi di service
	// (karena InitLogin untuk non-2FA sudah jalankan createSessionAndTokens)
	// Catatan: untuk non-2FA, kita perlu token. Ini di-handle via
	// pendekatan: service Login lama (yang return *LoginResult) dipanggil
	// kembali sebagai helper. Lihat catatan di bawah.
	loginResult, err := h.service.Login(req, ipAddress, deviceInfo)
	if err != nil {
		response.Error(c, err)
		return
	}

	cfg := h.cookieConfig()
	cookie.SetAuthCookies(c, cfg, loginResult.AccessToken, loginResult.RefreshToken,
		h.cfg.JWTAccessExpiryMinutes*60, h.cfg.JWTRefreshExpiryDays*24*60*60)
	cookie.SetCSRFCookie(c, cfg, csrfToken, h.cfg.JWTRefreshExpiryDays*24*60*60)

	userID := loginResult.User.ID
	h.auditSvc.Log(auditservice.LogParams{
		UserID:    &userID,
		Module:    "auth",
		Action:    "LOGIN_SUCCESS",
		IPAddress: ipAddress,
		UserAgent: deviceInfo,
	})

	response.Success(c, http.StatusOK, dto.LoginResponse{
		User:      loginResult.User,
		CSRFToken: csrfToken,
	}, "login successful")
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	result, err := h.service.VerifyOTPAndLogin(req, ipAddress, deviceInfo)
	if err != nil {
		response.Error(c, err)
		return
	}

	csrfToken, err := middleware.GenerateCSRFToken()
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to generate csrf token"))
		return
	}

	cfg := h.cookieConfig()
	cookie.SetAuthCookies(c, cfg, result.AccessToken, result.RefreshToken,
		h.cfg.JWTAccessExpiryMinutes*60, h.cfg.JWTRefreshExpiryDays*24*60*60)
	cookie.SetCSRFCookie(c, cfg, csrfToken, h.cfg.JWTRefreshExpiryDays*24*60*60)

	userID := result.User.ID
	h.auditSvc.Log(auditservice.LogParams{
		UserID:    &userID,
		Module:    "auth",
		Action:    "LOGIN_2FA_SUCCESS",
		IPAddress: ipAddress,
		UserAgent: deviceInfo,
	})

	response.Success(c, http.StatusOK, dto.LoginResponse{
		User:      result.User,
		CSRFToken: csrfToken,
	}, "login successful")
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req dto.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	if err := h.service.ResendOTP(req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "otp resent")
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshTokenFromCookie, err := c.Cookie(cookie.RefreshTokenCookie)
	if err != nil || refreshTokenFromCookie == "" {
		response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "missing refresh token"))
		return
	}

	result, err := h.service.RefreshToken(dto.RefreshTokenRequest{RefreshToken: refreshTokenFromCookie})
	if err != nil {
		// Refresh gagal (expired/revoked) -> bersihkan cookie agar frontend
		// tahu harus redirect ke login, bukan terus retry dengan cookie basi.
		cookie.ClearAuthCookies(c, h.cookieConfig())
		response.Error(c, err)
		return
	}

	csrfToken, err := middleware.GenerateCSRFToken()
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to generate csrf token"))
		return
	}

	cfg := h.cookieConfig()
	cookie.SetAuthCookies(
		c, cfg,
		result.AccessToken, result.RefreshToken,
		h.cfg.JWTAccessExpiryMinutes*60,
		h.cfg.JWTRefreshExpiryDays*24*60*60,
	)
	cookie.SetCSRFCookie(c, cfg, csrfToken, h.cfg.JWTRefreshExpiryDays*24*60*60)

	response.Success(c, http.StatusOK, dto.RefreshTokenResponse{CSRFToken: csrfToken}, "token refreshed")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshTokenFromCookie, _ := c.Cookie(cookie.RefreshTokenCookie)

	if refreshTokenFromCookie != "" {
		if err := h.service.Logout(dto.LogoutRequest{RefreshToken: refreshTokenFromCookie}); err != nil {
			response.Error(c, err)
			return
		}
	}

	cookie.ClearAuthCookies(c, h.cookieConfig())

	response.Success(c, http.StatusOK, nil, "logout successful")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.service.ChangePassword(userID, req); err != nil {
		response.Error(c, err)
		return
	}

	// Semua sesi lain direvoke server-side; sesi browser saat ini juga
	// sebaiknya di-clear agar user diarahkan login ulang dengan cookie baru.
	cookie.ClearAuthCookies(c, h.cookieConfig())

	h.auditSvc.Log(auditservice.LogParams{
		UserID:    &userID,
		Module:    "auth",
		Action:    "PASSWORD_CHANGED",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	response.Success(c, http.StatusOK, nil, "password changed successfully, please login again")
}

func (h *AuthHandler) RevokeAllSessions(c *gin.Context) {
	userID := c.GetUint64("user_id")

	if err := h.service.RevokeAllSessions(userID, "manual_revoke"); err != nil {
		response.Error(c, err)
		return
	}

	cookie.ClearAuthCookies(c, h.cookieConfig())

	h.auditSvc.Log(auditservice.LogParams{
		UserID: &userID,
		Module: "auth",
		Action: "SESSIONS_REVOKED",
	})

	response.Success(c, http.StatusOK, nil, "all sessions revoked, please login again")
}
