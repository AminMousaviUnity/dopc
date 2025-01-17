package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AminMousaviUnity/dopc/internal/services"
)

// DeliveryOrderPriceHandler holds the reference to the service we call.
type DeliveryOrderPriceHandler struct {
	dopcService services.DOPCService
}

// NewDeliveryOrderPriceHandler is a constructor for this handler,
// injecting the DOPCService
func NewDeliveryOrderPriceHandler(svc services.DOPCService) *DeliveryOrderPriceHandler {
	return &DeliveryOrderPriceHandler{
		dopcService: svc,
	}
}

// HandlerGetDeliveryOrderPrice processes GET /api/v1/delivery-order-price
func (h *DeliveryOrderPriceHandler) HandleGetDeliveryOrderPrice(w http.ResponseWriter, r *http.Request) {
	// 1. Parse query params
	query := r.URL.Query()

	venueSlug := query.Get("venue_slug")
	cartValueStr := query.Get("cart_value")
	userLatStr := query.Get("user_lat")
	userLonStr := query.Get("user_lon")

	// 2. Validate them
	if venueSlug == "" {
		http.Error(w, "missing venue_slug", http.StatusBadRequest)
		return
	}
	if cartValueStr == "" || userLatStr == "" || userLonStr == "" {
		http.Error(w, "missing one of the required params (cart_value, user_lat, user_lon)", http.StatusBadRequest)
		return
	}

	cartValue, err := strconv.Atoi(cartValueStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid cart_value: %v", err), http.StatusBadRequest)
		return
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user_lat: %v", err), http.StatusBadRequest)
		return
	}

	userLon, err := strconv.ParseFloat(userLonStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user_lon: %v", err), http.StatusBadRequest)
		return
	}

	// 3. Call the service
	resp, err := h.dopcService.CalculatePrice(
		venueSlug,
		cartValue,
		userLat,
		userLon,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("delivey isn't possible"), http.StatusBadRequest)
	}

	// 4. Return JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}
