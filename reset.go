package main

import "net/http"

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	// platform check
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	
	// Reset the hit counter
	cfg.fileserverHits.Store(0)

    // Delete all users
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete all users", err)
	}

	//Write Response
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
	w.WriteHeader(http.StatusOK)
}