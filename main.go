package main

import (
	"fmt"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8000"
	// Create an empty servemux
	mux := http.NewServeMux()

	// Build a fileserver
	fs := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", fs))

	// Add the readiness endpoint
	mux.HandleFunc("/healthz", handlerReadiness)

	svr := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %s on port:%s\n", filepathRoot, port)
	svr.ListenAndServe()
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
