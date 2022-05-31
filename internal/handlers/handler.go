package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AyratB/go-short-url/internal/app"
	custom_errors "github.com/AyratB/go-short-url/internal/errors"
	"github.com/AyratB/go-short-url/internal/middlewares"
	"github.com/AyratB/go-short-url/internal/repositories"
	"github.com/AyratB/go-short-url/internal/storage"
	"github.com/AyratB/go-short-url/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type PostURLRequest struct {
	URL string `json:"url"`
}

type PostURLResponse struct {
	Result string `json:"result"`
}

type Handler struct {
	configs *utils.Config
	sh      *shortener.Shortener
}

func NewHandler(configs *utils.Config) (*Handler, func() error, error) {

	var repo repositories.Repository
	var err error

	//configs.DatabaseDSN = "postgres://postgres:test@localhost:5432/postgres?sslmode=disable"

	if len(configs.DatabaseDSN) != 0 {
		repo, err = storage.NewDBStorage(configs.DatabaseDSN)
		if err != nil {
			return nil, nil, err
		}
	} else if len(configs.FileStoragePath) != 0 {
		repo, err = storage.NewFileStorage(configs.FileStoragePath)
		if err != nil {
			return nil, nil, err
		}
	} else {
		repo = storage.NewMemoryStorage()
	}

	return &Handler{
		sh:      shortener.GetNewShortener(repo),
		configs: configs,
	}, repo.CloseResources, nil
}

func (h *Handler) SaveJSONURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	p := PostURLRequest{}

	if err := json.Unmarshal(b, &p); err != nil {
		http.Error(w, "Incorrect body JSON format", http.StatusBadRequest)
		return
	}

	if len(p.URL) == 0 {
		http.Error(w, "URL can not be empty", http.StatusBadRequest)
		return
	}

	userGUID := fmt.Sprint(r.Context().Value(middlewares.CtxKey{}))

	shortURL, err := h.sh.MakeShortURL(p.URL, userGUID)

	if errors.Is(err, custom_errors.ErrDuplicateEntity) {
		shortURL, err = h.sh.GetExistingURLS(p.URL, userGUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusConflict)

	} else {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}

	resp, err := json.Marshal(PostURLResponse{Result: fmt.Sprintf("%s/%s", h.configs.BaseURL, shortURL)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func (h *Handler) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	err := h.sh.PingStorage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

type BatchRequest struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handler) SaveBatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	originalBatches := make([]BatchRequest, 0)
	shortenBatches := make([]BatchResponse, 0)

	if err := json.Unmarshal(b, &originalBatches); err != nil {
		http.Error(w, "Incorrect body JSON format", http.StatusBadRequest)
		return
	}

	userID := fmt.Sprint(r.Context().Value(middlewares.CtxKey{}))

	for _, originalButch := range originalBatches {

		shortenButch := BatchResponse{CorrelationId: originalButch.CorrelationId}

		shortenURL, err := h.sh.MakeShortURL(originalButch.OriginalUrl, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenButch.ShortURL = fmt.Sprintf("%s/%s", h.configs.BaseURL, shortenURL)
		shortenBatches = append(shortenBatches, shortenButch)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp, err := json.Marshal(shortenBatches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	shortenURL := chi.URLParam(r, "id")

	if len(shortenURL) == 0 {
		http.Error(w, "Need to set id", http.StatusBadRequest)
		return
	}

	originalURL, err := h.sh.GetOriginalURL(shortenURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", originalURL)
	w.Header().Set("content-type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SaveBodyURLHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	rawURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	userGUID := fmt.Sprint(r.Context().Value(middlewares.CtxKey{}).(string))

	var shortURL string

	shortURL, err = h.sh.MakeShortURL(string(rawURL), userGUID)

	if errors.Is(err, custom_errors.ErrDuplicateEntity) {
		shortURL, err = h.sh.GetExistingURLS(string(rawURL), userGUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
	} else {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte(fmt.Sprintf("%s/%s", h.configs.BaseURL, shortURL)))
}

func (h *Handler) GetAllSavedUserURLs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	userGUID := fmt.Sprint(r.Context().Value(middlewares.CtxKey{}).(string))

	// получаем все урлы
	urls, err := h.sh.GetAllSavedUserURLs(h.configs.BaseURL, userGUID)

	if err != nil {
		http.Error(w, "Errors happens when get all saved URLS!", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
