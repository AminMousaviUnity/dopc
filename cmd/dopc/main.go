package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// EchoResponse holds the query parameters we want to echo back.
type EchoResponse struct {
	VenueSlug string  `json:"venue_slug"`
	CartValue int     `json:"cart_value"`
	UserLat   float64 `json:"user_lat"`
	UserLon   float64 `json:"user_lon"`
}

func main() {
	// Register our handler function for the /api/v1/delivery-order-price route.
	http.HandleFunc("/api/v1/delivery-order-price", handleDeliveryOrderPrice)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleDeliveryOrderPrice is a minimal handler that:
// 1. Parses the query parameters: venue_slug, cart_value, user_lat, user_lon
// 2. Validates and converts them to the correct types
// 3. Returns them as JSON or errors if something is invalid
func handleDeliveryOrderPrice(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Extract the parameters as strings
	venueSlug := query.Get("venue_slug")
	cartValueStr := query.Get("cart_value")
	userLatStr := query.Get("user_lat")
	userLonStr := query.Get("user_lon")

	// -- Basic validation & parsing --

	// Check if venue_slug is provided
	if venueSlug == "" {
		http.Error(w, "missing venue_slug", http.StatusBadRequest)
	}

	// Convert cart_value from string to int
	cartValue, err := strconv.Atoi(cartValueStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid cart_value: %v", err), http.StatusBadRequest)
		return
	}

	// Convert user_lat from string to float64
	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user_lat: %v", err), http.StatusBadRequest)
		return
	}

	// Convert user_lon from string to float64
	userLon, err := strconv.ParseFloat(userLonStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user_lon: %v", err), http.StatusBadRequest)
		return
	}

	// -- Build our response (just echoing back what we parsed) --
	resp := EchoResponse{
		VenueSlug: venueSlug,
		CartValue: cartValue,
		UserLat:   userLat,
		UserLon:   userLon,
	}

	// -- Return as JSON --
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)

}
