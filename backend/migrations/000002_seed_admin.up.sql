-- Default admin user
-- Email: admin@ppms.local
-- Password: Admin@12345 (bcrypt hash di bawah, cost 12)
-- WAJIB diganti setelah first login!

INSERT INTO users (full_name, email, password_hash, system_role, is_active)
VALUES (
    'System Administrator',
    'admin@ppms.local',
    '$2a$12$B9Ukxq8z3QYxVTDhKksiwuLZkCblkvv5V1ovQspt6.MgQUQv5MQhC',
    'ADMIN',
    true
);