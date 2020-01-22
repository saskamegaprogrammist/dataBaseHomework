package models

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type Thread struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int  `json:"votes"`
	Date time.Time `json:"created"`
}

func (thread *Thread) CreateThread() (Thread, error) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
	}
	rows := transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", thread.User)
	var userId int
	var threadExists Thread
	err = rows.Scan(&userId)
	if err != nil || userId == 0 {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadExists, fmt.Errorf("can't find user with nickname %s", thread.User)
	}
	rows = transaction.QueryRow("SELECT id, slug FROM forum WHERE slug = $1", thread.Forum)
	var forumId int
	err = rows.Scan(&forumId, &thread.Forum)
	if err != nil || forumId == 0 {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadExists, fmt.Errorf("can't find forum ")
	}
	rows = transaction.QueryRow("SELECT * FROM thread WHERE (usernick, title, forumslug, message) = ($1, $2, $3, $4)",thread.User, thread.Title, thread.Forum, thread.Message)
	err = rows.Scan(&threadExists.Id, &threadExists.Slug, &threadExists.Date, &threadExists.Title, &threadExists.Message, &threadExists.Votes, &threadExists.Forum, &threadExists.User)
	if err != nil {
		log.Println(err)
	}
	if threadExists.Id != 0  {
		return threadExists, fmt.Errorf("thread exists")
	}
	if thread.Slug != "" {
		rows = transaction.QueryRow("SELECT * FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&threadExists.Id, &threadExists.Slug, &threadExists.Date, &threadExists.Title, &threadExists.Message, &threadExists.Votes, &threadExists.Forum, &threadExists.User)
		if err != nil {
			log.Println(err)
		}
		if threadExists.Id != 0 {
			return threadExists, fmt.Errorf("thread exists")
		}
	}
	rows = transaction.QueryRow("INSERT INTO thread (usernick, forumslug, created, slug, message, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		thread.User, thread.Forum, thread.Date, thread.Slug, thread.Message, thread.Title )
	err = rows.Scan(&thread.Id)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadExists, err
	}
	_, err = transaction.Exec("UPDATE forum SET threads = threads + 1 WHERE forum.slug = $1 ", thread.Forum)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadExists, err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return threadExists, nil
}

func GetThreadsByForum(params utils.SearchParams, forumSlug string) ([]Thread, error) {
	threadsFound := make([]Thread, 0)
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadsFound, err
	}
	row := transaction.QueryRow("SELECT id, slug FROM forum WHERE slug = $1", forumSlug)
	var forumId int
	err = row.Scan(&forumId, &forumSlug)
	if err != nil || forumId == 0 {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadsFound, fmt.Errorf("can't find forum with slug %s", forumSlug)
	}
	var rows *pgx.Rows
	since,  _ := time.Parse(time.RFC3339Nano, params.Since)
	if params.Decs {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 AND created <= $2 ORDER BY created DESC LIMIT $3", forumSlug, since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 ORDER BY created DESC LIMIT $2", forumSlug, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 AND created <= $2 ORDER BY created DESC", forumSlug, since)
			} else {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 ORDER BY created DESC  ", forumSlug)
			}
		}
	} else {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 AND created >= $2 ORDER BY created LIMIT $3", forumSlug, since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 ORDER BY created LIMIT $2", forumSlug, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 AND created >= $2 ORDER BY created", forumSlug, since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT * " +
					"FROM thread WHERE forumslug =  $1 ORDER BY created ", forumSlug)
			}
		}
	}
	if err!=nil {
		log.Println(err)
	}
	for rows.Next() {
		var threadFound Thread
		err = rows.Scan(&threadFound.Id, &threadFound.Slug, &threadFound.Date, &threadFound.Title, &threadFound.Message, &threadFound.Votes, &threadFound.Forum, &threadFound.User)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return threadsFound, err
		}
		threadsFound = append(threadsFound, threadFound)
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return threadsFound, nil
}


