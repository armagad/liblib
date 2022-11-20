package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
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

	mux := http.NewServeMux()
	mux.HandleFunc(COLLECTION_ROOT, booooooks)
	var srv http.Server
	srv.Handler = mux
	srv.Addr = port

	l, _ := net.Listen("tcp", port)
	defer l.Close()
	go srv.Serve(l)
	<-c
	srv.Shutdown(context.Background())
	println()
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
		switch r.Method {
		case http.MethodPost:
			// CREATE
			book, err := ReadBodyAsBook(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(err.Error()))
				return
			}
			book.Id = badkey
			book.URLid = BookId(fmt.Sprintf("%d", book.Id))
			badkey++

			// -- db
			libraryLock.Lock()
			library[book.URLid] = book
			libraryLock.Unlock()
		case http.MethodGet:
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
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(b)
		default:
			w.WriteHeader(http.StatusNotFound)
		}

	} else {
		bookId := BookId(r.URL.Path[len(COLLECTION_ROOT):])
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			// UPDATE
			upbook, err := ReadBodyAsBook(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(err.Error()))
				return
			}
			if book, ok := library[bookId]; ok {
				if upbook.Title != "" {
					book.Title = upbook.Title
				}
				if upbook.Author != "" {
					book.Author = upbook.Author
				}
				if book.Abridged {
					book.Abridged = upbook.Abridged
				}

				libraryLock.Lock()
				library[bookId] = book
				libraryLock.Unlock()
			}
		case http.MethodDelete:
			// DELETE
			libraryLock.RLock()
			_, ok := library[bookId]
			libraryLock.RUnlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
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
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				w.Write(b)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

type Book struct {
	Title  string
	Author string

	Abridged bool

	CopiesAvailable uint

	Id    uint64 `json:"-"`
	URLid BookId `json:"-"`
}
