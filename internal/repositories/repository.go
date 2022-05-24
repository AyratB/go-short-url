package repositories

type Repository interface {
	GetAll(userID string) (map[string]string, error)
	GetByKey(key, userID string) (string, error)
	Set(key, value, userID string) error
	CloseResources() error
}
