package main

// https://github.com/redis/go-redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func connect(ctx context.Context) *redis.Client {
	log.Print("connecting")
	rdb := redis.NewClient(&redis.Options{
		// Addr: "redis:6379",
		Addr:        "host.docker.internal:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		DialTimeout: 1 * time.Second,
		MaxRetries:  1,
	})

	pong, err := rdb.Echo(ctx, "ping").Result()
	if err == redis.Nil {
		log.Println("key2 does not exist")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Println("Echo: ", pong)
	}
	log.Print("connected")
	return rdb
}

func playSimple(ctx context.Context, rdb *redis.Client) {
	log.Println("set value")
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		log.Panic(err)
	}

	log.Print("get value")
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("key %s", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		log.Print("key2 does not exist")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("key2 -> %s", val2)
	}
}

func playQueue(ctx context.Context, rdb *redis.Client) {
	pushed, err := rdb.LPush(ctx, "queue1", 1, 2, 3).Result()
	if err == redis.Nil {
		log.Print("key2 does not exist")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("LPush: %+v", pushed)
	}

	for {
		poped, err := rdb.BRPop(ctx, 1*time.Second, "queue1").Result()
		if err == redis.Nil {
			log.Print("queue1 does not exist")
			break
		} else if err != nil {
			log.Panic(err)
			break
		} else {
			log.Printf("BRPop: %+v", poped)
		}
	}
}

func playStreams(ctx context.Context, rdb *redis.Client) {
	log.Printf("XAdd")
	argXAdd := redis.XAddArgs{
		Stream: "stream1",
		Values: []string{"key1", "value1", "key2", "value2"},
	}
	resultXAdd, err := rdb.XAdd(ctx, &argXAdd).Result()
	if err == redis.Nil {
		log.Print("XAdd failed")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("XAdd: %+v", resultXAdd)
	}

	// -----
	log.Printf("XGroupDestroy")
	resultXGroupDestroy, err := rdb.XGroupDestroy(ctx, "stream1", "group1").Result()
	if err == redis.Nil {
		log.Print("XGroupDestroy failed")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("XGroupDestroy: %+v", resultXGroupDestroy)
	}

	log.Printf("XGroupCreate")
	resultXGroupCreate, err := rdb.XGroupCreateMkStream(ctx, "stream1", "group1", "0").Result()
	if err == redis.Nil {
		log.Print("XGroupCreate failed")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("XGroupCreate: %+v", resultXGroupCreate)
	}

	log.Printf("XGroupCreateConsumer")
	resultXGroupCreateConsumer, err := rdb.XGroupCreateConsumer(ctx, "stream1", "group1", "consumer1").Result()
	if err == redis.Nil {
		log.Print("XGroupCreateConsumer failed")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("XGroupCreateConsumer: %+v", resultXGroupCreateConsumer)
	}

	for {
		log.Printf("XReadGroup")
		argXReadGroup := redis.XReadGroupArgs{
			Group:    "group1",
			Consumer: "consumer1",
			Streams:  []string{"stream1", ">"},
			Count:    1,
			Block:    1 * time.Second,
		}
		resultXReadGroup, err := rdb.XReadGroup(ctx, &argXReadGroup).Result()
		if err == redis.Nil {
			log.Print("XReadGroup failed")
			break
		} else if err != nil {
			log.Panic(err)
			break
		} else {
			log.Printf("XReadGroup: %+v", resultXReadGroup)
		}

		log.Print("XAck")
		for _, res := range resultXReadGroup {
			var ids = []string{}
			for _, m := range res.Messages {
				log.Printf("message - %+v", m)
				ids = append(ids, m.ID)
			}
			resultXAckGroup, err := rdb.XAck(ctx, res.Stream, "group1", ids...).Result()
			if err == redis.Nil {
				log.Print("XAck failed")
			} else if err != nil {
				log.Panic(err)
				break
			} else {
				log.Printf("XAck: %+v", resultXAckGroup)
			}
		}
	}

	// -----
	log.Printf("XGroupDestroy")
	resultXGroupDestroy, err = rdb.XGroupDestroy(ctx, "stream1", "group1").Result()
	if err == redis.Nil {
		log.Print("stream1 does not exist")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("XGroupDestroy: %+v", resultXGroupDestroy)
	}

	log.Printf("Del")
	resultXDel, err := rdb.Del(ctx, "stream1").Result()
	if err == redis.Nil {
		log.Print("stream1 does not exist")
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("Del: %+v", resultXDel)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	ctx := context.Background()

	rdb := connect(ctx)
	defer rdb.Close()
	// playSimple(ctx, rdb)
	// playQueue(ctx, rdb)
	playStreams(ctx, rdb)

	log.Print("done")
}
