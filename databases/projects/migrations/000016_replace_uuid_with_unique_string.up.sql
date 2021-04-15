drop table if exists tasks;
drop table if exists columns;
drop table if exists projects;
drop table if exists users;
drop table if exists ma_token;
drop table if exists products;

create table projects (
  project_id uuid primary key,
  user_id varchar(128) not null,
  open boolean not null,
  column_order text array[4],
  name varchar(36) not null,
  created timestamp without time zone default (now() at time zone 'utc')
);

create table columns (
 column_id uuid primary key,
 project_id uuid,
 title varchar(36) not null,
 column_name varchar(8) not null,
 task_ids text[],
 created timestamp without time zone default (now() at time zone 'utc'),
 constraint fk_project
     foreign key (project_id)
         references projects(project_id)
);

create table tasks (
   task_id uuid primary key,
   title varchar(48) not null,
   project_id uuid not null,
   content text,
   created timestamp without time zone default (now() at time zone 'utc')
);