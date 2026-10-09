package middleware

import (
	"net/http"
)

func ContentTypeValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if r.Header.Get("Content-Type") != "application/json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"bad_request","message":"Content-Type must be application/json"}`))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
