package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr string
	staticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// a new config struct
	var config config
	// define and parse command line flags to get the runtime values
	flag.StringVar(&config.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&config.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// a new application struct that contains the dependencies shared by the handlers and other files.
	app := application {
		errorLog: errorLog,
		infoLog: infoLog,
	}

	mux := http.NewServeMux()
	staticfileServer := http.FileServer(http.Dir(config.staticDir))

	// Staitc file server route
	mux.Handle("/static/", http.StripPrefix("/static", staticfileServer))

	// other application routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Configure the http.Server instance
	srv := http.Server{
		Addr: config.addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", config.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}