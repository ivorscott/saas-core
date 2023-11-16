ALTER TABLE users ENABLE ROW LEVEL SECURITY;

CREATE POLICY users_isolation_policy ON users
    USING (tenant_id = (SELECT current_setting('app.current_tenant')));