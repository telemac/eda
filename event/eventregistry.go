package event

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var (
	ErrEventNotRegistered = errors.New("event type not registered")
)

type EventRegistryIntf interface {
	Register(event Eventer)
	New(eventType string) (Eventer, error)
}

func Factory[T Eventer]() *T {
	return new(T)
}

type EventerFactory func() Eventer

type EventRegistry struct {
	events map[string]EventerFactory
	mutex  sync.RWMutex
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		events: make(map[string]EventerFactory),
	}
}

func (er *EventRegistry) Register(event Eventer) {
	er.mutex.Lock()
	defer er.mutex.Unlock()
	isPtr := reflect.ValueOf(event).Kind() == reflect.Ptr
	if !isPtr {
		//panic("event must be a pointer")
	}
	er.events[event.Type()] = func() Eventer {
		isPtr := reflect.ValueOf(event).Kind() == reflect.Ptr
		if isPtr {
			return event
		} else {
			// TODO : find a better way to do this
			value := reflect.ValueOf(event)
			iface := value.Interface()
			ptr := reflect.NewAt(value.Type(), unsafe.Pointer(&iface))
			return ptr.Interface().(Eventer)
		}
	}
}

func (er *EventRegistry) New(eventType string) (Eventer, error) {
	er.mutex.RLock()
	defer er.mutex.RUnlock()
	factory, ok := er.events[eventType]
	if !ok {
		return nil, fmt.Errorf("%w : %s", ErrEventNotRegistered, eventType)
	}
	return factory(), nil
}
