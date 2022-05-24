package storage

type MemoryStorage struct {
	shortUserURLs map[string]map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortUserURLs: make(map[string]map[string]string),
	}
}

func (ms *MemoryStorage) GetAll(userID string) (map[string]string, error) {
	return ms.shortUserURLs[userID], nil
}

func (ms *MemoryStorage) GetByKey(key, userID string) (string, error) {

	if ms.shortUserURLs[userID] != nil {
		return ms.shortUserURLs[userID][key], nil
	}
	return "", nil
}

func (ms *MemoryStorage) Set(key, value, userID string) error {

	if userURLs, ok := ms.shortUserURLs[userID]; ok {
		userURLs[key] = value
	} else {
		ms.shortUserURLs[userID] = make(map[string]string)
		ms.shortUserURLs[userID][key] = value
	}
	return nil
}

func (ms *MemoryStorage) CloseResources() error {
	return nil
}
