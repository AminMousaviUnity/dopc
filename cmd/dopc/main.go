package main

import (
	"log"
	"net/http"

	"github.com/yourname/dopc/internal/clients"
	"github.com/yourname/dopc/internal/handlers"
	"github.com/yourname/dopc/internal/services"
)

func main() {
	apiClient := clients.NewHomeAssignmentAPIClient(nil)

	docsService := services.NewDOPCService(apiClient)

	docsHandler := handlers.NewDeliveryOrderPriceHandler(docsService)

	http.HandleFunc("/api/v1/delivery-order-price", docsHandler.HandleGetDeliveryOrderPrice)

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
