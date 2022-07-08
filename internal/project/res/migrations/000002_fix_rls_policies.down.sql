REASSIGN OWNED BY user_a TO postgres;
DROP OWNED BY user_a;
DROP USER user_a;

DROP POLICY memberships_isolation_policy ON memberships;
DROP POLICY comments_isolation_policy ON comments;
DROP POLICY tasks_isolation_policy ON tasks;
DROP POLICY columns_isolation_policy ON columns;
DROP POLICY projects_isolation_policy ON projects;

-- create policies
CREATE POLICY memberships_isolation_policy ON memberships
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY comments_isolation_policy ON comments
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY tasks_isolation_policy ON tasks
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY columns_isolation_policy ON columns
    USING (tenant_id = current_setting('app.current_tenant'));

CREATE POLICY projects_isolation_policy ON projects
    USING (tenant_id = current_setting('app.current_tenant'));