ALTER TABLE projects
ADD COLUMN created timestamp without time zone default (now() at time zone 'utc');

ALTER TABLE columns
ADD COLUMN created timestamp without time zone default (now() at time zone 'utc');

ALTER TABLE tasks
ADD COLUMN created timestamp without time zone default (now() at time zone 'utc');