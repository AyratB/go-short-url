package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	DB *sql.DB

	shortUserURLs map[string]map[string]string
}

func NewDBStorage(dsn string) (*DBStorage, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &DBStorage{DB: db}, nil
}

func (d *DBStorage) CloseResources() error {
	return d.DB.Close()
}

func (d *DBStorage) GetAll() map[string]map[string]string {
	return nil
}

func (d *DBStorage) GetByOriginalURLForUser(originalURL, userID string) (string, error) {
	//urls := f.GetAll()
	//
	//if usersURLs, ok := urls[userID]; ok {
	//	return usersURLs[originalURL], nil
	//}

	return "", nil
}

func (d *DBStorage) Set(originalURL, shortenURL, userID string) error {

	//if userURLs, ok := f.shortUserURLs[userID]; ok {
	//	userURLs[originalURL] = shortenURL
	//} else {
	//	f.shortUserURLs[userID] = make(map[string]string)
	//	f.shortUserURLs[userID][originalURL] = shortenURL
	//}
	//
	//r := &service.Record{
	//	OriginalURL: originalURL,
	//	ShortenURL:  shortenURL,
	//	UserID:      userID,
	//}
	//if err := f.writer.Write(r); err != nil {
	//	return err
	//}
	return nil
}

func (d *DBStorage) PingStorage() error {
	return d.DB.Ping()
}
