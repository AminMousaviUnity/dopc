package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AminMousaviUnity/dopc/internal/models"
)

// HomeAssignmentAPIClient defines the contract for fetching
// venue data from the Home Assignment API.
type HomeAssignmentAPIClient interface {
	// Fetches static info (coordinates) for a given venue.
	GetVenueStatic(venueSlug string) (models.VenueStaticResponse, error)

	// Fetches dynamic info (delivery rules) for a given venue.
	GetVenueDynamic(venueSlug string) (models.VenueDynamicResponse, error)
}

// homeAssignmentAPIClient is the concrete implementation of HomeAssignmentClientAPI.
// It holds an http.Client (for making requests) and a baseURL (the root endpoint)
type homeAssignmentAPIClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewHomeAssignmentAPIClient creates a new client instance for talking to the Home Assignment API.
// - If no http.Client is provided, we default to a new one with a small timeout.
func NewHomeAssignmentAPIClient(client *http.Client) HomeAssignmentAPIClient {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	return &homeAssignmentAPIClient{
		httpClient: client,
		// This baseURL should be the root endpoint for your assignment API, e.g.:
		// https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues
		baseURL: "https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues",
	}
}

// GetVenueStatic fetches static data from the /static endpoint of the Home Assignment API.
func (c *homeAssignmentAPIClient) GetVenueStatic(venueSlug string) (models.VenueStaticResponse, error) {
	var result models.VenueStaticResponse

	// Build the URL: e.g. https://.../v1/venues/<slug>/static
	url := fmt.Sprintf("%s/%s/static", c.baseURL, venueSlug)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return result, fmt.Errorf("failed to call static endpoint: %w", err)
	}
	defer resp.Body.Close()

	// Expecting a 200 OK. If not, treat it as an error.
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("static endpoint returned status: %d", resp.StatusCode)
	}

	// Decode the JSON response into our models.VenueStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode static JSON: %w", err)
	}

	return result, nil
}

// GetVenueDynamic fetches dynamic data from the /dynamic endpoint of the Home Assignment APi.
func (c *homeAssignmentAPIClient) GetVenueDynamic(venueSlug string) (models.VenueDynamicResponse, error) {
	var result models.VenueDynamicResponse

	// Build the URL: e.g. https://.../v1/venues/<slug>/dynamic
	url := fmt.Sprintf("%s/%s/dynamic", c.baseURL, venueSlug)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return result, fmt.Errorf("failed to call dynamic endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("dynamic endpoint returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode dynamic JSON: %w", err)
	}

	return result, nil
}
