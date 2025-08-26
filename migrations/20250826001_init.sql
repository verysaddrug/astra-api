-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    login VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE documents (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mime VARCHAR(64) NOT NULL,
    file BOOLEAN NOT NULL,
    public BOOLEAN NOT NULL,
    owner UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    grants TEXT[],
    json_data JSONB
);
-- +goose Down
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS users;
