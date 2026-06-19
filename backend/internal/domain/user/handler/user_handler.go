package handler

import (
	"net/http"
	"strconv"

	"github.com/Kal-el21/backend/configs"
	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/user/dto"
	"github.com/Kal-el21/backend/internal/domain/user/entity"
	"github.com/Kal-el21/backend/internal/domain/user/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service  service.UserService
	cfg      *configs.Config
	auditSvc auditservice.AuditService
}

func NewUserHandler(service service.UserService, cfg *configs.Config, auditSvc auditservice.AuditService) *UserHandler {
	return &UserHandler{service: service, cfg: cfg, auditSvc: auditSvc}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	user, err := h.service.Create(req, h.cfg.BcryptCost)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "user",
		Action:     "CREATE_USER",
		EntityType: "USER",
		EntityID:   &user.ID,
		NewData:    user,
	})

	response.Success(c, http.StatusCreated, toUserResponse(user), "user created successfully")
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.service.GetAll(page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	result := make([]dto.UserResponse, len(users))
	for i, u := range users {
		result[i] = *toUserResponse(&u)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid user id"))
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, toUserResponse(user), "")
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid user id"))
		return
	}

	oldUser, _ := h.service.GetByID(id)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	user, err := h.service.Update(id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "user",
		Action:     "UPDATE_USER",
		EntityType: "USER",
		EntityID:   &user.ID,
		OldData:    oldUser,
		NewData:    user,
	})

	response.Success(c, http.StatusOK, toUserResponse(user), "user updated successfully")
}

func (h *UserHandler) AssignRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid user id"))
		return
	}

	oldUser, _ := h.service.GetByID(id)

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	user, err := h.service.AssignRole(id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "user",
		Action:     "ASSIGN_ROLE",
		EntityType: "USER",
		EntityID:   &user.ID,
		OldData:    map[string]string{"system_role": string(oldUser.SystemRole)},
		NewData:    map[string]string{"system_role": string(user.SystemRole)},
	})

	response.Success(c, http.StatusOK, toUserResponse(user), "role assigned successfully")
}

func (h *UserHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid user id"))
		return
	}

	actorID := c.GetUint64("user_id")

	if err := h.service.Deactivate(id, actorID); err != nil {
		response.Error(c, err)
		return
	}

	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "user",
		Action:     "DEACTIVATE_USER",
		EntityType: "USER",
		EntityID:   &id,
	})

	response.Success(c, http.StatusOK, nil, "user deactivated successfully")
}

func (h *UserHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid user id"))
		return
	}

	if err := h.service.Restore(id); err != nil {
		response.Error(c, err)
		return
	}

	actorID := c.GetUint64("user_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorID,
		Module:     "user",
		Action:     "RESTORE_USER",
		EntityType: "USER",
		EntityID:   &id,
	})

	response.Success(c, http.StatusOK, nil, "user restored successfully")
}

func toUserResponse(u *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:         u.ID,
		FullName:   u.FullName,
		Email:      u.Email,
		SystemRole: string(u.SystemRole),
		DivisionID: u.DivisionID,
		IsActive:   u.IsActive,
	}
}
