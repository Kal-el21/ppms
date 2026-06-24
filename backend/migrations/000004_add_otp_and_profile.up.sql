-- =========================================================
-- Phase V1.5: OTP (2FA via email), profile photo, email notification
-- =========================================================

-- Tabel OTP sementara untuk 2FA login via email.
-- Record dihapus otomatis setelah dipakai atau expired.
CREATE TABLE otp_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token_hash VARCHAR NOT NULL,             -- bcrypt hash dari 6-digit OTP
    purpose VARCHAR NOT NULL CHECK (purpose IN ('LOGIN', 'PASSWORD_RESET')),
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,                       -- NULL = belum dipakai
    ip_address VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_otp_tokens_user_id ON otp_tokens(user_id);
CREATE INDEX idx_otp_tokens_expires_at ON otp_tokens(expires_at);

-- Kolom tambahan di users: foto profile & flag 2FA aktif
ALTER TABLE users
    ADD COLUMN profile_photo_url VARCHAR,
    ADD COLUMN two_fa_enabled BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN email_notification_enabled BOOLEAN NOT NULL DEFAULT true;

-- Kolom channel di notification_preferences sudah ada sejak Phase 0
-- (channel VARCHAR NOT NULL DEFAULT 'IN_APP'), tapi kita pastikan
-- constraint-nya bisa handle EMAIL juga (sudah bisa sesuai Phase 0 schema).