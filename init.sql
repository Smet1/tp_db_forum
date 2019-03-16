CREATE EXTENSION IF NOT EXISTS CITEXT;

-- TRUNCATE TABLE forum_users;
-- TRUNCATE TABLE forum_forum;

CREATE TABLE IF NOT EXISTS forum_users (
  id SERIAL PRIMARY KEY,
  nickname CITEXT NOT NULL UNIQUE,
  fullname TEXT,
  email TEXT NOT NULL UNIQUE,
  about TEXT
);

CREATE TABLE IF NOT EXISTS forum_forum (
  posts INTEGER DEFAULT 0,
  slug citext NOT NULL UNIQUE,
  threads INTEGER DEFAULT 0,
  title TEXT NOT NULL ,
  "user" CITEXT REFERENCES forum_users(nickname) NOT NULL
);