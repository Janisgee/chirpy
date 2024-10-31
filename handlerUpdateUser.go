package main

import (
	"encoding/json"
	"net/http"

	"github.com/Janisgee/chirpy.git/internal/auth"
	"github.com/Janisgee/chirpy.git/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	params := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}

	// Get JWT from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to find token from 'Authorization'", err)
		return
	}

	//Validate JWT from header
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get validation from provided JWT", err)
		return
	}

	// Decodes the request JSON into Go struct
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil || params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// Hashed the new password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Update user info in databasse
	userInfo, err := cfg.db.UpdateUserInfo(r.Context(), database.UpdateUserInfoParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user information", err)
		return
	}

	// Maps the database user to your API User struct
	user := User{
		ID:        userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email:     userInfo.Email,
	}

	// Returns the correct 201 status code
	respondWithJSON(w, http.StatusOK, user)

}
