package main

import (
	
)

type Comentario struct {
	id string
	autor string
	texto string
	like uint64 
	data_publicacao string
	reply []Comentario
}

