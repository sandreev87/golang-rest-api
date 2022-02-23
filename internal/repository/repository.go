package repository

import (
	"github.com/coocood/freecache"
	"github.com/jmoiron/sqlx"
	"github.com/sandreev87/golang-rest-api/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Cache interface {
	Get(uuid []byte) ([]byte, error)
	Set(key []byte, val []byte, expireIn int) error
	Del(key []byte) (affected bool)
}

type Repository struct {
	Authorization
	Cache
}

func NewRepository(db *sqlx.DB, cache *freecache.Cache) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Cache:         NewFreeCache(cache),
	}
}
