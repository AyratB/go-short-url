package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetURLHandler)
		r.Post("/", SaveURLHandler)
	})
	return r
}

func TestSaveURLHandler(t *testing.T) {

	type want struct {
		statusCode  int
		redirectURL string
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
				redirectURL: "http://localhost:8080/rfBd67",
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
				redirectURL: "uncorrect URL format\n",
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
				redirectURL: "",
				contentType: "",
			},
			requestType: http.MethodDelete,
		},
		{
			name:    "simple test #4 with uncorrect url format",
			request: "/",
			body:    "/ya.ru",
			want: want{
				statusCode:  http.StatusBadRequest,
				redirectURL: "uncorrect URL format\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.requestType, tt.request, strings.NewReader(tt.body))

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.redirectURL, body)

			resp.Body.Close()
		})
	}
}

func TestGetURLHandler(t *testing.T) {

	type want struct {
		statusCode  int
		redirectURL string
		contentType string
		body        string
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
				redirectURL: "https://yatest.ru",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodGet,
		},
		{
			name:    "negative test #2 with wrong method type",
			request: "/test",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				redirectURL: "",
				contentType: "",
			},
			requestType: http.MethodDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.requestType, tt.request, nil)

			fmt.Println(body)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.redirectURL, resp.Header.Get("Location"))

			resp.Body.Close()
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, body)

	require.NoError(t, err)

	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := httpClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
