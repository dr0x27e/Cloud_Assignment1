package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StartTime Application start time
var StartTime = time.Now()

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Only GET method is supported:
	if r.Method != "GET" {
		fmt.Println("Method not supported: " + r.Method)
		return
	}

	// Creating a client with a timout:
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()

	// CountriesNow
	statusCountries, errCountries := client.Get(CountriesNowApi + "countries")
	if errCountries != nil {
		log.Println("CountriesNow", errCountries)
	}

	// REST Countries
	statusRest, errRest := client.Get(RESTCountriesApi + "all?fields=name,flags")
	if errRest != nil {
		log.Println("RestCountries", errRest)
	}

	// Formatting struct:
	response := statusResponse{
		CountriesNow:  statusCountries.StatusCode,
		RestCountries: statusRest.StatusCode,
		Version:       VERSION,
		Uptime:        int(time.Since(StartTime).Round(time.Second).Seconds()),
	}

	// Closing bodies.
	defer statusRest.Body.Close()
	defer statusCountries.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	// Send the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
