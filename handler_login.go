package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Janisgee/chirpy.git/internal/auth"
	"github.com/Janisgee/chirpy.git/internal/database"
)

// This endpoint will be used to give the user a token that they can use to make authenticated requests.
func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}

	type response struct {
		User
		Token         string `json:"token"`
		Refresh_token string `json:"refresh_token"`
	}

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

	// Create JSON Web Token
	accessToken, err := auth.MakeJWT(userdata.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	// Create Refresh Token
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token", err)
		return
	}

	//Store refresh token in the database
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:  refresh_token,
		UserID: userdata.ID,
	}
	refreshTokenData, err := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token database", err)
		return
	}

	// Maps the database user to your API User struct
	userResponse := response{
		User: User{
			ID:          userdata.ID,
			CreatedAt:   userdata.CreatedAt,
			UpdatedAt:   userdata.UpdatedAt,
			Email:       userdata.Email,
			IsChirpyRed: userdata.IsChirpyRed,
		},
		Token:         accessToken,
		Refresh_token: refreshTokenData.Token,
	}

	// Returns the correct 200 status code
	respondWithJSON(w, http.StatusOK, userResponse)

}
