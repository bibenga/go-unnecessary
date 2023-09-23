package main

import (
	"context"
	"log"
	"log/slog"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
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
	}
	defer db.Close()
	slog.Info("connected", "db", db)

	words := db.Collection("words").Limit(5).Documents(ctx)
	for {
		doc, err := words.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			slog.Error("Failed to iterate", "error", err)
		}
		log.Println("row", "data", doc.Data())
	}
}
