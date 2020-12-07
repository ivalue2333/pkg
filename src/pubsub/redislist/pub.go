package redislist

import (
	"context"
	"encoding/json"
	"github.com/ivalue2333/pkg/src/pubsub"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type pub struct {
	pool  Pool
	topic string
	Id    int64
}

func NewPublisher(pool Pool, topic string) pubsub.Publisher {
	return &pub{
		pool:  pool,
		topic: topic,
		Id:    0,
	}
}

func (p *pub) Close() {

}

func (p *pub) Send(ctx context.Context, msg pubsub.Message) (err error) {
	return p.send(ctx, msg)
}

// SendAsync fake async
func (p *pub) SendAsync(ctx context.Context, msg pubsub.Message, callback func(pubsub.Message, error)) {
	err := p.send(ctx, msg)

	if callback != nil {
		callback(msg, err)
	}
}

func (p *pub) send(ctx context.Context, msg pubsub.Message) (err error) {

	if msg.EventTime.IsZero() {
		msg.EventTime = time.Now()
	}

	if msg.SequenceId == 0 {
		atomic.AddInt64(&p.Id, 1)
		msg.SequenceId = p.Id
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrapf(err, "Marshal failed: %+v", msg)
	}

	conn := p.pool.Get()
	defer conn.Close()

	_, err = conn.Do("LPUSH", p.topic, data)

	return err
}
