package shortener

import "math/rand"

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

func GetShortURL(longUrl string) string {
	shortURL, ok := shortURLs[longUrl]
	if !ok {
		shortURL = getRandomURL(longUrl)
	}
	return shortURL
}
