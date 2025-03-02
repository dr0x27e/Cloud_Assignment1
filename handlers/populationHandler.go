package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
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

	// To fetch the Population we first need to get the Country name.
	// We fetch this from http://129.241.150.113:8080/v3.1/alpha/
	// Using the Iso2 code.

	// Fetching the iso2 value.
	iso2 := r.PathValue("country_code")

	// URL to invoke:
	url := RESTCountriesApi + "alpha/" + iso2 + "?fields=name"

	// Getting the response:
	resp, err := http.Get(url)
	if ServerError(err) {
		return
	}
	defer resp.Body.Close()

	// Read Response:
	body, err := io.ReadAll(resp.Body)
	if ServerError(err) {
		return
	}

	// Fetch the country name:
	var response NameInfo
	err = json.Unmarshal(body, &response)
	if ServerError(err) {
		return
	}

	countryName := response.Name.Common

	// Now we need to get the country and its population data.

	// Url:
	url = CountriesNowApi + "countries/population"

	// Creating payload:
	payload := strings.NewReader(fmt.Sprintf(`{"country": "%s"}`, countryName))

	// Creating request:
	req, err := http.NewRequest("POST", url, payload)
	if ServerError(err) {
		return
	}

	// Setting header type:
	req.Header.Set("Content-Type", "application/json")

	// Creating client:
	client := &http.Client{}

	// Issue Request:
	resp, cliErr := client.Do(req)
	if ServerError(cliErr) {
		return
	}
	defer resp.Body.Close()

	// Read Response:
	body, err = io.ReadAll(resp.Body)
	if ServerError(err) {
		return
	}

	// Get the query parameter 'limit':
	limitString := r.URL.Query().Get("limit")

	var finalResponse populationResponse

	if limitString != "" {
		// Split the string into 2:
		parts := strings.Split(limitString, "-")
		// Make sure its only 2 ints.
		if len(parts) != 2 {
			fmt.Println("Invalid range format in query limit, parts: ", parts)
			fmt.Println("Max 2 parts allowed")
			return
		}

		// Fetching the start and end year:
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil && err2 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}

		// Fetching the JSON body in an unfiltered Struct
		var unfilteredResponse populationResponse
		err = json.Unmarshal(body, &unfilteredResponse)
		if ServerError(err) {
			return
		}

		// Creating a mean
		mean := 0

		// Going through all the unfiltered data and only appending the right data into the final response.
		for _, pc := range unfilteredResponse.Data.Value {
			if pc.Year >= start && pc.Year <= end {
				finalResponse.Data.Value = append(finalResponse.Data.Value, pc)
				mean += pc.Value
			}
		}

		// Setting the mean.
		finalResponse.Mean = mean / len(finalResponse.Data.Value)

	} else {
		// Fetching the JSON body into finalResponse
		err = json.Unmarshal(body, &finalResponse)
		if ServerError(err) {
			return
		}

		for _, pc := range finalResponse.Data.Value {
			finalResponse.Mean += pc.Value
		}

		finalResponse.Mean = finalResponse.Mean / len(finalResponse.Data.Value)
	}

	// Formatting struct:
	type formatted struct {
		Mean   int `json:"mean"`
		Values []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"values"`
	}

	// Formatting so that it looks like in task.
	final := formatted{
		Mean:   finalResponse.Mean,
		Values: finalResponse.Data.Value,
	}

	// Return the combined result as a JSON response:
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(final); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
