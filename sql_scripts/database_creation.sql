CREATE EXTENSION IF NOT EXISTS citext;
DROP TABLE IF EXISTS forum_user CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS thread CASCADE;

CREATE TABLE forum_user (
    id SERIAL NOT NULL PRIMARY KEY,
    nickname citext NOT NULL UNIQUE,
    email varchar(100) NOT NULL UNIQUE,
    fullname varchar(100) NOT NULL,
    about varchar(200)
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
    slug varchar(200) NOT NULL UNIQUE,
    title varchar(200) NOT NULL,
    posts int DEFAULT 0,
    threads int DEFAULT 0,
    userid int NOT NULL,
    FOREIGN KEY (userid) REFERENCES forum_user (id),
    CONSTRAINT valid_slug CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')

);

CREATE UNIQUE INDEX forum_slug_idx ON forum (slug);

CREATE TABLE thread (
    id SERIAL NOT NULL PRIMARY KEY,
    slug varchar(200) NOT NULL UNIQUE,
    created timestamp,
    title varchar(100) NOT NULL,
    message text NOT NULL,
    votes int DEFAULT 0,
    forumid int NOT NULL,
    userid int NOT NULL,
    FOREIGN KEY (userid) REFERENCES forum_user (id),
    FOREIGN KEY (forumid) REFERENCES forum (id),
    CONSTRAINT valid_slug_thread CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);

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

CREATE TABLE post (
    id SERIAL NOT NULL PRIMARY KEY,
    message text NOT NULL,
    created timestamp,
    parent int DEFAULT 0,
    isEdited boolean DEFAULT false,
    forumid int NOT NULL,
    userid int NOT NULL,
    threadid int NOT NULL,
    FOREIGN KEY (userid) REFERENCES forum_user (id),
    FOREIGN KEY (forumid) REFERENCES forum (id),
    FOREIGN KEY (threadid) REFERENCES thread (id)
);