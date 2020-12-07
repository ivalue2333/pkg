package pubsub

import (
	"context"
	"time"
)

func NewMessage(ctx context.Context, key string, payload []byte) Message {
	return Message{
		Key:     key,
		Payload: payload,
	}
}

type Message struct {
	// Set the sequence id to assign to the current message
	SequenceId int64
	// Sets the key of the message for routing policy
	Key string
	// Payload for the message
	Payload []byte
	// Set the event time for a given message
	EventTime time.Time
	// traceId
	TraceId string
}
