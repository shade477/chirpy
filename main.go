package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const fileroot = "."
	const port = "8080"

	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	var apiCfg apiConfig

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileroot)))))
	mux.HandleFunc("/healthz", Handler)
	mux.HandleFunc("/metrics", apiCfg.Handler)
	mux.HandleFunc("/reset", apiCfg.ResetHandler)

	log.Printf("Server running on port: %s", port)
	log.Fatal(server.ListenAndServe())
}

func Handler(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) Handler(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}