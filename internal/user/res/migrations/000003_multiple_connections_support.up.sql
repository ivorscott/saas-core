-- Create new table and insert data
SELECT
    user_id, email, email_verified, first_name,
    last_name, picture, locale, created_at, updated_at
INTO user_profiles
FROM users;

ALTER TABLE user_profiles
ADD PRIMARY KEY (user_id);

DROP INDEX IF EXISTS idx_user_email, idx_user_last_name;
CREATE INDEX idx_user_profile_email ON user_profiles(email);
CREATE INDEX idx_user_profile_first_name ON user_profiles(first_name);
CREATE INDEX idx_user_profile_last_name ON user_profiles(last_name);

ALTER TABLE users
DROP COLUMN email,
DROP COLUMN email_verified,
DROP COLUMN first_name,
DROP COLUMN last_name,
DROP COLUMN picture,
DROP COLUMN locale;

-- We can't uniquely identified rows by user_id alone, so we create a composite primary key.
-- This allows multiple rows to have the same user_id but different tenant_id.
ALTER TABLE projects DROP CONSTRAINT projects_user_id_fkey;
ALTER TABLE invites DROP CONSTRAINT invites_user_id_fkey;
ALTER TABLE memberships DROP CONSTRAINT memberships_user_id_fkey;
ALTER TABLE teams DROP CONSTRAINT teams_user_id_fkey;
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (user_id, tenant_id);
ALTER TABLE teams ADD CONSTRAINT teams_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_profiles (user_id);
ALTER TABLE memberships ADD CONSTRAINT memberships_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_profiles (user_id);
ALTER TABLE invites ADD CONSTRAINT invites_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_profiles (user_id);
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_profiles (user_id);