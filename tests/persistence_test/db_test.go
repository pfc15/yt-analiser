package persistence_test

import (
	"database/sql"
	"log"
	"testing"
	"youtube_tracker/internal/persistence"

	_ "github.com/mattn/go-sqlite3"
)

func check_table(db *sql.DB, name_table string) bool {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?;`, name_table).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	return count>0
}


func TestStart_data_base(t *testing.T) {
	db := persistence.Start_data_base()

	if !check_table(db, "CANAL") || !check_table(db, "VIDEO") || !check_table(db, "METRICA") || !check_table(db, "COMENTARIO") {
		log.Println("não foi criado todos as tabelas")
		t.Fail()
	}

	if check_table(db, "nao tem") {
		log.Println("tabelas a mais foram criadas")
		t.Fail()
	}
}

