package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig mengizinkan origin frontend mengakses API.
// Di production, ganti allowedOrigins dengan domain resmi perusahaan
// (jangan gunakan AllowAllOrigins di production).
func CORSConfig(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Disposition"}, // perlu untuk file download (reporting)
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
