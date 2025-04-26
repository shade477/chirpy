package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/shade477/servers/internal/database"
)

// struct to store all runtime api configs
type apiConfig struct {
	fileserverHits atomic.Int32
	db	*database.Queries
	platform	string
}

// to create a user
func(cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	// to Stores the request parameters
	type parameters struct {
		Email string `json:"email"`
	}

	// to decode the request body
	decoder := json.NewDecoder(r.Body)

	// Create instance of the parameter struct
	params := parameters{}

	// Decode and store into the struct
	err := decoder.Decode(&params)
	if err != nil {
		log.Default()
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request", err)
		return
	}

	// Execute and store the db query
	dbUser,err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create User", err)
		return
	}

	// Map the db response to the user struct to Marshal as response
	user := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func(cfg *apiConfig) MetricHandler(w http.ResponseWriter, r *http.Request) {
	// Store the template
	template := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`
	// Set the Header
	w.Header().Set("Content-Type", "text/html")
	// Set the response code
	w.WriteHeader(http.StatusOK)
	// Set body
	w.Write([]byte(fmt.Sprintf(template, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	// Add a middle functon to keep track number of hits
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) Handler(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

