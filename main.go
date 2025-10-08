package main

import (
    
	"youtube_tracker/internal/ytclient"
	"youtube_tracker/internal/persistence"
	"youtube_tracker/internal/domain"
	// "fmt"
	"log"
	"os"
)

func main() {
	
	yt, _ := ytclient.NewYubeClient(domain.API_KEY)

	channel_id := "UCuAXFkgsw1L7xaCfnd5JJOw"
	db := persistence.Start_data_base()
	defer db.Close()

	canal, err := domain.GetChannel(yt, db, channel_id)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	canal.SaveCanalVideosData(yt, db)

	// yt.callCanalVideoList(channel_id, false)
	

	// video_dado, err := getVideoMetadata(yt, VIDEO_ID);
	// if  err != nil {
	// 	log.Fatalf("Error: %v", err)
	// 	os.Exit(1)
	// }
	
	// video_dado.saveVideoData(db)

}