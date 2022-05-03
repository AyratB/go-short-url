package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"log"
)

func main() {
	if err := server.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
