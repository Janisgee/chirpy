package main

import (
	"net/http"

	"github.com/Janisgee/chirpy.git/internal/auth"
	"github.com/Janisgee/chirpy.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpByUser(w http.ResponseWriter, r *http.Request) {

	// Retrieve Chirp ID from API path
	chirpID := r.PathValue("chirpID")
	if len(chirpID) == 0 {
		respondWithError(w, http.StatusBadRequest, "Empty string return as the request was not match with wildcard pattern from path", nil)
		return
	}
	// Convert string to a UUID
	parsedUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chip id into UUID type", err)
		return
	}
	// Get chirpdata userID
	reqChirpData, err := cfg.db.GetOneChirp(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't found chip data with the provided chirp id", err)
		return
	}

	// Get JWT from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to get refresh token from 'Authorization'", err)
		return
	}

	// Validate JWT from header and get userID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get validation from provided JWT", err)
		return
	}

	// Check if the chirpdata ID and userID are the same
	if reqChirpData.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User ID is not the same with chirp data ID", err)
		return
	}

	// If userID correct, delete chirp by userID
	_, err = cfg.db.DeleteOneChirpByUserId(r.Context(), database.DeleteOneChirpByUserIdParams{
		ID:     reqChirpData.ID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp with provided chirp id and user id", err)
		return
	}

	// response code 204 No Content - success but no body
	w.WriteHeader(http.StatusNoContent)

}
