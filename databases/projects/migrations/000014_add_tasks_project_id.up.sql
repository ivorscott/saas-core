ALTER TABLE tasks ADD COLUMN project_id UUID;

ALTER TABLE tasks
ADD CONSTRAINT fk_project
FOREIGN KEY(project_id)
REFERENCES projects(project_id);