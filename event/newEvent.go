package event

import "github.com/telemac/eda"

type Event[T eda.Validater] struct {
	data *T
}

func NewEvent[T eda.Validater]() *Event[T] {
	var e Event[T]
	e.data = e.Factory().(*T)
	return &e
}

func (e *Event[T]) Type() string {
	return GetTypeName(e)
}

func (e *Event[T]) Factory() interface{} {
	return new(T)
}

func (e *Event[T]) PublishTopic() string {
	return GetTypeNameCamelCase(e)
}

func (e *Event[T]) SubscribeTopic() string {
	return e.PublishTopic()
}

func (e *Event[T]) Data() *T {
	return e.data
}

// Validate
func (e *Event[T]) Validate() error {
	return eda.ValidateAll(e.data)
}
