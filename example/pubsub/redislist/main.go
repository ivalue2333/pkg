package main

import (
	"context"
	"fmt"
	"github.com/ivalue2333/pkg/src/pubsub"
	"github.com/ivalue2333/pkg/src/pubsub/redislist"
	"github.com/ivalue2333/pkg/src/storex/redisx"
	"os"
	"time"
)

var (
	topic = "topic_goods"
	pool  = redisx.CreatePool(os.Getenv("REDIS_URI"))
	pubc  = redislist.NewPublisher(pool, topic)
	subc = redislist.NewSubscriber(pool, topic)
	ctx = context.Background()
)

func init() {
}

func main()  {
	datas1 := []string{"msg1", "msg2", "msg3", "msg4"}

	go func() {
		err := subc.Start(ctx, func(ctx context.Context, message pubsub.Message) error {
			fmt.Println("message", message)
			fmt.Println("payload", string(message.Payload))
			return nil
		})
		fmt.Println("return from circle")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(3 * time.Second)

	for _, data := range datas1 {
		err := pubc.Send(ctx, pubsub.NewMessage(ctx, data, []byte(data)))
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(10 * time.Second)
}
