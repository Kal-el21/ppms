-- Down migration for 000005_sync_erd_and_missing_columns

-- Drop newly created tables
DROP TABLE IF EXISTS approval_levels CASCADE;
DROP TABLE IF EXISTS approval_workflows CASCADE;
DROP TABLE IF EXISTS project_snapshots CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS idx_approval_levels_workflow_id;
DROP INDEX IF EXISTS idx_project_snapshots_snapshot_date;
DROP INDEX IF EXISTS idx_project_snapshots_project_id;
DROP INDEX IF EXISTS idx_project_requests_request_number;
DROP INDEX IF EXISTS idx_projects_project_code;
DROP INDEX IF EXISTS idx_tasks_deleted_by;
DROP INDEX IF EXISTS idx_handovers_deleted_by;

-- Remove added columns from projects
ALTER TABLE projects
    DROP COLUMN IF EXISTS project_code,
    DROP COLUMN IF EXISTS completed_at,
    DROP COLUMN IF EXISTS updated_by;

-- Remove added columns from project_requests
ALTER TABLE project_requests
    DROP COLUMN IF EXISTS request_number,
    DROP COLUMN IF EXISTS current_revision,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS rejected_at;

-- Remove added columns from tasks
ALTER TABLE tasks
    DROP COLUMN IF EXISTS completed_at,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS deleted_by;

-- Remove added columns from budgets
ALTER TABLE budgets
    DROP COLUMN IF EXISTS updated_by;

-- Remove added columns from handovers
ALTER TABLE handovers
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS deleted_by;
