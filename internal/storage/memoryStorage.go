package storage

import "github.com/AyratB/go-short-url/internal/entities"

type MemoryStorage struct {
	shortUserURLs map[string]map[string]*entities.URLInfo // userGUID  : Original : shortURL
}

func (ms *MemoryStorage) DeleteURLS(batches []string, userID string) {
	//TODO implement me
	panic("implement me")
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortUserURLs: make(map[string]map[string]*entities.URLInfo),
	}
}

func (ms *MemoryStorage) GetAll() (map[string]map[string]*entities.URLInfo, error) {
	return ms.shortUserURLs, nil
}

func (ms *MemoryStorage) GetByOriginalURLForUser(originalURL, userGUID string) (string, error) {
	urls, _ := ms.GetAll()

	if usersURLs, ok := urls[userGUID]; ok {
		return usersURLs[originalURL].ShortenURL, nil
	}

	return "", nil
}

func (ms *MemoryStorage) Set(originalURL, shortenURL, userGUID string) error {

	if _, ok := ms.shortUserURLs[userGUID]; !ok {
		ms.shortUserURLs[userGUID] = make(map[string]*entities.URLInfo)
	}
	ms.shortUserURLs[userGUID][originalURL] = &entities.URLInfo{
		ShortenURL: shortenURL,
		IsDeleted:  false,
	}
	return nil
}

func (ms *MemoryStorage) CloseResources() error {
	return nil
}

func (ms *MemoryStorage) PingStorage() error {
	return nil
}
