-- =========================================================
-- Phase V1.5: Sinkronisasi ERD dengan Migration
-- - Tambah kolom yang belum ada di migration
-- - Buat tabel yang ada di ERD tapi belum ada di migration
-- =========================================================

-- =========================================================
-- PROJECTS: tambah kolom project_code, completed_at, updated_by
-- =========================================================
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS project_code VARCHAR UNIQUE,
    ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

CREATE INDEX IF NOT EXISTS idx_projects_project_code ON projects(project_code);

-- =========================================================
-- PROJECT_REQUESTS: tambah request_number, current_revision,
-- approved_at, rejected_at
-- =========================================================
ALTER TABLE project_requests
    ADD COLUMN IF NOT EXISTS request_number VARCHAR UNIQUE,
    ADD COLUMN IF NOT EXISTS current_revision INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS rejected_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_project_requests_request_number ON project_requests(request_number);

-- =========================================================
-- TASKS: tambah completed_at, updated_by, deleted_by
-- =========================================================
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

CREATE INDEX IF NOT EXISTS idx_tasks_deleted_by ON tasks(deleted_by);

-- =========================================================
-- BUDGETS: tambah updated_by
-- =========================================================
ALTER TABLE budgets
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- =========================================================
-- HANDOVERS: tambah updated_by, deleted_by
-- =========================================================
ALTER TABLE handovers
    ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id);

CREATE INDEX IF NOT EXISTS idx_handovers_deleted_by ON handovers(deleted_by);

-- =========================================================
-- PROJECT_SNAPSHOTS: untuk reporting berbasis snapshot
-- (sesuai SDD section 14, ERD)
-- =========================================================
CREATE TABLE project_snapshots (
    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL REFERENCES projects(id),

    snapshot_date DATE NOT NULL,

    project_status VARCHAR NOT NULL,

    project_progress DECIMAL,

    budget_allocated DECIMAL,
    budget_used DECIMAL,

    total_tasks INT,
    completed_tasks INT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_project_snapshots_project_id ON project_snapshots(project_id);
CREATE INDEX idx_project_snapshots_snapshot_date ON project_snapshots(snapshot_date);

-- =========================================================
-- APPROVAL_WORKFLOWS: untuk multi-level approval masa depan
-- =========================================================
CREATE TABLE approval_workflows (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================================================
-- APPROVAL_LEVELS: level-level dalam workflow
-- =========================================================
CREATE TABLE approval_levels (
    id BIGSERIAL PRIMARY KEY,

    workflow_id BIGINT NOT NULL REFERENCES approval_workflows(id),

    level_order INT NOT NULL,

    role_required VARCHAR NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_approval_levels_workflow_id ON approval_levels(workflow_id);
