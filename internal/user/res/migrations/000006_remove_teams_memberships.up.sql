DROP POLICY memberships_isolation_policy ON memberships;
DROP POLICY teams_isolation_policy ON teams;
DROP POLICY projects_isolation_policy ON projects;

ALTER TABLE invites DROP CONSTRAINT invites_team_id_fkey;

DROP TABLE memberships;
DROP TABLE teams;
DROP TABLE projects;

ALTER TABLE invites DROP COLUMN team_id;
