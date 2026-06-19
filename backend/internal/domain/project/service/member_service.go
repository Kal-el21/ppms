package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/project/dto"
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/project/errors"
	"github.com/Kal-el21/backend/internal/domain/project/repository"
	"gorm.io/gorm"
)

type MemberService interface {
	AddMember(projectID uint64, req dto.AddMemberRequest) (*entity.ProjectMember, error)
	ChangeMemberRole(memberID uint64, actorIsAdmin bool, actorProjectRole string, req dto.ChangeMemberRoleRequest) error
	RemoveMember(memberID uint64, actorIsAdmin bool, actorProjectRole string, actorUserID uint64, removedBy uint64) error
	GetActiveMembers(projectID uint64) ([]entity.ProjectMember, error)
	GetMemberRole(projectID, userID uint64) (entity.ProjectRole, error)
}

type memberService struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) MemberService {
	return &memberService{repo: repo}
}

func (s *memberService) AddMember(projectID uint64, req dto.AddMemberRequest) (*entity.ProjectMember, error) {
	existing, err := s.repo.FindByProjectAndUser(projectID, req.UserID)
	if err == nil && existing != nil {
		return nil, domainerrors.ErrMemberAlreadyExists
	}

	member := &entity.ProjectMember{
		ProjectID:   projectID,
		UserID:      req.UserID,
		ProjectRole: entity.ProjectRole(req.ProjectRole),
		Status:      entity.MemberActive,
	}

	if err := s.repo.Create(member); err != nil {
		return nil, err
	}

	return member, nil
}

// ChangeMemberRole menerapkan Special Rules dari Permission Matrix section 7:
// PROJECT_MANAGER boleh ubah MEMBER <-> OBSERVER, tapi TIDAK boleh ubah PROJECT_MANAGER lain.
// ADMIN selalu boleh (override).
func (s *memberService) ChangeMemberRole(memberID uint64, actorIsAdmin bool, actorProjectRole string, req dto.ChangeMemberRoleRequest) error {
	member, err := s.repo.FindByID(memberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainerrors.ErrMemberNotFound
		}
		return err
	}

	if !actorIsAdmin {
		// Hanya PROJECT_MANAGER yang boleh sampai sini (dijaga oleh project context middleware)
		if member.ProjectRole == entity.RoleProjectManager {
			return domainerrors.ErrCannotModifyPM
		}
		if entity.ProjectRole(req.ProjectRole) == entity.RoleProjectManager {
			// PM tidak boleh promote orang lain jadi PM melalui endpoint ini
			// (mengangkat PM baru tetap dimungkinkan, tapi sebagai keputusan desain
			// kita batasi di sini agar konsisten dengan "tidak boleh mengubah PROJECT_MANAGER")
			return domainerrors.ErrCannotModifyPM
		}
	}

	// Last-PM protection: jika member yang diubah adalah satu-satunya PM aktif, tolak
	if member.ProjectRole == entity.RoleProjectManager {
		count, err := s.repo.CountActiveProjectManagers(member.ProjectID)
		if err != nil {
			return err
		}
		if count <= 1 {
			return domainerrors.ErrLastPMProtection
		}
	}

	return s.repo.UpdateRole(memberID, entity.ProjectRole(req.ProjectRole))
}

// RemoveMember menerapkan last-PM protection (FR-05.04) + special rules PM.
func (s *memberService) RemoveMember(memberID uint64, actorIsAdmin bool, actorProjectRole string, actorUserID uint64, removedBy uint64) error {
	member, err := s.repo.FindByID(memberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainerrors.ErrMemberNotFound
		}
		return err
	}

	if member.ProjectRole == entity.RoleProjectManager {
		count, err := s.repo.CountActiveProjectManagers(member.ProjectID)
		if err != nil {
			return err
		}

		// Prevent removal of last PROJECT_MANAGER (FR-05.04)
		if count <= 1 {
			return domainerrors.ErrLastPMProtection
		}

		// PROJECT_MANAGER tidak boleh menghapus PROJECT_MANAGER lain (hanya ADMIN boleh)
		if !actorIsAdmin {
			return domainerrors.ErrCannotModifyPM
		}
	}

	// Prevent self-removal if user is last PROJECT_MANAGER (FR-05.04)
	if member.UserID == actorUserID && member.ProjectRole == entity.RoleProjectManager {
		count, err := s.repo.CountActiveProjectManagers(member.ProjectID)
		if err != nil {
			return err
		}
		if count <= 1 {
			return domainerrors.ErrLastPMProtection
		}
	}

	return s.repo.ChangeStatus(memberID, entity.MemberRemoved, removedBy)
}

func (s *memberService) GetActiveMembers(projectID uint64) ([]entity.ProjectMember, error) {
	return s.repo.FindActiveByProject(projectID)
}

func (s *memberService) GetMemberRole(projectID, userID uint64) (entity.ProjectRole, error) {
	member, err := s.repo.FindByProjectAndUser(projectID, userID)
	if err != nil {
		return "", err
	}
	return member.ProjectRole, nil
}
