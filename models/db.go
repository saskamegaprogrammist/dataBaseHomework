package models

import (
	"github.com/saskamegaprogrammist/dataBaseHomework/utils"
	"io/ioutil"
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
	return nil
}

func (db *DB) Clear() error {
	initScript, err := ioutil.ReadFile("./sql_scripts/database_creation.sql")
	if err != nil {
		log.Println(err)
		return err
	}
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
	_, err = dataBase.Exec(string(initScript))
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