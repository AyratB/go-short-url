package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"github.com/AyratB/go-short-url/internal/utils"
	"log"
)

func main() {

	configs, err := utils.GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	resourcesCloser, err := server.Run(configs)
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
