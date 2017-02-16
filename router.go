package main

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	// Routes for the User entity

	router = SetPublicRoutes(router)
	router = SetAdminRoutes(router)
	//Routes for authentication
	router = SetAuthRoutes(router)
	// Routes for logged-in user.
	router = SetUserRoutes(router)
	// Routes for Stripe Integration
	//router = SetStripeRoutes(router)
	return router

}
