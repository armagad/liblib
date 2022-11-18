package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"sync"
	"io/ioutil"
)

var port = ":4040"

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc(COLLECTION_ROOT, booooooks)

	go http.ListenAndServe(port, nil)
	<-c
}

var library map[BookId]Book
type BookId string

var libraryLock sync.RWMutex
var badkey uint64 = 1

func init() {
	library = map[BookId]Book{}
}


const COLLECTION_ROOT = "/books/"

func booooooks(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == COLLECTION_ROOT {
		if r.Method == http.MethodPost {
			// CREATE
			b, _ := ioutil.ReadAll(r.Body)
			book := Book{Title: string(b)}
			fmt.Println(badkey, book.Title)
			libraryLock.Lock()
			k := BookId(fmt.Sprintf("%d", badkey))
			library[k] = book
			badkey ++
			libraryLock.Unlock()
		} else {
			// LIST
			libraryLock.RLock()
			for id, book := range library {
				w.Write([]byte(fmt.Sprintf("%s %s\n",id, book.Title)))
			}
			libraryLock.RUnlock()
		}

	} else {
		bookId := BookId(r.URL.Path[len(COLLECTION_ROOT):])
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			// UPDATE
			b, _ := ioutil.ReadAll(r.Body)
			if book, ok := library[bookId]; ok {
				book.Title = string(b)
				libraryLock.Lock()
				library[bookId] = book
				libraryLock.Unlock()
			}
		case http.MethodDelete:
			// DELETE
			libraryLock.Lock()
			delete(library, bookId)
			libraryLock.Unlock()
		case http.MethodGet:
			// GET
			libraryLock.RLock()
			book, ok := library[bookId]
			libraryLock.RUnlock()
			if ok {
				w.Write([]byte(book.Title))
			}
		default:
		}
	}
}

type Book struct {
	Title  string
	Author string

	Abridged bool
	Edition  string
	ISBN     string

	CopiesAvailable uint

	Id    uint64
	URLid BookId
}
