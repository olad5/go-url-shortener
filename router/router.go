package router

import (
	"fmt"
	"net/http"

	"github.com/olad5/go-url-shortener/middleware"
)

func Home(http.ResponseWriter, *http.Request) {
	fmt.Printf("hello from Home")
}

func Initialize() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	return middleware.Json(mux)
}
