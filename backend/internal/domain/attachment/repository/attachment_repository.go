package repository

import (
	"github.com/Kal-el21/backend/internal/domain/attachment/entity"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(attachment *entity.Attachment) error
	FindByID(id uint64) (*entity.Attachment, error)
	FindByEntity(entityType entity.EntityType, entityID uint64) ([]entity.Attachment, error)
	FindVersionHistory(entityType entity.EntityType, entityID uint64, fileName string) ([]entity.Attachment, error)
	Delete(id uint64, deletedBy uint64) error
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(attachment *entity.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) FindByID(id uint64) (*entity.Attachment, error) {
	var attachment entity.Attachment
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

// FindByEntity menggunakan composite index (entity_type, entity_id) yang sudah
// didefinisikan di migration Phase 0 untuk query performant.
func (r *attachmentRepository) FindByEntity(entityType entity.EntityType, entityID uint64) ([]entity.Attachment, error) {
	var attachments []entity.Attachment
	err := r.db.Where("entity_type = ? AND entity_id = ? AND deleted_at IS NULL", entityType, entityID).
		Order("created_at DESC").
		Find(&attachments).Error
	return attachments, err
}

func (r *attachmentRepository) FindVersionHistory(entityType entity.EntityType, entityID uint64, fileName string) ([]entity.Attachment, error) {
	var attachments []entity.Attachment
	err := r.db.Where("entity_type = ? AND entity_id = ? AND original_name = ?", entityType, entityID, fileName).
		Order("version DESC").
		Find(&attachments).Error
	return attachments, err
}

func (r *attachmentRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.Attachment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}
