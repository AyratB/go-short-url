package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"github.com/AyratB/go-short-url/internal/utils"
	"log"
)

func main() {

	_, err := server.Run(utils.GetConfigs())
	//defer func() {
	//	if resourcesCloser != nil {
	//		resourcesCloser()
	//	}
	//}()

	if err != nil {
		//resourcesCloser()
		log.Fatal(err)
	}

}
