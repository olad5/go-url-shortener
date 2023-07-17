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
	"github.com/olad5/go-url-shortener/utils"
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
	tests := []struct {
		Name         string
		requestBody  utils.RequestBody
		ExpectedCode int
	}{
		{
			Name: "Should Shorten the url",
			requestBody: utils.RequestBody{
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

func TestRedirect(t *testing.T) {
	shortCode, err := createShortUrl()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		Name               string
		ShortCode          string
		ExpectedClickCount float64
		ExpectedCode       int
	}{
		{
			Name:               "Should redirect to the original url found in the database",
			ShortCode:          shortCode,
			ExpectedClickCount: 1,
			ExpectedCode:       307,
		},
	}
	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer

			url := "/api/v1/" + test.ShortCode
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
			clickCount, err := findShortUrl(test.ShortCode)
			if err != nil {
				t.Fatal(err)
			}
			if count := clickCount; count != test.ExpectedClickCount {
				t.Errorf("handler returned wrong status code: got %v want %v", count, test.ExpectedClickCount)
				return
			}
		}

		t.Run(test.Name, fn)
	}
}

func createShortUrl() (string, error) {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(utils.RequestBody{Url: "https://example.com"}); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "/api/v1/shorten", &b)
	if err != nil {
		return "", err
	}
	response := httptest.NewRecorder()
	router.Serve(response, req)
	responseBody := make(map[string]interface{})
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return "", err
	}

	shortCode := responseBody["data"].(map[string]interface{})["short_url"].(string)
	return shortCode, nil
}

func findShortUrl(shortCode string) (float64, error) {
	var b bytes.Buffer
	req, err := http.NewRequest(http.MethodGet, "/api/v1/info/"+shortCode, &b)
	if err != nil {
		return 0, err
	}
	response := httptest.NewRecorder()
	router.Serve(response, req)
	responseBody := make(map[string]interface{})
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return 0, err
	}

	clickCount := (responseBody["data"].(map[string]interface{})["click_count"].(float64))
	if err != nil {
		return 0, err
	}
	return clickCount, nil
}
