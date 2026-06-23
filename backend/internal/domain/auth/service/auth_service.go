package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/Kal-el21/backend/configs"
	"github.com/Kal-el21/backend/internal/domain/auth/dto"
	"github.com/Kal-el21/backend/internal/domain/auth/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/auth/errors"
	"github.com/Kal-el21/backend/internal/domain/auth/jwt"
	"github.com/Kal-el21/backend/internal/domain/auth/repository"
	userentity "github.com/Kal-el21/backend/internal/domain/user/entity"
	userrepo "github.com/Kal-el21/backend/internal/domain/user/repository"
	"golang.org/x/crypto/bcrypt"
	gormerrors "gorm.io/gorm"
)

type AuthService interface {
	Login(req dto.LoginRequest, ipAddress, deviceInfo string) (*LoginResult, error)
	RefreshToken(req dto.RefreshTokenRequest) (*RefreshResult, error)
	Logout(req dto.LogoutRequest) error
	ChangePassword(userID uint64, req dto.ChangePasswordRequest) error
	RevokeAllSessions(userID uint64, reason string) error
}

type authService struct {
	userRepo    userrepo.UserRepository
	sessionRepo repository.SessionRepository
	cfg         *configs.Config
}

func NewAuthService(userRepo userrepo.UserRepository, sessionRepo repository.SessionRepository, cfg *configs.Config) AuthService {
	return &authService{userRepo: userRepo, sessionRepo: sessionRepo, cfg: cfg}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
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
		return nil, domainerrors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domainerrors.ErrInvalidCredentials
	}

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
			return nil // already invalid/revoked, treat as success
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

	// Security best practice: revoke semua sesi lain setelah ganti password
	return s.sessionRepo.RevokeAllByUserID(userID, "password_changed")
}

func (s *authService) RevokeAllSessions(userID uint64, reason string) error {
	return s.sessionRepo.RevokeAllByUserID(userID, reason)
}

var _ = userentity.User{} // keep import referenced if unused elsewhere

// LoginResult membawa token mentah dari service ke handler. Token TIDAK
// pernah masuk ke dto.LoginResponse yang dikirim sebagai JSON body — hanya
// dipakai secara internal agar handler bisa menulis Set-Cookie.
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         dto.UserSummary
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}
