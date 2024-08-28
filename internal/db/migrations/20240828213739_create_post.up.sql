CREATE EXTENSION IF NOT EXISTS "citext";
CREATE TABLE posts (
                       id BIGSERIAL PRIMARY KEY,
                       content TEXT NOT NULL,
                       title TEXT NOT NULL,
                       user_id BIGINT NOT NULL,
                       tags TEXT[],
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       first_name VARCHAR(255) NOT NULL,
                       last_name VARCHAR(255) NOT NULL,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password bytea NOT NULL,
                       email CITEXT NOT NULL UNIQUE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE posts
    ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users (id);
CREATE INDEX idx_users_email ON users(email);