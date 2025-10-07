package main

import (
	"log"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeClient struct {
	service *youtube.Service
	apiKey  string
}

type VideoMetadata struct {
	Title        string
	Description  string
	ChannelTitle string
	PublishedAt  string
	ViewCount    uint64
	LikeCount    uint64
}

type CommentData struct {
	ID          string
	Author      string
	Text        string
	LikeCount   uint64
	PublishedAt string
	ReplyTo     string
}

// Update interface to use minimal types
type YouTubeClientInterface interface {
	callVideoData(videoID string) (*VideoMetadata, error)
	callCommentData(videoID string, maxResult int64) ([]CommentData, error)
}

func newYubeClient(apikey string) (*YouTubeClient, error) {
	ctx := context.Background()
	yt_client := &YouTubeClient{apiKey: apikey}

	service, err := youtube.NewService(ctx, option.WithAPIKey(yt_client.apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %v", err)
	}
	yt_client.service = service

	return yt_client, nil
}

func (yt *YouTubeClient) callVideoData(videoID string) (*VideoMetadata, error) {
	call := yt.service.Videos.List([]string{"snippet", "statistics"}).Id(videoID)
	response, err := call.Do()
	if err != nil || len(response.Items) == 0 {
		return nil, fmt.Errorf("error fetching video metadata: %v", err)
	}
	video := response.Items[0]
	return &VideoMetadata{
		Title:        video.Snippet.Title,
		Description:  video.Snippet.Description,
		ChannelTitle: video.Snippet.ChannelTitle,
		PublishedAt:  video.Snippet.PublishedAt,
		ViewCount:    video.Statistics.ViewCount,
		LikeCount:    video.Statistics.LikeCount,
	}, nil
}

func (yt *YouTubeClient) callCommentData(videoID string, maxResult int64) ([]CommentData, error) {
	call := yt.service.CommentThreads.List([]string{"snippet", "replies"}).
		VideoId(videoID).
		TextFormat("plainText").
		MaxResults(maxResult)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}
	var comments []CommentData
	for _, item := range response.Items {
		top := item.Snippet.TopLevelComment.Snippet
		comments = append(comments, CommentData{
			ID:          item.Id,
			Author:      top.AuthorDisplayName,
			Text:        top.TextDisplay,
			LikeCount:   uint64(top.LikeCount),
			PublishedAt: top.PublishedAt,
			ReplyTo:     "",
		})
		if item.Replies != nil {
			for _, reply := range item.Replies.Comments {
				r := reply.Snippet
				comments = append(comments, CommentData{
					ID:          reply.Id,
					Author:      r.AuthorDisplayName,
					Text:        r.TextDisplay,
					LikeCount:   uint64(r.LikeCount),
					PublishedAt: r.PublishedAt,
					ReplyTo:     item.Id,
				})
			}
		}
	}
	return comments, nil
}

func (yt *YouTubeClient) callCanalAllVideoList(canal string) ([]string){

    // 1. Get uploads playlist ID
    chResp, err := yt.service.Channels.List([]string{"contentDetails"}).
        Id(canal).
        Do()
    if err != nil {
        log.Fatalf("Channels.List error: %v", err)
    }
    if len(chResp.Items) == 0 {
        log.Fatalf("No channel found with ID %s", canal)
    }
    uploadsPlaylistID := chResp.Items[0].ContentDetails.RelatedPlaylists.Uploads
    fmt.Println("Uploads playlist ID:", uploadsPlaylistID)

    // 2. Iterate through the playlist items (videos)
    var allVideoIDs []string
    nextPageToken := ""
    for {
        plReq := yt.service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
            PlaylistId(uploadsPlaylistID).
            MaxResults(50)
        if nextPageToken != "" {
            plReq = plReq.PageToken(nextPageToken)
        }
        plResp, err := plReq.Do()
        if err != nil {
            log.Fatalf("PlaylistItems.List error: %v", err)
        }

        for _, item := range plResp.Items {
            // item.ContentDetails.VideoId is the video ID
            allVideoIDs = append(allVideoIDs, item.ContentDetails.VideoId)
            // You can also read snippet (title, publish date etc.)
            fmt.Printf("Video: %s — %s\n", item.ContentDetails.VideoId, item.Snippet.Title)
        }

        if plResp.NextPageToken == "" {
            break
        }
        nextPageToken = plResp.NextPageToken
    }

	return allVideoIDs
}
