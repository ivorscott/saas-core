CREATE TABLE projects (
    project_id UUID PRIMARY KEY,
    user_id UUID,
    open boolean not null,
    column_order text ARRAY[4],
    name varchar(36) not null,
    CONSTRAINT fk_creator
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);
