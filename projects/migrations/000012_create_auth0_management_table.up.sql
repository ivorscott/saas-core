CREATE TABLE ma_token (
    ma_token_id UUID PRIMARY KEY,
    token text NOT NULL,
    expiration text NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc')
);