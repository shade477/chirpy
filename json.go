package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	// log the error if any
	if err != nil {
		log.Println(err)
	}

	// log the response message
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	// Struct to store the json response
	type errorResponse struct {
		Error string `json:"error"`
	}

	// call the function to initate response
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	// Set Response Header
	w.Header().Set("Content-Type", "application/json")

	// Marshall the payload
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}