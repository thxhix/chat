package uuid

import googleUuid "github.com/google/uuid"

func NewUUID() (googleUuid.UUID, error) {
	return googleUuid.NewV7()
}
