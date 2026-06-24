package service

import (
	"fmt"

	"github.com/Kal-el21/backend/internal/domain/notification/entity"
	"github.com/Kal-el21/backend/internal/domain/notification/repository"
	userentity "github.com/Kal-el21/backend/internal/domain/user/entity"
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

// EmailSender interface kecil untuk menghindari circular import
type EmailSender interface {
	SendNotification(to string, subject string, htmlBody string) error
}

// UserEmailProvider interface untuk ambil email & preference user
type UserEmailProvider interface {
	FindByID(id uint64) (*userentity.User, error)
}

type NotificationService interface {
	Create(params CreateNotificationParams) error
	GetByUserID(userID uint64, page, limit int, unreadOnly bool) ([]entity.Notification, int64, error)
	CountUnread(userID uint64) (int64, error)
	MarkAsRead(id uint64, userID uint64) error
	MarkAllAsRead(userID uint64) error
}

type notificationService struct {
	repo         repository.NotificationRepository
	prefRepo     repository.PreferenceRepository
	emailSender  EmailSender
	userProvider UserEmailProvider
}

func NewNotificationService(
	repo repository.NotificationRepository,
	prefRepo repository.PreferenceRepository,
	emailSender EmailSender,
	userProvider UserEmailProvider,
) NotificationService {
	return &notificationService{
		repo:         repo,
		prefRepo:     prefRepo,
		emailSender:  emailSender,
		userProvider: userProvider,
	}
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

	// Simpan notifikasi in-app
	notification := &entity.Notification{
		UserID:         params.UserID,
		Type:           params.Type,
		Title:          params.Title,
		Message:        params.Message,
		EntityType:     params.EntityType,
		EntityID:       params.EntityID,
		ActionURL:      params.ActionURL,
		Channel:        entity.ChannelInApp,
		DeliveryStatus: entity.DeliverySent,
	}

	if err := s.repo.Create(notification); err != nil {
		return err
	}

	// Kirim email jika user mengaktifkan email notification
	go func() {
		user, err := s.userProvider.FindByID(params.UserID)
		if err != nil || !user.EmailNotificationEnabled {
			return
		}

		htmlBody := fmt.Sprintf(`
			<div style="font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px;">
				<div style="background:linear-gradient(135deg,#2563EB,#DC2626);padding:16px 24px;border-radius:8px 8px 0 0;">
					<p style="color:white;font-size:16px;font-weight:700;margin:0;">PPMS</p>
				</div>
				<div style="background:#fff;border:1px solid #E2E8F0;border-top:none;padding:24px;border-radius:0 0 8px 8px;">
					<p style="font-size:16px;font-weight:600;color:#0F172A;margin:0 0 8px;">%s</p>
					<p style="font-size:14px;color:#64748B;margin:0 0 20px;line-height:1.6;">%s</p>
					<a href="%s" style="display:inline-block;background:#2563EB;color:white;text-decoration:none;padding:10px 20px;border-radius:6px;font-size:13px;font-weight:600;">
						Lihat Detail
					</a>
				</div>
			</div>`,
			params.Title,
			params.Message,
			params.ActionURL,
		)

		_ = s.emailSender.SendNotification(user.Email, params.Title, htmlBody)
	}()

	return nil
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
