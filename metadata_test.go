package main

import (
	"log"
	"testing"
)


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


func TestSaveData(t *testing.T){
    db := start_data_base()
    meta := MetaDado{
        titulo: "video Teste",
        video_id: "12;DROP TABLE COMENTARIO;",
        descricao: "essa é uma descrição",
        canal: "asdmksad",
        data_publicacao: "11/09/2001",
        quant_view: 11,
        quant_like: 10,
        comentarios: make([]Comentario, 0),
    }

    err := meta.saveData(db)
    if err!= nil{
        t.Fail()
    }

    var id string
    err = db.QueryRow(`SELECT VIDEO.id FROM VIDEO WHERE VIDEO.id==?`, meta.video_id).Scan(&id)
    if err != nil{
        log.Println(err)
        t.Fail()
    }

    err = meta.saveData(db)
    if err!= nil{
        t.Fail()
    }
}
