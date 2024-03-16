package event

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/eda/edaentities"
	"github.com/telemac/eda/events"
	"testing"
)

func TestEventRegistry(t *testing.T) {
	assert := assert.New(t)
	eventRegistry := NewEventRegistry()
	assert.NotNil(eventRegistry)
	e := Factory[events.UserCreationRequested]()
	eventRegistry.Register(e)

	eventType := events.UserCreationRequested{}.Type()
	event, err := eventRegistry.New(eventType)
	assert.Nil(err)
	assert.NotNil(event)
	assert.Equal(eventType, event.Type())

	ucr := events.UserCreationRequested{
		User: edaentities.User{
			FirstName: "John",
			LastName:  "Doe",
		},
		Password: edaentities.Password{
			"password",
		},
	}

	err = ucr.Validate()
	assert.Nil(err)

	jsonUcr, err := json.Marshal(ucr)
	assert.Nil(err)
	assert.Equal("{\"first_name\":\"John\",\"last_name\":\"Doe\",\"password\":\"password\"}", string(jsonUcr))

	err = json.Unmarshal(jsonUcr, &event)
	assert.NoError(err)

	assert.Equal(&ucr, event)

	jsonEvent, err := json.Marshal(event)
	assert.NoError(err)
	assert.Equal(jsonUcr, jsonEvent)

}
