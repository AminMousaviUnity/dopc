package handlers

import (
	"encoding/json"
	"net/http"
)

func GetServiceDeliveryPrice(w http.ResponseWriter, r *http.Request) {
	var DOPC models.DOPCService
	if err := json.NewDecoder(r.Body).Decode(&task)
}
