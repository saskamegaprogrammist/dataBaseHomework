package models

import (
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
)

type User struct {
	Id int32 `json:"-"`
	Nickname string `json:"nickname"`
	Email string `json:"email"`
	Fullname string `json:"fullname"`
	About string `json:"about"`
}
type Exists struct {
	Exists bool
}
func (user *User) CreateUser() ([]User, error) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
	}
	var usersExists[] User
	rows := transaction.QueryRow("SELECT * FROM forum_user WHERE nickname = $1", user.Nickname)
	if err != nil {
		log.Println(err)
	}
	var userExists User
	_ = rows.Scan(&userExists.Id, &userExists.Nickname, &userExists.Email, &userExists.Fullname, &userExists.About)
	if userExists.Id != 0  {
		usersExists = append(usersExists, userExists)
	}
	var userExistsEmail User
	rows = transaction.QueryRow("SELECT * FROM forum_user WHERE email = $1", user.Email)
	if err != nil {
		log.Println(err)
	}
	_ = rows.Scan(&userExistsEmail.Id, &userExistsEmail.Nickname, &userExistsEmail.Email, &userExistsEmail.Fullname, &userExistsEmail.About)
	if userExistsEmail.Id != 0 {
		usersExists = append(usersExists, userExistsEmail)
	}
	if len(usersExists) > 0 {
		err = transaction.Commit()
		if err != nil {
			log.Fatalln(err)
		}
		return usersExists, nil
	}
	rows = transaction.QueryRow("INSERT INTO forum_user (nickname, email, fullname, about) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Nickname, user.Email, user.Fullname, user.About)
	err = rows.Scan(&user.Id)
	if err != nil {
		log.Fatalln(err)
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil, nil
}