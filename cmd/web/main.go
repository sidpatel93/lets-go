package main

import (
	"flag"
	"log"
	"net/http"
)
func main() {
	addr:=flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	mux := http.NewServeMux()
	staticfileServer := http.FileServer(http.Dir("./ui/static/"))

	// Staitc file server route
	mux.Handle("/static/", http.StripPrefix("/static", staticfileServer))

	// other application routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}