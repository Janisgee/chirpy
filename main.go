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
	jwtSecret      string
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

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		fmt.Printf("JWTSECRET must be set.\n")
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
		jwtSecret:      jwtSecret,
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
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)

	// ("/api/chirps") Get one chirp
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetOneChirp)

	// ("/api/chirps") Create chirp
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	// ("/api/users") allow users to be created
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

	// ("/api/login") allow user to login
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)

	// ("/api/refresh") allow to get refresh token
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)

	// ("/api/revoke") allow to revoke token
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	// ("/api/revoke") allow to update email and password their own
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	// ("/api/chirps/{chirpID}") delete user chip
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpByUser)

	// ("/api/revoke") set user to be is_chirpy_red
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhooks)

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
