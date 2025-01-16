package services

import (
	"github.com/AminMousaviUnity/DOPC/internal/models"
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
	apiClient HomeAssignmentAPIClient
}

// NewDOPCService is a constructor that returns a DOPCService interface.
// We inject the HomeAssignmentAPIClient so this service can fetch venue data.
func NewDOPCService(apiClient HomeAssignmentAPIClient) DOPCService {
	return &dopcService{
		apiClient: apiClient,
	}
}

func (s *dopsService) CalculatePrice(
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
	distanceRanges := dynamicResp.VenueRaw.DeliverySpecs.DeliveryPricing.distanceRanges

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
