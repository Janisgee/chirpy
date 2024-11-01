package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Janisgee/chirpy.git/internal/auth"
	"github.com/Janisgee/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {

	params := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}

	// Decodes the request JSON into Go struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil || params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Creates the user in the database
	userdata, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPw,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	// Maps the database user to your API User struct
	user := User{
		ID:          userdata.ID,
		CreatedAt:   userdata.CreatedAt,
		UpdatedAt:   userdata.UpdatedAt,
		Email:       userdata.Email,
		IsChirpyRed: userdata.IsChirpyRed,
	}

	// Returns the correct 201 status code
	respondWithJSON(w, http.StatusCreated, user)

}
