package middlewares

import (
	"context"
	"errors"
	"github.com/AyratB/go-short-url/internal/utils"
	"net/http"

	"github.com/google/uuid"
)

type CookieHandler struct {
	decoder *utils.Decoder
}

func NewCookieHandler(decoder *utils.Decoder) *CookieHandler {
	return &CookieHandler{decoder: decoder}
}

type CookieUserName string

const cookieUserName = CookieUserName("UserID")

type CtxKey struct{}

func (c *CookieHandler) CookieHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(string(cookieUserName))
		var currentUser = ""

		if errors.Is(err, http.ErrNoCookie) {

			userID := uuid.New().String()

			token := c.decoder.EnCode(userID)

			newCookie := &http.Cookie{
				Name:  string(cookieUserName),
				Value: token,
				Path:  "/",
			}

			http.SetCookie(w, newCookie)
			r.AddCookie(newCookie)

			currentUser = userID
		} else if err != nil {
			http.Error(w, "Cookie crumbled", http.StatusInternalServerError)
		} else {
			decoded, err := c.decoder.Decode(cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			if len(decoded) != 0 {
				currentUser = decoded
			}
		}
		ctx := context.WithValue(r.Context(), CtxKey{}, currentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
