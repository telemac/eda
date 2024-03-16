package broker

import (
	"context"
	"github.com/telemac/eda/event"
)

type Subscriber[T event.Eventer] interface {
	Subscribe(ctx context.Context, topic string, callback func(topic string, event T)) error
}
