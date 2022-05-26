package utils

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
}

func GetConfigs() *Config {
	config := &Config{}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	if len(config.ServerAddress) == 0 {
		flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "server address")
	}
	config.BaseURL = os.Getenv("BASE_URL")
	if len(config.BaseURL) == 0 {
		flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base url")
	}
	config.FileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	if len(config.FileStoragePath) == 0 {
		flag.StringVar(&config.FileStoragePath, "f", "", "file storage path")
	}
	config.DatabaseDSN = os.Getenv("DATABASE_DSN")
	if len(config.DatabaseDSN) == 0 {
		flag.StringVar(&config.DatabaseDSN, "d", "", "db storage path")
	}

	flag.Parse()
	return config
}
