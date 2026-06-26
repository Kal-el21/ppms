-- Down migration for 000006_add_snapshot_data
ALTER TABLE project_snapshots
    DROP COLUMN IF EXISTS snapshot_data;
