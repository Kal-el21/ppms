package middleware

import (
	"net/http"
	"strconv"

	projectrepo "github.com/Kal-el21/backend/internal/domain/project/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// ProjectContextMiddleware memvalidasi project membership dan project role,
// lalu inject project_id dan project_role ke context.
// ADMIN bypass: tetap lolos walau bukan member project.
func ProjectContextMiddleware(memberRepo projectrepo.MemberRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectIDStr := c.Param("id")
		if projectIDStr == "" {
			projectIDStr = c.Param("projectId")
		}

		projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		systemRole, _ := c.Get("system_role")
		userID := c.GetUint64("user_id")

		c.Set("project_id", projectID)

		// Admin bypass (Permission Evaluation Order #1)
		if systemRole == "ADMIN" {
			c.Set("project_role", "ADMIN_OVERRIDE")
			c.Next()
			return
		}

		member, err := memberRepo.FindByProjectAndUser(projectID, userID)
		if err != nil {
			response.Error(c, apperrors.New(apperrors.ErrNotProjectMember, "you are not an active member of this project"))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("project_role", string(member.ProjectRole))
		c.Set("project_member_id", member.ID)

		c.Next()
	}
}

// RequireProjectRole memvalidasi project_role dari context yang sudah di-set
// oleh ProjectContextMiddleware. ADMIN_OVERRIDE selalu lolos.
func RequireProjectRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("project_role")
		if !exists {
			response.Error(c, apperrors.New(apperrors.ErrUnauthorized, "project context not resolved"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		roleStr := role.(string)

		if roleStr == "ADMIN_OVERRIDE" {
			c.Next()
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, apperrors.New(apperrors.ErrInsufficientProjectRole, "insufficient project role permission"))
		c.AbortWithStatus(http.StatusForbidden)
	}
}
