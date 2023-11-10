package routes

import (
	"gocommerce/handlers"

	"github.com/gorilla/mux"
)

func RegisterAuthenticationRoutes(router *mux.Router, authHandler handlers.AuthenticationHandler) {
	// Setting up authentication routes
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")          // Endpoint for user login
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")    // Endpoint for user registration
	router.HandleFunc("/api/auth/refresh", authHandler.RefreshToken).Methods("POST") // Endpoint for refreshing the authentication token
}
