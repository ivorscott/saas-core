--- Enable uuid extension v9.1 and newer ( uuid_generate_v4() )
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
auth0_id VARCHAR(128) NOT NULL,
email VARCHAR(64) NOT NULL,
email_verified BOOLEAN DEFAULT FALSE,
first_name VARCHAR(255),
last_name VARCHAR(255),
picture VARCHAR(255),
locale VARCHAR(8),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE IF NOT EXISTS teams (
team_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,
name VARCHAR(32) NOT NULL,
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TYPE ROLE AS ENUM ('administrator', 'editor', 'commenter','viewer');

CREATE TABLE IF NOT EXISTS memberships (
membership_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,
team_id UUID NOT NULL,
role ROLE DEFAULT 'editor',
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY (user_id) REFERENCES users (user_id),
FOREIGN KEY (team_id) REFERENCES teams (team_id)
);

CREATE TABLE IF NOT EXISTS invites(
invite_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,
team_id UUID NOT NULL,
read BOOLEAN DEFAULT FALSE,
accepted BOOLEAN DEFAULT FALSE,
expiration TIMESTAMP WITHOUT TIME ZONE DEFAULT ((NOW() + '5 days') AT TIME ZONE 'utc'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY (user_id) REFERENCES users (user_id),
FOREIGN KEY (team_id) REFERENCES teams (team_id)
);

CREATE TABLE ma_token (
ma_token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
access_token TEXT NOT NULL,
scope TEXT NOT NULL,
expires_in INT NOT NULL,
token_type VARCHAR(16),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE projects (
project_id UUID PRIMARY KEY,
name VARCHAR(36) NOT NULL,
user_id UUID NOT NULL,
team_id VARCHAR(36),
active BOOLEAN DEFAULT TRUE,
public BOOLEAN DEFAULT FALSE,
column_order TEXT ARRAY[10],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY (user_id) REFERENCES users (user_id)
);

