package main

import (
	"context"
	"log"
	"time"
)

func genCtx(ctx context.Context) <-chan int {
	dst := make(chan int, 2)
	n := 1
	go func() {
		for {
			log.Print("gorutine")
			select {
			case <-ctx.Done():
				log.Print("done...")
				return // returning not to leak the goroutine
			case dst <- n:
				log.Printf("sent %d", n)
				time.Sleep(100 * time.Millisecond)
				n++
			}
		}
	}()
	return dst
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	log.Print("context")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := range genCtx(ctx) {
		log.Printf("recv %d", n)
		if n == 5 {
			break
		}
	}

	cancel()
	time.Sleep(100 * time.Millisecond)
}
