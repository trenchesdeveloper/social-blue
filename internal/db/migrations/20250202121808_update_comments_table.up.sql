-- Up migration: Add NOT NULL constraint to created_at and updated_at columns
ALTER TABLE comments
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;