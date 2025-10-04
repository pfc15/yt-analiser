package main

import "testing"


type MockYtClient struct{}

func (m *MockYtClient) callVideoData(videoID string) (*VideoMetadata, error) {
    return &VideoMetadata{
        Title:        "Test Title",
        Description:  "Test Description",
        ChannelTitle: "Test Channel",
        PublishedAt:  "2022-01-01T00:00:00Z",
        ViewCount:    123,
        LikeCount:    45,
    }, nil
}

func (m *MockYtClient) callCommentData(videoID string, maxResult int64) ([]CommentData, error) {
    return []CommentData{
        {
            ID:          "c1",
            Author:      "Author1",
            Text:        "Comment1",
            LikeCount:   1,
            PublishedAt: "2022-01-01T00:00:00Z",
            ReplyTo:     "",
        },
        {
            ID:          "r1",
            Author:      "Author2",
            Text:        "Reply1",
            LikeCount:   2,
            PublishedAt: "2022-01-01T01:00:00Z",
            ReplyTo:     "c1",
        },
    }, nil
}

func TestGetVideoMetadata(t *testing.T) {
	mock := &MockYtClient{}

	meta, _ := getVideoMetadata(mock, "1234")

	if meta.quant_view!=123 || meta.titulo!="Test Title"{
        t.Fail()
    }
}

