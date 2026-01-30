BEGIN;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name         text NOT NULL,
    email        text NOT NULL,
    password     text NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    archived_at  timestamptz
);

-- Unique email index
CREATE UNIQUE INDEX idx_users_unique_email_active
    ON users (email)
    WHERE archived_at IS NULL;

-- Index for archived lookups
CREATE INDEX idx_users_archived_at
    ON users (archived_at);


CREATE TABLE IF NOT EXISTS todos (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       text NOT NULL,
    body        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    valid_till  timestamptz NOT NULL,
    complete   boolean NOT NULL DEFAULT false
);

CREATE INDEX idx_todos_user_id
    ON todos (user_id);


CREATE TABLE IF NOT EXISTS sessions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    expires_at  timestamptz NOT NULL
);

CREATE INDEX idx_sessions_user_id
    ON sessions (user_id);

-- fast cleanup of expired sessions
CREATE INDEX idx_sessions_expires_at
    ON sessions (expires_at);

COMMIT
