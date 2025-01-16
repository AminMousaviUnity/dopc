package helpers

func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) int {
	...
}

func calculateDeliveryFee(distance int, basePrice int, ranges []models.DistanceRange) (int, error) {
	...
}