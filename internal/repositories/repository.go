package repositories

type Repository interface {
	GetAll() (map[string]map[string]string, error)
	GetByOriginalURLForUser(originalURL, userGUID string) (string, error)
	Set(originalURL, shortenURL, userGUID string) error
	CloseResources() error
	PingStorage() error
}
