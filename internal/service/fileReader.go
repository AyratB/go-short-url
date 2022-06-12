package service

import (
	"encoding/json"
	"github.com/AyratB/go-short-url/internal/entities"
	"os"
)

type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Reader) ReadAll() (shortURLs map[string]map[string]*entities.URLInfo, err error) {

	defer func() {
		err = r.file.Close()
	}()

	shortURLs = make(map[string]map[string]*entities.URLInfo)

	for r.decoder.More() {
		record := &Record{}
		if err = r.decoder.Decode(&record); err != nil {
			return nil, err
		}

		if _, ok := shortURLs[record.UserID]; !ok {
			shortURLs[record.UserID] = make(map[string]*entities.URLInfo)
		}

		shortURLs[record.UserID][record.OriginalURL] = &entities.URLInfo{
			ShortenURL: record.ShortenURL,
			IsDeleted:  false,
		}
	}

	return
}
