package shortener

import (
	"errors"
	"fmt"
	"github.com/AyratB/go-short-url/internal/utils"
	"math/rand"
)

var shortURLs map[string]string

const (
	letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterCount = 6
)

func getRandomURL(longURL string) string {

	if shortURLs == nil {
		shortURLs = make(map[string]string)
	}

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	res := string(b)
	shortURLs[longURL] = res

	return res
}

func MakeSHortURL(longURL string) (string, error) {

	if !utils.IsValidURL(longURL) {
		return "", errors.New("uncorrect URL format")
	}

	shortURL, ok := shortURLs[longURL]
	if !ok {
		shortURL = getRandomURL(longURL)
	}
	return shortURL, nil
}

func GetRawURL(shortURL string) (string, error) {

	for longValue, shortValue := range shortURLs {
		if shortValue == shortURL {
			return longValue, nil
		}
	}

	return "", fmt.Errorf("no URL for id = %s", shortURL)
}
