package server

import (
	"github.com/AyratB/go-short-url/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func Run(host string) error {
	router := mux.NewRouter()
	router.HandleFunc("/{id}", handlers.GetURLHandler)
	router.HandleFunc("/", handlers.SaveURLHandler)

	http.Handle("/", router)

	server := &http.Server{
		Addr: host,
	}

	return server.ListenAndServe()
}
