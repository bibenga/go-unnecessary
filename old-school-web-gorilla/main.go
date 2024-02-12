package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	store := sessions.NewCookieStore([]byte("32-byte-long-auth-key"))

	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler())

	fs := http.FileServer(http.Dir("old-school-web-gorilla/static"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "_gorilla_session")
		if err != nil {
			panic(err)
		}
		// fmt.Printf("ID=%v, IsNew=%v\n", session.ID, session.IsNew)
		session.Values["foo"] = "bar"
		if err := session.Save(r, w); err != nil {
			panic(err)
		}

		tmpl := template.Must(template.ParseFiles(
			"old-school-web-gorilla/templates/layout.html",
			"old-school-web-gorilla/templates/index.html",
		))
		if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
			panic(err)
		}
	})

	CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))

	srv := &http.Server{
		Handler: handlers.LoggingHandler(
			log.Writer(),
			CSRF(r),
		),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Ready: http://127.0.0.1:8000/")
	log.Fatal(srv.ListenAndServe())
}
