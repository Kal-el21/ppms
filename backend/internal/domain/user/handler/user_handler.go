package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kal-el21/backend/configs"
	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/user/dto"
	"github.com/Kal-el21/backend/internal/domain/user/entity"
	"github.com/Kal-el21/backend/internal/domain/user/service"
	"github.com/Kal-el21/backend/internal/infrastructure/minio"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service     service.UserService
	cfg         *configs.Config
	auditSvc    auditservice.AuditService
	minioClient *minio.Client
}

func NewUserHandler(service service.UserService, cfg *configs.Config, auditSvc auditservice.AuditService, minioClient *minio.Client) *UserHandler {
	return &UserHandler{service: service, cfg: cfg, auditSvc: auditSvc, minioClient: minioClient}
}

var allowedProfilePhotoTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
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
		ID:                       u.ID,
		FullName:                 u.FullName,
		Email:                    u.Email,
		SystemRole:               string(u.SystemRole),
		DivisionID:               u.DivisionID,
		IsActive:                 u.IsActive,
		ProfilePhotoURL:          u.ProfilePhotoURL,
		TwoFAEnabled:             u.TwoFAEnabled,
		EmailNotificationEnabled: u.EmailNotificationEnabled,
	}
}
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	user, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, toUserResponse(user), "profile updated successfully")
}

func (h *UserHandler) UploadProfilePhoto(c *gin.Context) {
	userID := c.GetUint64("user_id")

	file, err := c.FormFile("photo")
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "photo file is required"))
		return
	}

	// Validasi ukuran (max 5MB untuk foto profile)
	if file.Size > 5*1024*1024 {
		response.Error(c, apperrors.New(apperrors.ErrFileTooLarge, "photo size exceeds 5MB limit"))
		return
	}

	contentType := file.Header.Get("Content-Type")
	ext, ok := allowedProfilePhotoTypes[contentType]
	if !ok {
		response.Error(c, apperrors.New(apperrors.ErrUnsupportedFile, "photo must be JPG, PNG, or WebP"))
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to open file"))
		return
	}
	defer src.Close()

	objectName := fmt.Sprintf("profile-photos/%d/%s%s", userID, uuid.NewString(), ext)
	if err := h.minioClient.Upload(c.Request.Context(), objectName, src, file.Size, contentType); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to upload photo; please verify object storage is available"))
		return
	}

	photoURL, err := h.minioClient.GetPresignedDownloadURL(c.Request.Context(), objectName)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to get photo url"))
		return
	}

	if err := h.service.UpdateProfilePhoto(userID, photoURL); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"photo_url": photoURL}, "profile photo updated")
}

func (h *UserHandler) Toggle2FA(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req dto.Toggle2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	if err := h.service.Toggle2FA(userID, req.Enabled); err != nil {
		response.Error(c, err)
		return
	}

	msg := "two-factor authentication disabled"
	if req.Enabled {
		msg = "two-factor authentication enabled"
	}

	response.Success(c, http.StatusOK, nil, msg)
}

func (h *UserHandler) ToggleEmailNotification(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req dto.ToggleEmailNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	if err := h.service.ToggleEmailNotification(userID, req.Enabled); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "email notification preference updated")
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint64("user_id")

	user, err := h.service.GetByID(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, toUserResponse(user), "")
}
