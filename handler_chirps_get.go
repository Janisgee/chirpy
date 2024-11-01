package main

import (
	"net/http"

	"github.com/Janisgee/chirpy.git/internal/database"
	"github.com/google/uuid"
)

// Handler to get all chirps
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	// Grab the query parameters from the URL (author_id)
	authorIDString := r.URL.Query().Get("author_id")
	// (sort)
	sortMethod := r.URL.Query().Get("sort")

	var chirpsResult []database.Chirp
	var err error

	if authorIDString != "" {

		// Parse string into UUID
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Fail to convert author ID string into UUID format", err)
		}

		// Get all chirps by UserID
		if sortMethod == "desc" {
			chirpsResult, err = cfg.db.GetAllChipsByUserIDDesc(r.Context(), authorID)
		} else {
			chirpsResult, err = cfg.db.GetAllChipsByUserID(r.Context(), authorID)
		}
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't found chip data with the provided author id", err)
			return
		}

	} else {

		//Get all chirps data from database
		if sortMethod == "desc" {

			chirpsResult, err = cfg.db.GetAllChirpsDesc(r.Context())
		} else {
			chirpsResult, err = cfg.db.GetAllChirps(r.Context())
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrive all chirps from database", err)
			return
		}
		if len(chirpsResult) == 0 {
			respondWithJSON(w, http.StatusOK, Chirp{})
			return
		}
	}
	arrayChirps := []Chirp{}
	// Maps the database chirp to API chirp struct
	for i := range chirpsResult {
		chirp := Chirp{
			ID:        chirpsResult[i].ID,
			CreatedAt: chirpsResult[i].CreatedAt,
			UpdatedAt: chirpsResult[i].UpdatedAt,
			Body:      chirpsResult[i].Body,
			UserID:    chirpsResult[i].UserID,
		}

		arrayChirps = append(arrayChirps, chirp)

	}
	respondWithJSON(w, http.StatusOK, arrayChirps)
}

// Handler to get all chirps
func (cfg *apiConfig) handlerGetOneChirp(w http.ResponseWriter, r *http.Request) {

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

	reqChirpData, err := cfg.db.GetOneChirp(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't found chip data with the provided chirp id", err)
		return
	}
	chirp := Chirp{
		ID:        reqChirpData.ID,
		CreatedAt: reqChirpData.CreatedAt,
		UpdatedAt: reqChirpData.UpdatedAt,
		Body:      reqChirpData.Body,
		UserID:    reqChirpData.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
