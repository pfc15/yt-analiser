package persistence

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"os"
)

var ROOT = "/home/pfc15/Documents/aleatorio/youtube_tracker/"


func Start_data_base() (*sql.DB){
	db, err := sql.Open("sqlite3", ROOT+"/mydb.sqlite3")
    if err != nil {
        log.Fatal(err)
    }
	query, err := os.ReadFile(ROOT+"/sql/create.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}

	return db
}