package main

import (
	"net/http"

	"github.com/google/uuid"
)

// Handler to get all chirps
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	//Get all chirps data from database
	allChirpsData, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrive all chirps from database", err)
		return
	}
	if len(allChirpsData) == 0 {
		respondWithJSON(w, http.StatusOK, Chirp{})
		return
	}

	arrayChirps := []Chirp{}
	// Maps the database chirp to API chirp struct
	for i := range allChirpsData {
		chirp := Chirp{
			ID:        allChirpsData[i].ID,
			CreatedAt: allChirpsData[i].CreatedAt,
			UpdatedAt: allChirpsData[i].UpdatedAt,
			Body:      allChirpsData[i].Body,
			UserID:    allChirpsData[i].UserID,
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
