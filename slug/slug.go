package main

import (
	"log"

	"github.com/gosimple/slug"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	examples := []string{
		"Español",
		"Más",
		"Надежда",
	}

	for _, example := range examples {
		log.Print(slug.Make(example))
	}
}
