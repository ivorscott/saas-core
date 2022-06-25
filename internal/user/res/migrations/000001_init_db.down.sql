DROP POLICY memberships_isolation_policy ON memberships;
DROP POLICY invites_isolation_policy ON invites;
DROP POLICY users_isolation_policy ON users;
DROP POLICY teams_isolation_policy ON teams;
DROP POLICY projects_isolation_policy ON projects;

ALTER TABLE memberships DISABLE ROW LEVEL SECURITY;
ALTER TABLE invites DISABLE ROW LEVEL SECURITY;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
ALTER TABLE teams DISABLE ROW LEVEL SECURITY;
ALTER TABLE projects DISABLE ROW LEVEL SECURITY;

DROP INDEX idx_project_tenant, idx_project_team, idx_project_user;
DROP TABLE IF EXISTS projects;

DROP INDEX idx_invite_tenant, idx_invite_user;
DROP TABLE IF EXISTS invites;

DROP INDEX idx_membership_tenant, idx_memberships_user;
DROP TABLE IF EXISTS memberships;
DROP TYPE ROLE;

DROP INDEX idx_team_tenant, idx_team_name;
DROP TABLE IF EXISTS teams;

DROP INDEX idx_user_tenant, idx_user_email, idx_user_last_name;
DROP TABLE IF EXISTS users;