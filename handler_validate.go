package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Janisgee/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// Handler to connect to Chirps APU
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	// // Get user_id by login email
	// userdata, err := cfg.db.GetUser(r.Context(), cfg.CurrentUserEmail)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid input", err)
	// 	return
	// }

	// Decode request body (chirp body)
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Get clean body string
	cleanedString, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Create chirp params for building data
	chirpParams := database.CreateChirpParams{
		Body:   cleanedString,
		UserID: params.User_ID,
	}

	// Create Chirp in the database
	chirpData, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// Maps the database chirp to API chirp struct
	chirp := Chirp{
		ID:        chirpData.ID,
		CreatedAt: chirpData.CreatedAt,
		UpdatedAt: chirpData.UpdatedAt,
		Body:      chirpData.Body,
		UserID:    chirpData.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}

func validateChirp(body string) (string, error) {
	// Validate Chirps are under 140 characters
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	// Validate words
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// Finalized clean string
	cleanedString := getCleanedBody(body, badWords)

	return cleanedString, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		_, ok := badWords[loweredWord]
		if ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
