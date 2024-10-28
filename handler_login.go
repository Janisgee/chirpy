package main

import (
	"encoding/json"
	"net/http"

	"github.com/Janisgee/chirpy.git/internal/auth"
)

// This endpoint will be used to give the user a token that they can use to make authenticated requests.
func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
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

	// Get the user in the database
	userdata, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Check to see if password matches the stored hash password string
	err = auth.CheckPasswordHash(params.Password, userdata.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Maps the database user to your API User struct
	user := User{
		ID:        userdata.ID,
		CreatedAt: userdata.CreatedAt,
		UpdatedAt: userdata.UpdatedAt,
		Email:     userdata.Email,
	}

	// Returns the correct 204 status code
	respondWithJSON(w, http.StatusOK, user)

}
