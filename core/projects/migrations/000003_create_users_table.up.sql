CREATE TABLE users (
    id UUID not null unique,
    auth0id varchar(24) not null unique,
    email varchar(64) not null,
    firstname varchar(255),
    lastname varchar(255),
    picture varchar(255),
    created timestamp without time zone default (now() at time zone 'utc')
);