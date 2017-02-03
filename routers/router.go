package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	// Routes for the User entity
	// router = SetUserRoutes(router)
	// Routes for the Product entity
	// Routes for Auth checking
	//  router = SetAuthRoutes(router)
	router = SetOrderRoutes(router)
	// Routes for Stripe Integration
	router = SetStripeRoutes(router)
	return router

}
