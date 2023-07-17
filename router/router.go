package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/olad5/go-url-shortener/handlers"
	"github.com/olad5/go-url-shortener/middleware"
	"github.com/olad5/go-url-shortener/utils"
)

func Initialize() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Serve)
	return middleware.Json(mux)
}

var (
	baseUrl = "/api/v1"
	routes  = []route{
		newRoute("GET", baseUrl+"/healthcheck", handlers.Healthcheck),
		newRoute("POST", baseUrl+"/shorten", handlers.Shorten),
		newRoute("GET", baseUrl+"/info/([^/]+)", handlers.Info),
		newRoute("GET", baseUrl+"/([^/]+)", handlers.Redirect),
	}
)

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), utils.ParamsContextkey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len((allow)) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}
