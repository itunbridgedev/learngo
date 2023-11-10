// handlers/index.go

package handlers

import (
	"database/sql"
)

type Handlers struct {
	ProductHandler        ProductHandler
	UserHandler           UserHandler
	OrderHandler          OrderHandler
	ShoppingCartHandler   ShoppingCartHandler
	AuthenticationHandler AuthenticationHandler
	// Add other handlers as needed
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{
		ProductHandler:        ProductHandler{DB: db},
		UserHandler:           UserHandler{DB: db},
		OrderHandler:          OrderHandler{DB: db},
		ShoppingCartHandler:   ShoppingCartHandler{DB: db},
		AuthenticationHandler: AuthenticationHandler{DB: db},
	}
}
