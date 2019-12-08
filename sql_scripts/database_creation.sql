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


CREATE TABLE forum (
    id SERIAL NOT NULL PRIMARY KEY,
    slug varchar(200) NOT NULL UNIQUE,
    posts int DEFAULT 0,
    threads int DEFAULT 0,
    userid int NOT NULL,
    FOREIGN KEY (userid) REFERENCES forum_user (id),
    CONSTRAINT valid_slug CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')

);

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