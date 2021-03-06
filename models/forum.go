package models

import (
	"fmt"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
)

type Forum struct {
	Id int `json:"-"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"user"`
	Threads int `json:"threads"`
	Posts int `json:"posts"`
}

func (forum *Forum) CreateForum() (Forum, error) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	var forumExists Forum
	if err != nil {
		//log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return forumExists, err
	}
	rows := transaction.QueryRow("SELECT id, nickname FROM forum_user WHERE nickname = $1", forum.User)
	var userId int
	err = rows.Scan(&userId, &forum.User)
	if err != nil || userId == 0 {
		//log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return forumExists, fmt.Errorf("can't find user with nickname %s", forum.User)
	}
	rows = transaction.QueryRow("SELECT * FROM forum WHERE slug = $1", forum.Slug)
	_ = rows.Scan(&forumExists.Id, &forumExists.Slug, &forumExists.Title, &forumExists.Posts, &forumExists.Threads, &forumExists.User)
	if forumExists.Id != 0  {
		return forumExists, fmt.Errorf("forum exists")
	}

	rows = transaction.QueryRow("INSERT INTO forum (slug, title, usernick) VALUES ($1, $2, $3) RETURNING id, slug",
		forum.Slug, forum.Title, forum.User)
	err = rows.Scan(&forum.Id, &forum.Slug)
	if err != nil {
		//log.Println(err)
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
		//log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}
	rows := transaction.QueryRow("SELECT * FROM forum WHERE slug = $1", forumSlug)
	err = rows.Scan(&forum.Id, &forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.User)
	if err != nil {
		//log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find forum with slug %s", forumSlug)
	}
	err = transaction.Commit()
	if err != nil {
		//log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}
	return nil
}
