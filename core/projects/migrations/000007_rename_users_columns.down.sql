ALTER TABLE users
RENAME COLUMN email_verified to emailverified;

ALTER TABLE users
RENAME COLUMN users_id to id;

ALTER TABLE users
RENAME COLUMN first_name to firstname;

ALTER TABLE users
RENAME COLUMN last_name to lastname;

ALTER TABLE users
RENAME COLUMN auth0_id to auth0id;