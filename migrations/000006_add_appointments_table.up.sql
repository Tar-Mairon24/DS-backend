CREATE TABLE appointments (
id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
title VARCHAR(255) NOT NULL,
description TEXT,
start_date DATETIME NOT NULL,
end_date DATETIME NOT NULL,
status enum ('scheduled', 'completed', 'canceled') DEFAULT 'scheduled',
notes MEDIUMTEXT NULL,
id_client INT UNSIGNED NOT NULL,
id_property INT UNSIGNED NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
deleted_at TIMESTAMP NULL,
INDEX idx_id_client (id_client),
INDEX idx_id_property (id_property),
INDEX idx_deleted_at (deleted_at),
CONSTRAINT fk_appointments_client
FOREIGN KEY (id_client) REFERENCES users(id) ON DELETE CASCADE,
CONSTRAINT fk_appointments_property
FOREIGN KEY (id_property) REFERENCES properties(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE appointment_agents (
appointment_id INT UNSIGNED NOT NULL,
user_id INT UNSIGNED NOT NULL,
PRIMARY KEY (appointment_id, user_id),
INDEX idx_aa_user_id (user_id),
CONSTRAINT fk_appointment_agents_appointment
FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
CONSTRAINT fk_appointment_agents_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;