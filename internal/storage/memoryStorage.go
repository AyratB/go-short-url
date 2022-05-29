package storage

type MemoryStorage struct {
	shortUserURLs map[string]map[string]string // userGUID  : Original : shortURL
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortUserURLs: make(map[string]map[string]string),
	}
}

func (ms *MemoryStorage) GetAll() (map[string]map[string]string, error) {
	return ms.shortUserURLs, nil
}

func (ms *MemoryStorage) GetByOriginalURLForUser(originalURL, userGUID string) (string, error) {
	urls, _ := ms.GetAll()

	if usersURLs, ok := urls[userGUID]; ok {
		return usersURLs[originalURL], nil
	}

	return "", nil
}

func (ms *MemoryStorage) Set(originalURL, shortenURL, userGUID string) error {

	if userURLs, ok := ms.shortUserURLs[userGUID]; ok {
		userURLs[originalURL] = shortenURL
	} else {
		ms.shortUserURLs[userGUID] = make(map[string]string)
		ms.shortUserURLs[userGUID][originalURL] = shortenURL
	}
	return nil
}

func (ms *MemoryStorage) CloseResources() error {
	return nil
}

func (ms *MemoryStorage) PingStorage() error {
	return nil
}
