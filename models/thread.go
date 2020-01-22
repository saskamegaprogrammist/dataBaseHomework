package models

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
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
	var threadExistsForumId int
	rows = transaction.QueryRow("SELECT * FROM thread WHERE (usernick, title, forumid, message) = ($1, $2, $3, $4)",thread.User, thread.Title, forumId, thread.Message)
	err = rows.Scan(&threadExists.Id, &threadExists.Slug, &threadExists.Date, &threadExists.Title, &threadExists.Message, &threadExists.Votes, &threadExistsForumId, &threadExists.User)
	if err != nil {
		log.Println(err)
	}
	if threadExists.Id != 0  {
		rows = transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", threadExistsForumId)
		_ = rows.Scan(&threadExists.Forum)
		return threadExists, fmt.Errorf("thread exists")
	}
	if thread.Slug != "" {
		rows = transaction.QueryRow("SELECT * FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&threadExists.Id, &threadExists.Slug, &threadExists.Date, &threadExists.Title, &threadExists.Message, &threadExists.Votes, &threadExistsForumId, &threadExists.User)
		if err != nil {
			log.Println(err)
		}
		if threadExists.Id != 0 {
			rows = transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", threadExistsForumId)
			_ = rows.Scan(&threadExists.Forum)
			return threadExists, fmt.Errorf("thread exists")
		}
	}
	rows = transaction.QueryRow("INSERT INTO thread (usernick, forumid, created, slug, message, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		thread.User, forumId, thread.Date, thread.Slug, thread.Message, thread.Title )
	err = rows.Scan(&thread.Id)
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
	row := transaction.QueryRow("SELECT id FROM forum WHERE slug = $1", forumSlug)
	var forumId int
	err = row.Scan(&forumId)
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
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND created <= $2 ORDER BY created DESC LIMIT $3", forumSlug, since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created DESC LIMIT $2", forumSlug, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND created <= $2 ORDER BY created DESC", forumSlug, since)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created DESC  ", forumSlug)
			}
		}
	} else {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND created >= $2 ORDER BY created LIMIT $3", forumSlug, since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created LIMIT $2", forumSlug, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND created >= $2 ORDER BY created", forumSlug, since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created, title, message, votes, usernick, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created ", forumSlug)
			}
		}
	}
	if err!=nil {
		log.Println(err)
	}
	for rows.Next() {
		var threadFound Thread
		err = rows.Scan(&threadFound.Id, &threadFound.Slug, &threadFound.Date, &threadFound.Title, &threadFound.Message, &threadFound.Votes, &threadFound.User, &threadFound.Forum)
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
	var forumId int
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id, forumid FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&actualId, &forumId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return newPosts, fmt.Errorf("can't find thread with id %d", thread.Id),1
		}
	} else {
		rows := transaction.QueryRow("SELECT id, forumid FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&actualId, &forumId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return newPosts, fmt.Errorf("can't find thread with slug %s", thread.Slug),1
		}
		thread.Id = actualId
	}
	rows := transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", forumId)
	err = rows.Scan(&thread.Forum)
	if err != nil {
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return newPosts, fmt.Errorf("can't find forum"),1
	}
	timeNow := time.Now()
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


		rows = transaction.QueryRow("INSERT INTO post (usernick, forumid, created, parent, message, threadid, path) VALUES ($1, $2, $3, $4, $5, $6, (SELECT path FROM post WHERE id = $4) || (select nextval('post_id')::BIGINT)) RETURNING id, isEdited",
			newPosts[i].User, forumId, timeNow, newPosts[i].Parent, newPosts[i].Message, thread.Id)
		err = rows.Scan(&newPosts[i].Id, &newPosts[i].Edited)
		fmt.Println(newPosts[i].Id, thread.Id, forumId, len(newPosts))
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Fatalln(errRollback)
			}
			return newPosts, err, 2
		}
		newPosts[i].Date = timeNow
		newPosts[i].Thread = thread.Id
		newPosts[i].Forum = thread.Forum
	}
	if len(newPosts) != 0 {
		_, err = transaction.Exec("UPDATE forum SET posts = posts + $1 WHERE forum.id = $2 ",  len(newPosts), forumId)
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
		rows := transaction.QueryRow("SELECT id, slug, created, title, message, votes, usernick, forum FROM thread_full_view WHERE id = $1", thread.Id)
		err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.User, &thread.Forum)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with id %d", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id, slug, created, title, message, votes, usernick, forum FROM thread_full_view WHERE slug = $1",  thread.Slug)
		err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.User, &thread.Forum)
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
	if thread.Message != "" {
		_, err = transaction.Exec("UPDATE thread SET message = $2 WHERE id = $1",  thread.Id, thread.Message)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err
		}
	}
	if thread.Title != "" {
		_, err = transaction.Exec("UPDATE thread SET title = $2 WHERE id = $1",  thread.Id, thread.Title)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err
		}
	}
	rows := transaction.QueryRow("SELECT id, slug, created, title, message, votes, usernick, forum FROM thread_full_view WHERE id = $1",  thread.Id)
	err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.User, &thread.Forum)
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
	var actualId int
	var votes int
	if thread.Id != 0 {
		rows := transaction.QueryRow("SELECT id, votes FROM thread WHERE id = $1", thread.Id)
		err = rows.Scan(&actualId, &votes)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with id %s", thread.Id)
		}
	} else {
		rows := transaction.QueryRow("SELECT id, votes FROM thread WHERE slug = $1", thread.Slug)
		err = rows.Scan(&actualId, &votes)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("can't find thread with slug %s", thread.Slug)
		}
		thread.Id = actualId
	}
	var userId int
	rows := transaction.QueryRow("SELECT id  FROM forum_user WHERE nickname = $1", vote.Nickname)
	err = rows.Scan(&userId)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	var exists bool
	rows = transaction.QueryRow("SELECT EXISTS (SELECT * FROM votes WHERE userid = $1 AND threadid = $2)", userId, thread.Id)
	err = rows.Scan(&exists)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	if exists {
		var voteExists int
		rows = transaction.QueryRow("SELECT vote FROM votes WHERE userid = $1 AND threadid = $2", userId, thread.Id)
		err = rows.Scan(&voteExists)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Fatalln(errRollback)
			}
			return err
		}
		fmt.Println(votes)
		fmt.Println(vote.Voice)
		fmt.Println(voteExists)
		if voteExists!=vote.Voice {
			_, err = transaction.Exec("UPDATE thread SET votes = $2 WHERE id = $1",  thread.Id, votes + vote.Voice - voteExists)
			if err != nil {
				log.Println(err)
				err = transaction.Rollback()
				if err != nil {
					log.Fatalln(err)
				}
				return err
			}
			_, err = transaction.Exec("UPDATE votes SET vote = $3 WHERE userid = $1 AND threadid = $2", userId, thread.Id, vote.Voice)
			if err != nil {
				log.Println(err)
				err = transaction.Rollback()
				if err != nil {
					log.Fatalln(err)
				}
				return err
			}
		}
	} else {
		var id int
		rows = transaction.QueryRow("INSERT INTO votes (userid, vote, threadid) VALUES ($1, $2, $3) returning id",
			userId, vote.Voice, thread.Id)
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Fatalln(errRollback)
			}
			return err
		}
		_, err = transaction.Exec("UPDATE thread SET votes = $2 WHERE id = $1",  thread.Id, votes + vote.Voice)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err
		}
	}

	rows = transaction.QueryRow("SELECT id, slug, created, title, message, votes, usernick, forum FROM thread_full_view WHERE id = $1",  thread.Id)
	err = rows.Scan(&thread.Id, &thread.Slug, &thread.Date, &thread.Title, &thread.Message, &thread.Votes,  &thread.User, &thread.Forum)
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