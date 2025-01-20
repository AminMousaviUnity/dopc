package services

import (
	"errors"
	"math"

	"github.com/AminMousaviUnity/dopc/internal/models"
	"github.com/AminMousaviUnity/dopc/internal/clients"
)

// DOPCService is the contract for our business logic.
// It has one main method that calculates the delivery order price.
type DOPCService interface {
	CalculatePrice(
		venueSlug string,
		cartValue int,
		userLat float64,
		userLon float64,
	) (models.DeliveryPriceResponse, error)
}

type dopcService struct {
	apiClient clients.HomeAssignmentAPIClient
}

// NewDOPCService is a constructor that returns a DOPCService interface.
// We inject the HomeAssignmentAPIClient so this service can fetch venue data.
func NewDOPCService(apiClient clients.HomeAssignmentAPIClient) DOPCService {
	return &dopcService{
		apiClient: apiClient,
	}
}

func (s *dopcService) CalculatePrice(
	venueSlug string,
	cartValue int,
	userLat float64,
	userLon float64,
) (models.DeliveryPriceResponse, error) {
	// 1. Fetch the static data (coordinates)
	staticResp, err := s.apiClient.GetVenueStatic(venueSlug)
	if err != nil {
		// Return the error to the handler
		return models.DeliveryPriceResponse{}, err
	}

	// 2. Fetch the dynamic data (delivery fees, surcharges)
	dynamicResp, err := s.apiClient.GetVenueDynamic(venueSlug)
	if err != nil {
		return models.DeliveryPriceResponse{}, err
	}

	// 3. Extract needed fields
	venueLon := staticResp.VenueRaw.Location.Coordinates[0]
	venueLat := staticResp.VenueRaw.Location.Coordinates[1]

	orderMinNoSurcharge := dynamicResp.VenueRaw.DeliverySpecs.OrderMinimumNoSrucharge
	basePrice := dynamicResp.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice
	distanceRanges := dynamicResp.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges

	// 4. Calculate the distance between user & venue
	distanceMeters := calculateHaversineDistance(userLat, userLon, venueLat, venueLon)

	// 5. Determine if delivery is possible and compute the delivery fee
	deliveryFee, err := calculateDeliveryFee(distanceMeters, basePrice, distanceRanges)
	if err != nil {
		// e.g. if distance is too far
		return models.DeliveryPriceResponse{}, err
	}

	// 6. Compute small order surcharge
	smallOrderSurcharge := 0
	if cartValue < orderMinNoSurcharge {
		smallOrderSurcharge = orderMinNoSurcharge - cartValue
	}

	// 7. Sum up total price
	totalPrice := cartValue + smallOrderSurcharge + deliveryFee

	// 8. Return the models
	return models.DeliveryPriceResponse{
		TotalPrice:          totalPrice,
		SmallOrderSurcharge: smallOrderSurcharge,
		CartValue:           cartValue,
		Delivery: models.DeliveryDetails{
			Fee:      deliveryFee,
			Distance: distanceMeters,
		},
	}, nil
}

// calculateHaversineDistance returns the great-circle distance (in meters)
// between two lat/lon points using the Haversine formula.
func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) int {
    // Earth’s radius in meters
    const earthRadius = 6371000.0

    // Convert degrees to radians
    dLat := (lat2 - lat1) * math.Pi / 180.0
    dLon := (lon2 - lon1) * math.Pi / 180.0

    rLat1 := lat1 * math.Pi / 180.0
    rLat2 := lat2 * math.Pi / 180.0

    // Haversine formula
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    // Distance in meters
    distance := earthRadius * c

    // Round to nearest integer
    return int(math.Round(distance))
}

// calculateDeliveryFee uses distanceRanges to find the correct range for `distance`
// and calculates: basePrice + a + round(b * distance / 10).
// If the distance is beyond the max allowed range, returns an error.
func calculateDeliveryFee(distance int, basePrice int, ranges []models.DistanceRange) (int, error) {
    for _, dr := range ranges {
        // "max": 0 => means no delivery if distance >= dr.Min
        // Otherwise, the distance must fall between [dr.Min, dr.Max).

        if dr.Max == 0 {
            // This indicates that if distance >= dr.Min, it's not deliverable
            if distance >= dr.Min {
                return 0, errors.New("delivery not possible (distance too large)")
            }
            // If distance < dr.Min, keep checking next range
        } else {
            if distance >= dr.Min && distance < dr.Max {
                // We found our range
                variableFee := int(math.Round(float64(dr.B) * float64(distance) / 10.0))
                return basePrice + dr.A + variableFee, nil
            }
        }
    }

    // If no range matched, we assume it's out of range
    return 0, errors.New("delivery not possible (no matching distance range)")
}
