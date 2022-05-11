package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"github.com/AyratB/go-short-url/internal/utils"
	"log"
)

func main() {
	sa := utils.GetEnvOrDefault("SERVER_ADDRESS", utils.DefaultServerAddress)

	if err := server.Run(sa); err != nil {
		log.Fatal(err)
	}
}
