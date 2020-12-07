package main

import (
	"context"
	"fmt"
	"github.com/ivalue2333/pkg/src/jsonx"
	"github.com/ivalue2333/pkg/src/pubsub"
	"github.com/ivalue2333/pkg/src/pubsub/redisps"
	"github.com/ivalue2333/pkg/src/storex/redisx"
	"os"
	"time"
)

var (
	topics = []string{"topic_goos", "topic_user"}
	pool   = redisx.CreatePool(os.Getenv("REDIS_URI"))
	pubc   = []pubsub.Publisher{}

	subc = redisps.NewSubscriber(pool, topics...)

	ctx = context.Background()
)

func init() {
	for _, topic := range topics {
		pubc = append(pubc, redisps.NewPublisher(pool, topic))
	}
}

func main() {
	datas1 := []string{"msg1", "msg2", "msg3", "msg4"}
	datas2 := []struct {
		Name string
		Age  int
	}{
		{Name: "name1", Age: 12},
		{Name: "name2", Age: 15},
	}

	var err error

	go func() {
		err = subc.Start(ctx, func(ctx context.Context, message pubsub.Message) error {
			msg := jsonx.MarshalUnsafe(message)
			fmt.Println("msg:" + string(msg))
			fmt.Println(string(message.Payload))
			return nil
		})
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	for _, data := range datas1 {
		err = pubc[0].Send(ctx, pubsub.NewMessage(ctx, data, []byte(data)))
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(1 * time.Second)

	for _, d := range datas2 {
		data := jsonx.MarshalUnsafe(d)
		err = pubc[1].Send(ctx, pubsub.NewMessage(ctx, string(data), data))
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(2 * time.Second)

}
