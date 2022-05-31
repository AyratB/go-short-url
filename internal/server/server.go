package server

import (
	"github.com/AyratB/go-short-url/internal/handlers"
	"github.com/AyratB/go-short-url/internal/middlewares"
	"github.com/AyratB/go-short-url/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Run(configs *utils.Config) (func() error, error) {

	r := chi.NewRouter()

	decoder := utils.NewDecoder()
	cookieHandler := middlewares.NewCookieHandler(decoder)

	r.Use(middleware.RequestID)
	r.Use(middlewares.GzipHandle)
	r.Use(cookieHandler.CookieHandler)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler, resourcesCloser, err := handlers.NewHandler(configs)
	if err != nil {
		return resourcesCloser, err
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/api/user/urls", handler.GetAllSavedUserURLs)
		r.Get("/{id}", handler.GetURLHandler)
		r.Get("/ping", handler.PingDBHandler)
		r.Post("/api/shorten", handler.SaveJSONURLHandler)
		r.Post("/api/shorten/batch", handler.SaveBatchHandler)
		r.Post("/", handler.SaveBodyURLHandler)
	})

	return resourcesCloser, http.ListenAndServe(configs.ServerAddress, r)
}
