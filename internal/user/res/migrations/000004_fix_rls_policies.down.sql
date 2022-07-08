REASSIGN OWNED BY user_a TO postgres;
DROP OWNED BY user_a;
DROP USER user_a;

DROP POLICY projects_isolation_policy ON projects;
DROP POLICY teams_isolation_policy ON teams;
DROP POLICY users_isolation_policy ON users;
DROP POLICY invites_isolation_policy ON invites;
DROP POLICY memberships_isolation_policy ON memberships;

-- create policies
CREATE POLICY projects_isolation_policy ON projects
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY teams_isolation_policy ON teams
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY users_isolation_policy ON users
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY invites_isolation_policy ON invites
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY memberships_isolation_policy ON memberships
    USING (tenant_id = current_setting('app.current_tenant'));