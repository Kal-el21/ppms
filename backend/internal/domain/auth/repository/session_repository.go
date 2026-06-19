package repository

import (
	"time"

	"github.com/Kal-el21/backend/internal/domain/auth/entity"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *entity.UserSession) error
	FindByRefreshTokenHash(hash string) (*entity.UserSession, error)
	RevokeByID(id uint64, reason string) error
	RevokeAllByUserID(userID uint64, reason string) error
	CountFailedAttempts(email string, since time.Time) (int64, error)
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *entity.UserSession) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) FindByRefreshTokenHash(hash string) (*entity.UserSession, error) {
	var session entity.UserSession
	err := r.db.Where("refresh_token_hash = ? AND revoked_at IS NULL AND expires_at > now()", hash).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) RevokeByID(id uint64, reason string) error {
	return r.db.Model(&entity.UserSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"revoked_at":     gorm.Expr("now()"),
			"revoked_reason": reason,
		}).Error
}

func (r *sessionRepository) RevokeAllByUserID(userID uint64, reason string) error {
	return r.db.Model(&entity.UserSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]interface{}{
			"revoked_at":     gorm.Expr("now()"),
			"revoked_reason": reason,
		}).Error
}

// CountFailedAttempts: placeholder untuk rate limiting.
// Implementasi penuh memerlukan tabel log percobaan login terpisah
// atau in-memory store (Redis). Untuk Phase 1, rate limiting
// dilakukan di middleware level menggunakan in-memory counter (lihat 3.8).
func (r *sessionRepository) CountFailedAttempts(email string, since time.Time) (int64, error) {
	return 0, nil
}
