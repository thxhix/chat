package security

import (
	"github.com/thxhix/chat/internal/config"
	"github.com/thxhix/chat/internal/security/crypt"
	"github.com/thxhix/chat/internal/security/jwt"
	"github.com/thxhix/chat/internal/security/password"
)

type Security struct {
	JWT      jwt.IJWTManager
	Password password.IPasswordManager
	Crypt    crypt.ICryptManager
}

func NewSecurity(cfg *config.Config) *Security {
	jwt_manager := jwt.NewJWTManager(cfg.JWT)
	password_hasher := password.NewPasswordManager()
	crypt_manager := crypt.NewCryptManager()

	return &Security{
		JWT:      jwt_manager,
		Password: password_hasher,
		Crypt:    crypt_manager,
	}
}
