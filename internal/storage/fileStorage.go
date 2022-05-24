package storage

import (
	"github.com/AyratB/go-short-url/internal/service"
)

type FileStorage struct {
	writer        *service.Writer
	reader        *service.Reader
	shortUserURLs map[string]map[string]string
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

func (f *FileStorage) GetAll() map[string]map[string]string {
	return f.shortUserURLs
}

func (f *FileStorage) GetByOriginalURLForUser(originalURL, userID string) (string, error) {
	urls := f.GetAll()

	if usersURLs, ok := urls[userID]; ok {
		return usersURLs[originalURL], nil
	}

	return "", nil
}

func (f *FileStorage) Set(originalURL, shortenURL, userID string) error {

	if userURLs, ok := f.shortUserURLs[userID]; ok {
		userURLs[originalURL] = shortenURL
	} else {
		f.shortUserURLs[userID] = make(map[string]string)
		f.shortUserURLs[userID][originalURL] = shortenURL
	}

	r := &service.Record{
		OriginalURL: originalURL,
		ShortenURL:  shortenURL,
		UserID:      userID,
	}
	if err := f.writer.Write(r); err != nil {
		return err
	}
	return nil
}
