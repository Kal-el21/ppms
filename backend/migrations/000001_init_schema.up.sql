-- =========================================================
-- PPMS Initial Schema Migration
-- =========================================================

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- =========================================================
-- DIVISIONS
-- =========================================================
CREATE TABLE divisions (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR NOT NULL UNIQUE,
    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT
);

-- =========================================================
-- USERS
-- =========================================================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    division_id BIGINT REFERENCES divisions(id),

    full_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,

    system_role VARCHAR NOT NULL CHECK (system_role IN ('ADMIN', 'USER', 'VIEWER')),

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_users_division_id ON users(division_id);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

ALTER TABLE users
    ADD CONSTRAINT fk_users_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

ALTER TABLE divisions
    ADD CONSTRAINT fk_divisions_deleted_by FOREIGN KEY (deleted_by) REFERENCES users(id);

-- =========================================================
-- USER SESSIONS
-- =========================================================
CREATE TABLE user_sessions (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id),

    refresh_token_hash VARCHAR NOT NULL,
    device_info TEXT,
    ip_address VARCHAR,

    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    revoked_reason VARCHAR,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token_hash ON user_sessions(refresh_token_hash);

-- =========================================================
-- PROJECT REQUESTS
-- =========================================================
CREATE TABLE project_requests (
    id BIGSERIAL PRIMARY KEY,

    requester_id BIGINT NOT NULL REFERENCES users(id),

    title VARCHAR NOT NULL,
    description TEXT,

    business_goal TEXT,
    expected_outcome TEXT,

    estimated_budget DECIMAL(18,2),

    status VARCHAR NOT NULL DEFAULT 'DRAFT' CHECK (status IN (
        'DRAFT', 'SUBMITTED', 'UNDER_REVIEW', 'APPROVED', 'REJECTED', 'REVISED'
    )),

    submitted_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id),

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_project_requests_requester_id ON project_requests(requester_id);
CREATE INDEX idx_project_requests_status ON project_requests(status);
CREATE INDEX idx_project_requests_deleted_at ON project_requests(deleted_at);
CREATE INDEX idx_project_requests_title_trgm ON project_requests USING gin (title gin_trgm_ops);

-- =========================================================
-- PROJECT REQUEST REVISIONS
-- =========================================================
CREATE TABLE project_request_revisions (
    id BIGSERIAL PRIMARY KEY,

    project_request_id BIGINT NOT NULL REFERENCES project_requests(id),

    revision_number INTEGER NOT NULL,

    title VARCHAR,
    description TEXT,

    business_goal TEXT,
    expected_outcome TEXT,

    estimated_budget DECIMAL(18,2),

    revision_reason TEXT,

    revised_by BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT uq_request_revision_number UNIQUE (project_request_id, revision_number)
);

CREATE INDEX idx_project_request_revisions_request_id ON project_request_revisions(project_request_id);

-- =========================================================
-- PROJECT REQUEST APPROVALS
-- =========================================================
CREATE TABLE project_request_approvals (
    id BIGSERIAL PRIMARY KEY,

    project_request_id BIGINT NOT NULL REFERENCES project_requests(id),

    -- Traceability: jika action = REQUEST_REVISION dan request direvisi ulang,
    -- revision_id menghubungkan approval ini ke revisi yang dihasilkannya.
    revision_id BIGINT REFERENCES project_request_revisions(id),

    reviewed_by BIGINT NOT NULL REFERENCES users(id),

    action VARCHAR NOT NULL CHECK (action IN ('APPROVED', 'REJECTED', 'REQUEST_REVISION')),

    comment TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_project_request_approvals_request_id ON project_request_approvals(project_request_id);
CREATE INDEX idx_project_request_approvals_reviewed_by ON project_request_approvals(reviewed_by);

-- =========================================================
-- PROJECTS
-- =========================================================
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,

    project_request_id BIGINT REFERENCES project_requests(id),

    name VARCHAR NOT NULL,
    description TEXT,

    start_date DATE,
    end_date DATE,

    status VARCHAR NOT NULL DEFAULT 'PLANNED' CHECK (status IN (
        'PLANNED', 'ACTIVE', 'ON_HOLD', 'COMPLETED', 'CANCELLED'
    )),

    created_by BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id),

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX idx_projects_name_trgm ON projects USING gin (name gin_trgm_ops);

