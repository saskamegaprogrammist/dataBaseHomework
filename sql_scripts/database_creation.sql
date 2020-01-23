SET search_path TO docker;
SET search_path TO public;
CREATE EXTENSION IF NOT EXISTS citext;
ALTER EXTENSION citext SET SCHEMA public;
DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS forum_user_new CASCADE;

CREATE TABLE forum_user (
    id SERIAL NOT NULL PRIMARY KEY,
    nickname citext COLLATE "C" NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    fullname varchar(100) NOT NULL,
    about text
    CONSTRAINT valid_nickname CHECK (nickname ~* '^[A-Za-z0-9_.]+$')
);

CREATE INDEX forum_user_nickname_idx ON forum_user (nickname);
CREATE INDEX forum_user_email_idx ON forum_user (email);

CREATE OR REPLACE FUNCTION check_email() RETURNS TRIGGER
LANGUAGE  plpgsql
AS $check_forum_user_email$
BEGIN
    IF (OLD.email != NEW.email) THEN
        BEGIN
            IF EXISTS(
                SELECT FROM forum_user
                WHERE email = NEW.email) THEN
                RAISE EXCEPTION 'Cannot update, email exists';
            END IF;
        END;
    END IF;
    RETURN NEW;
END
$check_forum_user_email$;


CREATE TRIGGER check_forum_user_email
    BEFORE UPDATE ON forum_user
    FOR EACH ROW
    EXECUTE PROCEDURE check_email();


CREATE TABLE forum (
    id SERIAL NOT NULL PRIMARY KEY,
    slug citext NOT NULL UNIQUE,
    title varchar(200) NOT NULL,
    posts int DEFAULT 0,
    threads int DEFAULT 0,
    usernick citext NOT NULL,
    CONSTRAINT valid_slug CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')

);

CREATE INDEX forum_slug_idx ON forum (slug);

CREATE TABLE thread (
    id SERIAL NOT NULL PRIMARY KEY,
    slug citext NOT NULL,
    created TIMESTAMPTZ,
    title varchar(100) NOT NULL,
    message text NOT NULL,
    votes int DEFAULT 0,
    forumslug citext NOT NULL,
    usernick citext NOT NULL
--     CONSTRAINT valid_slug_thread CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);
CREATE INDEX thread_slug ON thread(slug);
CREATE INDEX thread_user_forum_idx ON thread (usernick, forumslug);

DROP SEQUENCE IF EXISTS post_id;
CREATE SEQUENCE post_id START 1;

CREATE TABLE post (
    id SERIAL NOT NULL PRIMARY KEY,
    message text NOT NULL,
    created TIMESTAMPTZ,
    parent int DEFAULT 0,
    path BIGINT[] NOT NULL,
    isEdited boolean DEFAULT false,
    forumslug citext NOT NULL,
    usernick citext NOT NULL,
    threadid int NOT NULL
);

CREATE INDEX IF NOT EXISTS post_thread_path_idx ON post (threadid, path);
CREATE INDEX IF NOT EXISTS post_thread_id_idx ON post(threadid, id);
CREATE INDEX IF NOT EXISTS post_thread_id0_idx ON post (threadid, id) WHERE parent = 0;
CREATE INDEX IF NOT EXISTS post_thread_id_created_idx ON post (id, created, threadid);
CREATE INDEX IF NOT EXISTS post_thread_path1_id_idx ON post (threadid, (path[1]), id);
CREATE INDEX post_user_forum_idx ON post (usernick, forumslug);

CREATE TABLE votes (
    usernick citext NOT NULL,
    vote int DEFAULT 0,
    threadid int NOT NULL
);

CREATE INDEX votes_user_thread_idx ON votes (usernick, threadid);

CREATE TABLE forum_user_new (
    usernick citext NOT NULL ,
    forumslug citext NOT NULL
);

CREATE INDEX forum_forum_idx ON forum_user_new (forumslug);