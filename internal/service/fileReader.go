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

func (c *Reader) Read() (*Record, error) {
	record := &Record{}
	if err := c.decoder.Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}

func (c *Reader) Close() error {
	return c.file.Close()
}
