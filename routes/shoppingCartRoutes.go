package routes

import (
	"gocommerce/handlers"

	"github.com/gorilla/mux"
)

func RegisterShoppingCartRoutes(router *mux.Router, cartHandler handlers.ShoppingCartHandler) {
	// Setting up routes for shopping cart
	router.HandleFunc("/api/cart", cartHandler.GetCart).Methods("GET")                          // Retrieves the current user's shopping cart
	router.HandleFunc("/api/cart/items", cartHandler.AddItemToCart).Methods("POST")             // Adds an item to the shopping cart
	router.HandleFunc("/api/cart/test", cartHandler.TestCartRoute).Methods("POST")              // Test route for debugging
	router.HandleFunc("/api/cart/items/{id}", cartHandler.UpdateCartItem).Methods("PUT")        // Updates the quantity of an item in the cart
	router.HandleFunc("/api/cart/items/{id}", cartHandler.RemoveItemFromCart).Methods("DELETE") // Removes an item from the shopping cart
}
