CREATE TABLE files
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id    INTEGER REFERENCES profiles(id) ON DELETE CASCADE,
    file_name  VARCHAR(255) NOT NULL,
    file_url   VARCHAR(255) NOT NULL
);

