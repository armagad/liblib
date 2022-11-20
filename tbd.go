package main

import (
	"encoding/json"
	"io"
)

func ReadBodyAsBook(body io.ReadCloser) (book Book, err error) {
	var b []byte
	b, err = io.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &book)
	return
}
