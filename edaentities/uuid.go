package edaentities

import "errors"

type Uuid struct {
	UUID string `json:"uuid"`
}

// Validate implements msh.Validator
func (u Uuid) Validate() error {
	if len(u.UUID) != 36 {
		return errors.New("uuid should have 36 characters")
	}
	return nil
}
