package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/token"
	"github.com/thxhix/chat/internal/domain/user"
	"github.com/thxhix/chat/internal/security/crypt"
	"github.com/thxhix/chat/internal/security/jwt"
	"github.com/thxhix/chat/internal/security/password"
	"strconv"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, login string, password string) (userId int64, accessToken string, refreshToken string, err error)
	Login(ctx context.Context, login string, password string) (userId int64, accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, incomingRefreshToken string) (userId int64, accessToken string, refreshToken string, err error)
}

type AuthService struct {
	userRepo  user.IUserRepository
	tokenRepo token.ITokenRepository

	passwordManager password.IPasswordManager
	tokenManager    jwt.IJWTManager
	cryptManager    crypt.ICryptManager
}

// NewService constructs a new AuthService with given dependencies.
func NewService(ur user.IUserRepository, tr token.ITokenRepository, pm password.IPasswordManager, tm jwt.IJWTManager, cm crypt.ICryptManager) IAuthService {
	return &AuthService{
		userRepo:  ur,
		tokenRepo: tr,

		passwordManager: pm,
		tokenManager:    tm,
		cryptManager:    cm,
	}
}

func (s *AuthService) Register(ctx context.Context, login string, password string) (userId int64, accessToken string, refreshToken string, err error) {
	if err := user.ValidateLogin(login); err != nil {
		return 0, "", "", err
	}

	if err := user.ValidatePassword(password); err != nil {
		return 0, "", "", err
	}

	// Generate password hash for security store in storage
	passwordHash, err := s.passwordManager.HashPassword(password)
	if err != nil {
		return 0, "", "", err
	}

	userId, err = s.userRepo.Create(ctx, login, passwordHash)
	if err != nil {
		return 0, "", "", err
	}

	accessToken, err = s.tokenManager.GenerateAccessToken(userId)
	if err != nil {
		return 0, "", "", err
	}

	refreshToken, refreshJTI, refreshTTL, err := s.tokenManager.GenerateRefreshToken(userId)
	if err != nil {
		return 0, "", "", err
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(refreshTTL)
	refreshHash := s.cryptManager.Sha256Hex(refreshToken)

	if err := s.tokenRepo.Create(ctx, userId, refreshJTI, refreshHash, issuedAt, expiresAt); err != nil {
		return 0, "", "", err
	}

	return userId, accessToken, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, login string, password string) (userId int64, accessToken string, refreshToken string, err error) {
	if err := user.ValidateLogin(login); err != nil {
		return 0, "", "", err
	}

	au, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return 0, "", "", err
	}

	if !s.passwordManager.CheckPasswordHash(password, au.Password) {
		return 0, "", "", user.ErrUserNotFound
	}

	userId = au.ID

	accessToken, err = s.tokenManager.GenerateAccessToken(userId)
	if err != nil {
		return 0, "", "", err
	}

	refreshToken, refreshJTI, refreshTTL, err := s.tokenManager.GenerateRefreshToken(userId)
	if err != nil {
		return 0, "", "", err
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(refreshTTL)
	refreshHash := s.cryptManager.Sha256Hex(refreshToken)

	if err := s.tokenRepo.Create(ctx, userId, refreshJTI, refreshHash, issuedAt, expiresAt); err != nil {
		return 0, "", "", err
	}

	return userId, accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, incomingRefreshToken string) (userId int64, accessToken string, refreshToken string, err error) {
	userIdStr, jtiStr, err := s.tokenManager.ParseRefreshToken(incomingRefreshToken)
	if err != nil {
		return 0, "", "", err
	}

	userId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, "", "", err
	}
	oldJTI, err := uuid.Parse(jtiStr)
	if err != nil {
		return 0, "", "", err
	}

	incomingTokenHash := s.cryptManager.Sha256Hex(incomingRefreshToken)

	tokenRecord, err := s.tokenRepo.GetByJTI(ctx, oldJTI)
	if err != nil {
		return 0, "", "", err
	}

	now := time.Now().UTC()

	if err := token.ValidateToken(now, userId, incomingTokenHash, tokenRecord); err != nil {
		return 0, "", "", err
	}

	accessToken, err = s.tokenManager.GenerateAccessToken(userId)
	if err != nil {
		return 0, "", "", err
	}
	newRefreshToken, newJTI, newTTL, err := s.tokenManager.GenerateRefreshToken(userId)
	if err != nil {
		return 0, "", "", err
	}

	newHash := s.cryptManager.Sha256Hex(newRefreshToken)

	if err := s.tokenRepo.Rotate(ctx, userId, oldJTI, newJTI, newHash, now, now.Add(newTTL)); err != nil {
		return 0, "", "", err
	}

	return userId, accessToken, newRefreshToken, nil
}
