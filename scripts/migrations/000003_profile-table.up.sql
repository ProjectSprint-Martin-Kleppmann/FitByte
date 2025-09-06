CREATE TABLE profiles
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    name       VARCHAR(255) DEFAULT '',
    image_uri  VARCHAR(500) DEFAULT '',
    preference VARCHAR(255) DEFAULT '',
    weight_unit VARCHAR(10) DEFAULT '',
    height_unit VARCHAR(10) DEFAULT '',
    weight     DECIMAL(5,2) DEFAULT 0,
    height     DECIMAL(5,2) DEFAULT 0
);

-- Add index for soft deletes (GORM requirement)
CREATE INDEX idx_profiles_deleted_at ON profiles(deleted_at);

-- Add index for email lookups
CREATE INDEX idx_profiles_email ON profiles(email);