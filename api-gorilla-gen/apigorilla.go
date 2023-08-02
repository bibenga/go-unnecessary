package main

import (
	"log"
	"net/http"
	"time"
	"unnecessary/api-gorilla-gen/server"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	// ------------------------------------------------------------------------------------
	// https://github.com/gorilla/mux
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("api"))
	r.PathPrefix("/docs2").Handler(http.StripPrefix("/docs2/", fs))

	r.Use(server.NewValidator())

	api := server.NewServer()
	server.HandlerFromMuxWithBaseURL(api, r, "/api")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Ready: http://127.0.0.1:8000/")
	log.Fatal(srv.ListenAndServe())
}
