-- Create users table (matching gorm.Model structure)
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL
);

-- Create index for deleted_at (used by GORM for soft deletes)
CREATE INDEX idx_users_deleted_at ON users(deleted_at);