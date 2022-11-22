package books

import (
	"encoding/gob"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/armagad/liblib/search"
)

type DecoratorService interface {
	Get(BookId)
}

type BooksApi struct {
	CollectionRoot string
	dbFilename     string
	library        map[BookId]Book
	libraryLock    sync.RWMutex
	badkey         uint64
}

type BookId string

type Book struct {
	Title  string
	Author string

	Abridged bool

	CopiesAvailable uint

	Id    uint64 `json:"-"`
	URLid BookId `json:"-"`
}

func NewApi(pattern, filename string, svcs []DecoratorService) *BooksApi {
	api := &BooksApi{
		CollectionRoot: "/" + pattern + "/",
		dbFilename:     filename,
		library:        make(map[BookId]Book),
		badkey:         1}

	if filename != "" {
		dataFile, err := os.OpenFile(filename, os.O_CREATE, 0644)
		if err != nil {
			println(err)
		}
		defer dataFile.Close()

		err = gob.NewDecoder(dataFile).Decode(&api.library)
		if err != io.EOF && err != nil {
			println(err)
		}
		ugh := api.badkey
		for k, v := range api.library {
			utmp, err := strconv.ParseUint(string(k), 10, 64)
			if err == nil && utmp > ugh {
				ugh = utmp
			}
			search.AddItems(v)
		}
		api.badkey = ugh
	}
	return api
}

func (api *BooksApi) Close() {

	if api.dbFilename != "" {
		dbf, err := os.OpenFile(api.dbFilename, os.O_WRONLY, 0644)
		if err != nil {
			println(err)
		}
		err = gob.NewEncoder(dbf).Encode(api.library)
		if err != nil {
			println(err)
		}

	}
}
