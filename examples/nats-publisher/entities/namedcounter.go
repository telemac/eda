package entities

import "fmt"

type NamedCounter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (c NamedCounter) PublishTopic() string {
	return "test.named_counter"
}

func (c NamedCounter) SubscribeTopic() string {
	return c.PublishTopic()
}

func (c NamedCounter) Type() string {
	return fmt.Sprintf("%T", c)
}

func (c NamedCounter) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("Name is empty")
	}
	return nil
}
