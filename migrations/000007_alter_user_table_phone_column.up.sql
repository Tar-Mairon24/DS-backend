ALTER TABLE users add column phone VARCHAR(20) after email;
ALTER TABLE users add column mfa_activated BOOLEAN DEFAULT FALSE after verified;
ALTER TABLE users CHANGE COLUMN role role ENUM('admin', 'agente', 'owner', 'client') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL;