package utils

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseUrl         string
	FileStoragePath string
}

func GetConfigs() (*Config, error) {
	config := &Config{}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	if len(config.ServerAddress) == 0 {
		flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "server address")
	}
	config.BaseUrl = os.Getenv("BASE_URL")
	if len(config.BaseUrl) == 0 {
		flag.StringVar(&config.BaseUrl, "b", "http://localhost:8080", "base url")
	}
	config.FileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	if len(config.FileStoragePath) == 0 {
		flag.StringVar(&config.FileStoragePath, "f", "", "file storage path")
	}

	flag.Parse()
	return config, nil
}
