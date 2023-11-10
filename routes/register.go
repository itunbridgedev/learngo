// routes/register.go

package routes

import (
	"gocommerce/handlers"

	"github.com/gorilla/mux"
)

func RegisterAll(router *mux.Router, handlers *handlers.Handlers) {
	RegisterProductRoutes(router, handlers.ProductHandler)
	RegisterUserRoutes(router, handlers.UserHandler)
	RegisterOrderRoutes(router, handlers.OrderHandler)
	RegisterShoppingCartRoutes(router, handlers.ShoppingCartHandler)
	RegisterAuthenticationRoutes(router, handlers.AuthenticationHandler)
}
