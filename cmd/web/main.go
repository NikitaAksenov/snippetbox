package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/NikitaAksenov/snippetbox/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "snippetbox.db", "sqlite3 data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlite3", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
