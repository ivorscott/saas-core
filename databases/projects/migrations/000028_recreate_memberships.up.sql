DROP TABLE IF EXISTS memberships;
DROP TYPE IF EXISTS ROLE CASCADE;
CREATE TYPE ROLE AS ENUM ('administrator', 'editor', 'commenter','viewer');

CREATE TABLE IF NOT EXISTS memberships (
   membership_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   user_id UUID NOT NULL,
   team_id UUID NOT NULL,
   role ROLE DEFAULT 'editor',
   created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
   updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);
