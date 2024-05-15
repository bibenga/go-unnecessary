package main

import (
	"log"
	"net/http"
	"time"
	"unnecessary/api-gorilla-gen/server"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ws: connect")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("ws: connected")

	wsw, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	_, err = wsw.Write([]byte("Olala"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	wsw.Close()
	log.Printf("hello has been sent")

	err = conn.WriteMessage(websocket.TextMessage, []byte("Olala2"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("hello2 has been sent")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	// ------------------------------------------------------------------------------------
	// https://github.com/gorilla/mux
	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler())

	fs := http.FileServer(http.Dir("api"))
	r.PathPrefix("/docs2").Handler(http.StripPrefix("/docs2/", fs))

	r.Path("/ws").HandlerFunc(wsHandler)

	rapi := r.NewRoute().Subrouter()
	rapi.Use(server.NewValidator())
	api := server.NewServer()
	server.HandlerWithOptions(api, server.GorillaServerOptions{
		BaseURL:    "/api",
		BaseRouter: rapi,
	})

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(log.Writer(), r),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Ready: http://127.0.0.1:8000/")
	log.Fatal(srv.ListenAndServe())
}
