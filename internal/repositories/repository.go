package repositories

type Repository interface {
	GetAll() map[string]map[string]string
	GetByOriginalURLForUser(originalURL, userID string) (string, error)
	Set(originalURL, shortenURL, userID string) error
	CloseResources() error
	PingStorage() error
}
