package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders menambahkan HTTP security headers standar.
// Tidak menggantikan CSP/HSTS yang idealnya diatur di reverse proxy (nginx/traefik)
// pada deployment production, tapi sebagai defense-in-depth tambahan di level app.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
