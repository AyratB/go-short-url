package shortener

import (
	"errors"
	"fmt"
	"github.com/AyratB/go-short-url/internal/storage"
	"github.com/AyratB/go-short-url/internal/utils"
	"math/rand"
)

var shortURLs map[string]string

func init() {
	s := storage.Storage{}
	shortURLs = s.GetDB().(map[string]string)
}

const (
	letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterCount = 6
	addressHead = "http://localhost:8080"
)

func getRandomURL(longURL string) string {

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	res := string(b)
	shortURLs[longURL] = res

	return res
}

func MakeShortURL(longURL string) (string, error) {

	if !utils.IsValidURL(longURL) {
		return "", errors.New("uncorrect URL format")
	}

	shortURL, ok := shortURLs[longURL]
	if !ok {
		shortURL = getRandomURL(longURL)
	}

	return fmt.Sprintf("%s/%s", addressHead, shortURL), nil
}

func GetRawURL(shortURL string) (string, error) {

	for longValue, shortValue := range shortURLs {
		if shortValue == shortURL {
			return longValue, nil
		}
	}

	return "", fmt.Errorf("no URL for id = %s", shortURL)
}
