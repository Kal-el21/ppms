package service

import (
	"github.com/Kal-el21/backend/internal/domain/notification/entity"
	"github.com/Kal-el21/backend/internal/domain/notification/repository"
)

type CreateNotificationParams struct {
	UserID     uint64
	Type       string
	Title      string
	Message    string
	EntityType string
	EntityID   *uint64
	ActionURL  string
}

type NotificationService interface {
	Create(params CreateNotificationParams) error
	GetByUserID(userID uint64, page, limit int, unreadOnly bool) ([]entity.Notification, int64, error)
	CountUnread(userID uint64) (int64, error)
	MarkAsRead(id uint64, userID uint64) error
	MarkAllAsRead(userID uint64) error
}

type notificationService struct {
	repo     repository.NotificationRepository
	prefRepo repository.PreferenceRepository
}

func NewNotificationService(repo repository.NotificationRepository, prefRepo repository.PreferenceRepository) NotificationService {
	return &notificationService{repo: repo, prefRepo: prefRepo}
}

// Create menghormati notification_preferences: jika user sudah mute type ini,
// notifikasi tidak dibuat sama sekali (bukan dibuat lalu disembunyikan).
func (s *notificationService) Create(params CreateNotificationParams) error {
	enabled, err := s.prefRepo.IsEnabled(params.UserID, params.Type)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}

	notification := &entity.Notification{
		UserID:         params.UserID,
		Type:           params.Type,
		Title:          params.Title,
		Message:        params.Message,
		EntityType:     params.EntityType,
		EntityID:       params.EntityID,
		ActionURL:      params.ActionURL,
		Channel:        entity.ChannelInApp, // EMAIL channel diaktifkan di v1.5
		DeliveryStatus: entity.DeliverySent,
	}

	return s.repo.Create(notification)
}

func (s *notificationService) GetByUserID(userID uint64, page, limit int, unreadOnly bool) ([]entity.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindByUserID(userID, page, limit, unreadOnly)
}

func (s *notificationService) CountUnread(userID uint64) (int64, error) {
	return s.repo.CountUnread(userID)
}

func (s *notificationService) MarkAsRead(id uint64, userID uint64) error {
	return s.repo.MarkAsRead(id, userID)
}

func (s *notificationService) MarkAllAsRead(userID uint64) error {
	return s.repo.MarkAllAsRead(userID)
}
