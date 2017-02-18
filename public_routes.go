package main

import (
	"./controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func SetPublicRoutes(router *mux.Router) *mux.Router {

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/about-us", aboutUsHandler).Methods("GET")
	router.HandleFunc("/sorry", sorryHandler).Methods("GET")
	router.HandleFunc("/privacy-policy", privacyPolicyHandler).Methods("GET")
	router.HandleFunc("/tos", tosHandler).Methods("GET")
	router.HandleFunc("/glass", glassHandler).Methods("GET")
	router.HandleFunc("/glass", postGlassHandler).Methods("POST")
	router.HandleFunc("/congratulations", congratulationsHandler).Methods("GET")
	router.HandleFunc("/terms/{id}", termsHandler).Methods("GET")
	router.HandleFunc("/decision/{id}", decisionHandler).Methods("GET")
	router.HandleFunc("/offer/{id}/charge", controllers.ChargeNewCustomerForOffer).Methods("POST")
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/logout", userLogout).Methods("GET")
	router.HandleFunc("/submit-email", emailHandler).Methods("POST")
	router.HandleFunc("/user/register", controllers.Register).Methods("POST")
	router.HandleFunc("/confirmation/{uuid}}", GetPaymentConfirmation).Methods("GET")

	router.HandleFunc("/about", aboutUsHandler)

	return router
}
