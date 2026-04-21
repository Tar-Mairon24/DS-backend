ALTER TABLE users DROP COLUMN phone;
ALTER TABLE users DROP COLUMN mfa_activated;
ALTER TABLE users CHANGE COLUMN role role ENUM('admin', 'agente') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL;