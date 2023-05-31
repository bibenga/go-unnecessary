package main

// https://github.com/segmentio/kafka-go

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// const broker = "kafka1:9092"
const broker = "host.docker.internal:9092"

func listTopics() {
	log.Print("list topics")
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	// i do not know but it seems that to read the list of topics we do not need to connect to the controller
	// controller, err := conn.Controller()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// var controllerConn *kafka.Conn
	// controllerConn, err = kafka.Dial(
	// 	"tcp",
	// 	net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	// )
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer controllerConn.Close()
	controllerConn := conn

	partitions, err := controllerConn.ReadPartitions()
	if err != nil {
		log.Panic(err)
	}

	// m := map[string]struct{}{}
	// for _, p := range partitions {
	// 	m[p.Topic] = struct{}{}
	// }
	// for k := range m {
	// 	log.Print(k)
	// }

	accept := func(t string) bool {
		switch t {
		case "__consumer_offsets",
			"__transaction_state":
			return false
		}
		return true
	}

	for _, p := range partitions {
		if !accept(p.Topic) {
			continue
		}
		log.Printf("%s, %d", p.Topic, p.ID)
	}
}

func createTopic(topic string) {
	log.Printf("create topic %s", topic)
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Panic(err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Panic(err)
	}
	defer controllerConn.Close()
	// controllerConn := conn

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     5,
			ReplicationFactor: 1,
			ConfigEntries: []kafka.ConfigEntry{
				{ConfigName: "retention.bytes", ConfigValue: "10000000"},
				{ConfigName: "retention.ms", ConfigValue: "3600000"},
			},
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("topic %s created", topic)
}

func deleteTopic(topic string) {
	log.Printf("delete topic %s", topic)

	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Panic(err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	)
	if err != nil {
		log.Panic(err)
	}
	defer controllerConn.Close()

	err = controllerConn.DeleteTopics(topic)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("topic %s deleted", topic)
}

type PredefinedPartitionBalancer struct {
	fallback kafka.RoundRobin
}

func (b *PredefinedPartitionBalancer) Balance(message kafka.Message, partitions ...int) int {
	if message.Partition < 0 {
		return b.fallback.Balance(message, partitions...)
	}
	return message.Partition
}

func readAndWrite(topic string) {
	ctx := context.Background()

	rlog := log.New(log.Default().Writer(), "[reader] - ", log.Default().Flags())
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			ClientID:  "go-writer",
		},
		Topic:          topic,
		GroupID:        "UnnecessaryID",
		Logger:         rlog,
		ErrorLogger:    rlog,
		IsolationLevel: kafka.ReadCommitted,
	})
	defer r.Close()

	wlog := log.New(log.Default().Writer(), "[writer] - ", log.Default().Flags())
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Dialer: &kafka.Dialer{
			Timeout:         10 * time.Second,
			DualStack:       true,
			ClientID:        "go-writer",
			TransactionalID: "go-writer",
		},
		Balancer:    &PredefinedPartitionBalancer{},
		Logger:      wlog,
		ErrorLogger: wlog,
	})
	w.AllowAutoTopicCreation = false
	err := w.WriteMessages(
		ctx,
		kafka.Message{Topic: topic, Partition: 0, Value: []byte("v0")},
		kafka.Message{Topic: topic, Partition: 1, Value: []byte("v1")},
		kafka.Message{Topic: topic, Partition: 4, Value: []byte("v4")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	defer w.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatal("failed to read message", err)
			break
		}
		log.Printf("message at offset Topic=(%s) Partition=%d Offset=%d Key=(%s) Value=(%s)",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	listTopics()
	// deleteTopic("t1")
	// createTopic("t1")
	readAndWrite("t1")
}