func (thread *Thread) CreatePosts(newPosts []Post) ([]Post,  error, int) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return newPosts, err, 3
	}
	var actualId int
	var forumSlug string
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id, forumslug FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&actualId, &forumSlug)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return newPosts, fmt.Errorf("can't find thread with id %d", thread.Id),1
		}
	} else {
		rows := transaction.QueryRow("SELECT id, forumslug FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&actualId, &forumSlug)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return newPosts, fmt.Errorf("can't find thread with slug %s", thread.Slug),1
		}
		thread.Id = actualId
	}
	rows := transaction.QueryRow("SELECT slug FROM forum WHERE slug = $1", forumSlug)
	err = rows.Scan(&thread.Forum)
	if err != nil {
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return newPosts, fmt.Errorf("can't find forum"),1
	}
	if len(newPosts) != 0 {
		timeNow := time.Now()
		vals := []interface{}{}
		query := "INSERT INTO post (usernick, forumslug, created, parent, message, threadid, path) VALUES"
		for i:=0; i<len(newPosts); i++ {
			var userId int
			rows = transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", newPosts[i].User)
			err = rows.Scan(&userId)
			if err != nil {
				err = transaction.Rollback()
				if err != nil {
					log.Fatalln(err)
				}
				return newPosts, fmt.Errorf("can't find user"), 1
			}

			if newPosts[i].Parent != 0 {
				var existingParent int
				rows = transaction.QueryRow("SELECT id FROM post WHERE id = $1 AND threadid = $2", newPosts[i].Parent, thread.Id)
				err = rows.Scan(&existingParent)
				if err != nil || existingParent==0 {
					errRollback := transaction.Rollback()
					if errRollback != nil {
						log.Fatalln(errRollback)
					}
					return newPosts, fmt.Errorf("Cannot insert post, parent doesnot exist"), 2
				}
			}

			query += " (?, ?, ?, ?, ?, ?, (SELECT path FROM post WHERE id = ?) || (select nextval('post_id')::BIGINT)),"
			vals = append(vals, newPosts[i].User, forumSlug, timeNow, newPosts[i].Parent, newPosts[i].Message, thread.Id, newPosts[i].Parent)

		}
		query = query[:len(query)-1]
		query += " RETURNING id, isEdited, created";

		count := strings.Count(query, "?")
		for k := 0; k < count; k++ {
			query = strings.Replace(query, "?", "$"+strconv.Itoa(k+1), 1)
		}

		rowsMany, err := transaction.Query(query, vals...)
		j:=0
		for rowsMany.Next() {
			var time time.Time
			err = rowsMany.Scan(&newPosts[j].Id, &newPosts[j].Edited, &time)
			if err != nil {
				log.Println(err)
				errRollback := transaction.Rollback()
				if errRollback != nil {
					log.Fatalln(errRollback)
				}
				return newPosts, err, 2
			}
			newPosts[j].Date = time
			newPosts[j].Thread = thread.Id
			newPosts[j].Forum = thread.Forum
			j++
		}
		_, err = transaction.Exec("UPDATE forum SET posts = posts + $1 WHERE slug = $2 ",  len(newPosts), forumSlug)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return newPosts, err, 2
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return newPosts, nil, 0
}

func (thread *Thread) GetThread() error {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT * FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.Forum, &thread.User)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with id %d", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT * FROM thread WHERE slug = $1",  thread.Slug)
		err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.Forum, &thread.User)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	return nil
}

func (thread *Thread) UpdateThread() error  {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	var actualId int
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&actualId)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with id %s", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&actualId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
		thread.Id = actualId
	}
	_, err = transaction.Exec("UPDATE thread SET message = coalesce(nullif($2, ''), message), title = coalesce(nullif($3, ''), title) WHERE id = $1",  thread.Id, thread.Message, thread.Title)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}
	rows := transaction.QueryRow("SELECT * FROM thread WHERE id = $1",  thread.Id)
	err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.Forum, &thread.User)
	if err != nil {
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	return nil
}


func (thread *Thread) Vote(vote *Vote) error  {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id, slug, created, title, message, votes, usernick, forumslug FROM thread WHERE id = $1",  thread.Id)
		err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.User, &thread.Forum)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with id %s", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id, slug, created, title, message, votes, usernick, forumslug FROM thread WHERE slug = $1",  thread.Slug)
		err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.User, &thread.Forum)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
	}
	var userId int
	rows := transaction.QueryRow("SELECT id  FROM forum_user WHERE nickname = $1", vote.Nickname)
	err = rows.Scan(&userId)
	if err != nil || userId == 0  {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return fmt.Errorf("can't find user with nickname %s", vote.Nickname)
	}
	var voteExists int
	rows = transaction.QueryRow("SELECT vote FROM votes WHERE usernick = $1 AND threadid = $2", vote.Nickname, thread.Id)
	err = rows.Scan(&voteExists)
	if voteExists != 0 {
		if voteExists!=vote.Voice {
			_, err = transaction.Exec("UPDATE thread SET votes = $2 WHERE id = $1",  thread.Id, thread.Votes + vote.Voice - voteExists)
			if err != nil {
				log.Println(err)
				err = transaction.Rollback()
				if err != nil {
					log.Fatalln(err)
				}
				return err
			}
			_, err = transaction.Exec("UPDATE votes SET vote = $3 WHERE usernick = $1 AND threadid = $2", vote.Nickname, thread.Id, vote.Voice)
			if err != nil {
				log.Println(err)
				err = transaction.Rollback()
				if err != nil {
					log.Fatalln(err)
				}
				return err
			}
			thread.Votes = thread.Votes + vote.Voice - voteExists
		}
	} else {
		_, err = transaction.Exec("INSERT INTO votes (usernick, vote, threadid) VALUES ($1, $2, $3)",
			vote.Nickname, vote.Voice, thread.Id)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Fatalln(errRollback)
			}
			return err
		}
		_, err = transaction.Exec("UPDATE thread SET votes = $2 WHERE id = $1",  thread.Id, thread.Votes + vote.Voice)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err
		}
		thread.Votes = thread.Votes + vote.Voice

	}
	err = transaction.Commit()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	return nil
}