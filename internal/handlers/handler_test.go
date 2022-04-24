package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveURLHandler(t *testing.T) {

	type want struct {
		statusCode  int
		resultURL   string
		contentType string
	}

	tests := []struct {
		name        string
		request     string
		body        string
		requestType string
		want        want
	}{
		{
			name:    "simple positive test #1",
			request: "/",
			body:    "https://ya.ru",
			want: want{
				statusCode:  http.StatusCreated,
				resultURL:   "http://localhost:8080/rfBd67",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
		{
			name:    "simple test #2 with empty URL",
			request: "/",
			body:    "",
			want: want{
				statusCode:  http.StatusBadRequest,
				resultURL:   "uncorrect URL format\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
		{
			name:    "simple test #3 with uncorrect request type",
			request: "/",
			body:    "https://ya.ru",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				resultURL:   "Only POST requests are allowed by this route!\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodDelete,
		},
		{
			name:    "simple test #4 with uncorrect url format",
			request: "/",
			body:    "/ya.ru",
			want: want{
				statusCode:  http.StatusBadRequest,
				resultURL:   "uncorrect URL format\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.requestType, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "text/plain; charset=utf-8")

			w := httptest.NewRecorder()

			h := http.HandlerFunc(SaveURLHandler)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			fmt.Println(result.Header)

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			shortenerResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.resultURL, string(shortenerResult))
		})
	}
}

func TestGetURLHandler(t *testing.T) {

	type want struct {
		statusCode  int
		resultURL   string
		contentType string
		errorText   string
	}

	tests := []struct {
		name        string
		request     string
		requestType string
		want        want
	}{
		{
			name:    "simple positive test #1",
			request: "/test",
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				resultURL:   "https://yatest.ru",
				errorText:   "",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodGet,
		},
		{
			name:    "negative test #2 with wrong method type",
			request: "/test",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				resultURL:   "",
				errorText:   "Only GET requests are allowed by this route!\n",
				contentType: "",
			},
			requestType: http.MethodDelete,
		},
		{
			name:    "negative test #3 with empty id",
			request: "/",
			want: want{
				statusCode:  http.StatusBadRequest,
				resultURL:   "",
				errorText:   "Need to set id\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(tt.requestType, tt.request, nil)

			w := httptest.NewRecorder()

			// NEED to mock gorilla/mux
			vars := map[string]string{
				"id": strings.TrimPrefix(tt.request, "/"),
			}

			request = mux.SetURLVars(request, vars)

			h := http.HandlerFunc(GetURLHandler)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.resultURL, result.Header.Get("Location"))
			assert.Equal(t, tt.want.errorText, fmt.Sprint(w.Body))

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
