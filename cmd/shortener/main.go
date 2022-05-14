package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"github.com/AyratB/go-short-url/internal/utils"
	"log"
)

func main() {
	sa := utils.GetEnvOrDefault("SERVER_ADDRESS", utils.DefaultServerAddress)

	resourcesCloser, err := server.Run(sa)
	defer func() {
		if resourcesCloser != nil {
			resourcesCloser()
		}
	}()

	if err != nil {
		resourcesCloser()
		log.Fatal(err)
	}
}
