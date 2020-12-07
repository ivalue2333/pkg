package pubsub

import (
	"context"
)

type Publisher interface {
	Close()
	Send(ctx context.Context, msg Message) (err error)
	SendAsync(ctx context.Context, msg Message, callback func(Message, error))
}
