package main

import (
	"fmt"
	"log"
	"os"
)




func main() {
	if API_KEY == "SET-API-KEY" {
		log.Fatal("⚠️ Please set your YouTube API key in the API_KEY constant.")
	}

	meta_dado, err := getVideoMetadata(API_KEY, VIDEO_ID);
	if  err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Título do Vídeo: %s\n", meta_dado.titulo)
	fmt.Printf("Canal: %s\n", meta_dado.canal)
	fmt.Printf("Descrição: %s\n", meta_dado.descricao)
	fmt.Printf("Data de Publicação: %s\n", meta_dado.data_publicacao)
	fmt.Printf("Visualizações: %d\n", meta_dado.quant_view)
	fmt.Printf("Curtidas: %d\n", meta_dado.quant_like)
}