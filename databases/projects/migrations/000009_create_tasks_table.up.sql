CREATE TABLE tasks (
    task_id UUID PRIMARY KEY,
    title varchar(48) not null,
    content text
);