package shortener

import (
	"errors"
	"fmt"
	"github.com/AyratB/go-short-url/internal/repositories"
	"github.com/AyratB/go-short-url/internal/utils"
	"math/rand"
	"time"
)

const (
	letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterCount = 8
)

func GetNewShortener(repo repositories.Repository) *Shortener {
	return &Shortener{repo: repo}
}

type Shortener struct {
	repo repositories.Repository
}

func (s *Shortener) MakeShortURL(longURL, userID string) (string, error) {

	if !utils.IsValidURL(longURL) {
		return "", errors.New("uncorrect URL format")
	}

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	shortURL := string(b)

	if err := s.repo.Set(longURL, shortURL, userID); err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *Shortener) GetOriginalURL(shortenURL string) (string, error) {

	urlsMap, err := s.repo.GetAll()
	if err != nil {
		return "", err
	}

	for _, shortenURLs := range urlsMap {
		for originalURL, shortValue := range shortenURLs {
			if shortValue == shortenURL {
				return originalURL, nil
			}
		}
	}

	return "", fmt.Errorf("no URL for id = %s", shortenURL)
}

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Shortener) GetExistingURLS(originalURL, userGUID string) (string, error) {
	return s.repo.GetByOriginalURLForUser(originalURL, userGUID)
}

func (s *Shortener) GetAllSavedUserURLs(baseURL, userID string) ([]*URL, error) {
	urlsMap, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	urls := make([]*URL, 0)

	for longURL, shortURL := range urlsMap[userID] {
		urls = append(urls, &URL{
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, shortURL),
			OriginalURL: longURL,
		})
	}

	return urls, nil
}

func (s *Shortener) PingStorage() error {
	return s.repo.PingStorage()
}
