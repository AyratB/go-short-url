package main

import (
	shortener "github.com/AyratB/go-short-url/internal/app"
	"io"
	"log"
	"net/http"
)

//func GetHandler(w http.ResponseWriter, r *http.Request) {
//	// этот обработчик принимает только запросы, отправленные методом GET
//	if r.Method != http.MethodGet {
//		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
//		return
//	}
//	// продолжаем обработку запроса
//	// ...
//}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	rawURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(rawURL) == 0 {
		http.Error(w, "Raw URL string length must be greater than 0", http.StatusBadRequest)
		return
	}

	shortURL := shortener.GetShortURL(string(rawURL))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func main() {

	//http.HandleFunc("/", HelloWorld)
	http.HandleFunc("/POST/", PostHandler)

	//http.HandleFunc("/GET/{id}", GetHandler)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	log.Fatal(server.ListenAndServe())
}
