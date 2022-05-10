package server

import (
	"github.com/AyratB/go-short-url/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Run(host string) error {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/api/shorten", handlers.PostShortenURLHandler)
		r.Get("/{id}", handlers.GetURLHandler)
		r.Post("/", handlers.SaveURLHandler)
	})

	return http.ListenAndServe(host, r)
}
