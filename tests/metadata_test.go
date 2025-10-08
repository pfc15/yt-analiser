package domain_test

import (
	"log"
	"testing"
    "youtube_tracker/internal/ytclient"
    "youtube_tracker/internal/persistence"
    "youtube_tracker/internal/domain"

)



func TestGetVideoMetadata(t *testing.T) {
    log.Println("teste GetVideoMetaData: ")
	mock := &ytclient.MockYtClient{}

	meta, _ := domain.GetVideoMetadata(mock, "1234")

	if meta.Quant_view!=123 || meta.Titulo!="Test Title"{
        t.Fail()
    }
    
}

func TestSaveData(t *testing.T){
    log.Println("teste GetVideoMetaData: ")
    db := persistence.Start_data_base()
    meta := domain.MetaDado{
        Titulo: "video Teste",
        Video_id: "12;DROP TABLE COMENTARIO;",
        Descricao: "essa é uma descrição",
        Canal: "asdmksad",
        Data_publicacao: "11/09/2001",
        Quant_view: 11,
        Quant_like: 10,
        Comentarios: make([]domain.Comentario, 0),
    }

    err := meta.SaveVideoData(db)
    if err!= nil{
        t.Fail()
    }

    var id string
    err = db.QueryRow(`SELECT VIDEO.id FROM VIDEO WHERE VIDEO.id==?`, meta.Video_id).Scan(&id)
    if err != nil{
        log.Println(err)
        t.Fail()
    }

    err = meta.SaveVideoData(db)
    if err!= nil{
        t.Fail()
    }
}
