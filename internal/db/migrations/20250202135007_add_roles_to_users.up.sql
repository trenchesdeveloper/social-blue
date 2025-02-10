CREATE TABLE IF NOT EXISTS roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  level INT NOT NULL DEFAULT 0,
  description TEXT
);

INSERT INTO roles (name, level, description)
VALUES ('admin', 3, 'Admin role with all permissions'),
       ('user', 1, 'User role with limited permissions'),
       ('moderator', 2, 'Guest role with no permissions');

ALTER TABLE users
ADD COLUMN role_id BIGINT REFERENCES roles(id) DEFAULT 1;

UPDATE users
SET role_id = (
    SELECT id
    FROM roles
    WHERE name = 'user'
);

ALTER TABLE users
ALTER COLUMN role_id DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;
