package entities

import (
	"fmt"
	"github.com/telemac/eda/event"
)

type NamedCounter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (c NamedCounter) Type() string {
	return event.GetTypeName(c)
}
func (c NamedCounter) PublishTopic() string {
	return fmt.Sprintf("test.%s.%d", c.Type(), c.Count)
}
func (c NamedCounter) SubscribeTopic() string {
	return fmt.Sprintf("test.%s.*", c.Type())
}
func (c NamedCounter) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("Name is empty")
	}
	return nil
}
