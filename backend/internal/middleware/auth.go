package middleware

import (
	"net/http"
	"strings"

	"github.com/Kal-el21/backend/configs"
	"github.com/Kal-el21/backend/internal/domain/auth/jwt"
	userrepo "github.com/Kal-el21/backend/internal/domain/user/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT and injects user_id, system_role, division_id into context.
func AuthMiddleware(cfg *configs.Config, userRepo userrepo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "missing or invalid authorization header"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateAccessToken(tokenString, cfg.JWTAccessSecret)
		if err != nil {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "invalid or expired access token"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Validate user still active (in case deactivated after token issued)
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
