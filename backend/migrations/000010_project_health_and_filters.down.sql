-- Revert 000010_project_health_and_filters
DROP INDEX IF EXISTS idx_projects_health;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_end_date;

ALTER TABLE projects DROP COLUMN IF EXISTS health;