-- =========================================================
-- PROJECT MEMBERS
-- =========================================================
CREATE TABLE project_members (
    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL REFERENCES projects(id),
    user_id BIGINT NOT NULL REFERENCES users(id),

    project_role VARCHAR NOT NULL CHECK (project_role IN (
        'PROJECT_MANAGER', 'MEMBER', 'OBSERVER'
    )),

    status VARCHAR NOT NULL DEFAULT 'ACTIVE' CHECK (status IN (
        'ACTIVE', 'SUSPENDED', 'LEFT', 'REMOVED'
    )),

    joined_at TIMESTAMP NOT NULL DEFAULT now(),
    left_at TIMESTAMP,

    status_changed_by BIGINT REFERENCES users(id),
    status_changed_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT uq_project_member UNIQUE (project_id, user_id)
);

CREATE INDEX idx_project_members_project_id ON project_members(project_id);
CREATE INDEX idx_project_members_user_id ON project_members(user_id);
CREATE INDEX idx_project_members_status ON project_members(status);

-- =========================================================
-- MILESTONES
-- =========================================================
CREATE TABLE milestones (
    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL REFERENCES projects(id),

    name VARCHAR NOT NULL,
    description TEXT,

    order_index INTEGER NOT NULL DEFAULT 0,

    start_date DATE,
    end_date DATE,

    status VARCHAR NOT NULL DEFAULT 'PLANNED' CHECK (status IN (
        'PLANNED', 'ACTIVE', 'COMPLETED', 'CANCELLED'
    )),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id),

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_milestones_project_id ON milestones(project_id);
CREATE INDEX idx_milestones_status ON milestones(status);
CREATE INDEX idx_milestones_deleted_at ON milestones(deleted_at);

-- =========================================================
-- TASKS
-- =========================================================
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL REFERENCES projects(id),
    milestone_id BIGINT REFERENCES milestones(id),

    title VARCHAR NOT NULL,
    description TEXT,

    priority VARCHAR NOT NULL DEFAULT 'MEDIUM' CHECK (priority IN (
        'LOW', 'MEDIUM', 'HIGH', 'URGENT'
    )),

    status VARCHAR NOT NULL DEFAULT 'TODO' CHECK (status IN (
        'TODO', 'IN_PROGRESS', 'DONE', 'CANCELLED'
    )),

    progress INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),

    order_index INTEGER NOT NULL DEFAULT 0,

    start_date DATE,
    due_date DATE,

    created_by BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id),

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_milestone_id ON tasks(milestone_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);
CREATE INDEX idx_tasks_title_trgm ON tasks USING gin (title gin_trgm_ops);

-- =========================================================
-- TASK ASSIGNEES
-- =========================================================
CREATE TABLE task_assignees (
    id BIGSERIAL PRIMARY KEY,

    task_id BIGINT NOT NULL REFERENCES tasks(id),
    user_id BIGINT NOT NULL REFERENCES users(id),

    assigned_by BIGINT NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMP NOT NULL DEFAULT now(),

    -- Histori assignment tetap utuh untuk audit/reporting,
    -- bukan physical delete saat unassign.
    unassigned_at TIMESTAMP,
    unassigned_by BIGINT REFERENCES users(id)
);

CREATE INDEX idx_task_assignees_task_id ON task_assignees(task_id);
CREATE INDEX idx_task_assignees_user_id ON task_assignees(user_id);

-- Hanya satu assignment aktif (belum unassigned) per task-user
CREATE UNIQUE INDEX uq_task_assignee_active
    ON task_assignees(task_id, user_id)
    WHERE unassigned_at IS NULL;

-- =========================================================
-- TASK COMMENTS
-- =========================================================
CREATE TABLE task_comments (
    id BIGSERIAL PRIMARY KEY,

    task_id BIGINT NOT NULL REFERENCES tasks(id),
    user_id BIGINT NOT NULL REFERENCES users(id),

    comment TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);

-- =========================================================
-- BUDGETS
-- =========================================================
CREATE TABLE budgets (
    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL UNIQUE REFERENCES projects(id),

    allocated_budget DECIMAL(18,2) NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id),

    version INTEGER NOT NULL DEFAULT 1
);

