package routes

import (
	"gocommerce/handlers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, userHandler handlers.UserHandler) {
	// Setting up routes for users
	router.HandleFunc("/api/users/{id}", userHandler.GetUser).Methods("GET")       // To get a specific user by ID
	router.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")        // To create a new user
	router.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")    // To update an existing user by ID
	router.HandleFunc("/api/users/{id}", userHandler.DeleteUser).Methods("DELETE") // To delete a user by ID
}
