package storage

type MemoryStorage struct {
	shortUserURLs map[string]map[string]string // user ID  : Original : shortURL
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortUserURLs: make(map[string]map[string]string),
	}
}

func (ms *MemoryStorage) GetAll() map[string]map[string]string {
	return ms.shortUserURLs
}

func (ms *MemoryStorage) GetByOriginalURLForUser(originalURL, userID string) (string, error) {
	urls := ms.GetAll()

	if usersURLs, ok := urls[userID]; ok {
		return usersURLs[originalURL], nil
	}

	return "", nil
}

func (ms *MemoryStorage) Set(originalURL, shortenURL, userID string) error {

	if userURLs, ok := ms.shortUserURLs[userID]; ok {
		userURLs[originalURL] = shortenURL
	} else {
		ms.shortUserURLs[userID] = make(map[string]string)
		ms.shortUserURLs[userID][originalURL] = shortenURL
	}
	return nil
}

func (ms *MemoryStorage) CloseResources() error {
	return nil
}