-- =========================================================
-- BUDGET TRANSACTIONS
-- (IMMUTABLE: tidak ada UPDATE/DELETE setelah insert)
-- =========================================================
CREATE TABLE budget_transactions (
    id BIGSERIAL PRIMARY KEY,

    budget_id BIGINT NOT NULL REFERENCES budgets(id),

    type VARCHAR NOT NULL CHECK (type IN ('EXPENSE', 'REFUND', 'ADJUSTMENT')),

    adjustment_type VARCHAR CHECK (adjustment_type IN (
        'ERROR_CORRECTION', 'BUDGET_REALLOCATION', 'AUDIT_CORRECTION', 'MANUAL_OVERRIDE'
    )),

    amount DECIMAL(18,2) NOT NULL,

    reason TEXT,
    description TEXT,

    transaction_date TIMESTAMP NOT NULL DEFAULT now(),

    -- Mencegah duplicate transaction akibat retry / double-click
    idempotency_key VARCHAR UNIQUE,

    created_by BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP NOT NULL DEFAULT now(),

    -- FR-08.04: ADJUSTMENT wajib memiliki adjustment_type dan reason
    CONSTRAINT chk_adjustment_requires_type_reason CHECK (
        type != 'ADJUSTMENT'
        OR (adjustment_type IS NOT NULL AND reason IS NOT NULL)
    )
);

CREATE INDEX idx_budget_transactions_budget_id ON budget_transactions(budget_id);
CREATE INDEX idx_budget_transactions_type ON budget_transactions(type);

-- =========================================================
-- ATTACHMENTS
-- =========================================================
CREATE TABLE attachments (
    id BIGSERIAL PRIMARY KEY,

    uploaded_by BIGINT NOT NULL REFERENCES users(id),

    entity_type VARCHAR NOT NULL CHECK (entity_type IN (
        'PROJECT_REQUEST', 'PROJECT', 'MILESTONE', 'TASK', 'BUDGET_TRANSACTION', 'HANDOVER'
    )),
    entity_id BIGINT NOT NULL,

    version INTEGER NOT NULL DEFAULT 1,

    file_name VARCHAR NOT NULL,
    original_name VARCHAR NOT NULL,

    file_path VARCHAR NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id)
);

CREATE INDEX idx_attachments_entity ON attachments(entity_type, entity_id);
CREATE INDEX idx_attachments_deleted_at ON attachments(deleted_at);

-- =========================================================
-- HANDOVERS
-- =========================================================
CREATE TABLE handovers (
    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL REFERENCES projects(id),

    sender_id BIGINT NOT NULL REFERENCES users(id),
    sender_division_id BIGINT REFERENCES divisions(id),

    receiver_id BIGINT REFERENCES users(id),

    description TEXT,

    delivery_date DATE,
    delivery_time TIME,

    received_at TIMESTAMP,

    status VARCHAR NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING', 'RECEIVED', 'RETURNED'
    )),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    deleted_at TIMESTAMP,
    deleted_by BIGINT REFERENCES users(id),

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_handovers_project_id ON handovers(project_id);
CREATE INDEX idx_handovers_status ON handovers(status);
CREATE INDEX idx_handovers_deleted_at ON handovers(deleted_at);

-- =========================================================
-- NOTIFICATIONS
-- =========================================================
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id),

    type VARCHAR NOT NULL,

    title VARCHAR NOT NULL,
    message TEXT,

    entity_type VARCHAR,
    entity_id BIGINT,

    action_url VARCHAR,

    -- Delivery channel & status (mempersiapkan v1.5 email notification
    -- tanpa migration besar di kemudian hari)
    channel VARCHAR NOT NULL DEFAULT 'IN_APP' CHECK (channel IN ('IN_APP', 'EMAIL')),
    delivery_status VARCHAR NOT NULL DEFAULT 'SENT' CHECK (delivery_status IN (
        'PENDING', 'SENT', 'FAILED'
    )),

    is_read BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_entity ON notifications(entity_type, entity_id);

-- =========================================================
-- NOTIFICATION PREFERENCES
-- =========================================================
CREATE TABLE notification_preferences (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id),

    type VARCHAR NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT uq_notification_preference UNIQUE (user_id, type)
);

CREATE INDEX idx_notification_preferences_user_id ON notification_preferences(user_id);

-- =========================================================
-- AUDIT LOGS
-- (Immutable: hanya INSERT, tidak ada UPDATE/DELETE endpoint)
-- =========================================================
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT REFERENCES users(id),

    module VARCHAR NOT NULL,
    action VARCHAR NOT NULL,

    entity_type VARCHAR,
    entity_id BIGINT,

    old_data JSON,
    new_data JSON,

    ip_address VARCHAR,
    user_agent TEXT,

    request_id VARCHAR,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_module ON audit_logs(module);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);