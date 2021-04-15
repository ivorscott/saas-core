ALTER TABLE users
ADD COLUMN emailverified boolean not null,
ADD COLUMN locale varchar(8);