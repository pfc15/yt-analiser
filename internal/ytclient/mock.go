package ytclient



type MockYtClient struct{}

func (m *MockYtClient) CallVideoData(videoID string) (*VideoMetadata, error) {
    return &VideoMetadata{
        Title:        "Test Title",
        Description:  "Test Description",
        Channel_id: "1234",
        PublishedAt:  "2022-01-01T00:00:00Z",
        ViewCount:    123,
        LikeCount:    45,
    }, nil
}

func (m *MockYtClient) CallCommentData(videoID string, maxResult int64) ([]CommentData, error) {
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


func (m *MockYtClient) CallCanalVideoList(id string, paginacao bool) ([]string) {
    return []string{
        "1234",
        "4321",
    }
}

func (m *MockYtClient) CallCanal(canalId string) (*CanalMetadata, error) {
    return &CanalMetadata{
        Id:"c1234",
        Nome: "canal teste",
    }, nil
}


