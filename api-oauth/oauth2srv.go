package main

// curl -v 'http://127.0.0.1:8000/authorize?client_id=000000&response_type=code&redirect_uri=http://127.0.0.1:8000/callback'
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" "http://127.0.0.1:8000/token" -d "client_id=000000" -d "code=NMY4ODC4ZGQTZMUZZI0ZMJK1LTG3MZETMTCWMTNHNTEYYWY4" -d "redirect_uri=http://127.0.0.1:8000/callback" -d "grant_type=authorization_code"
// or
// curl -X POST -H "Content-Type: application/x-www-form-urlencoded" "http://127.0.0.1:8000/token" -d "client_id=000000" -d "username=a" -d "password=a" -d "redirect_uri=http://127.0.0.1:8000/callback" -d "grant_type=password"

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
)

func main() {
	manager := manage.NewDefaultManager()

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("SECRET"), jwt.SigningMethodHS512))

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Public: true,
		Domain: "http://127.0.0.1:8000",
		// Secret: "999999",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if username == "a" && password == "a" {
			userID = "a"
		}
		return
	})

	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return "000000", nil
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	})
	// http.HandleFunc("/auth", authHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		slog.Error("Internal Error", "error", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		slog.Error("Response Error", "error", re.Error.Error())
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("/authorize")

		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("/token")
		srv.HandleTokenRequest(w, r)
	})

	slog.Info("Listen on 127.0.0.1:8000")
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		slog.Error("Error", "error", err)
	}
}
