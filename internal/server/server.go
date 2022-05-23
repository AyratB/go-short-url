package server

import (
	"github.com/AyratB/go-short-url/internal/handlers"
	"github.com/AyratB/go-short-url/internal/middlewares"
	"github.com/AyratB/go-short-url/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Run(configs *utils.Config) (*handlers.Handler, error) {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.CookieHandler)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler := handlers.NewHandler(configs)

	r.Route("/", func(r chi.Router) {
		r.Post("/api/shorten", handler.PostShortenURLHandler)
		r.Get("/api/user/urls", handler.GetAllSavedURLs)
		r.Get("/{id}", handler.GetURLHandler)
		r.Post("/", handler.SaveURLHandler)
	})

	return handler, http.ListenAndServe(configs.ServerAddress, r)
}
