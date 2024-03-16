package broker

import "github.com/telemac/eda/event"

type Broker[T event.Eventer] interface {
	Publisher[T]
	Subscriber[T]
	Requester[T]
}
