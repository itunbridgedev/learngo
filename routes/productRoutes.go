package routes

import (
	"gocommerce/handlers"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router, productHandler handlers.ProductHandler) {
	// Setting up routes for products
	router.HandleFunc("/api/products", productHandler.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", productHandler.GetProduct).Methods("GET")
	router.HandleFunc("/api/products", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/api/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/api/products/{id}", productHandler.DeleteProduct).Methods("DELETE")
}
