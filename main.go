package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Janisgee/chirpy.git/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error in connecting database: %s\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	// Create an empty servemux
	mux := http.NewServeMux()

	// ("/app") Build a fileserver
	fs := http.FileServer(http.Dir(filepathRoot))
	handler := http.StripPrefix("/app", fs)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	// apiCfg.middlewareMetricsInc(handler)
	// ("/healthz") Add the readiness endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// ("/metrics") Add the requset count endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerRequestCount)

	// ("/reset") Add the requset count endpoint
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteAllUsers)

	// ("/api/validate_chirp") connect to Chirpy API
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	// ("/api/users") allow users to be created
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)

	svr := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %s on port:%s\n", filepathRoot, port)
	svr.ListenAndServe()
}
