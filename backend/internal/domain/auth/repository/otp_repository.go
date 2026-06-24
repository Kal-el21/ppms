package repository

import (
	"github.com/Kal-el21/backend/internal/domain/auth/entity"
	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(otp *entity.OTPToken) error
	FindActiveByUserAndPurpose(userID uint64, purpose entity.OTPPurpose) (*entity.OTPToken, error)
	MarkUsed(id uint64) error
	DeleteExpired() error
}

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(otp *entity.OTPToken) error {
	// Invalidate previous OTP milik user yang sama sebelum buat baru
	r.db.Model(&entity.OTPToken{}).
		Where("user_id = ? AND purpose = ? AND used_at IS NULL", otp.UserID, otp.Purpose).
		Update("used_at", gorm.Expr("now()"))
	return r.db.Create(otp).Error
}

func (r *otpRepository) FindActiveByUserAndPurpose(userID uint64, purpose entity.OTPPurpose) (*entity.OTPToken, error) {
	var otp entity.OTPToken
	err := r.db.Where(
		"user_id = ? AND purpose = ? AND used_at IS NULL AND expires_at > now()",
		userID, purpose,
	).Order("created_at DESC").First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) MarkUsed(id uint64) error {
	return r.db.Model(&entity.OTPToken{}).
		Where("id = ?", id).
		Update("used_at", gorm.Expr("now()")).Error
}

func (r *otpRepository) DeleteExpired() error {
	return r.db.Where("expires_at < now()").Delete(&entity.OTPToken{}).Error
}
