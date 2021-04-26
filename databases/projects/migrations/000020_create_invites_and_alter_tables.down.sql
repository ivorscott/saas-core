ALTER TABLE memberships
DROP COLUMN role;

DROP TYPE role;

ALTER TABLE memberships
RENAME COLUMN membership_id TO member_id;

ALTER TABLE memberships
RENAME TO team_member;

ALTER TABLE team_member
ADD COLUMN is_leader boolean default false,
ADD COLUMN invite_sent boolean default false,
ADD COLUMN invite_accepted boolean default false;

DROP TABLE IF EXISTS invites;

ALTER TABLE teams
RENAME TO team;

ALTER TABLE team
RENAME COLUMN user_id TO leader_id;

ALTER TABLE team
ADD COLUMN projects text ARRAY[100];

ALTER TABLE projects
DROP COLUMN team_id;

ALTER TABLE projects
DROP COLUMN public;

ALTER TABLE projects
RENAME active To open;

DROP EXTENSION IF EXISTS "uuid-ossp";