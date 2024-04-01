package event

import (
	"github.com/stretchr/testify/assert"
	"github.com/telemac/eda/edaentities"
	"testing"
)

type UserDeleted struct {
	edaentities.User
	edaentities.Uuid
	Err error
}

func (ud UserDeleted) Validate() error {
	if err := ud.User.Validate(); err != nil {
		return err
	}
	if err := ud.Uuid.Validate(); err != nil {
		return err
	}
	return nil
}

func TestNewEvent(t *testing.T) {
	assert := assert.New(t)
	userDeletedEvent := NewEvent[UserDeleted]()

	userDeleted := userDeletedEvent.Data()
	userDeleted.User.FirstName = "Alexandre"
	userDeleted.User.LastName = "HEIM"
	userDeleted.Uuid.UUID = "123456789012345678901234567890123456"

	err := userDeletedEvent.Validate()
	assert.NoError(err)

	assert.NotNil(userDeletedEvent)
	assert.Equal("UserDeleted", userDeletedEvent.Type())
	assert.Equal("user.deleted", userDeletedEvent.PublishTopic())
}
