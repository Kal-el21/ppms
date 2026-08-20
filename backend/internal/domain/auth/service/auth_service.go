package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Kal-el21/backend/configs"
	"github.com/Kal-el21/backend/internal/domain/auth/dto"
	"github.com/Kal-el21/backend/internal/domain/auth/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/auth/errors"
	"github.com/Kal-el21/backend/internal/domain/auth/jwt"
	"github.com/Kal-el21/backend/internal/domain/auth/repository"
	userentity "github.com/Kal-el21/backend/internal/domain/user/entity"
	userrepo "github.com/Kal-el21/backend/internal/domain/user/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/logger"
	"golang.org/x/crypto/bcrypt"
	gormerrors "gorm.io/gorm"
)

// EmailSenderIface adalah interface kecil untuk dependency injection email
// ke dalam auth service, tanpa circular import ke infrastructure/email package.
// Nama sengaja menggunakan suffix "Iface" untuk membedakan dari package name.
type EmailSenderIface interface {
	SendOTP(to string, recipientName string, otp string) error
}

type AuthService interface {
	Login(req dto.LoginRequest, ipAddress, deviceInfo string) (*LoginResult, error)
	InitLogin(req dto.LoginRequest, ipAddress, deviceInfo string) (*dto.LoginInitResponse, error)
	VerifyOTPAndLogin(req dto.VerifyOTPRequest, ipAddress, deviceInfo string) (*LoginResult, error)
	ResendOTP(req dto.ResendOTPRequest) error
	RefreshToken(req dto.RefreshTokenRequest) (*RefreshResult, error)
	Logout(req dto.LogoutRequest) error
	ChangePassword(userID uint64, req dto.ChangePasswordRequest) error
	RevokeAllSessions(userID uint64, reason string) error
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         dto.UserSummary
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type authService struct {
	userRepo    userrepo.UserRepository
	sessionRepo repository.SessionRepository
	otpRepo     repository.OTPRepository
	otpSessions *OTPSessionStore
	emailSvc    EmailSenderIface // gunakan tipe lokal, bukan package.Type
	cfg         *configs.Config
}

func NewAuthService(
	userRepo userrepo.UserRepository,
	sessionRepo repository.SessionRepository,
	otpRepo repository.OTPRepository,
	otpSessions *OTPSessionStore,
	emailSvc EmailSenderIface,
	cfg *configs.Config,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		otpRepo:     otpRepo,
		otpSessions: otpSessions,
		emailSvc:    emailSvc,
		cfg:         cfg,
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateOTPCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", (int(b[0])<<16|int(b[1])<<8|int(b[2]))%1000000), nil
}

func (s *authService) Login(req dto.LoginRequest, ipAddress, deviceInfo string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gormerrors.ErrRecordNotFound) {
			return nil, domainerrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		logger.Log.Warn().
			Uint64("user_id", user.ID).
			Str("email", user.Email).
			Msg("auth: login attempt for inactive user")
		return nil, domainerrors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Log.Warn().
			Uint64("user_id", user.ID).
			Str("email", user.Email).
			Err(err).
			Msg("auth: invalid password")
		return nil, domainerrors.ErrInvalidCredentials
	}

	return s.createSessionAndTokens(user, ipAddress, deviceInfo)
}

func (s *authService) InitLogin(req dto.LoginRequest, ipAddress, deviceInfo string) (*dto.LoginInitResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gormerrors.ErrRecordNotFound) {
			return nil, domainerrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		logger.Log.Warn().
			Uint64("user_id", user.ID).
			Str("email", user.Email).
			Msg("auth: init-login attempt for inactive user")
		return nil, domainerrors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Log.Warn().
			Uint64("user_id", user.ID).
			Str("email", user.Email).
			Err(err).
			Msg("auth: invalid password on init-login")
		return nil, domainerrors.ErrInvalidCredentials
	}

	if !user.TwoFAEnabled {
		result, err := s.createSessionAndTokens(user, ipAddress, deviceInfo)
		if err != nil {
			return nil, err
		}
		summary := result.User
		return &dto.LoginInitResponse{
			TwoFARequired: false,
			User:          &summary,
		}, nil
	}

	otpCode, err := generateOTPCode()
	if err != nil {
		return nil, err
	}

	otpHash, err := bcrypt.GenerateFromPassword([]byte(otpCode), 10)
	if err != nil {
		return nil, err
	}

	otpRecord := &entity.OTPToken{
		UserID:    user.ID,
		TokenHash: string(otpHash),
		Purpose:   entity.OTPLogin,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IPAddress: ipAddress,
	}
	if err := s.otpRepo.Create(otpRecord); err != nil {
		return nil, err
	}

	go func() {
		_ = s.emailSvc.SendOTP(user.Email, user.FullName, otpCode)
	}()

	sessionToken, err := s.otpSessions.Create(user.ID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &dto.LoginInitResponse{
		TwoFARequired:   true,
		OTPSessionToken: sessionToken,
	}, nil
}

