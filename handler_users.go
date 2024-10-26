package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {

	params := struct {
		Email string `json:"email"`
	}{}

	// Decodes the request JSON
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil || params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Creates the user in the database
	userdata, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Maps the database user to your API User struct
	user := User{
		ID:        userdata.ID,
		CreatedAt: userdata.CreatedAt,
		UpdatedAt: userdata.UpdatedAt,
		Email:     userdata.Email,
	}

	// Returns the correct 201 status code
	respondWithJSON(w, http.StatusCreated, user)

}

// func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
// 	// Get the user in the database
// 	userdata, err := cfg.db.GetUser(r.Context(), cfg.CurrentUserEmail)
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid input", err)
// 		return
// 	}

// }
