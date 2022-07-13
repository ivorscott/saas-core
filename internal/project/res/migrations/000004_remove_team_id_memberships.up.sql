DROP POLICY memberships_isolation_policy ON memberships;

DROP TABLE memberships;

ALTER TABLE projects DROP COLUMN team_id;
