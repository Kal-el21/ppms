package dto

import "time"

type UploadAttachmentResponse struct {
	ID           uint64    `json:"id"`
	EntityType   string    `json:"entity_type"`
	EntityID     uint64    `json:"entity_id"`
	OriginalName string    `json:"original_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	Version      int       `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
}

type AttachmentDownloadResponse struct {
	DownloadURL string `json:"download_url"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
}
