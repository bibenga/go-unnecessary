package main

import (
	"context"
	"log"
	"log/slog"
	"time"
)

func sender(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		slog.Info("> sender")
		defer func() {
			close(dst)
			slog.Info("< sender")
		}()
		for {
			select {
			case <-ctx.Done():
				slog.Info("sender", "err", ctx.Err())
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

func listener(ctx context.Context, gen <-chan int) {
	slog.Info("> listener")
	defer func() {
		slog.Info("< listener")
	}()
	for {
		slog.Info("listener")
		select {
		case <-ctx.Done():
			slog.Info("listener", "err", ctx.Err())
			return // returning not to leak the goroutine
		case n := <-gen:
			slog.Info("listener recv", "n", n)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	slog.Info(">")
	gen := sender(ctx)

	go listener(ctx, gen)

out:
	for {
		slog.Info("main")
		select {
		case <-ctx.Done():
			slog.Info("main done", "err", ctx.Err())
			break out
		case n := <-gen:
			slog.Info("main recv", "n", n)
		}
	}
	slog.Info("<")

	// wait a gorutines, it is bad
	time.Sleep(1 * time.Second)
}
