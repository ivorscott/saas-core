CREATE TABLE IF NOT EXISTS invite(
  invite_id uuid not null unique,
  user_id varchar(128) not null unique,
  team_id uuid not null unique,
  read boolean default false,
  accepted boolean default false,
  expiration timestamp without time zone default ((now() + '5 days') at time zone 'utc'),
  created timestamp without time zone default (now() at time zone 'utc')
);

ALTER TABLE projects
DROP COLUMN open;

ALTER TABLE projects
ADD COLUMN active boolean default true;

ALTER TABLE projects
ADD COLUMN public boolean default false;

--- add uuid extension v9.1 and newer ( uuid_generate_v4() )
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--- add column
ALTER TABLE projects
ADD COLUMN team_id uuid;

--- update existing records where data is null
UPDATE projects SET team_id = uuid_generate_v4() WHERE team_id IS NULL;

--- then add the not null constraint
ALTER TABLE projects
ALTER COLUMN team_id SET NOT NULL;

ALTER TABLE team
RENAME TO teams;

ALTER TABLE team_member
RENAME TO memberships;

ALTER TABLE memberships
DROP COLUMN is_leader, DROP COLUMN invite_sent, DROP COLUMN invite_accepted;

CREATE TYPE role AS ENUM ('admin', 'editor', 'commenter','viewer');

ALTER TABLE memberships
ADD COLUMN role role default 'editor';

ALTER TABLE memberships
RENAME COLUMN member_id TO membership_id;



