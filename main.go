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

	http.HandleFunc("/books/", booooooks)

	go http.ListenAndServe(port, nil)
	<-c
}

func booooooks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sup"))
}
