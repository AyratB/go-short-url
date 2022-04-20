package main

import (
	shortener "github.com/AyratB/go-short-url/internal/app"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		http.Error(w, "Need to set id", http.StatusMethodNotAllowed)
		return
	}

	longURL, err := shortener.GetRawURL(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	rawURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	shortURL, err := shortener.MakeSHortURL(string(rawURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/{id}", GetHandler)
	router.HandleFunc("/", PostHandler)

	http.Handle("/", router)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	log.Fatal(server.ListenAndServe())
}
