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
	assert.NotNil(err)

	password.Password = "password"
	err = password.Validate()
	assert.Nil(err)

	jsonStr, err := json.Marshal(password)
	assert.Nil(err)
	assert.Equal(`{"password":"password"}`, string(jsonStr))

}
