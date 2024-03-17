package events

import (
	"github.com/telemac/eda/edaentities"
)

type UserCreationRequested struct {
	edaentities.User
	edaentities.Password
}

// Type implements msh.Eventer
func (ucr UserCreationRequested) Type() string {
	return "UserCreationRequested"
}

func (ucr UserCreationRequested) Factory() any { return new(UserCreationRequested) }

// Topic implements msh.Eventer
func (ucr UserCreationRequested) PublishTopic() string {
	return "user.creation.requested." + ucr.LastName
}

// SubscribeTopic implements msh.Eventer
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
	return "UserCreationDone"
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
