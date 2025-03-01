package main

import (
	"Assignment1/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Fetching the Port:
	port := os.Getenv("PORT")

	// Override port if there is no port specified:
	if port == "" {
		log.Println("PORT not set, using default 8080")
		port = "8080"
	}

	// Initializing the router:
	router := http.NewServeMux()

	// Routing .../
	router.HandleFunc(handlers.EMPTY, handlers.EmptyHandler)

	// Routing .../info/
	router.HandleFunc(handlers.INFO, handlers.InfoHandler)

	// Routing .../population/
	router.HandleFunc(handlers.POPULATION, handlers.PopulationHandler)

	// Routing .../status
	router.HandleFunc(handlers.STATUS, handlers.StatusHandler)

	// Starting server:
	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
