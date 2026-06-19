package repository

import (
	"github.com/Kal-el21/backend/internal/domain/milestone/entity"
	"gorm.io/gorm"
)

type MilestoneRepository interface {
	Create(milestone *entity.Milestone) error
	FindByID(id uint64) (*entity.Milestone, error)
	FindByProjectID(projectID uint64) ([]entity.Milestone, error)
	UpdateWithVersionCheck(milestone *entity.Milestone, expectedVersion int) (int64, error)
	UpdateOrderIndex(id uint64, orderIndex int) error
	Delete(id uint64, deletedBy uint64) error
}

type milestoneRepository struct {
	db *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &milestoneRepository{db: db}
}

func (r *milestoneRepository) Create(milestone *entity.Milestone) error {
	return r.db.Create(milestone).Error
}

func (r *milestoneRepository) FindByID(id uint64) (*entity.Milestone, error) {
	var milestone entity.Milestone
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&milestone).Error
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *milestoneRepository) FindByProjectID(projectID uint64) ([]entity.Milestone, error) {
	var milestones []entity.Milestone
	err := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID).
		Order("order_index ASC").
		Find(&milestones).Error
	return milestones, err
}

func (r *milestoneRepository) UpdateWithVersionCheck(milestone *entity.Milestone, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Milestone{}).
		Where("id = ? AND version = ?", milestone.ID, expectedVersion).
		Updates(map[string]interface{}{
			"name":        milestone.Name,
			"description": milestone.Description,
			"start_date":  milestone.StartDate,
			"end_date":    milestone.EndDate,
			"status":      milestone.Status,
			"version":     gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}

func (r *milestoneRepository) UpdateOrderIndex(id uint64, orderIndex int) error {
	return r.db.Model(&entity.Milestone{}).
		Where("id = ?", id).
		Update("order_index", orderIndex).Error
}

func (r *milestoneRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.Milestone{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}
