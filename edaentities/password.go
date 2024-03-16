package edaentities

import (
	"errors"
)

type Password struct {
	Password string `json:"password"`
}

func (p Password) Validate() error {
	if p.Password == "" {
		return errors.New("Password is required")
	}
	return nil
}
