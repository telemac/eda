package broker

import (
	"context"
	"github.com/telemac/eda/event"
)

type Publisher[T event.Eventer] interface {
	Publish(ctx context.Context, topic string, event T) error
	PublishEvent(ctx context.Context, event T) error
}
