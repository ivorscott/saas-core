DROP TYPE role;

ALTER TABLE memberships
DROP COLUMN role;

ALTER TABLE memberships
RENAME COLUMN membership_id TO member_id;

ALTER TABLE memberships
ADD COLUMN is_leader boolean default false,
ADD COLUMN invite_sent boolean default false,
ADD COLUMN invite_accepted boolean default false;

ALTER TABLE memberships
RENAME TO team_member;

ALTER TABLE teams
RENAME TO team;

ALTER TABLE projects
DROP COLUMN team_id;

ALTER TABLE projects
DROP COLUMN public;

ALTER TABLE projects
RENAME active To open;

DROP TABLE IF EXISTS invites;
