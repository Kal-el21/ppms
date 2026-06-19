package middleware

import (
	"net/http"

	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// RequireSystemRole checks if the authenticated user's system_role
// is within the allowed roles. ADMIN always passes (handled by caller
// including ADMIN explicitly if needed, or via override logic later).
func RequireSystemRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("system_role")
		if !exists {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "unauthenticated"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		roleStr := role.(string)

		// ADMIN Override (Permission Evaluation Order #1)
		if roleStr == "ADMIN" {
			c.Next()
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, apperrors.New(apperrors.ErrInsufficientSystemRole, "insufficient system role permission"))
		c.AbortWithStatus(http.StatusForbidden)
	}
}
