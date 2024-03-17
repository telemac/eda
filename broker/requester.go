package broker

import (
	"context"
	"github.com/telemac/eda/event"
	"time"
)

type EventHandler func(requestEvent any) (any, error)

// TODO : implement on natsbroker
type Requester[T event.Eventer] interface {
	Request(ctx context.Context, topic string, event T, timeout time.Duration) (any, error)
	RequestEvent(ctx context.Context, event T, timeout time.Duration) (any, error)
}
