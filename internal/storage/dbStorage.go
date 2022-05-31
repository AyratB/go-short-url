package storage

import (
	"context"
	"database/sql"
	"fmt"
	customerrors "github.com/AyratB/go-short-url/internal/errors"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"time"
)

type DBStorage struct {
	DB            *sql.DB
	shortUserURLs map[string]map[string]string
}

type DBEntity struct {
	originalURL string
	shortenURL  string
	UserData
}

type UserData struct {
	userID   int
	userGUID string
}

func NewDBStorage(dsn string) (*DBStorage, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	dbStorage := &DBStorage{DB: db}

	if err = dbStorage.initTables(); err != nil {
		return nil, err
	}

	return dbStorage, nil
}

func (d *DBStorage) initTables() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	initQuery := `
		CREATE TABLE IF NOT EXISTS users (
    		id 				SERIAL PRIMARY KEY,
    		guid        	text NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS user_urls (
    		id 				SERIAL PRIMARY KEY,
    		original_url 	TEXT NOT NULL,
			shorten_url		TEXT NOT NULL,
			user_id			INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON user_urls (original_url);
	`
	if _, err := d.DB.ExecContext(ctx, initQuery); err != nil {
		return err
	}
	return nil
}

func (d *DBStorage) CloseResources() error {
	return d.DB.Close()
}

func (d *DBStorage) GetAll() (map[string]map[string]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	urls := make([]DBEntity, 0)

	query := `
		SELECT uu.original_url, uu.shorten_url, u.id, u.guid 
			FROM user_urls as uu
		JOIN users as u ON uu.user_id = u.id
	`

	rows, err := d.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var e DBEntity
		if err = rows.Scan(&e.originalURL, &e.shortenURL, &e.userID, &e.userGUID); err != nil {
			return nil, err
		}
		urls = append(urls, e)
	}

	result := make(map[string]map[string]string)

	for _, urlInfo := range urls {
		if userData, ok := result[urlInfo.userGUID]; ok {
			userData[urlInfo.originalURL] = urlInfo.shortenURL
		} else {
			result[urlInfo.userGUID] = make(map[string]string)
			result[urlInfo.userGUID][urlInfo.originalURL] = urlInfo.shortenURL
		}
	}

	return result, nil
}

func (d *DBStorage) getUserByGUID(userGUID string) (userID int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = d.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE guid = $1", userGUID).Scan(&userID)
	defer cancel()
	return
}

func (d *DBStorage) saveUser(userGUID string) (userID int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = d.DB.QueryRowContext(ctx, "INSERT INTO users(guid) VALUES($1) RETURNING id", userGUID).
		Scan(&userID)
	defer cancel()
	return
}

func (d *DBStorage) GetByOriginalURLForUser(originalURL, userGUID string) (string, error) {
	urls, err := d.GetAll()
	if err != nil {
		return "", err
	}

	if usersURLs, ok := urls[userGUID]; ok {
		return usersURLs[originalURL], nil
	}

	return "", nil
}

func (d *DBStorage) Set(originalURL, shortenURL, userGUID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	urls, err := d.GetAll()
	if err != nil {
		return err
	}
	var userID int

	if userData, ok := urls[userGUID]; ok {
		if _, ok := userData[originalURL]; !ok {
			if userID, err = d.getUserByGUID(userGUID); err != nil {
				return err
			}
		}
	} else {
		if userID, err = d.saveUser(userGUID); err != nil {
			return err
		}
	}

	result, err := d.DB.ExecContext(ctx, "INSERT INTO user_urls (original_url, shorten_url, user_id) VALUES ($1, $2, $3)", originalURL, shortenURL, userID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pgerrcode.UniqueViolation {
			return customerrors.ErrDuplicateEntity
		}
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (d *DBStorage) PingStorage() error {
	return d.DB.Ping()
}
