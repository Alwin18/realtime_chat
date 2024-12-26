package utils

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

func StringToUUID(input string) (uuid.UUID, error) {
	if input == "" {
		return uuid.UUID{}, errors.New("input string is empty")
	}

	u, err := uuid.Parse(input)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid UUID format: " + err.Error())
	}

	return u, nil
}

func BytesToStruct[T any](data []byte) (T, error) {
	var result T

	if len(data) == 0 {
		return result, errors.New("input data is empty")
	}

	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
