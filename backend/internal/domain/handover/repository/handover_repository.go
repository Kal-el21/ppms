package repository

import (
	"github.com/Kal-el21/backend/internal/domain/handover/entity"
	"gorm.io/gorm"
)

type HandoverRepository interface {
	Create(handover *entity.Handover) error
	FindByID(id uint64) (*entity.Handover, error)
	FindByProjectID(projectID uint64, statusFilter string) ([]entity.Handover, error)
	UpdateWithVersionCheck(handover *entity.Handover, expectedVersion int) (int64, error)
}

type handoverRepository struct {
	db *gorm.DB
}

func NewHandoverRepository(db *gorm.DB) HandoverRepository {
	return &handoverRepository{db: db}
}

func (r *handoverRepository) Create(handover *entity.Handover) error {
	return r.db.Create(handover).Error
}

func (r *handoverRepository) FindByID(id uint64) (*entity.Handover, error) {
	var handover entity.Handover
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&handover).Error
	if err != nil {
		return nil, err
	}
	return &handover, nil
}

func (r *handoverRepository) FindByProjectID(projectID uint64, statusFilter string) ([]entity.Handover, error) {
	var handovers []entity.Handover
	query := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID)
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}
	err := query.Order("created_at DESC").Find(&handovers).Error
	return handovers, err
}

func (r *handoverRepository) UpdateWithVersionCheck(handover *entity.Handover, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Handover{}).
		Where("id = ? AND version = ?", handover.ID, expectedVersion).
		Updates(map[string]interface{}{
			"receiver_id":   handover.ReceiverID,
			"description":   handover.Description,
			"delivery_date": handover.DeliveryDate,
			"delivery_time": handover.DeliveryTime,
			"received_at":   handover.ReceivedAt,
			"status":        handover.Status,
			"version":       gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}
