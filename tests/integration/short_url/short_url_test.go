//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/olad5/go-url-shortener/router"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../config/.test.env")
	if err != nil {
		panic("Error Loading .env file")
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestShorten(t *testing.T) {
	type requestBody struct {
		Url string `json:"url"`
	}
	tests := []struct {
		Name         string
		requestBody  requestBody
		ExpectedCode int
	}{
		{
			Name: "Should Shorten the url",
			requestBody: requestBody{
				Url: "https://example.com",
			},
			ExpectedCode: 200,
		},
	}
	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.requestBody)

			url := "/api/v1/shorten"

			req, err := http.NewRequest(http.MethodPost, url, &b)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			rr := httptest.NewRecorder()
			router.Serve(rr, req)

			if status := rr.Code; status != test.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.ExpectedCode)
				return
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		Name         string
		ShortCode    string
		ExpectedCode int
	}{
		{
			Name:         "Should return a 404 error for the shortcode not found",
			ShortCode:    "ksksksksk",
			ExpectedCode: 404,
		},
	}
	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer

			url := "/api/v1/info/" + test.ShortCode
			req, err := http.NewRequest(http.MethodGet, url, &b)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			rr := httptest.NewRecorder()
			router.Serve(rr, req)

			if status := rr.Code; status != test.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.ExpectedCode)
				return
			}
		}

		t.Run(test.Name, fn)
	}
}
