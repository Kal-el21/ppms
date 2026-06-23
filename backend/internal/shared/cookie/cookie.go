package cookie

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookie  = "ppms_access_token"
	RefreshTokenCookie = "ppms_refresh_token"
	CSRFTokenCookie    = "ppms_csrf_token"
)

type Config struct {
	Domain   string
	Secure   bool // wajib true di production (HTTPS only)
	SameSite http.SameSite
}

// SetAuthCookies menulis access & refresh token sebagai httpOnly cookie.
// httpOnly=true berarti document.cookie di JavaScript TIDAK BISA membaca
// cookie ini sama sekali — inilah yang mencegah XSS mencuri token, beda
// dengan localStorage yang bisa diakses script apa pun yang berhasil
// inject ke halaman.
func SetAuthCookies(c *gin.Context, cfg Config, accessToken string, refreshToken string, accessMaxAgeSeconds int, refreshMaxAgeSeconds int) {
	c.SetSameSite(cfg.SameSite)

	c.SetCookie(
		AccessTokenCookie,
		accessToken,
		accessMaxAgeSeconds,
		"/",
		cfg.Domain,
		cfg.Secure,
		true, // httpOnly
	)

	c.SetCookie(
		RefreshTokenCookie,
		refreshToken,
		refreshMaxAgeSeconds,
		"/api/v1/auth", // refresh token HANYA dikirim browser ke endpoint /auth/*,
		cfg.Domain,     // mempersempit blast radius jika ada kebocoran path lain
		cfg.Secure,
		true,
	)
}

// SetCSRFCookie menulis CSRF token sebagai cookie BUKAN httpOnly (sengaja),
// karena frontend JavaScript perlu membacanya untuk disertakan di header
// X-CSRF-Token pada setiap request mutasi (POST/PUT/PATCH/DELETE).
// Ini pola "double submit cookie" standar untuk proteksi CSRF.
func SetCSRFCookie(c *gin.Context, cfg Config, csrfToken string, maxAgeSeconds int) {
	c.SetSameSite(cfg.SameSite)
	c.SetCookie(
		CSRFTokenCookie,
		csrfToken,
		maxAgeSeconds,
		"/",
		cfg.Domain,
		cfg.Secure,
		false, // NOT httpOnly — frontend perlu baca ini
	)
}

func ClearAuthCookies(c *gin.Context, cfg Config) {
	c.SetSameSite(cfg.SameSite)
	c.SetCookie(AccessTokenCookie, "", -1, "/", cfg.Domain, cfg.Secure, true)
	c.SetCookie(RefreshTokenCookie, "", -1, "/api/v1/auth", cfg.Domain, cfg.Secure, true)
	c.SetCookie(CSRFTokenCookie, "", -1, "/", cfg.Domain, cfg.Secure, false)
}
