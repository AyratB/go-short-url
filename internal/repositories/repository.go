package repositories

import "github.com/AyratB/go-short-url/internal/entities"

type Repository interface {
	GetAll() (map[string]map[string]*entities.URLInfo, error)
	GetByOriginalURLForUser(originalURL, userGUID string) (string, error)
	Set(originalURL, shortenURL, userGUID string) error
	CloseResources() error
	PingStorage() error
}
