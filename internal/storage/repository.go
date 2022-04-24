package storage

type Repository interface {
	GetDB() interface{}
}

var shortURLs = map[string]string{
	"https://yatest.ru": "test",
}

type Storage struct{}

func (s *Storage) GetDB() interface{} {
	return shortURLs
}
