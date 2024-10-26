package main

import (
	"fmt"
	"net/http"
)

// Handler that writes the number of requests that have been counted to HTTP response
func (cfg *apiConfig) handlerRequestCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()

	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)
}

// middleware method on *apiConfig
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hitCheck := cfg.fileserverHits.Add(1)
		fmt.Printf("hitCheck: %d\n", hitCheck)
		next.ServeHTTP(w, r)
	})
}
