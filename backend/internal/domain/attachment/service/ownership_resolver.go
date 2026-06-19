package service

import (
	"errors"

	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

// OwnershipResolver menentukan project_id yang terasosiasi dengan sebuah
// entity attachment, lalu memvalidasi apakah user adalah member aktif
// project tersebut (atau ADMIN). Ini menutup gap keamanan yang ditandai
// di Phase 4: sebelumnya endpoint /attachments tidak memvalidasi project
// membership sama sekali.
type OwnershipResolver interface {
	ResolveProjectID(entityType string, entityID uint64) (uint64, error)
	ValidateAccess(entityType string, entityID uint64, userID uint64, isAdmin bool) error
}

type ownershipResolver struct {
	db MembershipChecker
}

// MembershipChecker adalah interface kecil untuk menghindari import langsung
// ke domain project (circular dependency risk). Diimplementasikan oleh
// adapter yang membungkus *gorm.DB di main.go.
type MembershipChecker interface {
	GetProjectIDByEntity(entityType string, entityID uint64) (uint64, error)
	IsActiveMember(projectID uint64, userID uint64) (bool, error)
}

func NewOwnershipResolver(checker MembershipChecker) OwnershipResolver {
	return &ownershipResolver{db: checker}
}

func (r *ownershipResolver) ResolveProjectID(entityType string, entityID uint64) (uint64, error) {
	return r.db.GetProjectIDByEntity(entityType, entityID)
}

func (r *ownershipResolver) ValidateAccess(entityType string, entityID uint64, userID uint64, isAdmin bool) error {
	if isAdmin {
		return nil // ADMIN Override (Permission Evaluation Order #1)
	}

	projectID, err := r.db.GetProjectIDByEntity(entityType, entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New(apperrors.ErrNotFound, "related entity not found")
		}
		return err
	}

	isMember, err := r.db.IsActiveMember(projectID, userID)
	if err != nil {
		return err
	}

	if !isMember {
		return apperrors.New(apperrors.ErrNotProjectMember, "you are not an active member of the project that owns this resource")
	}

	return nil
}
