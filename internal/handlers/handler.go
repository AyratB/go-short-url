package handlers

import (
	"encoding/json"
	shortener "github.com/AyratB/go-short-url/internal/app"
	"github.com/AyratB/go-short-url/internal/repositories"
	"github.com/AyratB/go-short-url/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
)

type PostURLRequest struct {
	URL string `json:"url"`
}

type PostURLResponse struct {
	Result string `json:"result"`
}

type Handler struct {
	sh *shortener.Shortener
}

func NewHandler() (*Handler, error) {
	var repo repositories.Repository
	var err error

	filePath := os.Getenv("FILE_STORAGE_PATH")
	if len(filePath) == 0 {
		repo = storage.NewMemoryStorage()
	} else {
		repo, err = storage.NewFileStorage(filePath)
		if err != nil {
			return nil, err
		}
	}
	return &Handler{sh: shortener.GetNewShortener(repo)}, nil
}

func (h *Handler) PostShortenURLHandler(w http.ResponseWriter, r *http.Request) {
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

	shortURL, err := h.sh.MakeShortURL(p.URL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp, err := json.Marshal(PostURLResponse{Result: shortURL})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(resp)
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	if len(id) == 0 {
		http.Error(w, "Need to set id", http.StatusBadRequest)
		return
	}

	longURL, err := h.sh.GetRawURL(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", longURL)
	w.Header().Set("content-type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	rawURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	shortURL, err := h.sh.MakeShortURL(string(rawURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
