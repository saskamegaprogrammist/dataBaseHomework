package models

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
)

type User struct {
	Id int `json:"-"`
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
	err = rows.Scan(&userExists.Id, &userExists.Nickname, &userExists.Email, &userExists.Fullname, &userExists.About)
	if err != nil {
		log.Println(err)
	}

	if userExists.Id != 0  {
		usersExists = append(usersExists, userExists)
	}
	var userExistsEmail User
	rows = transaction.QueryRow("SELECT * FROM forum_user WHERE email = $1", user.Email)
	if err != nil {
		log.Println(err)
	}
	err = rows.Scan(&userExistsEmail.Id, &userExistsEmail.Nickname, &userExistsEmail.Email, &userExistsEmail.Fullname, &userExistsEmail.About)
	if err != nil {
		log.Println(err)
	}
	if userExistsEmail.Id != 0 && userExists.Id != userExistsEmail.Id  {
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
	rows := transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", user.Nickname)
	err = rows.Scan(&user.Id)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find user with nickname %s", user.Nickname), 1
	}
	rows = transaction.QueryRow("UPDATE forum_user SET email = coalesce(nullif($2, ''), email)," +
			" about = coalesce(nullif($3, ''), about), fullname = coalesce(nullif($4, ''), fullname)" +
		"WHERE nickname = $1 RETURNING email, fullname, about",  user.Nickname, user.Email, user.About, user.Fullname)
	err = rows.Scan(&user.Email, &user.Fullname, &user.About)
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

func GetUsersByForum(limit int, since string, desc bool, forumSlug string) ([]User, error) {
	usersFound := make([]User, 0)
	dataBase := utils.GetDataBase()
	transaction, err := dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if err != nil {
			log.Fatalln(errRollback)
		}
		return usersFound, err
	}
	var forumId int
	row := transaction.QueryRow("SELECT id, slug FROM forum WHERE slug = $1", forumSlug)
	err = row.Scan(&forumId, &forumSlug)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return usersFound, fmt.Errorf("can't find forum with slug %s", forumSlug)
	}

	//sqlSelect := "SELECT about, fullname, nickname, email FROM forum_user " +
	//				"JOIN (SELECT COALESCE(p_usernick, t_usernick) as merge_nick FROM ( " +
	//						 "SELECT DISTINCT usernick as p_usernick FROM post WHERE forumslug = $1) as p " +
	//	"FULL OUTER JOIN ( " +
	//	"SELECT DISTINCT usernick as t_usernick  FROM thread WHERE forumslug = $1) " +
	//	"as t ON p.p_usernick = t.t_usernick) " +
	//				"as u ON u.merge_nick = forum_user.nickname"

	sqlSelect := "SELECT about, fullname, nickname, email FROM forum_user " +
		"JOIN (SELECT DISTINCT usernick as merge_nick FROM forum_user_new "+
			"WHERE forumslug = $1 ) " +
			"as u ON u.merge_nick = forum_user.nickname"
	var rows *pgx.Rows

	if desc {
		if limit != -1 {
			if since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname   < $2 ORDER BY (nickname ) DESC LIMIT $3", forumSlug, since, limit)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname ) DESC LIMIT $2", forumSlug, limit)
			}
		} else {
			if since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname  < $2 ORDER BY (nickname ) DESC", forumSlug, since)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname ) DESC", forumSlug)
			}
		}
	} else {
		if limit != -1 {
			if since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname   > $2 ORDER BY (nickname ) LIMIT $3", forumSlug, since, limit)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname ) LIMIT $2", forumSlug, limit)
			}
		} else {
			if since != "" {
				rows, err = transaction.Query(sqlSelect+" WHERE nickname   > $2 ORDER BY (nickname ) ", forumSlug, since)
			} else {
				rows, err = transaction.Query(sqlSelect+" ORDER BY (nickname ) ", forumSlug)
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