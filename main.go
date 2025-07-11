package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	// _ means that it is imported but only for passive use
	_ "github.com/lib/pq"
	"github.com/shade477/servers/internal/database"
)

func main() {
	//import .env file
	godotenv.Load()
	// load the db string
	dbURL := os.Getenv("DB_URL")

	//Connect to DB
	db,err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error: Unable to connect to db")
	}

	dbQueries := database.New(db)

	// Assign constants for file path and port.
	// fileroot is the directory that will be used to serve static files.
	const fileroot = "."
	const port = "8080"

	// Create a new HTTP multiplexer (router) to handle different URL paths.
	mux := http.NewServeMux()

	
	// Create and configure the HTTP server.
	server := &http.Server{
		Addr:    ":" + port, // The port the server will listen on (e.g., ":8080")
		Handler: mux,        // The request multiplexer (router) to use
	}


	// Declare a variable of type apiConfig.
	// This will likely hold handler methods and middleware logic for your API.
	var apiCfg apiConfig

	apiCfg.db = dbQueries
	apiCfg.platform = os.Getenv("PLATFORM")

	// Serve static files under the "/app/" route.
	// - Strip "/app/" from the request path so files are resolved correctly.
	// - Serve files from the 'fileroot' directory.
	// - Wrap the file server with a middleware to track request metrics.
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileroot)))))

	// Handle GET request for "/admin/metrics".
	// This could return server metrics like request counts, etc.
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricHandler)

	// Handle GET request for "/api/healthz".
	// This is a health check endpoint to verify the server is running.
	mux.HandleFunc("GET /api/healthz", HealthHandler)

	// Handle POST request for "/admin/reset".
	// This will reset internal server state, metrics, etc.
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)


	// Handle POST request to "/api/validate_chirp"
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirpy)

	// Handle POST request to "/api/users"
	mux.HandleFunc("POST /api/users", apiCfg.createUser)

	// Print a log message showing the port the server is running on.
	log.Printf("Server running on port: %s", port)
	// Start the server and log any fatal errors if it crashes or fails to start.
	log.Fatal(server.ListenAndServe())
}

