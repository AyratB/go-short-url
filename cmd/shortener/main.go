package main

import (
	"github.com/AyratB/go-short-url/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/{id}", handlers.GetURLHandler)
	router.HandleFunc("/", handlers.SaveURLHandler)

	http.Handle("/", router)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	log.Fatal(server.ListenAndServe())
}
