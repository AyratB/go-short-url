package storage

//var shortURLs = map[string]string{
//	"https://yatest.ru": "test",
//}

import (
	"github.com/AyratB/go-short-url/internal/service"
)

type FileStorage struct {
	writer *service.Writer
	reader *service.Reader
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

	return &FileStorage{
		writer: w,
		reader: r,
	}, nil
}

func (f *FileStorage) CloseResources() error {
	err := f.writer.Close()
	if err != nil {
		return err
	}
	return f.reader.Close()
}

func (f *FileStorage) GetAll() (map[string]string, error) {
	return nil, nil
}

func (f *FileStorage) GetByKey(key string) (string, error) {
	return "", nil
}

func (f *FileStorage) Set(key string, value string) error {
	r := &service.Record{Key: key, Value: value}
	if err := f.writer.Write(r); err != nil {
		return err
	}
	return nil
}
