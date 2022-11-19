package main

import (
	"encoding/json"
	"io"
	"log"
)

func ReadBodyAsBook(body io.ReadCloser) (book Book) {
	var b []byte
	var err error
	b, err = io.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(b, &book)
	if err != nil {
		log.Println(err)
	}
	return book
}
