package models

type DeliveryPriceResponse struct {
	TotalPrice          int             `json:"total_price"`
	SmallOrderSurcharge int             `json:"small_order_surcharge"`
	CartValue           int             `json:"cart_value"`
	Delivery            DeliveryDetails `json:"delivery"`
}

type DeliveryDetails struct {
	Fee      int `json:"fee"`
	Distance int `json:"distance"`
}

type VenueStaticResponse struct {
	VenueRaw struct {
		Location struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"location"`
	} `json:"venue_raw"`
}

type VenueDynamicResponse struct {
	VenueRaw struct {
		DeliverySpecs struct {
			OrderMinimumNoSrucharge int             `json:"order_minimum_no_surcharge"`
			DeliveryPricing         DeliveryPricing `json:"delivery_pricing"`
		} `json:"delivery_specs"`
	} `json:"venue_raw"`
}

type DeliveryPricing struct {
	BasePrice      int             `json:"base_price"`
	DistanceRanges []DistanceRange `json:"distance_ranges"`
}

type DistanceRange struct {
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	A    int    `json:"a"`
	B    int    `json:"b"`
	Flag string `json:"flag"`
}
