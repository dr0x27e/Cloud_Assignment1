package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	// Only GET method is supported:
	if r.Method != "GET" {
		fmt.Println("Method not supported: " + r.Method)
		return
	}

	ServerError := func(err error) bool {
		if err != nil {
			fmt.Println(err)
			return true
		}
		return false
	}

	// Fetching the iso2 value.
	iso := r.PathValue("country_code")

	// URL to invoke:
	url := RESTCountriesApi + "alpha/" + iso

	// Get the query parameter 'limit':
	limitString := r.URL.Query().Get("limit")

	// If limit is empty, default to 10:
	limit := 10
	if limitString != "" {
		// convert limit to integer and check if its 0 -> 10:
		limitInt, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Invalid limit amount", http.StatusBadRequest)
			return
		} else if limitInt > 0 && limitInt < 10 {
			limit = limitInt
		}
	}

	// Create a new request:
	req, reqErr := http.NewRequest("GET", url, nil)
	if ServerError(reqErr) {
		return
	}

	// Setting header type content to json
	req.Header.Add("content-type", "application/json")

	// Instantiate the client:
	client := &http.Client{}

	// Issue request.
	response, cliErr := client.Do(req)
	if ServerError(cliErr) {
		return
	}

	// Read Response and turn it go readable:
	body, err := io.ReadAll(response.Body)
	if ServerError(err) {
		return
	}

	// Fetching response:
	var writer []infoResponse
	unMarshalErr := json.Unmarshal(body, &writer)
	if ServerError(unMarshalErr) {
		return
	}

	// Url for second API:
	url = CountriesNowApi + "countries/cities"

	// Creating payload:
	payload := strings.NewReader(fmt.Sprintf(`{"iso2": "%s"}`, iso))

	// Creating new request:
	req, err = http.NewRequest("POST", url, payload)
	if ServerError(err) {
		return
	}

	// Setting header type content to json
	req.Header.Add("content-type", "application/json")

	// Issue Request:
	response, cliErr = client.Do(req)
	if ServerError(cliErr) {
		return
	}

	// Read Response:
	body, err = io.ReadAll(response.Body)
	if ServerError(err) {
		return
	}

	// Fetching response:
	err = json.Unmarshal(body, &writer[0].Cities)
	if ServerError(err) {
		return
	}

	// Formatting struct:
	type formatted struct {
		Name       string            `json:"name"`
		Continents []string          `json:"continents"`
		Population int               `json:"population"`
		Languages  map[string]string `json:"languages"`
		Borders    []string          `json:"borders"`
		Flag       string            `json:"flag"`
		Capital    []string          `json:"capital"`
		Cities     []string          `json:"cities"`
	}

	// Formatting so that it looks like in task and also limiting amount of cities:
	finalResponse := formatted{
		Name:       writer[0].Name.Common,
		Continents: writer[0].Continents,
		Population: writer[0].Population,
		Languages:  writer[0].Languages,
		Borders:    writer[0].Borders,
		Flag:       writer[0].Flag.Png,
		Capital:    writer[0].Capital,
		Cities:     writer[0].Cities.Data[:limit],
	}

	// Return the combined result as a JSON response:
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(finalResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
