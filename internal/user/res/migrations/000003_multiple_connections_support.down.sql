ALTER TABLE projects DROP CONSTRAINT projects_user_id_fkey;
ALTER TABLE invites DROP CONSTRAINT invites_user_id_fkey;
ALTER TABLE memberships DROP CONSTRAINT memberships_user_id_fkey;
ALTER TABLE teams DROP CONSTRAINT teams_user_id_fkey;
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (user_id);
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id);
ALTER TABLE invites ADD CONSTRAINT invites_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id);
ALTER TABLE memberships ADD CONSTRAINT memberships_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id);
ALTER TABLE teams ADD CONSTRAINT teams_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id);

ALTER TABLE users
    ADD COLUMN email VARCHAR(64),
    ADD COLUMN email_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN first_name VARCHAR(255),
    ADD COLUMN last_name VARCHAR(255),
    ADD COLUMN picture VARCHAR(255),
    ADD COLUMN locale VARCHAR(8);

UPDATE users
 SET email = user_profiles.email,
     email_verified = user_profiles.email_verified,
     first_name = user_profiles.first_name,
     last_name = user_profiles.last_name,
     picture = user_profiles.picture,
     locale = user_profiles.locale
FROM users u
INNER JOIN user_profiles USING(user_id);

ALTER TABLE users
ALTER COLUMN email SET NOT NULL;

CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_user_last_name ON users(last_name);

DROP INDEX idx_user_profile_email, idx_user_profile_first_name, idx_user_profile_last_name;
DROP TABLE IF EXISTS user_profiles;
