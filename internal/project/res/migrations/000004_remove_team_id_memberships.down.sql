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

ALTER TABLE memberships ENABLE ROW LEVEL SECURITY;

CREATE POLICY memberships_isolation_policy ON memberships
    USING (tenant_id = current_setting('app.current_tenant'));

ALTER TABLE projects ADD COLUMN team_id VARCHAR(36);