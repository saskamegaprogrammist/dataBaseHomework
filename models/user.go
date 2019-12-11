package models

import (
	"fmt"
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
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return usersExists, nil
	}
	rows = transaction.QueryRow("INSERT INTO forum_user (nickname, email, fullname, about) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Nickname, user.Email, user.Fullname, user.About)
	err = rows.Scan(&user.Id)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return nil, err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil, nil
}

func (user *User) GetUser(userNickname string) error {
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
	rows := transaction.QueryRow("SELECT * FROM forum_user WHERE nickname = $1", userNickname)
	err = rows.Scan(&user.Id, &user.Nickname, &user.Email, &user.Fullname, &user.About)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find user with nickname %s", userNickname)
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

func (user *User) UpdateUser() (error, int) {
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
	}
	var userExistsNickname User
	rows := transaction.QueryRow("SELECT * FROM forum_user WHERE nickname = $1", user.Nickname)
	err = rows.Scan(&userExistsNickname.Id, &userExistsNickname.Nickname, &userExistsNickname.Email, &userExistsNickname.Fullname, &userExistsNickname.About)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find user with nickname %s", user.Nickname), 1
	}
	_, err = transaction.Exec("UPDATE forum_user SET (email, fullname, about) = ($2, $3, $4) WHERE nickname = $1;",  user.Nickname, user.Email, user.Fullname, user.About)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("email exists %s", user.Email), 2
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil, 0
}