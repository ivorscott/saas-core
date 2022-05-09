INSERT INTO users (user_id, auth0_id, email, email_verified, first_name, last_name, picture, locale, created_at, updated_at) VALUES
('a4b54ec1-57f9-4c39-ab53-d936dbb6c177','auth0|60a666916089a00069b2a773','testuser@devpie.io',false,'testuser','','https://s.gravatar.com/avatar/8e199cf73a009f2065e5b4dd7d5353e3?s=480\u0026r=pg\u0026d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fte.png','','2021-05-21T10:42:36.815181Z','2021-05-21T10:42:36.815181Z')
ON CONFLICT DO NOTHING;

