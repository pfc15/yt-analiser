package persistence

import (
	_ "github.com/mattn/go-sqlite3"
	"youtube_tracker/internal/domain"
	"database/sql"
	"log"
	"os"
)



func Start_data_base() (*sql.DB){
	db, err := sql.Open("sqlite3", domain.ROOT+"/mydb.sqlite3")
    if err != nil {
        log.Fatal(err)
    }
	query, err := os.ReadFile(domain.ROOT+"/sql/create.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}

	return db
}
