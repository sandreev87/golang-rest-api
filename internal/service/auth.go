package service

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sandreev87/golang-rest-api/internal/models"
	"github.com/sandreev87/golang-rest-api/internal/repository"
	"time"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo  repository.Authorization
	cache repository.Cache
}

func NewAuthService(repo repository.Authorization, cache repository.Cache) *AuthService {
	return &AuthService{repo: repo, cache: cache}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) SignIn(username, password string) (map[string]interface{}, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(user)
}

func (s *AuthService) UpdateRefreshToken(refreshToken string) (map[string]interface{}, error) {
	defer s.cache.Del([]byte(refreshToken))

	userBytes, err := s.cache.Get([]byte(refreshToken))
	if err != nil {
		return nil, err
	}
	var user models.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(user)
}

func (s *AuthService) GenerateTokens(user models.User) (map[string]interface{}, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	accessToken, _ := token.SignedString([]byte(signingKey))
	newRefreshToken := uuid.New()
	userBytes, _ := json.Marshal(user)

	err := s.cache.Set([]byte(newRefreshToken.String()), userBytes, 0)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token":         accessToken,
		"refresh_token": newRefreshToken.String(),
	}, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("Token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
