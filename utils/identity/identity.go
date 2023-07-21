package identity

import (
	"errors"

	"github.com/google/uuid"
)

func GenerateID() (string, error) {
	id := uuid.New().String()
	if id == "" {
		return "", errors.New("error on generating uuid")
	}
	return id, nil
}
