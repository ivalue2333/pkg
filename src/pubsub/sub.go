package pubsub

import (
	"context"
)

type MessageHandler func(ctx context.Context, message Message) error

type Subscriber interface {
	//Close stop receive message, and close any connection to the broker
	Close() error
	//Start start a loop to receive message from broker, and dispath to handler
	Start(ctx context.Context, handler MessageHandler) error
}
