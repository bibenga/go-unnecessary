package main

// https://github.com/rabbitmq/amqp091-go

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	rlog := log.New(log.Default().Writer(), "rabbitmq - ", log.Default().Flags())
	amqp.SetLogger(rlog)

	log.Print("is connecting...")
	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()
	log.Print("is connected")

	log.Print("is opening channel...")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	ch.Qos(1, 1024*1024, false)
	log.Print("channel opened")

	log.Print("will be declared queue...")
	q, err := ch.QueueDeclare(
		"q1",  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "failed to declare a queue")
	log.Print("queue was declared...")

	log.Print("will be consumed...")
	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "failed to register a consumer")
	log.Print("was consumed")

	errors := conn.NotifyClose(make(chan *amqp.Error))

	// forever := make(chan int)
	// go func() {
	// 	log.Print("waiting errors...")
	// 	for err := range errors {
	// 		log.Printf("error received: %+v", err)
	// 		forever <- 1
	// 	}
	// }()
	// go func() {
	// 	log.Print("waiting messages...")
	// 	for d := range msgs {
	// 		log.Printf("received a message: %+v", d)
	// 	}
	// 	log.Print("stop processing messages...")
	// }()
	// log.Printf("waiting for messages. To exit press CTRL+C")
	// <-forever
	// log.Printf("Exit")

	log.Printf("waiting for messages. To exit press CTRL+C")
out:
	for {
		select {
		case message := <-messages:
			log.Printf("received a message: %+v", message)
		case err := <-errors:
			log.Printf("error received: %+v", err)
			break out
		}
	}
	log.Printf("Exit")
}
