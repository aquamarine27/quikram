package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"quikram-service/internal/config"
	"quikram-service/internal/models"
	"quikram-service/internal/repository"
	"quikram-service/internal/validator"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	tokenRepo *repository.TokenRepo
	cfg       *config.Config
}

func NewAuthService(userRepo *repository.UserRepo, tokenRepo *repository.TokenRepo, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo, cfg: cfg}
}

type AuthResult struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func (s *AuthService) Register(email, password, name string) (*AuthResult, error) {
	if err := validator.Email(email); err != nil {
		return nil, err
	}
	if err := validator.Password(password); err != nil {
		return nil, err
	}

	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Plan:         "free",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.generateTokens(user)
}

func (s *AuthService) Login(email, password string) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return s.generateTokens(user)
}

func (s *AuthService) Refresh(refreshToken string) (*AuthResult, error) {
	hash := hashToken(refreshToken)
	stored, err := s.tokenRepo.FindByHash(hash)
	if err != nil || stored == nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(stored.ExpiresAt) {
		s.tokenRepo.DeleteByHash(hash)
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.FindByID(stored.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	s.tokenRepo.DeleteByHash(hash)
	return s.generateTokens(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.tokenRepo.DeleteByHash(hash)
}

func (s *AuthService) generateTokens(user *models.User) (*AuthResult, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(s.cfg.JWTAccessExpire).Unix(),
	})

	accessStr, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	refreshRaw := uuid.New().String() + "-" + uuid.New().String()
	refreshHash := hashToken(refreshRaw)

	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshExpire),
	}
	if err := s.tokenRepo.Create(rt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessStr,
		RefreshToken: refreshRaw,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func generatePassword(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
