package books

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/armagad/liblib/search"
)

func (api *BooksApi) Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == api.CollectionRoot {
		switch r.Method {
		case http.MethodPost:
			// CREATE
			book, err := ReadBodyAsBook(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(err.Error()))
				return
			}
			book.Id = api.badkey
			book.URLid = BookId(fmt.Sprintf("%d", book.Id))
			api.badkey++

			// -- db
			go search.AddItems(book)
			api.libraryLock.Lock()
			api.library[book.URLid] = book
			api.libraryLock.Unlock()
		case http.MethodGet:
			// LIST
			c := make(chan Book)

			// -- db
			go func() {
				api.libraryLock.RLock()
				for _, book := range api.library {
					c <- book
				}
				api.libraryLock.RUnlock()
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
		bookId := BookId(r.URL.Path[len(api.CollectionRoot):])
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			// UPDATE
			upbook, err := ReadBodyAsBook(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(err.Error()))
				return
			}
			if book, ok := api.library[bookId]; ok {
				if upbook.Title != "" {
					book.Title = upbook.Title
				}
				if upbook.Author != "" {
					book.Author = upbook.Author
				}
				if book.Abridged {
					book.Abridged = upbook.Abridged
				}

				api.libraryLock.Lock()
				api.library[bookId] = book
				api.libraryLock.Unlock()
			}
		case http.MethodDelete:
			// DELETE
			api.libraryLock.RLock()
			book, ok := api.library[bookId]
			api.libraryLock.RUnlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			go search.PurgeItems(book)
			api.libraryLock.Lock()
			delete(api.library, bookId)
			api.libraryLock.Unlock()
		case http.MethodGet:
			// GET
			api.libraryLock.RLock()
			book, ok := api.library[bookId]
			api.libraryLock.RUnlock()
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

func ReadBodyAsBook(body io.ReadCloser) (book Book, err error) {
	var b []byte
	b, err = io.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &book)
	return
}
