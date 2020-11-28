CREATE TABLE users (
    user_id UUID PRIMARY KEY not null,
    auth0_id varchar(24) not null unique,
    email varchar(64) not null,
    email_verified boolean not null,
    first_name varchar(255),
    last_name varchar(255),
    picture varchar(255),
    created timestamp without time zone default (now() at time zone 'utc')
);