-- Index untuk query FindBySystemRole (notify-all-admin) dan filter role lainnya
CREATE INDEX idx_users_system_role ON users(system_role) WHERE deleted_at IS NULL;

-- Index tambahan untuk overdue tasks query (Dashboard Phase 6)
CREATE INDEX idx_tasks_due_date_status ON tasks(due_date, status) WHERE deleted_at IS NULL;