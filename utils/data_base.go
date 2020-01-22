package utils

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
)

var dataBasePool *pgx.ConnPool

func CreateAddress(user, password, host, name string) string {
	return  fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s",
	user, password, host, name)
}

func CreateDataBaseConnection(user, password, host, name string, maxConn int) {
	dataBaseConfig := CreateAddress(user, password, host, name)
	connectionConfig, err := pgx.ParseConnectionString(dataBaseConfig)
	if err != nil {
		log.Println(err);
		return
	}
	dataBasePool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: connectionConfig,
		MaxConnections: maxConn,
	})
	if err != nil {
		log.Println(err);
		return
	}
}

func InitDataBase() {
	_, err := dataBasePool.Exec(`
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
    created TIMESTAMPTZ(9),
    title varchar(100) NOT NULL,
    message text NOT NULL,
    votes int DEFAULT 0,
    forumslug citext NOT NULL,
    usernick citext NOT NULL
--     CONSTRAINT valid_slug_thread CHECK (slug ~* '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);
CREATE INDEX thread_slug ON thread(slug);
CREATE INDEX thread_createdtitle_idx ON thread (created, title);
CREATE INDEX thread_forumslug_idx ON thread (forumslug);

DROP SEQUENCE IF EXISTS post_id;
CREATE SEQUENCE post_id START 1;

CREATE TABLE post (
    id SERIAL NOT NULL PRIMARY KEY,
    message text NOT NULL,
    created TIMESTAMPTZ(9), 
    parent int DEFAULT 0,
    path BIGINT[] NOT NULL,
    isEdited boolean DEFAULT false,
    forumslug citext NOT NULL,
    usernick citext NOT NULL,
    threadid int NOT NULL
);
CREATE INDEX post_thread_idx ON post (id, threadid);
CREATE INDEX post_path ON post (path);
CREATE INDEX post_created ON post (created);
CREATE INDEX post_level ON post (array_length(path, 1));
CREATE INDEX post_forumslug_idx ON post (forumslug);

CREATE TABLE votes (
    usernick citext NOT NULL,
    vote int DEFAULT 0,
    threadid int NOT NULL
);

CREATE INDEX votes_user_thread_idx ON votes (usernick, threadid);
`)
	if err != nil {
		log.Println(err)
	}

}

func GetDataBase() *pgx.ConnPool {
	return dataBasePool
}