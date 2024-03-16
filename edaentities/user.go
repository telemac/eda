package edaentities

import (
	"errors"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (u User) Validate() error {
	if u.FirstName == "" {
		return errors.New("first_name is required")
	}
	if u.LastName == "" {
		return errors.New("last_name is required")
	}
	return nil
}