func (s *authService) VerifyOTPAndLogin(req dto.VerifyOTPRequest, ipAddress, deviceInfo string) (*LoginResult, error) {
	userID, ok := s.otpSessions.Consume(req.OTPSessionToken)
	if !ok {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "otp session expired or invalid, please login again")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, domainerrors.ErrInvalidCredentials
	}

	otpRecord, err := s.otpRepo.FindActiveByUserAndPurpose(userID, entity.OTPLogin)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "otp expired or already used")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(otpRecord.TokenHash), []byte(req.OTPCode)); err != nil {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "invalid otp code")
	}

	if err := s.otpRepo.MarkUsed(otpRecord.ID); err != nil {
		return nil, err
	}

	return s.createSessionAndTokens(user, ipAddress, deviceInfo)
}

func (s *authService) ResendOTP(req dto.ResendOTPRequest) error {
	userID, ok := s.otpSessions.Peek(req.OTPSessionToken)
	if !ok {
		return apperrors.New(apperrors.ErrUnauthorized, "otp session expired, please login again")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	otpCode, err := generateOTPCode()
	if err != nil {
		return err
	}

	otpHash, err := bcrypt.GenerateFromPassword([]byte(otpCode), 10)
	if err != nil {
		return err
	}

	otpRecord := &entity.OTPToken{
		UserID:    userID,
		TokenHash: string(otpHash),
		Purpose:   entity.OTPLogin,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.otpRepo.Create(otpRecord); err != nil {
		return err
	}

	go func() {
		_ = s.emailSvc.SendOTP(user.Email, user.FullName, otpCode)
	}()

	return nil
}

func (s *authService) RefreshToken(req dto.RefreshTokenRequest) (*RefreshResult, error) {
	claims, err := jwt.ValidateRefreshToken(req.RefreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return nil, domainerrors.ErrInvalidToken
	}

	tokenHash := hashToken(req.RefreshToken)
	session, err := s.sessionRepo.FindByRefreshTokenHash(tokenHash)
	if err != nil {
		if errors.Is(err, gormerrors.ErrRecordNotFound) {
			return nil, domainerrors.ErrSessionRevoked
		}
		return nil, err
	}

	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, domainerrors.ErrInvalidToken
	}

	if !user.IsActive {
		return nil, domainerrors.ErrUserInactive
	}

	if err := s.sessionRepo.RevokeByID(session.ID, "rotated"); err != nil {
		return nil, err
	}

	newAccessToken, err := jwt.GenerateAccessToken(user.ID, string(user.SystemRole), user.DivisionID, s.cfg.JWTAccessSecret, s.cfg.JWTAccessExpiryMinutes)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshExpiryDays)
	if err != nil {
		return nil, err
	}

	newSession := &entity.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: hashToken(newRefreshToken),
		DeviceInfo:       session.DeviceInfo,
		IPAddress:        session.IPAddress,
		ExpiresAt:        time.Now().AddDate(0, 0, s.cfg.JWTRefreshExpiryDays),
	}

	if err := s.sessionRepo.Create(newSession); err != nil {
		return nil, err
	}

	_ = claims

	return &RefreshResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) Logout(req dto.LogoutRequest) error {
	tokenHash := hashToken(req.RefreshToken)
	session, err := s.sessionRepo.FindByRefreshTokenHash(tokenHash)
	if err != nil {
		if errors.Is(err, gormerrors.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return s.sessionRepo.RevokeByID(session.ID, "logout")
}

func (s *authService) ChangePassword(userID uint64, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return domainerrors.ErrWrongOldPassword
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), s.cfg.BcryptCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(newHash)
	user.Version++

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return s.sessionRepo.RevokeAllByUserID(userID, "password_changed")
}

func (s *authService) RevokeAllSessions(userID uint64, reason string) error {
	return s.sessionRepo.RevokeAllByUserID(userID, reason)
}

func (s *authService) createSessionAndTokens(user *userentity.User, ipAddress, deviceInfo string) (*LoginResult, error) {
	accessToken, err := jwt.GenerateAccessToken(user.ID, string(user.SystemRole), user.DivisionID, s.cfg.JWTAccessSecret, s.cfg.JWTAccessExpiryMinutes)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshExpiryDays)
	if err != nil {
		return nil, err
	}

	session := &entity.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: hashToken(refreshToken),
		DeviceInfo:       deviceInfo,
		IPAddress:        ipAddress,
		ExpiresAt:        time.Now().AddDate(0, 0, s.cfg.JWTRefreshExpiryDays),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserSummary{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			SystemRole: string(user.SystemRole),
			DivisionID: user.DivisionID,
		},
	}, nil
}
