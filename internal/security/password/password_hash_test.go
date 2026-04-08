package password

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var pm = NewPasswordManager()

func TestHasher_HashAndCheckPassword(t *testing.T) {
	password := "mySecret123!"

	// Генерация хэша
	hash, err := pm.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Проверка правильного пароля
	match := pm.CheckPasswordHash(password, hash)
	require.True(t, match, "expected password to match hash")

	// Проверка неправильного пароля
	wrongMatch := pm.CheckPasswordHash("wrongPassword", hash)
	require.False(t, wrongMatch, "expected wrong password not to match hash")
}

func TestHasher_HashPassword_UniqueHashes(t *testing.T) {
	password := "samePassword"

	hash1, err := pm.HashPassword(password)
	require.NoError(t, err)
	hash2, err := pm.HashPassword(password)
	require.NoError(t, err)

	// Два хэша одного и того же пароля должны быть разными из-за соли
	require.NotEqual(t, hash1, hash2)
}

func TestHasher_CheckPasswordHash_InvalidHash(t *testing.T) {
	// Некорректный хэш должен вернуть false
	match := pm.CheckPasswordHash("password", "invalidHash")
	require.False(t, match)
}
