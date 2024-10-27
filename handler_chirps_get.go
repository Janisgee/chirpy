package main

import "net/http"

// Handler to get all chirps
func (cfg *apiConfig) handlerAllChirpsGet(w http.ResponseWriter, r *http.Request) {

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
