package middleware

import "net/http"

func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		rw.WriteHeader(http.StatusOK)
		next.ServeHTTP(rw, r)
	})
}
