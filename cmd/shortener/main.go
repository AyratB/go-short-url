package main

import (
	"github.com/AyratB/go-short-url/internal/server"
	"log"
)

func main() {
	log.Fatal(server.Run("localhost:8080"))
}
