package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type MetaDado struct {
	titulo string
	video_id string
	descricao string
	canal string
	data_publicacao string 
	quant_view uint64
	quant_like uint64
	comentarios []Comentario
}


func (meta *MetaDado) getComments(yt_client YouTubeClientInterface, videoID string) error {
	response, err := yt_client.callCommentData(videoID, 20)
	if err!=nil{
		return err
	}
	if len(response) > 0 {
		fmt.Printf("Comments for Video ID: %s\n", videoID)

		for _, item := range response {
			comentario := Comentario{
				id: item.ID,
				autor: item.Author,
				texto: item.Text,
				like: uint64(item.LikeCount),
				data_publicacao: item.PublishedAt,
				reply: item.ReplyTo,
			}
			meta.comentarios = append(meta.comentarios, comentario)
		}
	} else {
		fmt.Printf("No comments found for Video ID: %s\n", videoID)
	}

	return nil
}

// func getVideoId(yt_client YouTubeClientInterface, canal string) ([]string, error) {
	
// }

func getVideoMetadata(yt_client YouTubeClientInterface, videoID string) (*MetaDado,error) {
	meta := MetaDado{}
	meta.video_id = videoID
	response, err := yt_client.callVideoData(videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching video metadata: %v", err)
	}

	if response !=nil{

		meta.titulo = response.Title
		meta.descricao = response.Description
		meta.canal = response.ChannelTitle
		meta.data_publicacao = response.PublishedAt
		meta.quant_view = response.ViewCount
		meta.quant_like = response.LikeCount
		


		err := meta.getComments(yt_client, videoID)
		if err != nil {
			return nil, fmt.Errorf("error fetching comments: %v", err)
		}

		return &meta, nil

	} else {
		fmt.Printf("Vídeo com ID '%s' não encontrado.\n", videoID)
		return nil, errors.New("vídeo com ID '%s' não encontrado")
	}
}

func (meta *MetaDado) saveData(db *sql.DB) ( error) {
	var isOnDb bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM VIDEO WHERE VIDEO.id = ?)`, meta.video_id).Scan(&isOnDb)
	
	if err != nil {
		fmt.Println("Query error:", err)
		return err
	}
	
	if !isOnDb {
		insert_video, err := db.Prepare(
		`INSERT INTO VIDEO (id, titulo, descricao, canal, data_publicacao) VALUES (?, ?, ?, ?, ?);`)
		if err !=nil{
			return err
		}
		_, err = insert_video.Exec(meta.video_id, meta.titulo, meta.descricao,meta.canal, meta.data_publicacao)
		if err !=nil{
			return err
		}
		insert_video.Close()
	}
	
	insert_metrica, err := db.Prepare(`INSERT INTO METRICA(video_id, data_coleta, quant_view, quant_like) VALUES(?,?,?,?)`)
	if err!=nil{return err}
	
	insert_metrica.Exec(meta.video_id, time.Now().Format("SS:MM:HH DD-MM-YYYY"),meta.quant_view, meta.quant_like)

	insert_comentario, err := db.Prepare(`INSERT INTO COMENTARIO(id, video_id, autor, texto, data_publicacao, reply) VALUES (?,?,?,?,?,?);`)
	for _, comentario := range meta.comentarios{
		if err!=nil{
			return  err
		}
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM COMENTARIO WHERE COMENTARIO.id = ?)`, comentario.id).Scan(&isOnDb)
		if err!=nil{
			return  err
		}
		if !isOnDb{
			if comentario.reply == ""{
			_, err = insert_comentario.Exec(comentario.id, meta.video_id, comentario.autor, comentario.texto, comentario.data_publicacao, nil)
			} else{
				_, err = insert_comentario.Exec(comentario.id, meta.video_id, comentario.autor, comentario.texto, comentario.data_publicacao, comentario.reply)
			}
			if err!=nil{
				return err
			}
		}
	}
    return nil
}
