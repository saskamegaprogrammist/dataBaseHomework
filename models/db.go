package models

import (
	"fmt"
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"log"
)

type DB struct {
	Forum int `json:"forum"`
	Post int `json:"post"`
	Thread int `json:"thread"`
	User int `json:"user"`
}

func (db *DB) GetStatus() error {
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
	rows := transaction.QueryRow("SELECT COUNT(*) FROM forum")
	err = rows.Scan(&db.Forum)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	rows = transaction.QueryRow("SELECT COUNT(*) FROM forum_user ")
	err = rows.Scan(&db.User)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	rows = transaction.QueryRow("SELECT COUNT(*) FROM thread ")
	err = rows.Scan(&db.Thread)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	rows = transaction.QueryRow("SELECT COUNT(*) FROM post ")
	err = rows.Scan(&db.Post)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func (db *DB) Clear() error {
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
	_, err = transaction.Exec(`DELETE FROM votes;`)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("votes")
	transaction, err = dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	_, err = transaction.Exec(`DELETE FROM post;`)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("post")
	transaction, err = dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	_, err = transaction.Exec(`DELETE FROM thread;`)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("thread")
	transaction, err = dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	_, err = transaction.Exec(`DELETE FROM forum;`)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("forum")
	transaction, err = dataBase.Begin()
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	_, err = transaction.Exec(`DELETE FROM forum_user;`)
	if err != nil {
		log.Println(err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			log.Fatalln(errRollback)
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("forum_user")
	return nil
}