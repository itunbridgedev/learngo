package models

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        int `json:"id"`         // Unique identifier for the cart item
	ProductID int `json:"product_id"` // Identifier for the product
	UserID    int `json:"user_id"`    // Identifier for the user who owns the cart
	Quantity  int `json:"quantity"`   // Quantity of the product in the cart
	// Add other relevant fields if necessary, such as price
}
