package models

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
	time2 "time"
)

type Thread struct {
	Id int32 `json:"id"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int32  `json:"votes"`
	Date string `json:"created"`
}

func (thread *Thread) CreateThread() (Thread, error) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
	}
	rows := transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", thread.User)
	var userId int32
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
	rows = transaction.QueryRow("SELECT id FROM forum WHERE slug = $1", thread.Forum)
	var forumId int32
	err = rows.Scan(&forumId)
	if err != nil || forumId == 0 {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadExists, fmt.Errorf("can't find forum with slug %s", thread.Forum)
	}
	var threadExistsUserId int32
	var threadExistsForumId int32
	rows = transaction.QueryRow("SELECT * FROM thread WHERE (created, title) = ($1, $2)", thread.Date, thread.Title)
	err = rows.Scan(&threadExists.Id, &threadExists.Slug, &threadExists.Date, &threadExists.Title, &threadExists.Message, &threadExists.Votes, &threadExistsForumId, &threadExistsUserId)
	if err != nil {
		log.Println(err)
	}
	if threadExists.Id != 0  {
		rows = transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", threadExistsUserId)
		_ = rows.Scan(&threadExists.User)
		rows = transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", threadExistsForumId)
		_ = rows.Scan(&threadExists.Forum)
		return threadExists, fmt.Errorf("thread exists")
	}
	rows = transaction.QueryRow("INSERT INTO thread (userid, forumid, created, slug, message, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		userId, forumId, thread.Date, thread.Slug, thread.Message, thread.Title )
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
	var threadsFound []Thread
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
	var rows *pgx.Rows
	if params.Decs {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND thread_full_view.created >= $2 ORDER BY created DESC LIMIT $3", forumSlug, params.Since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created DESC LIMIT $2", forumSlug, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND thread_full_view.created >= $2 ORDER BY created DESC ", forumSlug, params.Since)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created DESC  ", forumSlug)
			}
		}
	} else {
		if params.Limit != -1 {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND thread_full_view.created >= $2 ORDER BY created LIMIT $3", forumSlug, params.Since, params.Limit)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created LIMIT $2", forumSlug, params.Limit)
			}
		} else {
			if params.Since != "" {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 AND thread_full_view.created >= $2 ORDER BY created ", forumSlug, params.Since)
			} else {
				rows, err = transaction.Query("SELECT id, slug, created::text, title, message, votes, user_forum, forum " +
					"FROM thread_full_view WHERE forum =  $1 ORDER BY created ", forumSlug)
			}
		}
	}


	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadsFound, fmt.Errorf("can't find threads with forum %s", forumSlug)
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
	if len(threadsFound) == 0 {
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return threadsFound, fmt.Errorf("can't find threads with forum %s", forumSlug)
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return threadsFound, nil
}


func (thread *Thread) CreatePosts(newPosts []Post) ([]Post,  error, int32) {
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
	var actualId int32
	var forumId int32
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
	time := time2.Now().Format(time2.RFC3339)
	for i:=0; i<len(newPosts); i++ {
		var userId int32
		rows = transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", newPosts[i].User)
		err = rows.Scan(&userId)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return newPosts, fmt.Errorf("can't find user"), 1
		}
		rows = transaction.QueryRow("INSERT INTO post (userid, forumid, created, parent, message, threadid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, isEdited, created::text",
			userId, forumId, time, newPosts[i].Parent, newPosts[i].Message, thread.Id)
		err = rows.Scan(&newPosts[i].Id, &newPosts[i].Edited, &newPosts[i].Date)
		if err != nil {
			log.Println(err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				log.Fatalln(errRollback)
			}
			return newPosts, err, 2
		}
		newPosts[i].Thread = thread.Id
		newPosts[i].Forum = thread.Forum
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
		rows := transaction.QueryRow("SELECT id, slug, created::text, title, message, votes, user_forum, forum FROM thread_full_view WHERE id = $1", thread.Id)
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
		rows := transaction.QueryRow("SELECT id, slug, created::text, title, message, votes, user_forum, forum FROM thread_full_view WHERE slug = $1",  thread.Slug)
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
	var actualId int32
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
	_, err = transaction.Exec("UPDATE thread SET (message, title) = ($1, $2) WHERE id = $3 ",  thread.Message, thread.Title, thread.Id)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	rows := transaction.QueryRow("SELECT id, slug, created::text, title, message, votes, user_forum, forum FROM thread_full_view WHERE slug = $1",  thread.Slug)
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