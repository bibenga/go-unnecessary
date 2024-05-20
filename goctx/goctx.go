package main

import (
	"context"
	"log"
	"log/slog"
	"time"
)

func genCtx(ctx context.Context) <-chan int {
	dst := make(chan int, 2)
	n := 1
	go func() {
		for {
			slog.Info("gorutine")
			select {
			case <-ctx.Done():
				err := ctx.Err()
				slog.Info("done", "err", err)
				return // returning not to leak the goroutine
			case dst <- n:
				slog.Info("sent", "n", n)
				time.Sleep(500 * time.Millisecond)
				n++
			}
		}
	}()
	return dst
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	slog.Info(">")
	// for n := range genCtx(ctx) {
	// 	slog.Info("recv", "n", n)
	// 	if n >= 100 {
	// 		break
	// 	}
	// }
	gen := genCtx(ctx)
out:
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			slog.Info("main done", "err", err)
			break out
		case n := <-gen:
			slog.Info("main recv", "n", n)
			time.Sleep(500 * time.Millisecond)
			n++
		}
	}
	slog.Info("<")
}
