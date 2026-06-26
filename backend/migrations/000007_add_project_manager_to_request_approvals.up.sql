ALTER TABLE project_request_approvals
    ADD COLUMN IF NOT EXISTS project_manager_id BIGINT REFERENCES users(id);

CREATE INDEX IF NOT EXISTS idx_project_request_approvals_project_manager_id
    ON project_request_approvals(project_manager_id);
