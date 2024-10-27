package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Janisgee/chirpy.git/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	//Load env variables
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Printf("DB_URL must be set.\n")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		fmt.Printf("PLATFORM must be set.\n")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error in connecting database: %s\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	// Create an empty servemux
	mux := http.NewServeMux()

	// ("/app") Build a fileserver
	fs := http.FileServer(http.Dir(filepathRoot))
	handler := http.StripPrefix("/app", fs)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	// apiCfg.middlewareMetricsInc(handler)

	// ("/reset") Add the requset count endpoint
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteAllUsers)

	// ("/api/chirps") Get all chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerAllChirpsGet)

	// ("/api/chirps") Create chirp
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	// ("/api/users") allow users to be created
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

	// ("/healthz") Add the readiness endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// ("/metrics") Add the requset count endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerRequestCount)

	svr := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %s on port:%s\n", filepathRoot, port)
	svr.ListenAndServe()
}
