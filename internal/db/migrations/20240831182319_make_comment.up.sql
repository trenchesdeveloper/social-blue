CREATE  TABLE  IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    post_id BIGSERIAL NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()

);

ALTER TABLE comments
    ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users (id);

ALTER TABLE comments
    ADD CONSTRAINT fk_post
    FOREIGN KEY (post_id)
    REFERENCES posts (id);