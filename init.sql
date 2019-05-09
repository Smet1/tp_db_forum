CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users CASCADE;
DROP TABLE IF EXISTS forum_forum CASCADE;
DROP TABLE IF EXISTS forum_thread CASCADE;
DROP TABLE IF EXISTS forum_post CASCADE;
DROP TABLE IF EXISTS forum_vote CASCADE;

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
    created  timestamptz,
    forum    citext REFERENCES forum_forum (slug),
    id       SERIAL PRIMARY KEY,
    isEdited BOOLEAN   DEFAULT FALSE,
    message  TEXT                                     NOT NULL,
    parent   INTEGER   DEFAULT 0,
    thread   INTEGER REFERENCES forum_thread (id)     NOT NULL,
    path     INTEGER[] DEFAULT array []::INT[]
);


-- CREATE INDEX IF NOT EXISTS posts_id_index ON forum_post (id);
-- CREATE INDEX IF NOT EXISTS posts_forum_index ON forum_post (forum);
-- CREATE INDEX IF NOT EXISTS posts_main_index ON forum_post (thread, parent);
-- CREATE INDEX IF NOT EXISTS posts_thread_index ON forum_post (thread);

-- CREATE INDEX IF NOT EXISTS forum_post_path_id ON forum_post (path, id);
-- -- GetSortedPosts
-- CREATE INDEX IF not exists forum_post_path ON forum_post (path);
-- -- FlatSort
CREATE INDEX IF not exists forum_post_thread_id ON forum_post (thread, id);
-- -- TreeSort, ParentTreeSort
CREATE INDEX IF not exists forum_post_thread_path_id ON forum_post (thread, path, id);
-- GetForumUsers (хз)
CREATE INDEX IF not exists forum_post_author_forum ON forum_post (author, forum);


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
    slug  citext REFERENCES forum_forum (slug)  NOT NULL,
    UNIQUE (nickname, slug)
);