-- =========================================================
-- Phase 1: Project creation portfolio schema
-- - Extend project request metadata from portfolio form
-- - Extend project and budget portfolio fields
-- - Add annual CAPEX/OPEX ceiling table
-- =========================================================

-- =========================================================
-- PROJECT REQUESTS
-- =========================================================
ALTER TABLE project_requests
    ADD COLUMN IF NOT EXISTS category VARCHAR(100),
    ADD COLUMN IF NOT EXISTS initiation_type VARCHAR(32) CHECK (
        initiation_type IS NULL OR initiation_type IN ('NEW_INITIATIVE', 'RENEWAL', 'ENHANCEMENT')
    ),
    ADD COLUMN IF NOT EXISTS priority VARCHAR(16) NOT NULL DEFAULT 'MEDIUM' CHECK (
        priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')
    ),
    ADD COLUMN IF NOT EXISTS proposed_start_date DATE,
    ADD COLUMN IF NOT EXISTS proposed_end_date DATE,
    ADD COLUMN IF NOT EXISTS budget_type VARCHAR(16) CHECK (
        budget_type IS NULL OR budget_type IN ('CAPEX', 'OPEX')
    ),
    ADD COLUMN IF NOT EXISTS budget_name VARCHAR(200),
    ADD COLUMN IF NOT EXISTS notes TEXT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_project_requests_proposed_dates'
    ) THEN
        ALTER TABLE project_requests
            ADD CONSTRAINT chk_project_requests_proposed_dates CHECK (
                proposed_start_date IS NULL
                OR proposed_end_date IS NULL
                OR proposed_end_date >= proposed_start_date
            );
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_project_requests_initiation_type ON project_requests(initiation_type);
CREATE INDEX IF NOT EXISTS idx_project_requests_priority ON project_requests(priority);
CREATE INDEX IF NOT EXISTS idx_project_requests_budget_type ON project_requests(budget_type);
CREATE INDEX IF NOT EXISTS idx_project_requests_proposed_end_date ON project_requests(proposed_end_date);

-- =========================================================
-- PROJECT REQUEST REVISIONS
-- =========================================================
ALTER TABLE project_request_revisions
    ADD COLUMN IF NOT EXISTS category VARCHAR(100),
    ADD COLUMN IF NOT EXISTS initiation_type VARCHAR(32) CHECK (
        initiation_type IS NULL OR initiation_type IN ('NEW_INITIATIVE', 'RENEWAL', 'ENHANCEMENT')
    ),
    ADD COLUMN IF NOT EXISTS priority VARCHAR(16) NOT NULL DEFAULT 'MEDIUM' CHECK (
        priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')
    ),
    ADD COLUMN IF NOT EXISTS proposed_start_date DATE,
    ADD COLUMN IF NOT EXISTS proposed_end_date DATE,
    ADD COLUMN IF NOT EXISTS budget_type VARCHAR(16) CHECK (
        budget_type IS NULL OR budget_type IN ('CAPEX', 'OPEX')
    ),
    ADD COLUMN IF NOT EXISTS budget_name VARCHAR(200),
    ADD COLUMN IF NOT EXISTS notes TEXT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_project_request_revisions_proposed_dates'
    ) THEN
        ALTER TABLE project_request_revisions
            ADD CONSTRAINT chk_project_request_revisions_proposed_dates CHECK (
                proposed_start_date IS NULL
                OR proposed_end_date IS NULL
                OR proposed_end_date >= proposed_start_date
            );
    END IF;
END $$;

-- =========================================================
-- PROJECTS
-- =========================================================
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS category VARCHAR(100),
    ADD COLUMN IF NOT EXISTS initiation_type VARCHAR(32) CHECK (
        initiation_type IS NULL OR initiation_type IN ('NEW_INITIATIVE', 'RENEWAL', 'ENHANCEMENT')
    ),
    ADD COLUMN IF NOT EXISTS priority VARCHAR(16) NOT NULL DEFAULT 'MEDIUM' CHECK (
        priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')
    ),
    ADD COLUMN IF NOT EXISTS notes TEXT;

CREATE INDEX IF NOT EXISTS idx_projects_category ON projects(category);
CREATE INDEX IF NOT EXISTS idx_projects_initiation_type ON projects(initiation_type);
CREATE INDEX IF NOT EXISTS idx_projects_priority ON projects(priority);

-- =========================================================
-- BUDGETS
-- =========================================================
ALTER TABLE budgets
    ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS deleted_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS budget_type VARCHAR(16) CHECK (
        budget_type IS NULL OR budget_type IN ('CAPEX', 'OPEX')
    ),
    ADD COLUMN IF NOT EXISTS budget_name VARCHAR(200);

CREATE INDEX IF NOT EXISTS idx_budgets_budget_type ON budgets(budget_type);
CREATE INDEX IF NOT EXISTS idx_budgets_deleted_by ON budgets(deleted_by);

-- =========================================================
-- PORTFOLIO ANNUAL BUDGET CEILINGS
-- =========================================================
CREATE TABLE IF NOT EXISTS portfolio_budget_years (
    id BIGSERIAL PRIMARY KEY,

    year INTEGER NOT NULL UNIQUE,

    capex_ceiling DECIMAL(18,2) NOT NULL DEFAULT 0,
    opex_ceiling DECIMAL(18,2) NOT NULL DEFAULT 0,

    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    version INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT chk_portfolio_budget_years_year CHECK (year BETWEEN 2000 AND 2100),
    CONSTRAINT chk_portfolio_budget_years_non_negative CHECK (
        capex_ceiling >= 0 AND opex_ceiling >= 0
    )
);
