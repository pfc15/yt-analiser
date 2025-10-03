package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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

func (meta *MetaDado) getComments(service *youtube.Service, videoID string) error {
	call := service.CommentThreads.List([]string{"snippet", "replies"}).
		VideoId(videoID).
		TextFormat("plainText").
		MaxResults(20) // adjust as needed

	response, err := call.Do()
	if err != nil {
		return err
	}

	if len(response.Items) > 0 {
		fmt.Printf("Comments for Video ID: %s\n", videoID)

		for _, item := range response.Items {
			topLevel := item.Snippet.TopLevelComment.Snippet
			comentario := Comentario{
				id: item.Id,
				autor: topLevel.AuthorDisplayName,
				texto: topLevel.TextDisplay,
				like: uint64(topLevel.LikeCount),
				data_publicacao: topLevel.PublishedAt,
				reply: make([]Comentario, 0),
			}

			
			if item.Replies != nil && len(item.Replies.Comments) > 0 {
				fmt.Println("  Replies:")

				for _, reply := range item.Replies.Comments {
					replySnippet := reply.Snippet
					reply := Comentario{
						id: reply.Id,
						autor: replySnippet.AuthorDisplayName,
						texto: replySnippet.TextDisplay,
						like: uint64(replySnippet.LikeCount),
						data_publicacao: replySnippet.PublishedAt,
						reply: make([]Comentario, 0),
					}
					comentario.reply = append(comentario.reply, reply)
				}
			}
			meta.comentarios = append(meta.comentarios, comentario)
		}


		if response.NextPageToken != "" {
			fmt.Println("\nThere are more comments. Use NextPageToken to fetch them.")
			// you can implement pagination here if you want
		}
	} else {
		fmt.Printf("No comments found for Video ID: %s\n", videoID)
	}

	return nil
}

func getVideoMetadata(apiKey, videoID string) (*MetaDado,error) {
	meta := MetaDado{}
	meta.video_id = videoID
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %v", err)
	}

	call := service.Videos.List([]string{"snippet", "statistics"}).Id(videoID)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error fetching video metadata: %v", err)
	}

	if len(response.Items) > 0 {
		video := response.Items[0]
		snippet := video.Snippet
		stats := video.Statistics

		meta.titulo = snippet.Title
		meta.descricao = snippet.Description
		meta.canal = snippet.ChannelTitle
		meta.data_publicacao = snippet.PublishedAt
		meta.quant_view = stats.ViewCount
		meta.quant_like = stats.LikeCount
		


		err := meta.getComments(service, videoID)
		if err != nil {
			return nil, fmt.Errorf("error fetching comments: %v", err)
		}

		return &meta, nil

	} else {
		fmt.Printf("Vídeo com ID '%s' não encontrado.\n", videoID)
		return nil, errors.New("vídeo com ID '%s' não encontrado")
	}
}