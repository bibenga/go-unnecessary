package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")
	log.Printf("------------")
	log.Printf("playStdLog: %v", 1)

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(true, false),
	)

	messages, err := pubSub.Subscribe(context.Background(), "example.topic")
	if err != nil {
		panic(err)
	}

	go process(messages)
	go publish(pubSub)

	slog.Info("sleep")
	time.Sleep(3 * time.Second)
	slog.Info("terminate")
}

func publish(pubSub *gochannel.GoChannel) {
	for i := range 4 {
		msg := message.NewMessage(watermill.NewUUID(), []byte(fmt.Sprintf("message %v", i)))
		slog.Info("Send message", "ID", msg.UUID, "Payload", string(msg.Payload))
		// slog.Info("Send message", "msg", msg)
		err := pubSub.Publish("example.topic", msg)
		if err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		slog.Info("Received message", "ID", msg.UUID, "Payload", string(msg.Payload))
		msg.Ack()
	}
}
