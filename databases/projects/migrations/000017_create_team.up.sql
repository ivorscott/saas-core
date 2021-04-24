CREATE TABLE IF NOT EXISTS team (
  team_id UUID not null unique,
  leader_id varchar(128) not null unique,
  name varchar(68) not null unique,
  projects text ARRAY[100],
  created timestamp without time zone default (now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS team_member (
  member_id uuid not null unique,
  user_id varchar(128) not null unique,
  team_id uuid not null unique,
  is_leader boolean default false,
  invite_sent boolean default false,
  invite_accepted boolean default false,
  created timestamp without time zone default (now() at time zone 'utc')
);