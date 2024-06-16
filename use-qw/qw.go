package main

import (
	"log"
	"log/slog"

	"github.com/bibenga/gomod1"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")
	slog.Info("main")

	gomod1.SayHello()
}
