CREATE EXTENSION IF NOT EXISTS CITEXT;

TRUNCATE TABLE forum_users CASCADE;
TRUNCATE TABLE forum_forum CASCADE;
TRUNCATE TABLE forum_thread CASCADE;

DROP TABLE IF EXISTS forum_users CASCADE;
DROP TABLE IF EXISTS forum_forum CASCADE;
DROP TABLE IF EXISTS forum_thread CASCADE;
-- DROP TABLE IF EXISTS post;
-- DROP TABLE IF EXISTS vote;
-- DROP TABLE IF EXISTS thread;
-- DROP TABLE IF EXISTS forum_users;

CREATE TABLE IF NOT EXISTS forum_users
(
  id       SERIAL PRIMARY KEY,
  nickname CITEXT NOT NULL UNIQUE,
  fullname TEXT,
  email    CITEXT NOT NULL UNIQUE,
  about    TEXT
);

CREATE TABLE IF NOT EXISTS forum_forum
(
  posts   INTEGER DEFAULT 0,
  slug    citext                                   NOT NULL UNIQUE,
  threads INTEGER DEFAULT 0,
  title   TEXT                                     NOT NULL,
  "user"  CITEXT REFERENCES forum_users (nickname) NOT NULL
);

CREATE TABLE IF NOT EXISTS forum_thread
(
  author  CITEXT REFERENCES forum_users (nickname) NOT NULL,
  created TIMESTAMP,
  forum   CITEXT REFERENCES forum_forum (slug)     NOT NULL,
  id      SERIAL PRIMARY KEY,
  message TEXT                                     NOT NULL,
  slug    CITEXT                                            DEFAULT NULL UNIQUE,
  title   TEXT,
  votes   INTEGER                                  NOT NULL DEFAULT 0
)