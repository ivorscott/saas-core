DROP POLICY memberships_isolation_policy ON memberships;
DROP POLICY comments_isolation_policy ON comments;
DROP POLICY tasks_isolation_policy ON tasks;
DROP POLICY columns_isolation_policy ON columns;
DROP POLICY projects_isolation_policy ON projects;

ALTER TABLE projects DISABLE ROW LEVEL SECURITY;
ALTER TABLE columns DISABLE ROW LEVEL SECURITY;
ALTER TABLE tasks DISABLE ROW LEVEL SECURITY;
ALTER TABLE comments DISABLE ROW LEVEL SECURITY;
ALTER TABLE memberships DISABLE ROW LEVEL SECURITY;

DROP INDEX idx_membership_tenant, idx_membership_user, idx_membership_team;
DROP TABLE IF EXISTS memberships;
DROP TYPE ROLE;

DROP INDEX idx_comment_tenant, idx_comment_task;
DROP TABLE IF EXISTS comments;

DROP INDEX idx_task_tenant, idx_task_project, idx_task_key, idx_task_title;
DROP TABLE IF EXISTS tasks;

DROP INDEX idx_column_tenant, idx_column_project;
DROP TABLE IF EXISTS columns;

DROP INDEX idx_project_tenant, idx_project_team, idx_project_name, idx_project_public;
DROP TABLE IF EXISTS projects;