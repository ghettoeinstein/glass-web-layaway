package main

import (
	"./controllers"
	"github.com/gorilla/mux"
)

func SetAuthRoutes(router *mux.Router) *mux.Router {
	//Routes that deal with login & authentication.
	router.HandleFunc("/login", SMSLogin).Methods("GET")
	router.HandleFunc("/login/verify", GETVerifySMSLogin).Methods("GET")
	router.HandleFunc("/login/verify", POSTVerifySMSLogin).Methods("POST")
	router.HandleFunc("/team/login", controllers.GetLogin).Methods("GET")
	router.HandleFunc("/team/login", controllers.AdminLogin).Methods("POST")
	return router
}
