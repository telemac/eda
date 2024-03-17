package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telemac/eda/event"
	"sync"
)

var (
	ErrEventLareadyRegistered = errors.New("event already registered")
	ErrEventNotFound          = errors.New("event not found")
)

type EventFactory func() any

type RegistryEntry struct {
	EventType    string
	EventFactory EventFactory
}

type Registry[T event.Eventer] interface {
	// Register a new event
	Register(entry RegistryEntry) error
	New(eventType string) (any, error)
	IsRegistered(event event.Eventer) bool
}

/* implementation */

type registry[T event.Eventer] struct {
	eventFactoryMap map[string]RegistryEntry
	mutex           sync.RWMutex
}

// New creates a new registry
func New() Registry[event.Eventer] {
	return &registry[event.Eventer]{
		eventFactoryMap: make(map[string]RegistryEntry),
	}
}

func (r *registry[T]) Register(entry RegistryEntry) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, ok := r.eventFactoryMap[entry.EventType]
	if ok {
		return ErrEventLareadyRegistered
	}
	r.eventFactoryMap[entry.EventType] = entry
	return nil
}

func Register[T event.Eventer](r Registry[event.Eventer]) error {
	var ev T
	return r.Register(RegistryEntry{
		EventType:    ev.Type(),
		EventFactory: func() any { return new(T) },
	})
}

func (r *registry[T]) New(eventType string) (any, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entry, ok := r.eventFactoryMap[eventType]
	if !ok {
		return nil, fmt.Errorf("event %s : %w", eventType, ErrEventNotFound)
	}
	event := entry.EventFactory()
	return event, nil
}

func GetType[T event.Eventer]() string {
	var e T
	eventType := e.Type()
	return eventType
}

func UnmarshalEvent[T event.Eventer](r Registry[event.Eventer], eventTipe string, data []byte) (*T, error) {
	ev, err := r.New(eventTipe)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalEvent new event %s : %w", eventTipe, ErrEventNotFound)
	}
	err = json.Unmarshal(data, &ev)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalEvent unmarshal event %s : %w", eventTipe, err)
	}
	ucr, ok := ev.(*T)
	if !ok {
		return nil, fmt.Errorf("UnmarshalEvent cast event %s : %w", eventTipe, err)
	}
	return ucr, nil
}

func (r *registry[T]) IsRegistered(event event.Eventer) bool {
	var eventType = event.Type()
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	_, ok := r.eventFactoryMap[eventType]
	return ok
}
