package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AyratB/go-short-url/internal/app"
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

	userID := fmt.Sprint(r.Context().Value(utils.CurrentUser))

	shortURL, err := h.sh.MakeShortURL(p.URL, h.configs.BaseURL, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp, err := json.Marshal(PostURLResponse{Result: shortURL})
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

func (h *Handler) BatchHandler(w http.ResponseWriter, r *http.Request) {
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

	userID := fmt.Sprint(r.Context().Value(utils.CurrentUser))

	for _, originalButch := range originalBatches {

		shortenButch := BatchResponse{CorrelationId: originalButch.CorrelationId}

		if shortenButch.ShortURL, err = h.sh.MakeShortURL(originalButch.OriginalUrl, h.configs.BaseURL, userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

	userID := fmt.Sprint(r.Context().Value(utils.CurrentUser))

	shortURL, err := h.sh.MakeShortURL(string(rawURL), h.configs.BaseURL, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Handler) GetAllSavedUserURLs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	userID := fmt.Sprint(r.Context().Value(utils.CurrentUser))

	// получаем все урлы
	urls, err := h.sh.GetAllSavedUserURLs(h.configs.BaseURL, userID)

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
