package edaentities

import (
	"errors"
)

type String string

type Password struct {
	Password string `json:"password"`
}

func (p Password) String() string {
	return p.Password
}

func (p Password) Validate() error {
	if p.Password == "" {
		return errors.New("Password is required")
	}
	return nil
}
