package service

import (
	"github.com/sandreev87/golang-rest-api/internal/models"
	"github.com/sandreev87/golang-rest-api/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	SignIn(username, password string) (map[string]interface{}, error)
	GenerateTokens(user models.User) (map[string]interface{}, error)
	ParseToken(token string) (int, error)
	UpdateRefreshToken(token string) (map[string]interface{}, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, repos.Cache),
	}
}
