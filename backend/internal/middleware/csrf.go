package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/Kal-el21/backend/internal/shared/cookie"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CSRFProtection menerapkan pola "double submit cookie":
// 1. Saat login, server set CSRF token ke cookie (non-httpOnly, bisa dibaca JS)
// 2. Frontend membaca cookie itu, mengirimkannya kembali via header X-CSRF-Token
// 3. Server membandingkan nilai cookie vs header — HARUS SAMA
//
// Mengapa ini mencegah CSRF: situs jahat di domain lain BISA membuat browser
// korban mengirim request (cookie auth otomatis ikut), TAPI situs jahat itu
// TIDAK BISA membaca cookie CSRF milik domain PPMS (browser memblokir akses
// cross-origin ke cookie), sehingga tidak bisa menyertakan header yang cocok.
//
// Hanya diterapkan untuk method yang mengubah state (POST/PUT/PATCH/DELETE);
// GET/HEAD/OPTIONS dibiarkan lewat karena idempotent dan tidak mengubah data.
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		cookieToken, err := c.Cookie(cookie.CSRFTokenCookie)
		if err != nil || cookieToken == "" {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "missing csrf cookie"))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		headerToken := c.GetHeader("X-CSRF-Token")
		if headerToken == "" || headerToken != cookieToken {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "csrf token mismatch"))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
