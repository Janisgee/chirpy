package main

import (
	"net/http"
)

// Handler that count reset
func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not allowed to be accessed with production environment", nil)
		return
	}

	// Delete the user in the database
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't delete all users", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respondWithJSON(w, http.StatusOK, struct{}{})

}
