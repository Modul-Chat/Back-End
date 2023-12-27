package utils

import "github.com/google/uuid"

func GenerateUUIDs() uuid.UUID{
	newUUID:= uuid.New()

	return newUUID
}

