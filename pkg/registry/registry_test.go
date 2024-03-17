package registry

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/eda/events"
	"testing"
)

func Test_registry_Register(t *testing.T) {
	assert := assert.New(t)

	registry := New()

	// check GetType
	eventType := "UserCreationRequested"
	getType := GetType[events.UserCreationRequested]()
	assert.Equal(eventType, getType)

	// register the events.UserCreationRequested event
	err := registry.Register(RegistryEntry{
		EventType:    eventType,
		EventFactory: func() any { return new(events.UserCreationRequested) },
	})
	assert.NoError(err)

	// get the events.UserCreationRequested event instance
	ucrIntf, err := registry.New(eventType)
	// cast the interface to events.UserCreationRequested
	ucr, ok := ucrIntf.(*events.UserCreationRequested)
	assert.True(ok)
	assert.NoError(err)
	assert.NotNil(ucr)

	ucr.FirstName = "John"
	ucr.LastName = "Doe"
	ucr.Password.Password = "secret"

	ucrJson, err := json.Marshal(ucr)
	assert.NoError(err)
	assert.NotNil(ucrJson)

	// unmarshal the events.UserCreationRequested event using the UnmarshalEvent helper function
	ucr2, err := UnmarshalEvent[events.UserCreationRequested](registry, eventType, ucrJson)
	assert.NoError(err)
	assert.Equal(ucr, ucr2)

	var ucd events.UserCreationDone
	// check if the events.UserCreationDone event is not yet registered
	registered := registry.IsRegistered(&ucd)
	assert.False(registered)

	// register the events.UserCreationDone event using the Register helper function
	err = Register[events.UserCreationDone](registry)
	assert.NoError(err)

	// verify that the events.UserCreationDone event is registered
	registered = registry.IsRegistered(&ucd)
	assert.True(registered)

}
