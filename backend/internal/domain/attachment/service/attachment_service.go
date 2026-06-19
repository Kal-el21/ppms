package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/Kal-el21/backend/internal/domain/attachment/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/attachment/errors"
	"github.com/Kal-el21/backend/internal/domain/attachment/repository"
	"github.com/Kal-el21/backend/internal/domain/attachment/validator"
	"github.com/Kal-el21/backend/internal/infrastructure/minio"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/google/uuid"
)

type AttachmentService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, entityType string, entityID uint64, uploadedBy uint64, isAdmin bool) (*entity.Attachment, error)
	GetByEntity(entityType string, entityID uint64, userID uint64, isAdmin bool) ([]entity.Attachment, error)
	GetDownloadURL(ctx context.Context, id uint64, userID uint64, isAdmin bool) (string, *entity.Attachment, error)
	GetVersionHistory(id uint64, userID uint64, isAdmin bool) ([]entity.Attachment, error)
	Delete(id uint64, deletedBy uint64, isAdmin bool) error
}

type attachmentService struct {
	repo           repository.AttachmentRepository
	minioClient    *minio.Client
	resolver       OwnershipResolver
	requestChecker RequestOwnershipChecker
}

// RequestOwnershipChecker interface kecil untuk validasi ownership PROJECT_REQUEST
// tanpa circular import ke domain project_request secara langsung.
type RequestOwnershipChecker interface {
	IsRequestOwner(requestID uint64, userID uint64) (bool, error)
}

func NewAttachmentService(
	repo repository.AttachmentRepository,
	minioClient *minio.Client,
	resolver OwnershipResolver,
	requestChecker RequestOwnershipChecker,
) AttachmentService {
	return &attachmentService{repo: repo, minioClient: minioClient, resolver: resolver, requestChecker: requestChecker}
}

// checkAccess memvalidasi akses, dengan exception untuk PROJECT_REQUEST
// yang pakai ownership berbasis requester (bukan project membership),
// karena request belum punya project sampai approved.
func (s *attachmentService) checkAccess(entityType string, entityID uint64, userID uint64, isAdmin bool) error {
	if isAdmin {
		return nil
	}

	if entityType == "PROJECT_REQUEST" {
		isOwner, err := s.requestChecker.IsRequestOwner(entityID, userID)
		if err != nil {
			return err
		}
		if !isOwner {
			return apperrors.New(apperrors.ErrResourceNotOwned, "you are not the owner of this project request")
		}
		return nil
	}

	return s.resolver.ValidateAccess(entityType, entityID, userID, isAdmin)
}

func (s *attachmentService) Upload(ctx context.Context, file *multipart.FileHeader, entityTypeStr string, entityID uint64, uploadedBy uint64, isAdmin bool) (*entity.Attachment, error) {
	if err := validator.ValidateEntityType(entityTypeStr); err != nil {
		return nil, err
	}

	if err := s.checkAccess(entityTypeStr, entityID, uploadedBy, isAdmin); err != nil {
		return nil, err
	}

	if err := validator.ValidateFileSize(file.Size); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Magic-bytes validation (item Phase 7 #2) — baca 512 byte pertama
	// untuk deteksi MIME type aktual, bukan hanya percaya Content-Type header
	// yang bisa dipalsukan client.
	detectedMime, err := validator.DetectAndValidateMimeType(src)
	if err != nil {
		return nil, err
	}

	// Reset reader ke awal setelah dibaca untuk deteksi (DetectAndValidateMimeType
	// membaca beberapa byte pertama, file perlu di-seek balik sebelum upload)
	if _, err := src.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file reader: %w", err)
	}

	existingVersions, _ := s.repo.FindVersionHistory(
		entity.EntityType(entityTypeStr), entityID, file.Filename,
	)
	newVersion := len(existingVersions) + 1

	objectName := fmt.Sprintf("%s/%d/%s_%s", entityTypeStr, entityID, uuid.New().String(), file.Filename)

	if err := s.minioClient.Upload(ctx, objectName, src, file.Size, detectedMime); err != nil {
		return nil, err
	}

	attachment := &entity.Attachment{
		UploadedBy:   uploadedBy,
		EntityType:   entity.EntityType(entityTypeStr),
		EntityID:     entityID,
		Version:      newVersion,
		FileName:     objectName,
		OriginalName: file.Filename,
		FilePath:     objectName,
		FileSize:     file.Size,
		MimeType:     detectedMime,
	}

	if err := s.repo.Create(attachment); err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *attachmentService) GetByEntity(entityTypeStr string, entityID uint64, userID uint64, isAdmin bool) ([]entity.Attachment, error) {
	if err := validator.ValidateEntityType(entityTypeStr); err != nil {
		return nil, err
	}

	if err := s.checkAccess(entityTypeStr, entityID, userID, isAdmin); err != nil {
		return nil, err
	}

	return s.repo.FindByEntity(entity.EntityType(entityTypeStr), entityID)
}

func (s *attachmentService) GetDownloadURL(ctx context.Context, id uint64, userID uint64, isAdmin bool) (string, *entity.Attachment, error) {
	attachment, err := s.repo.FindByID(id)
	if err != nil {
		return "", nil, domainerrors.ErrAttachmentNotFound
	}

	if err := s.checkAccess(string(attachment.EntityType), attachment.EntityID, userID, isAdmin); err != nil {
		return "", nil, err
	}

	url, err := s.minioClient.GetPresignedDownloadURL(ctx, attachment.FilePath)
	if err != nil {
		return "", nil, err
	}

	return url, attachment, nil
}

func (s *attachmentService) GetVersionHistory(id uint64, userID uint64, isAdmin bool) ([]entity.Attachment, error) {
	attachment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, domainerrors.ErrAttachmentNotFound
	}

	if err := s.checkAccess(string(attachment.EntityType), attachment.EntityID, userID, isAdmin); err != nil {
		return nil, err
	}

	return s.repo.FindVersionHistory(attachment.EntityType, attachment.EntityID, attachment.OriginalName)
}

func (s *attachmentService) Delete(id uint64, deletedBy uint64, isAdmin bool) error {
	attachment, err := s.repo.FindByID(id)
	if err != nil {
		return domainerrors.ErrAttachmentNotFound
	}

	if err := s.checkAccess(string(attachment.EntityType), attachment.EntityID, deletedBy, isAdmin); err != nil {
		return err
	}

	if err := s.minioClient.Delete(context.Background(), attachment.FilePath); err != nil {
		_ = err
	}

	return s.repo.Delete(id, deletedBy)
}
