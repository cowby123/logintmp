package main

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "uploaddb"
)

var db *sql.DB

//SystemTime 系統時間
var SystemTime int64

func main() {

	log.Println("main")
	INIT()
	router := initRouter()
	router.Run(":8080")
}

//INIT 初始化
func INIT() {

	//連接到DB
	ConnectDB()

	//初始化android vmhost 控制器

}

//ConnectDB 連接到PostgresDB
func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

}
