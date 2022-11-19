package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
			var book Book
			b, _ := io.ReadAll(r.Body)
			// fmt.Println(string(b))
			err := json.Unmarshal(b, &book)
			if err != nil {
				fmt.Println(err)
			}
			book.Id = badkey
			book.URLid = BookId(fmt.Sprintf("%d", book.Id))
			badkey++

			fmt.Println(book.Id, book.Title)

			// -- db
			libraryLock.Lock()
			library[book.URLid] = book
			libraryLock.Unlock()
		} else {
			// LIST
			c := make(chan Book)

			// -- db
			go func() {
				libraryLock.RLock()
				for _, book := range library {
					c <- book
				}
				libraryLock.RUnlock()
				close(c)
			}()

			list := []Book{}
			for book := range c {
				list = append(list, book)
			}
			b, err := json.Marshal(list)
			if err != nil {
				fmt.Println("error:", err)
			}
			w.Write(b)
		}

	} else {
		bookId := BookId(r.URL.Path[len(COLLECTION_ROOT):])
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			// UPDATE
			b, _ := io.ReadAll(r.Body)
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
				b, err := json.Marshal(book)
				if err != nil {
					fmt.Println("error:", err)
				}
				w.Write(b)
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

	Id    uint64 `json:"-"`
	URLid BookId `json:"-"`
}
