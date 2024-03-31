package events

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/eda/edaentities"
	"github.com/telemac/eda/event"
	"testing"
)

func TestUserCreationRequested(t *testing.T) {
	assert := assert.New(t)
	ucr := UserCreationRequested{
		User: edaentities.User{
			FirstName: "Alexandre",
			LastName:  "HEIM",
		},
		Password: edaentities.Password{
			Password: "password",
		},
	}
	jsonStr, err := json.Marshal(ucr)
	assert.NoError(err)

	var ucr2 UserCreationRequested
	err = json.Unmarshal(jsonStr, &ucr2)
	assert.NoError(err)
	assert.Equal(ucr, ucr2)

	evType := ucr2.Type()
	assert.Equal("UserCreationRequested", evType)

	userCreationDone := UserCreationDone{
		UserCreationRequested: &ucr2,
		Uuid: edaentities.Uuid{
			UUID: "2228b248-ef94-11ee-8e16-de5687fee50d",
		},
	}
	jsonStr, err = json.Marshal(userCreationDone)
	assert.NoError(err)
	assert.Equal(`{"first_name":"Alexandre","last_name":"HEIM","password":"password","uuid":"2228b248-ef94-11ee-8e16-de5687fee50d"}`, string(jsonStr))

	evType = userCreationDone.Type()
	assert.Equal("UserCreationDone", evType)
	evTypeCamelCase := event.GetTypeNameCamelCase(userCreationDone)
	assert.Equal("user.creation.done", evTypeCamelCase)
	err = userCreationDone.Validate()
	assert.NoError(err)
}
