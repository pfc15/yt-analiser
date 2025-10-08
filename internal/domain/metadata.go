package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"youtube_tracker/internal/ytclient"
)

type MetaDado struct {
	Titulo string
	Video_id string
	Descricao string
	Canal string
	Data_publicacao string 
	Quant_view uint64
	Quant_like uint64
	Comentarios []Comentario
}


func (meta *MetaDado) GetComments(yt ytclient.YouTubeClientInterface, videoID string) error {
	response, err := yt.CallCommentData(videoID, BATCH_COMENTARIO)
	if err!=nil{
		return err
	}
	if len(response) > 0 {
		fmt.Printf("Comments for Video ID: %s\n", videoID)

		for _, item := range response {
			comentario := Comentario{
				Id: item.ID,
				Autor: item.Author,
				Texto: item.Text,
				Like: uint64(item.LikeCount),
				Data_publicacao: item.PublishedAt,
				Reply: item.ReplyTo,
			}
			meta.Comentarios = append(meta.Comentarios, comentario)
		}
	} else {
		fmt.Printf("No comments found for Video ID: %s\n", videoID)
	}

	return nil
}


func GetVideoMetadata(yt ytclient.YouTubeClientInterface, videoID string) (*MetaDado,error) {
	meta := MetaDado{}
	meta.Video_id = videoID
	response, err := yt.CallVideoData(videoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching video metadata: %v", err)
	}

	if response !=nil{

		meta.Titulo = response.Title
		meta.Descricao = response.Description
		meta.Canal = response.Channel_id
		meta.Data_publicacao = response.PublishedAt
		meta.Quant_view = response.ViewCount
		meta.Quant_like = response.LikeCount
		


		err := meta.GetComments(yt, videoID)
		if err != nil {
			return nil, fmt.Errorf("error fetching comments: %v", err)
		}

		return &meta, nil

	} else {
		fmt.Printf("Vídeo com ID '%s' não encontrado.\n", videoID)
		return nil, errors.New("vídeo com ID '%s' não encontrado")
	}
}

func (meta *MetaDado) SaveVideoData(db *sql.DB) ( error) {
	var isOnDb bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM VIDEO WHERE VIDEO.id = ?)`, meta.Video_id).Scan(&isOnDb)
	
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
		_, err = insert_video.Exec(meta.Video_id, meta.Titulo, meta.Descricao,meta.Canal, meta.Data_publicacao)
		if err !=nil{
			return err
		}
		insert_video.Close()
	}
	
	insert_metrica, err := db.Prepare(`INSERT INTO METRICA(video_id, data_coleta, quant_view, quant_like) VALUES(?,?,?,?)`)
	if err!=nil{return err}
	
	insert_metrica.Exec(meta.Video_id, time.Now().Format("SS:MM:HH DD-MM-YYYY"),meta.Quant_view, meta.Quant_like)

	insert_comentario, err := db.Prepare(`INSERT INTO COMENTARIO(id, video_id, autor, texto, data_publicacao, reply) VALUES (?,?,?,?,?,?);`)
	for _, comentario := range meta.Comentarios{
		if err!=nil{
			return  err
		}
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM COMENTARIO WHERE COMENTARIO.id = ?)`, comentario.Id).Scan(&isOnDb)
		if err!=nil{
			return  err
		}
		if !isOnDb{
			if comentario.Reply == ""{
			_, err = insert_comentario.Exec(comentario.Id, meta.Video_id, comentario.Autor, comentario.Texto, comentario.Data_publicacao, nil)
			} else{
				_, err = insert_comentario.Exec(comentario.Id, meta.Video_id, comentario.Autor, comentario.Texto, comentario.Data_publicacao, comentario.Reply)
			}
			if err!=nil{
				return err
			}
		}
	}
    return nil
}
