CREATE TABLE IF NOT EXISTS user_profiles (
    user_id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(64) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    picture VARCHAR(255),
    locale VARCHAR(8),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE INDEX idx_user_email ON user_profiles(email);
CREATE INDEX idx_user_last_name ON user_profiles(last_name);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(36),
    tenant_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    PRIMARY KEY (user_id, tenant_id)
);

CREATE TABLE IF NOT EXISTS seats (
    tenant_id VARCHAR(36) PRIMARY KEY,
    seats_used INTEGER DEFAULT 0,
    max_seats INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS invites(
    invite_id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    accepted BOOLEAN DEFAULT FALSE,
    expiration TIMESTAMP WITHOUT TIME ZONE DEFAULT ((NOW() + '5 days') AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (user_id) REFERENCES user_profiles (user_id)
);

CREATE INDEX idx_invite_tenant ON invites(tenant_id);
CREATE INDEX idx_invite_user ON invites(user_id);

-- enable RLS
ALTER TABLE invites ENABLE ROW LEVEL SECURITY;

-- create policies
CREATE POLICY invites_isolation_policy ON invites
    USING (tenant_id = (SELECT current_setting('app.current_tenant')));

CREATE USER user_a WITH PASSWORD 'postgres';
GRANT ALL ON ALL TABLES IN SCHEMA "public" TO user_a;