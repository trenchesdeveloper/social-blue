CREATE TABLE IF NOT EXISTS followers (
                                         user_id BIGINT NOT NULL,
                                         follower_id BIGINT NOT NULL,
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                         PRIMARY KEY (user_id, follower_id),
                                         FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
                                         FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
);