package service

import (
	"encoding/json"
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

func (r *Reader) ReadAll() (map[string]map[string]string, error) {

	defer r.file.Close()

	shortURLs := make(map[string]map[string]string)
	var err error

	for r.decoder.More() {
		record := &Record{}
		if err = r.decoder.Decode(&record); err != nil {
			return nil, err
		}

		if userData, ok := shortURLs[record.UserID]; ok {
			userData[record.Key] = record.Value
		} else {
			shortURLs[record.UserID] = make(map[string]string)
			shortURLs[record.UserID][record.Key] = record.Value
		}
	}

	return shortURLs, err
}
