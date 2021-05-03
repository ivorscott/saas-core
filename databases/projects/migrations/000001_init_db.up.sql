--- Enable uuid extension v9.1 and newer ( uuid_generate_v4() )
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE projects (
project_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
name VARCHAR(36) NOT NULL,
prefix VARCHAR(4) NOT NULL,
description TEXT,
user_id UUID,
team_id VARCHAR(36), ---> CAN BE NULL, represents uuid field in users service
active BOOLEAN DEFAULT TRUE,
public BOOLEAN DEFAULT FALSE,
column_order TEXT ARRAY[10],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE columns (
column_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
project_id UUID,
title VARCHAR(36) NOT NULL,
column_name VARCHAR(8) NOT NULL,
task_ids TEXT[],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY(project_id) REFERENCES projects (project_id)
);

CREATE TABLE tasks (
task_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
project_id UUID,
key VARCHAR(10),
seq SERIAL,
title VARCHAR(48) NOT NULL,
points INT DEFAULT 0,
content TEXT,
assigned_to VARCHAR(36), ---> CAN BE NULL, represents uuid field in users service
attachments TEXT[],
comments TEXT[],
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
FOREIGN KEY(project_id) REFERENCES projects (project_id)
);

CREATE TABLE comments (
comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
content TEXT,
likes INT,
user_id UUID,
edited BOOLEAN DEFAULT FALSE,
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

-- redundant data
CREATE TYPE ROLE AS ENUM ('administrator', 'editor', 'commenter','viewer');
CREATE TABLE IF NOT EXISTS memberships (
membership_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID,
team_id UUID,
role ROLE DEFAULT 'editor',
created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);
