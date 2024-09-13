package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr string
	dsn string
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
	flag.StringVar(&config.dsn, "dsn", "root:@/snippetbox?parseTime=true", "MySQL DB connection string")
	flag.StringVar(&config.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create a connection pool
	db, err := openDB(config.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// defer the close of the db connection
	defer db.Close()

	// a new application struct that contains the dependencies shared by the handlers and other files.
	app := application {
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Configure the http.Server instance
	srv := http.Server{
		Addr: config.addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", config.addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}