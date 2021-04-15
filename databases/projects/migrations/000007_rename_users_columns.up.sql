ALTER TABLE users
RENAME COLUMN emailverified to email_verified;

ALTER TABLE users
RENAME COLUMN id to user_id;

ALTER TABLE users
RENAME COLUMN firstname to first_name;

ALTER TABLE users
RENAME COLUMN lastname to last_name;

ALTER TABLE users
RENAME COLUMN auth0id to auth0_id;