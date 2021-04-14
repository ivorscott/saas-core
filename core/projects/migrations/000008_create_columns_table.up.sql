CREATE TABLE columns (
    column_id UUID PRIMARY KEY,
    project_id UUID,
    title varchar(36) not null,
    columns varchar(8) not null,
    task_ids text[],
    CONSTRAINT fk_project
        FOREIGN KEY(project_id)
            REFERENCES projects(project_id)
);