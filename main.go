package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/armagad/liblib/books"
	"github.com/armagad/liblib/search"
)

var gob_filename = "library.dev.gob"
var port = ":4040"

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// inv := inventory.NewService()
	// img := image.NewService()

	api := books.NewApi("books", gob_filename, []books.DecoratorService{})
	defer api.Close()

	mux := http.NewServeMux()
	mux.HandleFunc(api.CollectionRoot, api.Handler)
	mux.HandleFunc("/search/", search.HandleHttp())

	srv := http.Server{Handler: mux, Addr: port}

	l, _ := net.Listen("tcp", port)
	defer l.Close()
	go srv.Serve(l)
	<-c
	srv.Shutdown(context.Background())

	println()
}
