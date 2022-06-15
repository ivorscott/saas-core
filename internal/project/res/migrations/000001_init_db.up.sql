CREATE TABLE projects (
    project_id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
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

CREATE INDEX idx_project_tenant ON projects(tenant_id);
CREATE INDEX idx_project_team ON projects(team_id);
CREATE INDEX idx_project_name ON projects(name);
CREATE INDEX idx_project_public ON projects(public);


CREATE TABLE columns (
    column_id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    title VARCHAR(36) NOT NULL,
    column_name VARCHAR(8) NOT NULL,
    task_ids TEXT[],
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY(project_id) REFERENCES projects (project_id)
);
CREATE INDEX idx_column_tenant ON columns(tenant_id);
CREATE INDEX idx_column_project ON columns(project_id);

CREATE TABLE tasks (
    task_id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
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
CREATE INDEX idx_task_tenant ON tasks(tenant_id);
CREATE INDEX idx_task_project ON tasks(project_id);
CREATE INDEX idx_task_key ON tasks(key);
CREATE INDEX idx_task_title ON tasks(title);

CREATE TABLE comments (
    comment_id VARCHAR(36) PRIMARY KEY,
    task_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    content TEXT,
    likes INT,
    user_id VARCHAR(36) NOT NULL,
    edited BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);
CREATE INDEX idx_comment_tenant ON comments(tenant_id);
CREATE INDEX idx_comment_task ON comments(task_id);

-- role and memberships are redundant copies of user service data
CREATE TYPE ROLE AS ENUM ('administrator', 'editor', 'commenter','viewer');
CREATE TABLE IF NOT EXISTS memberships (
   membership_id VARCHAR(36) PRIMARY KEY,
   tenant_id VARCHAR(36) NOT NULL,
   user_id VARCHAR(36) NOT NULL,
   team_id VARCHAR(36) NOT NULL,
   role ROLE DEFAULT 'editor',
   created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
   updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);
CREATE INDEX idx_membership_tenant ON memberships(tenant_id);
CREATE INDEX idx_membership_user ON memberships(user_id);
CREATE INDEX idx_membership_team ON memberships(team_id);

-- enable RLS
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE columns ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE comments ENABLE ROW LEVEL SECURITY;
ALTER TABLE memberships ENABLE ROW LEVEL SECURITY;

-- create policies
CREATE POLICY projects_isolation_policy ON projects
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY columns_isolation_policy ON columns
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY tasks_isolation_policy ON tasks
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY comments_isolation_policy ON comments
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY memberships_isolation_policy ON memberships
    USING (tenant_id = current_setting('app.current_tenant'));