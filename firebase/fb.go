package main

import (
	"context"
	"log"

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

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	db, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing database client: %v", err)
	}
	defer db.Close()
	log.Print(db)

	words := db.Collection("words").Documents(ctx)
	for {
		doc, err := words.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		log.Println(doc.Data())
	}
}
