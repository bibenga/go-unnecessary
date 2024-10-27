package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	llama "github.com/go-skynet/go-llama.cpp"
)

var (
	threads   = 2
	tokens    = 128
	gpulayers = 0
	seed      = -1
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	l := slog.New(slog.Default().Handler())
	slog.SetDefault(l)

	model := ""

	l.Info("hello")
	_, err := llama.New(
		model,
		llama.EnableF16Memory,
		llama.SetContext(128),
		llama.EnableEmbeddings,
		llama.SetGPULayers(gpulayers),
	)
	if err != nil {
		l.Error("Loading the model failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	fmt.Printf("Model loaded successfully.\n")

}
