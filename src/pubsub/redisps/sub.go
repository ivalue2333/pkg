package redisps

import (
	"context"
	"encoding/json"
	"github.com/ivalue2333/pkg/src/logx"
	"github.com/ivalue2333/pkg/src/pubsub"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const DelayAfterErr = 5 * time.Second

type sub struct {
	pool     Pool
	topics   []string
	quitChan chan uint8

	once sync.Once
}

func NewSubscriber(pool Pool, topics ...string) pubsub.Subscriber {
	return &sub{
		pool:     pool,
		topics:   topics,
		quitChan: make(chan uint8),
	}
}

//Close stop receive message, and close any connection to the broker
func (s *sub) Close() (err error) {
	s.once.Do(func() {
		close(s.quitChan)
	})

	return
}

//Start start a loop to receive message from broker, and dispath to handler
func (s *sub) Start(ctx context.Context, handler pubsub.MessageHandler) (err error) {
	var psc *redis.PubSubConn
	psc = &redis.PubSubConn{Conn: s.pool.Get()}

	for {
		select {
		case <-s.quitChan:
			return
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = s.listenMessage(ctx, psc, handler)
			if !strings.Contains(err.Error(), "i/o timeout") {
				time.Sleep(DelayAfterErr)
			}

			if psc != nil {
				psc.Close()
				psc = nil
			}

			psc = &redis.PubSubConn{Conn: s.pool.Get()}
		}
	}
}

func (s *sub) listenMessage(ctx context.Context, psc *redis.PubSubConn, handler pubsub.MessageHandler) (err error) {
	channels := []interface{}{}
	for _, topic := range s.topics {
		channels = append(channels, topic)
	}

	err = psc.Subscribe(channels...)
	if err != nil {
		time.Sleep(DelayAfterErr)
		return err
	}

	for {
		select {
		case <-s.quitChan:
			return
		case <-ctx.Done():
			return ctx.Err()
		default:
			switch v := psc.Receive().(type) {
			case redis.Message:
				if err := s.processMessage(v.Channel, v.Data, handler); err != nil {
					logx.Warn(ctx, "process message failed.", err)
				}
			case redis.Subscription:
			case error:
				_ = psc.Unsubscribe(channels...)

				return v
			}
		}
	}
}

func (s *sub) processMessage(topic string, data []byte, handler pubsub.MessageHandler) (err error) {
	ctx := context.Background()
	logx.Debugf(ctx, "topic:%s, data:%s", topic, data)
	msg := pubsub.Message{}
	if err = json.Unmarshal(data, &msg); err != nil {
		return err
	}
	return handler(ctx, msg)
}
