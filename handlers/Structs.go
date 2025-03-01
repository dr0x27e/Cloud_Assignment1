package handlers

// Struct for InfoHandler:
type infoResponse struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       struct {
		Png string `json:"png"`
	} `json:"flags"`
	Capital []string `json:"capital"`
	Cities  struct {
		Data []string `json:"data"`
	}
}

// NameInfo Structs for PopulationHandler:
type NameInfo struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
}

type populationResponse struct {
	Mean int
	Data struct {
		Value []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

// Struct for status response:
type statusResponse struct {
	CountriesNow  int
	RestCountries int
	Version       string
	Uptime        int
}
