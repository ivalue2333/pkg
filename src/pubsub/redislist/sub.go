package redislist

import (
	"context"
	"encoding/json"
	"github.com/ivalue2333/pkg/src/logx"
	"github.com/ivalue2333/pkg/src/pubsub"
	"github.com/pkg/errors"
	"sync"

	"github.com/gomodule/redigo/redis"
)

type sub struct {
	pool  Pool
	topic string

	quitChan chan uint8
	once     sync.Once

	pullTimeout int
}

func NewSubscriber(pool Pool, topic string, opts ...SubOption) pubsub.Subscriber {
	s := &sub{
		pool:        pool,
		topic:       topic,
		quitChan:    make(chan uint8),
		pullTimeout: 3,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *sub) Close() error {
	s.once.Do(func() {
		close(s.quitChan)
	})

	return nil
}

func (s *sub) Start(ctx context.Context, handler pubsub.MessageHandler) (err error) {
	for {
		select {
		case <-s.quitChan:
			return
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := s.Pull(s.topic, handler)
			if err != nil {
				logx.Warn(ctx, "process message failed.", err)
			}
		}
	}
}

func (s *sub) omitErr(err error) bool {
	return err == redis.ErrNil
}

func (s *sub) Pull(topic string, handler pubsub.MessageHandler) error {

	data, err := s.pullOnce(topic)

	if err != nil && !s.omitErr(err) {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	msg := pubsub.Message{}
	if err = json.Unmarshal(data, &msg); err != nil {
		return errors.Wrapf(err, "unmarshal error: data(%s)", data)
	}
	ctx := context.Background()

	return handler(ctx, msg)
}

func (s *sub) pullOnce(topic string) ([]byte, error) {
	conn := s.pool.Get()
	defer conn.Close()
	arr, err := redis.Values(conn.Do("BRPOP", topic, s.pullTimeout))
	if err != nil {
		return nil, err
	}
	if len(arr) != 2 {
		return nil, errors.New("bad values")
	}
	data, ok := arr[1].([]byte)
	if !ok {
		return nil, errors.New("bad values")
	}
	return data, nil
}
