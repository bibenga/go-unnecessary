package main

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"
	"unnecessary/api-gorilla-gen/server"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Ok   bool
	Time time.Time
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("ws: connect")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade error", "error", err)
		return
	}
	defer conn.Close()
	slog.Info("ws: connected")

	m := Message{Ok: true, Time: time.Now()}
	mb, err := json.Marshal(m)
	if err != nil {
		slog.Error("json error", "error", err)
		panic(err)
	}

	wsw, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		slog.Error("writer error", "error", err)
		return
	}
	_, err = wsw.Write(mb)
	if err != nil {
		slog.Error("hello error", "error", err)
		return
	}
	wsw.Close()
	slog.Info("hello has been sent")

	err = conn.WriteMessage(websocket.TextMessage, mb)
	if err != nil {
		slog.Error("hello2 error", "error", err)
		return
	}
	slog.Info("hello2 has been sent")

	err = conn.WriteJSON(m)
	if err != nil {
		slog.Error("hello3 error", "error", err)
		return
	}
	slog.Info("hello3 has been sent")

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			slog.Error("reader error", "error", err)
			return
		}
		if err := conn.WriteMessage(messageType, payload); err != nil {
			slog.Error("writer error", "error", err)
			return
		}
	}
}

func writeLog(writer io.Writer, params handlers.LogFormatterParams) {
	slog.Info("access",
		"method", params.Request.Method,
		"url", params.URL.Path,
		"status", params.StatusCode,
		"size", params.Size,
		"time", params.TimeStamp,
	)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	// ------------------------------------------------------------------------------------
	// https://github.com/gorilla/mux
	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler())
	r.Use(handlers.ProxyHeaders)

	fs := http.FileServer(http.Dir("api"))
	r.PathPrefix("/docs2").Handler(http.StripPrefix("/docs2/", fs))

	r.Path("/ws").HandlerFunc(wsHandler)

	rapi := r.NewRoute().Subrouter()
	rapi.Use(server.NewValidator())
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))
	rapi.Use(csrfMiddleware)

	api := server.NewServer()
	server.HandlerWithOptions(api, server.GorillaServerOptions{
		BaseURL:    "/api",
		BaseRouter: rapi,
	})

	srv := &http.Server{
		Handler:      handlers.CustomLoggingHandler(nil, r, writeLog),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	slog.Info("Ready: http://127.0.0.1:8000/")
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("error", "error", err)
	}
}
