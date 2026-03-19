CREATE TABLE IF NOT EXISTS images (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    property_id INT UNSIGNED NOT NULL,
    url VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    main_image BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    main_image_property_key INT UNSIGNED GENERATED ALWAYS AS (
        CASE
            WHEN main_image = TRUE AND deleted_at IS NULL THEN property_id
            ELSE NULL
        END
    ) STORED,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    INDEX idx_property_id (property_id),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX uq_images_one_active_main_per_property (main_image_property_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;