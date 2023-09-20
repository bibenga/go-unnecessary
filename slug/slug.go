package main

import (
	"log"

	"github.com/gosimple/slug"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	examples := []string{
		"надежда",
		"fe",
		"love",
		"español",
		"más o menos",
	}

	for _, example := range examples {
		log.Print(slug.Make(example))
	}
}
