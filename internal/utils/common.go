package utils

import "os"

func GetEnvOrDefault(key string, defaultValue string) string {
	v := os.Getenv(key)
	if len(v) != 0 {
		return v
	}
	return defaultValue
}
