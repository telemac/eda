package event

import "github.com/telemac/eda"

const EDATypeHeader = "EDA-Type"

type Eventer interface {
	eda.Validater
	Type() string         // event type
	PublishTopic() string // topic on which the event is published
	SubscribeTopic() string
}
