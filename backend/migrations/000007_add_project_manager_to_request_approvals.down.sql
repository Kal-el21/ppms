DROP INDEX IF EXISTS idx_project_request_approvals_project_manager_id;

ALTER TABLE project_request_approvals
    DROP COLUMN IF EXISTS project_manager_id;
