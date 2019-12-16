package models

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"time"
)

type Post struct {
	Id int32 `json:"id"`
	Message string `json:"message"`
	Date time.Time `json:"created"`
	Parent int32 `json:"parent"`
	Edited bool `json:"isEdited"`
	User string `json:"author"`
	Forum string `json:"forum"`
	Thread int32 `json:"thread"`
}

func GetPostsByThread(params utils.SearchParams, thread Thread) ([]Post, error) {
	postsFound := make([]Post, 0)
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if err != nil {
			log.Fatalln(errRollback)
		}
		return postsFound, err
	}
	var actualId int32
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&actualId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return postsFound, fmt.Errorf("can't find thread with id %d", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&actualId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return postsFound, fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
		thread.Id = actualId
	}

	sqlSelect := "SELECT DISTINCT ON (nickname COLLATE \"C\") about, fullname, nickname, email FROM forum_user " +
		"JOIN (SELECT COALESCE(p_userid, t_userid) as merge_id FROM ( " +
		"SELECT DISTINCT userid as p_userid FROM post WHERE forumid = $1) as p " +
		"FULL OUTER JOIN ( " +
		"SELECT DISTINCT userid as t_userid  FROM thread WHERE forumid = $1) " +
		"as t ON p.p_userid = t.t_userid) " +
		"as u ON u.merge_id = forum_user.id"
	var rows *pgx.Rows

	if params.Decs {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname COLLATE \"C\"  < $2 ORDER BY (nickname COLLATE \"C\") DESC LIMIT $3", forumId, params.Since, params.Limit)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname COLLATE \"C\") DESC LIMIT $2", forumId, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname COLLATE \"C\" < $2 ORDER BY (nickname COLLATE \"C\") DESC", forumId, params.Since)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname COLLATE \"C\") DESC", forumId)
			}
		}
	} else {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname COLLATE \"C\"  > $2 ORDER BY (nickname COLLATE \"C\") LIMIT $3", forumId, params.Since, params.Limit)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname COLLATE \"C\") LIMIT $2", forumId, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname COLLATE \"C\"  > $2 ORDER BY (nickname COLLATE \"C\") ", forumId, params.Since)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname COLLATE \"C\") ", forumId)
			}
		}
	}

	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return usersFound, fmt.Errorf("can't find users with forum %s", forumSlug)
	}
	for rows.Next() {
		var userFound User
		err = rows.Scan(&userFound.About, &userFound.Fullname, &userFound.Nickname, &userFound.Email)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if err != nil {
				log.Fatalln(errRollback)
			}
			return usersFound, err
		}
		usersFound = append(usersFound, userFound)
	}

	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return usersFound, nil
}