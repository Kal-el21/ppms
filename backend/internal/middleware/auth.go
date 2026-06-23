package middleware

import (
	"net/http"

	"github.com/Kal-el21/backend/configs"
	"github.com/Kal-el21/backend/internal/domain/auth/jwt"
	userrepo "github.com/Kal-el21/backend/internal/domain/user/repository"
	"github.com/Kal-el21/backend/internal/shared/cookie"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware membaca access token dari httpOnly cookie (bukan lagi
// Authorization header). Fallback ke header tetap disediakan untuk
// kompatibilitas (misal kebutuhan testing via curl/Postman tanpa cookie jar),
// namun browser production akan selalu memakai jalur cookie.
func AuthMiddleware(cfg *configs.Config, userRepo userrepo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cookie.AccessTokenCookie)

		if err != nil || tokenString == "" {
			// Fallback: Authorization header (dipakai untuk testing manual / API client non-browser)
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "missing access token"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString, cfg.JWTAccessSecret)
		if err != nil {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "invalid or expired access token"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := userRepo.FindByID(claims.UserID)
		if err != nil || !user.IsActive {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "user account is inactive"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("system_role", claims.SystemRole)
		c.Set("division_id", claims.DivisionID)

		c.Next()
	}
}
