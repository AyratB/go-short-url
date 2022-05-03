package repositories

type Repository interface {
	GetAll() map[string]string
	GetByKey(key string) (string, bool)
	Set(key string, value string)
}
