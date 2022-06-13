package storage

import (
	"github.com/AyratB/go-short-url/internal/entities"
	"github.com/AyratB/go-short-url/internal/service"
)

type FileStorage struct {
	writer        *service.Writer
	reader        *service.Reader
	shortUserURLs map[string]map[string]*entities.URLInfo
}

func (f *FileStorage) DeleteURLS(batches []string, userID string) {
	//TODO implement me
	panic("implement me")
}

func NewFileStorage(filePath string) (*FileStorage, error) {

	w, err := service.NewWriter(filePath)
	if err != nil {
		return nil, err
	}

	r, err := service.NewReader(filePath)
	if err != nil {
		return nil, err
	}

	// читаем один раз, потом работаем в памяти
	shortUserURLs, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		writer:        w,
		reader:        r,
		shortUserURLs: shortUserURLs,
	}, nil
}

func (f *FileStorage) CloseResources() error {
	return f.writer.Close()
}

func (f *FileStorage) PingStorage() error {
	return nil
}

func (f *FileStorage) GetAll() (map[string]map[string]*entities.URLInfo, error) {
	return f.shortUserURLs, nil
}

func (f *FileStorage) GetByOriginalURLForUser(originalURL, userGUID string) (string, error) {
	urls, _ := f.GetAll()

	if usersURLs, ok := urls[userGUID]; ok {
		return usersURLs[originalURL].ShortenURL, nil
	}

	return "", nil
}

func (f *FileStorage) Set(originalURL, shortenURL, userGUID string) error {

	if _, ok := f.shortUserURLs[userGUID]; !ok {
		f.shortUserURLs[userGUID] = make(map[string]*entities.URLInfo)
	}
	f.shortUserURLs[userGUID][originalURL] = &entities.URLInfo{
		ShortenURL: shortenURL,
		IsDeleted:  false,
	}

	r := &service.Record{
		OriginalURL: originalURL,
		ShortenURL:  shortenURL,
		UserID:      userGUID,
	}
	if err := f.writer.Write(r); err != nil {
		return err
	}
	return nil
}
