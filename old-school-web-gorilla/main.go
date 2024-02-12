package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"old-school-web-gorilla/templates/layout.html",
		"old-school-web-gorilla/templates/index.html",
	))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func main() {
	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler())

	fs := http.FileServer(http.Dir("old-school-web-gorilla/static"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", IndexHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(log.Writer(), r),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Ready: http://127.0.0.1:8000/")
	log.Fatal(srv.ListenAndServe())
}
