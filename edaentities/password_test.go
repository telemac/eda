package edaentities

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPassword(t *testing.T) {
	assert := assert.New(t)

	password := Password{"Alexandre"}
	err := password.Validate()
	assert.NoError(err)

	password.Password = "password"
	err = password.Validate()
	assert.NoError(err)

	jsonStr, err := json.Marshal(password)
	assert.Nil(err)
	assert.Equal(`{"password":"password"}`, string(jsonStr))

}
