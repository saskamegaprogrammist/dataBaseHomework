SET search_path TO docker;
SET search_path TO public;
CREATE EXTENSION IF NOT EXISTS citext;
ALTER EXTENSION citext SET SCHEMA public;
DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE TABLE forum_user (
    id SERIAL NOT NULL PRIMARY KEY,
    nickname citext NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    fullname varchar(100) NOT NULL,
    about text
    CONSTRAINT valid_nickname CHECK (nickname ~* '^[A-Za-z0-9_.]+$')
);

CREATE UNIQUE INDEX forum_user_nickname_idx ON forum_user (nickname);
CREATE UNIQUE INDEX forum_user_email_idx ON forum_user (email);

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
    FOREIGN KEY (usernick) REFERENCES forum_user (nickname),
    CONSTRAINT valid_slug CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')

);

CREATE INDEX forum_slug_idx ON forum (slug);

CREATE TABLE thread (
    id SERIAL NOT NULL PRIMARY KEY,
    slug citext NOT NULL,
    created TIMESTAMP WITH TIME ZONE,
    title varchar(100) NOT NULL,
    message text NOT NULL,
    votes int DEFAULT 0,
    forumid int NOT NULL,
    usernick citext NOT NULL,
    FOREIGN KEY (usernick) REFERENCES forum_user (nickname),
    FOREIGN KEY (forumid) REFERENCES forum (id)
--     CONSTRAINT valid_slug_thread CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);
CREATE INDEX thread_slug ON thread(slug);
CREATE UNIQUE INDEX thread_createdtitle_idx ON thread (created, title);

CREATE OR REPLACE FUNCTION add_thread() RETURNS TRIGGER
LANGUAGE  plpgsql
AS $add_forum_thread$
BEGIN
    UPDATE forum
        SET threads = threads + 1
        WHERE forum.id = NEW.forumid;
    RETURN NEW;
END
$add_forum_thread$;


CREATE TRIGGER add_forum_thread
    AFTER INSERT ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE add_thread();

CREATE INDEX thread_forumid_idx ON thread (forumid);

DROP VIEW IF EXISTS thread_full_view ;

CREATE VIEW thread_full_view  AS
    SELECT id, slug, created, title, message, votes, usernick, forum FROM thread
    JOIN (SELECT slug as forum, id as fid FROM forum) as forum_t ON forum_t.fid = thread.forumid;

DROP SEQUENCE IF EXISTS post_id;
CREATE SEQUENCE post_id START 1;

CREATE TABLE post (
    id SERIAL NOT NULL PRIMARY KEY,
    message text NOT NULL,
    created TIMESTAMP WITH TIME ZONE,
    parent int DEFAULT 0,
    path BIGINT[] NOT NULL,
    isEdited boolean DEFAULT false,
    forumid int NOT NULL,
    usernick citext NOT NULL,
    threadid int NOT NULL,
    FOREIGN KEY (usernick) REFERENCES forum_user (nickname),
    FOREIGN KEY (forumid) REFERENCES forum (id),
    FOREIGN KEY (threadid) REFERENCES thread (id)
);

CREATE INDEX post_path_desc ON post (path DESC, id);
CREATE INDEX post_path ON post (path);
CREATE INDEX post_level ON post (array_length(path, 1));
CREATE INDEX post_forumid_idx ON post (forumid);

CREATE TABLE votes (
    id SERIAL NOT NULL PRIMARY KEY,
    userid int NOT NULL ,
    vote int DEFAULT 0,
    threadid int NOT NULL ,
    FOREIGN KEY (userid) REFERENCES forum_user (id),
    FOREIGN KEY (threadid) REFERENCES thread (id)
);

CREATE INDEX votes_thread_idx ON votes (userid);
CREATE INDEX votes_user_idx ON votes (threadid);