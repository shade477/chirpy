package main

import (
	"log"
	"net/http"
)

func main() {
	//Assgin const for file path and port 
	const fileroot = "."
	const port = "8080"

	//Create HTTP multiplexer
	mux := http.NewServeMux()

	//create http server struct
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	
	var apiCfg apiConfig

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileroot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricHandler)
	mux.HandleFunc("GET /api/healthz", HealthHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)

	log.Printf("Server running on port: %s", port)
	log.Fatal(server.ListenAndServe())
}

