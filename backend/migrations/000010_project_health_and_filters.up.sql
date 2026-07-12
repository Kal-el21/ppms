-- =========================================================
-- Phase 5: Project health column + list filter indexes
-- Add stored `health` (GREEN/YELLOW/RED) computed from status,
-- end_date and progress. Backfilled using status + overdue logic
-- (project progress is derived from milestones, not stored here).
-- =========================================================

ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS health VARCHAR(16) CHECK (
        health IS NULL OR health IN ('GREEN', 'YELLOW', 'RED')
    );

CREATE INDEX IF NOT EXISTS idx_projects_health ON projects(health);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_end_date ON projects(end_date);

-- Backfill existing rows. Progress is not stored on the project row, so we
-- approximate health from status and overdue end_date (matches calculateProjectHealth
-- for the no-progress-known case).
UPDATE projects
SET health = CASE
    WHEN status = 'COMPLETED' THEN 'GREEN'
    WHEN status = 'CANCELLED' THEN 'RED'
    WHEN status = 'ON_HOLD' THEN 'YELLOW'
    WHEN end_date IS NOT NULL AND end_date < CURRENT_DATE THEN 'RED'
    ELSE 'GREEN'
END
WHERE health IS NULL;
