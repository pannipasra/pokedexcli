package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is resiposible for making request to PokeAPI
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
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
