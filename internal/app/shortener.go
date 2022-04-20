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

func getRandomURL(longUrl string) string {

	if shortURLs == nil {
		shortURLs = make(map[string]string, 0)
	}

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	res := string(b)
	shortURLs[longUrl] = res

	return res
}

func MakeSHortURL(longUrl string) (string, error) {

	if !utils.IsValidURL(longUrl) {
		return "", errors.New("uncorrect URL format")
	}

	shortURL, ok := shortURLs[longUrl]
	if !ok {
		shortURL = getRandomURL(longUrl)
	}
	return shortURL, nil
}

func GetRawURL(shortURL string) (string, error) {

	if shortURLs != nil {
		for longValue, shortValue := range shortURLs {
			if shortValue == shortURL {
				return longValue, nil
			}
		}
	}
	return "", fmt.Errorf("no URL for id = %s", shortURL)
}
