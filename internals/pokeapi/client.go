package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pannipasra/pokedexcli/internals/pokecache"
)

// Client is resiposible for making request to PokeAPI
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Cache      *pokecache.Cache
}

// Config stores current pagination state
type Config struct {
	Next     *string
	Previous *string
}

// NewClient create a new PokeAPI client
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://pokeapi.co/api/v2",
		HTTPClient: &http.Client{},
		Cache:      pokecache.NewCache(5 * time.Minute),
	}
}

// =====================================================
// =====================================================

// ListLocationAreas retrieves the list of location areas
func (c *Client) ListLocationAreas(config *Config) (*LocationAreaResp, error) {
	url := fmt.Sprintf("%s/location-area", c.BaseURL)

	if config.Next != nil {
		url = *config.Next
	}

	// Check if we have this URL cached
	if cachedData, found := c.Cache.Get(url); found {
		// Use cached data
		var locationAreaResp LocationAreaResp
		err := json.Unmarshal(cachedData, &locationAreaResp)
		if err != nil {
			// Update config with pagination links
			config.Next = locationAreaResp.Next
			config.Previous = locationAreaResp.Previous
			return &locationAreaResp, nil
		}
	}

	// Make HTTP request
	res, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add to cache
	c.Cache.Add(url, body)

	// Parse JSON
	var locationArea LocationAreaResp
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return nil, err
	}

	// Update config
	config.Next = locationArea.Next
	config.Previous = locationArea.Previous

	return &locationArea, nil
}

// ListPreviousLocationAreas retrieves the previous list of location areas
func (c *Client) ListPreviousLocationAreas(config *Config) (*LocationAreaResp, error) {
	url := fmt.Sprintf("%s/location-area", c.BaseURL)

	if config.Previous != nil {
		url = *config.Previous
	}

	// Check if we have this URL in cached
	if cachedData, found := c.Cache.Get(url); found {
		// Use cached data
		var locationAreaResp LocationAreaResp
		err := json.Unmarshal(cachedData, &locationAreaResp)
		if err != nil {
			return nil, err
		}
		return &locationAreaResp, nil
	}

	// Make HTTP request
	res, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add to cache
	c.Cache.Add(url, body)

	// Parse JSON
	var locationArea LocationAreaResp
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return nil, err
	}

	// Update config
	config.Next = locationArea.Next
	config.Previous = locationArea.Previous

	return &locationArea, nil
}

func (c *Client) Explore(locationName string) (*ExploreAreaEncounter, error) {
	url := fmt.Sprintf("%s/location-area/%s", c.BaseURL, locationName)

	// Make a request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read a byte
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse to JSON
	var exploreEncouter ExploreAreaEncounter
	err = json.Unmarshal(body, &exploreEncouter)
	if err != nil {
		return nil, err
	}

	return &exploreEncouter, nil
}
