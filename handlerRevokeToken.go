package main

import (
	"net/http"

	"github.com/Janisgee/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to get refresh token from 'Authorization'", err)
		return
	}
	// Update the time of (updated_at and revoked_at)
	_, err = cfg.db.UpdateRefreshTokenByUser(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}
	// response code 204 No Content - success but no body
	w.WriteHeader(http.StatusNoContent)
}
