package models

import (
	"fmt"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
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
	rows = transaction.QueryRow("SELECT * FROM forum WHERE (created, title) = ($1, $2)", thread.Date, thread.Title)
	_ = rows.Scan(&threadExists.Id, &threadExists.Slug, &threadExists.Date, &threadExists.Title, &threadExists.Message, &threadExists.Votes, &threadExistsForumId, &threadExistsUserId)
	if threadExists.Id != 0  {
		rows = transaction.QueryRow("SELECT nickname FROM forum_user WHERE id = $1", threadExistsUserId)
		_ = rows.Scan(&threadExists.User)
		rows = transaction.QueryRow("SELECT slug FROM forum WHERE id = $1", threadExistsForumId)
		_ = rows.Scan(&threadExists.Forum)
		return threadExists, fmt.Errorf("thread exists")
	}
	rows = transaction.QueryRow("INSERT INTO thread (slug, created, title, message, forumid, userid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		thread.Slug, thread.Date, thread.Title, thread.Message, forumId, userId)
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