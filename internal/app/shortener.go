package shortener

import (
	"errors"
	"fmt"
	"github.com/AyratB/go-short-url/internal/repositories"
	"github.com/AyratB/go-short-url/internal/utils"
	"math/rand"
)

const (
	letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterCount = 6
)

func GetNewShortener(repo repositories.Repository) *Shortener {
	return &Shortener{repo: repo}
}

type Shortener struct {
	repo repositories.Repository
}

func (s *Shortener) getShortenURL(longURL, userID string) (string, error) {

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	shortURL := string(b)
	err := s.repo.Set(longURL, shortURL, userID)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *Shortener) MakeShortURL(longURL, baseURL, userID string) (string, error) {

	if !utils.IsValidURL(longURL) {
		return "", errors.New("uncorrect URL format")
	}

	// check if already exists
	shortURL, err := s.repo.GetByOriginalURLForUser(longURL, userID)

	if err != nil {
		return "", err
	}

	if len(shortURL) == 0 {
		shortURL, err = s.getShortenURL(longURL, userID)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s/%s", baseURL, shortURL), nil
}

func (s *Shortener) GetOriginalURL(shortenURL string) (string, error) {

	urlsMap := s.repo.GetAll()

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

func (s *Shortener) GetAllSavedUserURLs(baseURL, userID string) ([]*URL, error) {
	urlsMap := s.repo.GetAll()

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
