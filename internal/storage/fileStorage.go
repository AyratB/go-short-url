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

func (f *FileStorage) GetAll(userID string) (map[string]string, error) {
	return f.shortUserURLs[userID], nil
}

func (f *FileStorage) GetByKey(key, userID string) (string, error) {
	records, err := f.GetAll(userID)
	if err != nil {
		return "", err
	}
	return records[key], nil
}

func (f *FileStorage) Set(key, value, userID string) error {

	if userURLs, ok := f.shortUserURLs[userID]; ok {
		userURLs[key] = value
	} else {
		f.shortUserURLs[userID] = make(map[string]string)
		f.shortUserURLs[userID][key] = value
	}

	r := &service.Record{
		Key:    key,
		Value:  value,
		UserID: userID,
	}
	if err := f.writer.Write(r); err != nil {
		return err
	}
	return nil
}
