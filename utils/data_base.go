package utils

import (
	"fmt"
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
)

var dataBasePool *pgx.ConnPool

func CreateAddress(user, password, host, name string) string {
	return  fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s",
	user, password, host, name)
}

func CreateDataBaseConnection(user, password, host, name string, maxConn int) {
	dataBaseConfig := CreateAddress(user, password, host, name)
	connectionConfig, err := pgx.ParseConnectionString(dataBaseConfig)
	if err != nil {
		log.Println(err);
		return
	}
	dataBasePool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: connectionConfig,
		MaxConnections: maxConn,
	})
	if err != nil {
		log.Println(err);
		return
	}
}

func InitDataBase() {
	initScript, err := ioutil.ReadFile("./sql_scripts/database_creation.sql")
	if err != nil {
		log.Println(err)
	}
	_, err = dataBasePool.Exec(string(initScript))
	if err != nil {
		log.Println(err)
	}

}

func GetDataBase() *pgx.ConnPool {
	return dataBasePool
}