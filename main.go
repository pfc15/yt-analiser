package main

import (
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	// "fmt"
	"log"
	"os"
)


func start_data_base() (*sql.DB){
	db, err := sql.Open("sqlite3", "mybd.sqlite3")
    if err != nil {
        log.Fatal(err)
    }
	query, err := os.ReadFile("sql/create.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}

	return db
}

func main() {
	if API_KEY == "SET-API-KEY" {
		log.Fatal("⚠️ Please set your YouTube API key in the API_KEY constant.")
	}

	
	yt, _ := newYubeClient(API_KEY)


	meta_dado, err := getVideoMetadata(yt, VIDEO_ID);
	if  err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
	db := start_data_base()
	meta_dado.saveData(db)

}