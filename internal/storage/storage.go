package storage

var shortURLs = map[string]string{
	"https://yatest.ru": "test",
}

type Storage struct{}

func (s *Storage) GetAll() map[string]string {
	return shortURLs
}

func (s *Storage) GetByKey(key string) (string, bool) {

	shortURL, ok := shortURLs[key]
	if !ok {
		return "", false
	}
	return shortURL, true
}

func (s *Storage) Set(key string, value string) {
	shortURLs[key] = value
}
