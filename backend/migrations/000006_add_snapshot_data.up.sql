-- =========================================================
-- Phase V1.5: Add snapshot_data to project_snapshots
-- untuk reporting berbasis snapshot lengkap
-- =========================================================
ALTER TABLE project_snapshots
    ADD COLUMN IF NOT EXISTS snapshot_data jsonb;
