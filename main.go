package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World!")
	// Create an empty servemux
	mux := http.NewServeMux()
	svr := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	svr.ListenAndServe()
}
