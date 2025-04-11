package pokeapi

// LocationAreaResp represents the response from the location-area endpoint
type LocationAreaResp struct {
	Count    int                  `json:"count"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []LocationAreaResult `json:"results"`
}

// LocationAreaResult represents a single location area in the response
type LocationAreaResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
