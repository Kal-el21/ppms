package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/project/dto"
	"github.com/Kal-el21/backend/internal/domain/project/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	service  service.MemberService
	auditSvc auditservice.AuditService
}

func NewMemberHandler(service service.MemberService, auditSvc auditservice.AuditService) *MemberHandler {
	return &MemberHandler{service: service, auditSvc: auditSvc}
}

func (h *MemberHandler) GetList(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	members, err := h.service.GetActiveMembers(projectID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, members, "")
}

func (h *MemberHandler) Add(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid project id"))
		return
	}

	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	member, err := h.service.AddMember(projectID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, member, "member added successfully")
}

func (h *MemberHandler) ChangeRole(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid member id"))
		return
	}

	var req dto.ChangeMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	role, _ := c.Get("system_role")
	isAdmin := role == "ADMIN"
	projectRole, _ := c.Get("project_role")

	if err := h.service.ChangeMemberRole(memberID, isAdmin, projectRole.(string), req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "member role updated successfully")
}

func (h *MemberHandler) Remove(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid member id"))
		return
	}

	role, _ := c.Get("system_role")
	isAdmin := role == "ADMIN"
	projectRole, _ := c.Get("project_role")
	actorUserID := c.GetUint64("user_id")

	if err := h.service.RemoveMember(memberID, isAdmin, projectRole.(string), actorUserID, actorUserID); err != nil {
		response.Error(c, err)
		return
	}

	projectID := c.GetUint64("project_id")
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &actorUserID,
		Module:     "project",
		Action:     "REMOVE_MEMBER",
		EntityType: "PROJECT",
		EntityID:   &projectID,
		NewData:    map[string]uint64{"removed_member_id": memberID},
	})

	response.Success(c, http.StatusOK, nil, "member removed successfully")
}
