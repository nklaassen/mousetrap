package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	hub := newHub()

	fs := http.FileServer(http.Dir("www"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { hub.serveWs(w, r) })

	portNum := "8080"
	if len(os.Args) > 1 {
		portNum = os.Args[1]
	}

	log.Printf("starting mousetrap on port %v\n", portNum)
	log.Fatal(http.ListenAndServe(":"+portNum, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
