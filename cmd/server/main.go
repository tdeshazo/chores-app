package main

import (
	"log"
	"net/http"

	"chores-app/internal/db"
	"chores-app/internal/handlers"
	"chores-app/internal/web"
)

const (
	defaultDBPath = "chores.db"
	defaultAddr   = ":5000"
)

func main() {
	conn, err := db.Open(defaultDBPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer conn.Close()

	tmpl, err := web.ParseTemplates()
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	pages := handlers.Pages{DB: conn, Tmpl: tmpl}
	api := handlers.API{DB: conn}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", pages.Index)
	mux.HandleFunc("/kid/", pages.Kid)
	mux.HandleFunc("/api/update_status", api.UpdateStatus)

	log.Printf("Starting server on %s ...", defaultAddr)
	if err := http.ListenAndServe(defaultAddr, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
