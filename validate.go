package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ValidateChirpy(w http.ResponseWriter, r *http.Request) {
	// Create struct to store request json
	type parameters struct {
		Body string `json:"body"`
	}

	// Create struct to store the response json
	type returnVals struct {
		Cleaned string `json:"cleaned_body"`
	}

	// Using decoder to decode the request body
	decoder := json.NewDecoder(r.Body)

	// creating instance of the parameters struct
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameter", err)
		return
	}

	// Create validation rules

	// Handle error if chirp length exceeds 140 characters
	// if length >140 {
		// Set Content-Type header
		// w.Header().Set("Content-Type", "application/json")
		// Set Response Code
		// w.WriteHeader(400)
		// Write struct to json body
		// json.NewEncoder(w).Encode(map[string]string{"error": "Chirp is too long"})
		// return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(200)
	// json.NewEncoder(w).Encode(map[string]bool{"valid": true})
	// return

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned: removeProfanity(params.Body),
	})

}

func removeProfanity(msg string) string {
	profane := map[string]bool{
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	words := strings.Split(msg, " ")
	for i,w:= range words {
		for p := range profane {
			if strings.ToLower(w) == p {
				words[i] = "****"
				break
			}
		}
	}
	puremsg := strings.Join(words, " ")
	return puremsg
}