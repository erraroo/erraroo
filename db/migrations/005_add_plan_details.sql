ALTER TABLE plans ADD COLUMN name text NOT NULL DEFAULT 'default';
ALTER TABLE plans ADD COLUMN price_in_cents int NOT NULL DEFAULT 2900;
ALTER TABLE plans ADD COLUMN rate_limit_duration_seconds int NOT NULL DEFAULT 60;
ALTER TABLE plans ADD COLUMN rate_limit int NOT NULL DEFAULT 60;
ALTER TABLE plans ADD COLUMN projects_limit int NOT NULL DEFAULT 3;
