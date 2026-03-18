CREATE DATABASE IF NOT EXISTS inmosoftDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE inmosoftDB;

-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'agente') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Properties table
CREATE TABLE IF NOT EXISTS properties (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    listing_date TIMESTAMP NULL,
    address VARCHAR(500) NOT NULL,
    neighborhood VARCHAR(255),
    city VARCHAR(255) NOT NULL,
    zone VARCHAR(255),
    reference VARCHAR(500),
    price DECIMAL(15, 2) NOT NULL,
    construction_m2 INT DEFAULT 0,
    land_m2 INT DEFAULT 0,
    is_occupied BOOLEAN DEFAULT FALSE,
    is_furnished BOOLEAN DEFAULT FALSE,
    floors INT DEFAULT 1,
    bedrooms INT DEFAULT 0,
    bathrooms INT DEFAULT 0,
    garage_size INT DEFAULT 0,
    garden_m2 INT DEFAULT 0,
    gas_types JSON,
    amenities JSON,
    extras JSON,
    utilities JSON,
    notes LONGTEXT,
    owner_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    property_type VARCHAR(255) NOT NULL,
    transaction_type VARCHAR(255) NOT NULL,
    status VARCHAR(255) DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_user_id (user_id),
    INDEX idx_owner_id (owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create RefreshTokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id VARCHAR(255) PRIMARY KEY,
    token VARCHAR(500) UNIQUE NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create VerificationTokens table (for email verification)
CREATE TABLE IF NOT EXISTS verification_tokens (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    used BOOLEAN DEFAULT FALSE,
    resends INT DEFAULT 0,
    motive VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_properties join table (many-to-many)
CREATE TABLE IF NOT EXISTS user_properties (
    user_id INT UNSIGNED NOT NULL,
    property_id INT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, property_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    INDEX idx_property_id (property_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;