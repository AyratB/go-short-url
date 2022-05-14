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

func (s *Shortener) getRandomURL(longURL string) (string, error) {

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	res := string(b)
	err := s.repo.Set(longURL, res)

	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *Shortener) MakeShortURL(longURL string) (string, error) {

	if !utils.IsValidURL(longURL) {
		return "", errors.New("uncorrect URL format")
	}

	// check if already exists
	shortURL, err := s.repo.GetByKey(longURL)

	if err != nil {
		return "", err
	}

	if len(shortURL) == 0 {
		shortURL, err = s.getRandomURL(longURL)
		if err != nil {
			return "", err
		}
	}

	ah := utils.GetEnvOrDefault("BASE_URL", utils.DefaultBaseURL)

	return fmt.Sprintf("%s/%s", ah, shortURL), nil
}

func (s *Shortener) GetRawURL(shortURL string) (string, error) {

	shortURLs, err := s.repo.GetAll()

	if err != nil {
		return "", err
	}

	for longValue, shortValue := range shortURLs {
		if shortValue == shortURL {
			return longValue, nil
		}
	}

	return "", fmt.Errorf("no URL for id = %s", shortURL)
}
