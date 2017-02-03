package routers

import (
	"../common"
	"../controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetStripeRoutes(router *mux.Router) *mux.Router {

	stripeRouter := mux.NewRouter()
	stripeRouter.HandleFunc("/customer", controllers.GetCustomerForUser).Methods("GET")
	stripeRouter.HandleFunc("/customer/sources", controllers.AddSourceToCustomer).Methods("POST")
	stripeRouter.HandleFunc("customer/default_source", controllers.SetDefaultSource).Methods("POST")
	stripeRouter.HandleFunc("/customer/charge", controllers.ChargeCustomer).Methods("POST")

	router.PathPrefix("/customer").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(stripeRouter),
	))
	return router
}
