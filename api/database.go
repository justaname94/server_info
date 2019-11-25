package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbhost = "dbhost"
	dbport = "dbport"
	dbuser = "dbuser"
	dbname = "dbname"
)

func InitDb() *sql.DB {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbname])

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database!")
	return db
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	conf[dbhost] = "localhost"
	conf[dbport] = "26257"
	conf[dbuser] = "go_admin"
	conf[dbname] = "website"
	return conf
}
