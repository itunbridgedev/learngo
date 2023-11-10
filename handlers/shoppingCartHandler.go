// shoppingCartHandler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gocommerce/constants"
	"gocommerce/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ShoppingCartHandler struct {
	DB *sql.DB
}

func (h *ShoppingCartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the context set by the AuthenticationMiddleware
	ctx := r.Context()
	userID, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		// Handle the case where the user ID is not set or is of the wrong type
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}

	var cartItems []models.CartItem
	rows, err := h.DB.Query("SELECT id, product_id, user_id, quantity FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error fetching cart items: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.UserID, &item.Quantity); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Printf("Error reading cart item: %v", err)
			return
		}
		cartItems = append(cartItems, item)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error iterating cart items: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
}

func (h *ShoppingCartHandler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the context
	ctx := r.Context()
	userID, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}

	// Decode the request body to get cart item details
	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the cart item
	if err := h.validateCartItemDetails(item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new cart item into the database
	_, err := h.DB.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)", userID, item.ProductID, item.Quantity)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error adding item to cart: %v", err)
		return
	}

	// Send a successful response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart successfully"})
}

// validateCartItemDetails validates the details of a cart item.
func (h *ShoppingCartHandler) validateCartItemDetails(item models.CartItem) error {
	// Check if the product exists
	exists, err := h.productExists(item.ProductID)
	if err != nil {
		return fmt.Errorf("error checking product existence: %w", err)
	}
	if !exists {
		return errors.New("product does not exist")
	}

	// Validate the quantity
	if !isValidQuantity(item.Quantity) {
		return errors.New("invalid quantity")
	}

	// Add any additional validation rules as needed

	return nil
}

// productExists checks if a product with the given ID exists in the database.
func (h *ShoppingCartHandler) productExists(productID int) (bool, error) {
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	return exists, err
}

// isValidQuantity checks if the provided quantity is positive and within acceptable limits.
func isValidQuantity(quantity int) bool {
	return quantity > 0 // Add more complex logic as needed, such as maximum limit checks
}

func (h *ShoppingCartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		// Handle the case where the user ID is not set or is of the wrong type
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("UPDATE cart_items SET quantity = $1 WHERE id = $2 AND user_id = $3", item.Quantity, itemID, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error updating cart item: %v", err)
		return
	}

	// Send a successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart item updated successfully"})
}

func (h *ShoppingCartHandler) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the context
	ctx := r.Context()
	userID, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}

	// Extract the item ID from the URL parameters
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Delete the item from the database
	_, err = h.DB.Exec("DELETE FROM cart_items WHERE id = $1 AND user_id = $2", itemID, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error deleting cart item: %v", err)
		return
	}

	// Send a successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart item removed successfully"})
}
