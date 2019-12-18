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

			row, err := transaction.Query("SELECT extname::text FROM pg_extension")

			if err != nil {
				log.Println(err)
			}

			for row.Next() {
				var name string
				err = rows.Scan(&name)
				fmt.Println(name)

				if err != nil {
					log.Println(err)
				}

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
	fmt.Println(user.Email, userExists)
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
	fmt.Println(user.Id)
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
	fmt.Println(user.Nickname)
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
	var userExistsId int
	rows := transaction.QueryRow("SELECT id FROM forum_user WHERE nickname = $1", user.Nickname)
	err = rows.Scan(&userExistsId)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find user with nickname %s", user.Nickname), 1
	}
	if user.Email != "" {
		_, err = transaction.Exec("UPDATE forum_user SET email = $2 WHERE nickname = $1",  user.Nickname, user.Email)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return fmt.Errorf("email exists %s", user.Email), 2
		}
	}
	if user.About != "" {
		_, err = transaction.Exec("UPDATE forum_user SET about = $2 WHERE nickname = $1",  user.Nickname, user.About)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err, 2
		}
	}
	if user.Fullname != "" {
		_, err = transaction.Exec("UPDATE forum_user SET fullname = $2 WHERE nickname = $1",  user.Nickname, user.Fullname)
		if err != nil {
			log.Println(err)
			err = transaction.Rollback()
			if err != nil {
				log.Fatalln(err)
			}
			return err, 2
		}
	}
	rows = transaction.QueryRow("SELECT * FROM forum_user WHERE nickname = $1", user.Nickname)
	err = rows.Scan(&user.Id, &user.Nickname, &user.Email, &user.Fullname, &user.About)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return fmt.Errorf("can't find user with nickname %s", user.Nickname), 2
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil, 0
}

func GetUsersByForum(params utils.SearchParams, forumSlug string) ([]User, error) {
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
	row := transaction.QueryRow("SELECT id FROM forum WHERE slug = $1", forumSlug)
	err = row.Scan(&forumId)
	if err != nil {
		log.Println(err)
		err = transaction.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return usersFound, fmt.Errorf("can't find forum with slug %s", forumSlug)
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