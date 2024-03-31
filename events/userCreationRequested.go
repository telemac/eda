package events

import (
	"github.com/telemac/eda/edaentities"
	"github.com/telemac/eda/event"
)

type UserCreationRequested struct {
	edaentities.User
	edaentities.Password
}

// Type implements event.Eventer
func (ucr UserCreationRequested) Type() string {
	return event.GetTypeName(ucr)
}

func (ucr UserCreationRequested) Factory() any { return new(UserCreationRequested) }

// Topic implements event.Eventer
func (ucr UserCreationRequested) PublishTopic() string {
	return "user.creation.requested." + ucr.LastName
}

// SubscribeTopic implements event.Eventer
func (ucr UserCreationRequested) SubscribeTopic() string {
	return "user.creation.requested.*"
}

// Validate implements msh.Validator
func (ucr UserCreationRequested) Validate() error {
	var err error
	err = ucr.User.Validate()
	if err != nil {
		return err
	}
	err = ucr.Password.Validate()
	if err != nil {
		return err
	}
	return err
}

type UserCreationDone struct {
	*UserCreationRequested
	edaentities.Uuid
}

func (ucd UserCreationDone) Type() string {
	return event.GetTypeName(ucd)
}

func (ucd UserCreationDone) Topic() string {
	return "user.creation.done"
}

func (ucd UserCreationDone) Validate() error {
	var err error
	err = ucd.UserCreationRequested.Validate()
	if err != nil {
		return err
	}
	err = ucd.Uuid.Validate()
	if err != nil {
		return err
	}
	return err
}
