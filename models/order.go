// models/order.go

package models

import (
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
