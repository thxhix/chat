package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

type ICryptManager interface {
	Sha256Hex(s string) string
}

type CryptManager struct {
}

func NewCryptManager() *CryptManager {
	return &CryptManager{}
}

func (cm *CryptManager) Sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
