package main

import (
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	// "fmt"
	"log"
	"os"
	"net/http"
)


func start_data_base() (*sql.DB){
	db, err := sql.Open("sqlite3", "mydb.sqlite3")
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
	

	http.HandleFunc("/youtube/callback", youtubeCallback)

    go subscribe("UCuAXFkgsw1L7xaCfnd5JJOw") // replace with your channel ID

    log.Println("Server started on :8002")
    log.Fatal(http.ListenAndServe(":8002", nil))
	

	// yt, _ := newYubeClient(API_KEY)

	// channel_id := "UCuAXFkgsw1L7xaCfnd5JJOw"

	// yt.callCanalAllVideoList(channel_id)
	

	// meta_dado, err := getVideoMetadata(yt, VIDEO_ID);
	// if  err != nil {
	// 	log.Fatalf("Error: %v", err)
	// 	os.Exit(1)
	// }
	// db := start_data_base()
	// meta_dado.saveData(db)

}