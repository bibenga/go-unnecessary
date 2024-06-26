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
	// publisher watermill.Pun

	go process(pubSub)
	time.Sleep(100 * time.Millisecond)
	go publish(pubSub)

	slog.Info("sleep")
	time.Sleep(3 * time.Second)

	pubSub.Close()
	slog.Info("terminate")
}

func publish(publisher message.Publisher) {
	slog.Info("> publish")
	for i := range 4 {
		msg := message.NewMessage(watermill.NewUUID(), []byte(fmt.Sprintf("message %v", i)))
		slog.Info("Send message", "ID", msg.UUID, "Payload", string(msg.Payload))
		// slog.Info("Send message", "msg", msg)
		err := publisher.Publish("example.topic", msg)
		if err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	slog.Info("< publish")
}

func process(subscriper message.Subscriber) {
	slog.Info("> process")
	messages, err := subscriper.Subscribe(context.Background(), "example.topic")
	if err != nil {
		panic(err)
	}
	for msg := range messages {
		slog.Info("Received message", "ID", msg.UUID, "Payload", string(msg.Payload))
		msg.Ack()
	}
	slog.Info("< process")
}
