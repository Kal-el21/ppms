-- Down migration for 000009_project_creation_phase1_schema

DROP TABLE IF EXISTS portfolio_budget_years;

DROP INDEX IF EXISTS idx_budgets_budget_type;
DROP INDEX IF EXISTS idx_budgets_deleted_by;
DROP INDEX IF EXISTS idx_projects_priority;
DROP INDEX IF EXISTS idx_projects_initiation_type;
DROP INDEX IF EXISTS idx_projects_category;
DROP INDEX IF EXISTS idx_project_requests_proposed_end_date;
DROP INDEX IF EXISTS idx_project_requests_budget_type;
DROP INDEX IF EXISTS idx_project_requests_priority;
DROP INDEX IF EXISTS idx_project_requests_initiation_type;

ALTER TABLE project_requests
    DROP CONSTRAINT IF EXISTS chk_project_requests_proposed_dates;

ALTER TABLE project_request_revisions
    DROP CONSTRAINT IF EXISTS chk_project_request_revisions_proposed_dates;

ALTER TABLE budgets
    DROP COLUMN IF EXISTS budget_name,
    DROP COLUMN IF EXISTS budget_type,
    DROP COLUMN IF EXISTS deleted_by,
    DROP COLUMN IF EXISTS created_by;

ALTER TABLE projects
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS initiation_type,
    DROP COLUMN IF EXISTS category;

ALTER TABLE project_request_revisions
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS budget_name,
    DROP COLUMN IF EXISTS budget_type,
    DROP COLUMN IF EXISTS proposed_end_date,
    DROP COLUMN IF EXISTS proposed_start_date,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS initiation_type,
    DROP COLUMN IF EXISTS category;

ALTER TABLE project_requests
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS budget_name,
    DROP COLUMN IF EXISTS budget_type,
    DROP COLUMN IF EXISTS proposed_end_date,
    DROP COLUMN IF EXISTS proposed_start_date,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS initiation_type,
    DROP COLUMN IF EXISTS category;
