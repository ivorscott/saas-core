CREATE TABLE users (
    id UUID PRIMARY KEY not null,
    auth0_id varchar(128) not null unique,
    email varchar(64) not null,
    email_verified boolean not null,
    first_name varchar(255),
    last_name varchar(255),
    picture varchar(255),
    locale varchar(8),
    created timestamp without time zone default (now() at time zone 'utc')
);