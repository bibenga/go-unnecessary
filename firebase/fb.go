package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	ctx := context.Background()

	opt := option.WithCredentialsFile("firebase/firebase.json")
	conf := &firebase.Config{
		DatabaseURL: "https://aprende-palabras-b61c4.firebaseio.com",
	}
	slog.Info("try connect to google", "DatabaseURL", conf.DatabaseURL)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		slog.Error("error initializing app", "error", err)
		panic(err)
	}

	db, err := app.Firestore(ctx)
	if err != nil {
		slog.Error("Error initializing database client", "error", err)
		panic(err)
	}
	defer db.Close()
	slog.Info("connected", "db", db)

	// words := db.Collection("words").Limit(5).Documents(ctx)
	// for {
	// 	doc, err := words.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		slog.Error("Failed to iterate", "error", err)
	// 		panic(err)
	// 	}
	// 	log.Println("row", "data", doc.Data())
	// }

	auth, err := app.Auth(ctx)
	if err != nil {
		slog.Error("Failed create auth client", "error", err)
		panic(err)
	}

	customToken, err := auth.CustomToken(ctx, "")
	if err != nil {
		slog.Error("error minting custom token", "error", err)
		panic(err)
	}
	slog.Info("Got custom token", "customToken", customToken)

	idToken, err := signInWithCustomToken(customToken)
	if err != nil {
		slog.Error("error minting id token", "error", err)
		panic(err)
	}
	slog.Info("Got custom token", "idToken", idToken)

	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		slog.Error("error verifying ID token", "error", err)
	}
	slog.Info("Verified ID token", "token", token)
}

func signInWithCustomToken(token string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	verifyCustomTokenURL := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
	apiKey := ""
	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, apiKey), req)
	if err != nil {
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}
	return respBody.IDToken, err
}

func postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
