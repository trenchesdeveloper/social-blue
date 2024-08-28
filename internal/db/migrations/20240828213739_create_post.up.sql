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
                       username TEXT NOT NULL,
                       password TEXT NOT NULL,
                       email TEXT NOT NULL UNIQUE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);