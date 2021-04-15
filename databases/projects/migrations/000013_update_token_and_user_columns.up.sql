ALTER TABLE ma_token
DROP COLUMN expiration;

ALTER TABLE users
ALTER COLUMN auth0_id SET DATA TYPE varchar(128);