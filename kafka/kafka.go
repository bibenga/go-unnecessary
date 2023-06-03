package main

// https://github.com/segmentio/kafka-go

import (
	"context"
	"log"
	"net"
	"sort"
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

	topics := make(map[string][]int)
	for _, p := range partitions {
		if !accept(p.Topic) {
			continue
		}
		partitions, ok := topics[p.Topic]
		if ok {
			partitions = append(partitions, p.ID)
		} else {
			partitions = []int{p.ID}
		}
		sort.Ints(partitions)
		topics[p.Topic] = partitions
	}
	if len(topics) == 0 {
		log.Printf("there are no topics")
	} else {
		for topic, partitions := range topics {
			log.Printf("%s - %v", topic, partitions)
		}
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

func writer(ctx context.Context, topic string) {
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
	defer w.Close()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	log.Print("writing...")
	func() {
		for {
			select {
			case <-ctx.Done():
				wlog.Print("break")
				return
			case <-t.C:
				wlog.Print("write messages")
				err := w.WriteMessages(
					ctx,
					kafka.Message{Topic: topic, Partition: 0, Value: []byte("v0")},
					kafka.Message{Topic: topic, Partition: 1, Value: []byte("v1")},
					kafka.Message{Topic: topic, Partition: 4, Value: []byte("v4")},
				)
				if err != nil {
					wlog.Fatal("failed to write messages:", err)
				}
			}
		}
	}()

	t.Stop()

	if err := w.Close(); err != nil {
		wlog.Fatal("failed to close reader:", err)
	}
}

func reader(ctx context.Context, topic string) {
	rlog := log.New(log.Default().Writer(), "[reader] - ", log.Default().Flags())

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			ClientID:  "go-reader",
		},
		Topic:          topic,
		GroupID:        "UnnecessaryID",
		Logger:         rlog,
		ErrorLogger:    rlog,
		IsolationLevel: kafka.ReadCommitted,
	})
	defer r.Close()

	rlog.Print("reading...")
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			rlog.Fatal("failed to read message", err)
			break
		}
		rlog.Printf("message at offset Topic=(%s) Partition=%d Offset=%d Key=(%s) Value=(%s)",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		rlog.Fatal("failed to close reader:", err)
	}
}

func readAndWrite(topic string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go reader(ctx, topic)
	time.Sleep(3 * time.Second)
	go writer(ctx, topic)

	log.Printf("working...")
	for i := 0; i < 10; i++ {
		log.Printf("tik %d", i)
		time.Sleep(1 * time.Second)
	}

	log.Printf("cancel...")
	cancel()
	time.Sleep(1 * time.Second)
	log.Printf("terminated")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	listTopics()
	// deleteTopic("t1")
	// createTopic("t1")
	// readAndWrite("t1")
}
