package uuid

import googleUuid "github.com/google/uuid"

type IUUIDManager interface {
	NewUUIDv7() (googleUuid.UUID, error)
	NewUUIDv5(key string) googleUuid.UUID
}

type UUIDManager struct {
	V5Salt googleUuid.UUID
}

func NewUUIDManager(s string) *UUIDManager {
	appNamespace := googleUuid.NewSHA1(googleUuid.NameSpaceDNS, []byte(s))

	return &UUIDManager{
		V5Salt: appNamespace,
	}
}

func (m *UUIDManager) NewUUIDv7() (googleUuid.UUID, error) {
	return googleUuid.NewV7()
}

func (m *UUIDManager) NewUUIDv5(key string) googleUuid.UUID {
	return googleUuid.NewSHA1(m.V5Salt, []byte(key))
}
