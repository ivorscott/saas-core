ALTER TABLE seats ENABLE ROW LEVEL SECURITY;

CREATE POLICY seats_isolation_policy ON seats
    USING (tenant_id = (SELECT current_setting('app.current_tenant')));
