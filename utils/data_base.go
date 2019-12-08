package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	_ "github.com/jackc/pgx/pgxpool"
	"io/ioutil"
	"log"
)

var dataBasePool *pgxpool.Pool

func createAddress(user, password, host, name string, maxConn int) string {
	return  fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s pool_max_conns=%d",
	user, password, host, name, maxConn)
}

func CreateDataBaseConnection(user, password, host, name string, maxConn int) {
	dataBaseConfig := createAddress(user, password, host, name, maxConn)
	config, err := pgxpool.ParseConfig(dataBaseConfig)
	if err != nil {
		log.Println(err);
		return
	}
	dataBasePool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Println(err);
		return
	}
}

func InitDataBase() {
	initScript, err := ioutil.ReadFile("sql_sripts/database_creation.sql")
	if err != nil {
		log.Println(err)
	}
	_, err = dataBasePool.Exec(context.Background(), string(initScript));
	if err != nil {
		log.Println(err)
	}

}