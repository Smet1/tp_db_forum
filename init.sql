CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users CASCADE;
DROP TABLE IF EXISTS forum_forum CASCADE;
DROP TABLE IF EXISTS forum_thread CASCADE;
DROP TABLE IF EXISTS forum_post CASCADE;
DROP TABLE IF EXISTS forum_vote CASCADE;
DROP TABLE IF EXISTS forum_users_forum CASCADE;

DROP INDEX IF EXISTS forum_users_nickname;

DROP INDEX IF EXISTS forum_thread_slug;
DROP INDEX IF EXISTS forum_thread_forum_created;
DROP INDEX IF EXISTS forum_thread_author_forum;

DROP INDEX IF EXISTS forum_post_path_id;
DROP INDEX IF EXISTS forum_post_thread_id;
DROP INDEX IF EXISTS forum_post_thread;
DROP INDEX IF EXISTS forum_post_thread_path_id;
DROP INDEX IF EXISTS forum_post_thread_id_path_parent;
DROP INDEX IF EXISTS forum_post_author_forum;

DROP INDEX IF EXISTS forum_vote_nickname_thread;

CREATE TABLE IF NOT EXISTS forum_users
(
--     id       SERIAL PRIMARY KEY,
    nickname CITEXT NOT NULL UNIQUE,
    fullname TEXT,
    email    CITEXT NOT NULL UNIQUE,
    about    TEXT
);

-- GetPostDetails (хз)
-- CREATE INDEX IF not exists forum_users_nickname_email ON forum_users (nickname, email);
CREATE INDEX IF not exists forum_users_nickname ON forum_users (nickname);

CREATE TABLE IF NOT EXISTS forum_forum
(
    posts   INTEGER DEFAULT 0,
    slug    citext                                   NOT NULL UNIQUE,
    threads INTEGER DEFAULT 0,
    title   TEXT                                     NOT NULL,
    "user"  CITEXT REFERENCES forum_users (nickname) NOT NULL
);

CREATE INDEX IF not exists forum_forum_slug ON forum_forum (slug);


CREATE TABLE IF NOT EXISTS forum_thread
(
    author  CITEXT REFERENCES forum_users (nickname) NOT NULL,
    created timestamptz,
    forum   CITEXT REFERENCES forum_forum (slug)     NOT NULL,
    id      SERIAL PRIMARY KEY,
    message TEXT                                     NOT NULL,
    slug    CITEXT                                            DEFAULT NULL UNIQUE,
    title   TEXT,
    votes   INTEGER                                  NOT NULL DEFAULT 0
);

-- CREATE INDEX IF not exists forum_thread_author_title_slug ON forum_thread (author, title, slug);
-- CREATE INDEX IF NOT EXISTS forum_thread_id_slug ON forum_thread (id, slug);
-- CREATE INDEX IF NOT EXISTS forum_thread_forum ON forum_thread (forum);

-- GetForumBySlug
CREATE INDEX IF not exists forum_thread_slug ON forum_thread (slug);
-- GetForumThreads
-- CREATE INDEX IF NOT EXISTS forum_thread_created ON forum_thread (created);
CREATE INDEX IF NOT EXISTS forum_thread_forum_created ON forum_thread (forum, created);
-- GetForumUsers (но хз)
CREATE INDEX IF not exists forum_thread_author_forum ON forum_thread (author, forum);



CREATE TABLE IF NOT EXISTS forum_post
(
    author   citext REFERENCES forum_users (nickname) NOT NULL,
    created  TIMESTAMPTZ DEFAULT current_timestamp, -- сделать current_timestamp (посоветовал ник)
    forum    citext REFERENCES forum_forum (slug),
    id       SERIAL PRIMARY KEY,
    isEdited BOOLEAN     DEFAULT FALSE,
    message  TEXT                                     NOT NULL,
    parent   INTEGER     DEFAULT 0,
    thread   INTEGER REFERENCES forum_thread (id)     NOT NULL,
    path     INTEGER[]   DEFAULT array []::INT[]
);

CREATE INDEX IF NOT EXISTS forum_post_path_id ON forum_post (id, (path [1]));
CREATE INDEX IF NOT EXISTS forum_post_path ON forum_post (path);
CREATE INDEX IF NOT EXISTS forum_post_path_1 ON forum_post ((path [1]));
CREATE INDEX IF not exists forum_post_thread_id ON forum_post (thread, id);
CREATE INDEX IF not exists forum_post_thread ON forum_post (thread);
CREATE INDEX IF not exists forum_post_thread_path_id ON forum_post (thread, path, id);
CREATE INDEX IF NOT EXISTS forum_post_thread_id_path_parent ON forum_post (thread, id, (path[1]), parent);
CREATE INDEX IF not exists forum_post_author_forum ON forum_post (author, forum); -- если нигде не используется удалить нахуй


CREATE TABLE IF NOT EXISTS forum_vote
(
    nickname CITEXT REFERENCES forum_users (nickname) NOT NULL,
    voice    SMALLINT CHECK ( voice IN (-1, 1) ),
    thread   INTEGER REFERENCES forum_thread (id)     NOT NULL,
    UNIQUE (nickname, thread)
);

-- GetVoteByNicknameAndThreadID
CREATE INDEX IF not exists forum_vote_nickname_thread ON forum_vote (nickname, thread);

CREATE TABLE IF NOT EXISTS forum_users_forum
(
    nickname citext REFERENCES forum_users (nickname) NOT NULL,
    slug     citext REFERENCES forum_forum (slug)     NOT NULL,
    UNIQUE (nickname, slug)
);