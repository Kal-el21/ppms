package repository

import (
	"github.com/Kal-el21/backend/internal/domain/notification/entity"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	FindByUserID(userID uint64, page, limit int, unreadOnly bool) ([]entity.Notification, int64, error)
	CountUnread(userID uint64) (int64, error)
	MarkAsRead(id uint64, userID uint64) error
	MarkAllAsRead(userID uint64) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByUserID(userID uint64, page, limit int, unreadOnly bool) ([]entity.Notification, int64, error) {
	var notifications []entity.Notification
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = false")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, total, err
}

func (r *notificationRepository) CountUnread(userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(id uint64, userID uint64) error {
	return r.db.Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", id, userID). // pastikan user hanya bisa mark notifikasi miliknya sendiri
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uint64) error {
	return r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}
