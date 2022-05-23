package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"github.com/AyratB/go-short-url/internal/utils"
	"log"
)

func main() {

	handler, err := server.Run(utils.GetConfigs())

	defer func() {
		for _, closers := range handler.ReposClosers {
			err = closers()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
}
