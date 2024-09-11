package main

import (
	"flag"
	"log"
	"net/http"
)

type config struct {
	addr string
	staticDir string
}

func main() {
	// a new config struct
	var config config
	// define and parse command line flags to get the runtime values
	flag.StringVar(&config.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&config.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	
	mux := http.NewServeMux()
	staticfileServer := http.FileServer(http.Dir(config.staticDir))

	// Staitc file server route
	mux.Handle("/static/", http.StripPrefix("/static", staticfileServer))

	// other application routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", config.addr)
	err := http.ListenAndServe(config.addr, mux)
	log.Fatal(err)
}