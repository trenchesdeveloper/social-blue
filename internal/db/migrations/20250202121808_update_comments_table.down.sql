-- Down migration: Remove NOT NULL constraint from created_at and updated_at columns
ALTER TABLE comments
ALTER COLUMN created_at DROP NOT NULL,
ALTER COLUMN updated_at DROP NOT NULL;