CREATE TABLE projects (
project_id VARCHAR(36) PRIMARY KEY,
name VARCHAR(36) NOT NULL,
prefix VARCHAR(4) NOT NULL,
description TEXT,
user_id VARCHAR(36) NOT NULL,
team_id VARCHAR(36),
active BOOLEAN DEFAULT TRUE,
public BOOLEAN DEFAULT FALSE,
column_order TEXT ARRAY[10],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE columns (
column_id VARCHAR(36) PRIMARY KEY,
project_id VARCHAR(36) NOT NULL,
title VARCHAR(36) NOT NULL,
column_name VARCHAR(8) NOT NULL,
task_ids TEXT[],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY(project_id) REFERENCES projects (project_id)
);

CREATE TABLE tasks (
task_id VARCHAR(36) PRIMARY KEY,
project_id VARCHAR(36) NOT NULL,
key VARCHAR(10),
seq SERIAL,
title VARCHAR(48) NOT NULL,
points INT DEFAULT 0,
content TEXT,
assigned_to VARCHAR(36),
attachments TEXT[],
comments TEXT[],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY(project_id) REFERENCES projects (project_id)
);

CREATE TABLE comments (
comment_id VARCHAR(36) PRIMARY KEY,
content TEXT,
likes INT,
user_id VARCHAR(36) NOT NULL,
edited BOOLEAN DEFAULT FALSE,
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

-- redundant data
CREATE TYPE ROLE AS ENUM ('administrator', 'editor', 'commenter','viewer');
CREATE TABLE IF NOT EXISTS memberships (
membership_id VARCHAR(36) PRIMARY KEY,
user_id VARCHAR(36) NOT NULL,
team_id VARCHAR(36) NOT NULL,
role ROLE DEFAULT 'editor',
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);
