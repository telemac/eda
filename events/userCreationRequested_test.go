package events

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/telemac/eda/edaentities"
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
}
