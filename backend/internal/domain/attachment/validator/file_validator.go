package validator

import (
	"io"
	"net/http"

	domainerrors "github.com/Kal-el21/backend/internal/domain/attachment/errors"
)

const maxFileSize = 25 * 1024 * 1024 // 25 MB

// allowedMimeTypes: whitelist final, dicocokkan terhadap hasil deteksi
// magic-bytes (http.DetectContentType), BUKAN terhadap Content-Type header
// dari client yang bisa dipalsukan.
var allowedMimeTypes = map[string]bool{
	"application/pdf":    true,
	"image/jpeg":         true,
	"image/png":          true,
	"image/webp":         true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"text/csv":   true,
	"text/plain": true,
	// application/zip muncul untuk docx/xlsx modern (OOXML adalah ZIP container)
	// jika http.DetectContentType tidak bisa membaca signature spesifik OOXML;
	// kita tetap whitelist agar tidak false-reject file docx/xlsx valid.
	"application/zip": true,
}

func ValidateFileSize(size int64) error {
	if size > maxFileSize {
		return domainerrors.ErrFileTooLarge
	}
	return nil
}

// DetectAndValidateMimeType membaca maksimal 512 byte pertama file (cukup
// untuk http.DetectContentType bekerja akurat), mendeteksi MIME type
// sebenarnya berdasarkan magic bytes, lalu memvalidasi terhadap whitelist.
// Mengembalikan MIME type yang TERDETEKSI (bukan yang diklaim client) agar
// disimpan ke database — mencegah penyimpanan metadata yang menyesatkan.
func DetectAndValidateMimeType(reader io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)

	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	detectedMime := http.DetectContentType(buffer[:n])

	// http.DetectContentType mengembalikan parameter tambahan terkadang
	// (misal "text/plain; charset=utf-8"), kita ambil bagian utama saja
	// untuk dicocokkan ke whitelist.
	mainType := detectedMime
	for i, ch := range detectedMime {
		if ch == ';' {
			mainType = detectedMime[:i]
			break
		}
	}

	if !allowedMimeTypes[mainType] {
		return "", domainerrors.ErrUnsupportedMime
	}

	return mainType, nil
}

func ValidateEntityType(entityType string) error {
	valid := map[string]bool{
		"PROJECT_REQUEST": true, "PROJECT": true, "MILESTONE": true,
		"TASK": true, "BUDGET_TRANSACTION": true, "HANDOVER": true,
	}
	if !valid[entityType] {
		return domainerrors.ErrInvalidEntityType
	}
	return nil
}
