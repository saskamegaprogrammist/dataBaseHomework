package models

import (
	"fmt"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
)

type Forum struct {
	Id int32 `json:"-"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"user"`
	Threads int32 `json:"threads"`
	Posts int32 `json:"posts"`
}

func (forum *Forum) CreateForum() (Forum, error) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	var forumExists Forum
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return forumExists, err
	}
	rows := transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", forum.User)
	var userId int32
	err = rows.Scan(&userId)
	if err != nil || userId == 0 {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return forumExists, fmt.Errorf("can't find user with nickname %s", forum.User)
	}
	rows = transaction.QueryRow("SELECT * FROM forum WHERE slug = $1", forum.Slug)
	var forumExistsUserId int32
	_ = rows.Scan(&forumExists.Id, &forumExists.Slug, &forumExists.Title, &forumExists.Posts, &forumExists.Threads, &forumExistsUserId)
	if forumExists.Id != 0  {
		rows = transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", forumExistsUserId)
		_ = rows.Scan(&forumExists.User)
		return forumExists, fmt.Errorf("forum exists")
	}

	rows = transaction.QueryRow("INSERT INTO forum (slug, title, userid) VALUES ($1, $2, $3) RETURNING id",
		forum.Slug, forum.Title, userId)
	err = rows.Scan(&forum.Id)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return forumExists, err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return forumExists, nil
}

func (forum *Forum) GetForum(forumSlug string) error {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}
	var forumUserId int32
	rows := transaction.QueryRow("SELECT * FROM forum WHERE slug = $1", forumSlug)
	err = rows.Scan(&forum.Id, &forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forumUserId)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find forum with slug %s", forumSlug)
	}
	rows = transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", forumUserId)
	err = rows.Scan(&forum.User)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}
	return nil
}
