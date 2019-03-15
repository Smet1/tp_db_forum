CREATE EXTENSION IF NOT EXISTS CITEXT;

TRUNCATE TABLE forum_users;

CREATE TABLE IF NOT EXISTS forum_users (
  id SERIAL PRIMARY KEY,
  nickname CITEXT NOT NULL UNIQUE,
  fullname TEXT,
  email TEXT NOT NULL UNIQUE,
  about TEXT
)