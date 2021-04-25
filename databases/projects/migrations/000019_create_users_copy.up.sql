CREATE TABLE IF NOT EXISTS users (
    user_id UUID not null unique,
    auth0_id varchar(128) not null unique,
    email varchar(64) not null,
    created timestamp without time zone default (now() at time zone 'utc')
);