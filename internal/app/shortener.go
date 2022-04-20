package shortener

import "math/rand"

var ShortURLs map[string]string

const (
	letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterCount = 6
)

func getRandomURL() string {
	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetShortURL(longUrl string) string {
	shortURL, ok := ShortURLs[longUrl]
	if !ok {
		shortURL = getRandomURL()
	}
	return shortURL
}
