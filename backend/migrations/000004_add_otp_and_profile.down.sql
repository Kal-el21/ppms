ALTER TABLE users
    DROP COLUMN IF EXISTS profile_photo_url,
    DROP COLUMN IF EXISTS two_fa_enabled,
    DROP COLUMN IF EXISTS email_notification_enabled;

DROP TABLE IF EXISTS otp_tokens;