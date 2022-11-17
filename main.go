package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var port = ":4040"

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc(COLLECTION_ROOT, booooooks)

	go http.ListenAndServe(port, nil)
	<-c
}

const COLLECTION_ROOT = "/books/"

func booooooks(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == COLLECTION_ROOT {
		if r.Method == http.MethodPost {
			// CREATE

		} else {
			// LIST

		}

	} else {
		// bookId := r.URL.Path[len(COLLECTION_ROOT):]
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			// UPDATE
		case http.MethodDelete:
			// DELETE
		case http.MethodGet:
			// GET
		default:

		}
	}

	w.Write([]byte("Sup"))
}

type Book struct {
	Title  string
	Author string

	Abridged bool
	Edition  string
	ISBN     string

	CopiesAvailable uint

	Id    uint64
	URLid string
}
