package routes

import (
	"gocommerce/handlers"

	"github.com/gorilla/mux"
)

func RegisterOrderRoutes(router *mux.Router, orderHandler handlers.OrderHandler) {
	// Setting up routes for orders
	router.HandleFunc("/api/orders", orderHandler.GetOrders).Methods("GET")           // Endpoint to get all orders
	router.HandleFunc("/api/orders/{id}", orderHandler.GetOrder).Methods("GET")       // Endpoint to get a specific order by ID
	router.HandleFunc("/api/orders", orderHandler.CreateOrder).Methods("POST")        // Endpoint to create a new order
	router.HandleFunc("/api/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")    // Endpoint to update an existing order by ID
	router.HandleFunc("/api/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE") // Endpoint to delete an order by ID
}
