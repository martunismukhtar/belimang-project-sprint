-- Create role enum type
CREATE TYPE user_role AS ENUM ('user', 'admin');

-- Users table without FK to other services (default schema)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(30) NOT NULL CHECK (char_length(username) >= 5),
    email VARCHAR(255) NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    password VARCHAR(255) NOT NULL CHECK (char_length(password) >= 5 AND char_length(password) <= 30),
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Helpful indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_role_unique ON users (email, role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
