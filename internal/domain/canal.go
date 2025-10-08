package domain

import (
	"database/sql"
	"fmt"
	"youtube_tracker/internal/ytclient"
)

type canal struct {
	Id string
	Nome string
	IsSubscribed bool
	Videos []string
}

func GetChannel(yt ytclient.YouTubeClientInterface, db *sql.DB,id string) (*canal, error) {
	
	c := canal{}
	err := db.QueryRow(`SELECT id, nome, isSubscribed  FROM CANAL WHERE CANAL.id==?`, id).Scan(&c.Id, &c.Nome, &c.IsSubscribed)
	
	if err == sql.ErrNoRows{
		info, err := yt.CallCanal(id)
		if err != nil {
			return nil, err
		}
		c.Id = info.Id
		c.Nome = info.Nome
		c.IsSubscribed = true

		db.Exec(`INSERT INTO CANAL(id, nome, isSusbcribed) VALUES (?,?,?)`, c.Id,c.Nome,c.IsSubscribed)
	}

	if c.IsSubscribed {
		c.Videos = yt.CallCanalVideoList(id, true)
	}

	return &c, nil
}

func (c *canal) SaveCanalVideosData(yt ytclient.YouTubeClientInterface, db *sql.DB) (error) {
	if len(c.Videos) == 0 {
		return fmt.Errorf("channel not subscribed or no videos on channel")
	}

	for _, video_id := range c.Videos {
		video_dado, err := GetVideoMetadata(yt, video_id)
		if err != nil {
			return err
		}
		video_dado.SaveVideoData(db)
	}
	return nil
}
