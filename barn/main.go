package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bibenga/barn-go/examples"
	humanhash "github.com/bibenga/humanhash-go"
)

func main() {
	examples.Setup(true)

	db := examples.InitDb(false)
	defer db.Close()

	_, name, err := humanhash.NewUuid()
	if err != nil {
		panic(err)
	}

	slog.Info("name: %s", name)

	_, cancel := context.WithCancel(context.Background())

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)
	s := <-osSignal
	slog.Info("os signal received", "signal", s)

	cancel()

}
