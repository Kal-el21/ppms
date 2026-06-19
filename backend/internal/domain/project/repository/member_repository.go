package repository

import (
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(member *entity.ProjectMember) error
	FindByID(id uint64) (*entity.ProjectMember, error)
	FindByProjectAndUser(projectID, userID uint64) (*entity.ProjectMember, error)
	FindActiveByProject(projectID uint64) ([]entity.ProjectMember, error)
	CountActiveProjectManagers(projectID uint64) (int64, error)
	UpdateRole(id uint64, role entity.ProjectRole) error
	ChangeStatus(id uint64, status entity.MemberStatus, changedBy uint64) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(member *entity.ProjectMember) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) FindByID(id uint64) (*entity.ProjectMember, error) {
	var member entity.ProjectMember
	err := r.db.Where("id = ?", id).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindByProjectAndUser(projectID, userID uint64) (*entity.ProjectMember, error) {
	var member entity.ProjectMember
	err := r.db.Where("project_id = ? AND user_id = ? AND status = 'ACTIVE'", projectID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindActiveByProject(projectID uint64) ([]entity.ProjectMember, error) {
	var members []entity.ProjectMember
	err := r.db.Where("project_id = ? AND status = 'ACTIVE'", projectID).
		Order("joined_at ASC").
		Find(&members).Error
	return members, err
}

func (r *memberRepository) CountActiveProjectManagers(projectID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.ProjectMember{}).
		Where("project_id = ? AND project_role = 'PROJECT_MANAGER' AND status = 'ACTIVE'", projectID).
		Count(&count).Error
	return count, err
}

func (r *memberRepository) UpdateRole(id uint64, role entity.ProjectRole) error {
	return r.db.Model(&entity.ProjectMember{}).
		Where("id = ?", id).
		Update("project_role", role).Error
}

func (r *memberRepository) ChangeStatus(id uint64, status entity.MemberStatus, changedBy uint64) error {
	return r.db.Model(&entity.ProjectMember{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            status,
			"status_changed_by": changedBy,
			"status_changed_at": gorm.Expr("now()"),
			"left_at":           gorm.Expr("now()"),
		}).Error
}
