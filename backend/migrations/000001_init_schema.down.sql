-- =========================================================
-- PPMS Initial Schema Rollback
-- Drop in reverse dependency order
-- =========================================================

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS notification_preferences;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS handovers;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS budget_transactions;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS task_comments;
DROP TABLE IF EXISTS task_assignees;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS milestones;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS project_request_approvals;
DROP TABLE IF EXISTS project_request_revisions;
DROP TABLE IF EXISTS project_requests;
DROP TABLE IF EXISTS user_sessions;

ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS fk_users_deleted_by;
ALTER TABLE IF EXISTS divisions DROP CONSTRAINT IF EXISTS fk_divisions_deleted_by;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS divisions;

DROP EXTENSION IF EXISTS pg_trgm;