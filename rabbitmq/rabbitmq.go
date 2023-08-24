package main

// https://github.com/rabbitmq/amqp091-go

import (
	"log"
	"time"

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

	log.Print("try to connect")
	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()
	log.Print("the connection was established")

	log.Print("try to open a channel...")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	ch.Qos(1, 1024*1024, false)
	log.Print("the channel was opened")

	log.Print("try to declare a queue...")
	q, err := ch.QueueDeclare(
		"q1",  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "failed to declare a queue")
	log.Print("the queue was created or updated...")

	log.Print("try to consume...")
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
	log.Print("the consumer was connected")

	errors := conn.NotifyClose(make(chan *amqp.Error))

	log.Printf("the program is waiting for messages. To exit press CTRL+C")
	stop := time.NewTimer(5 * time.Second)
	defer stop.Stop()

	err = process(messages, errors, stop)
	failOnError(err, "failed on process messages")

	stop.Stop()

	// out:
	// 	for {
	// 		select {
	// 		case message := <-messages:
	// 			log.Printf("received a message: %+v", message)
	// 		case err := <-errors:
	// 			log.Printf("error received: %+v", err)
	// 			break out
	// 		}
	// 	}
	log.Printf("Exit")
}

func process(messages <-chan amqp.Delivery, errors chan *amqp.Error, stop *time.Timer) error {
	// message loop
	for {
		select {
		case message := <-messages:
			err := processOne(&message)
			if err != nil {
				return err
			}
		case <-stop.C:
			log.Printf("exit by timer")
			return nil
		case err := <-errors:
			log.Printf("error received: %+v", err)
			return err
		}
	}
}

func processOne(message *amqp.Delivery) error {
	// process concrete message
	log.Printf("received a message: %+v", message)
	return nil
}
