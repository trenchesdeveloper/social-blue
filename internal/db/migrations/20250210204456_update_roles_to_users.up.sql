UPDATE roles
SET description = ''
WHERE description IS NULL;

-- 2. Alter the 'description' column to set it as NOT NULL.
ALTER TABLE roles
ALTER COLUMN description SET NOT NULL;