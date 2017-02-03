package routers

import (
	"../common"
	"../controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetupOrderRoutes(router *mux.Router) *mux.Router {
	orderRouter := mux.NewRouter()
	//Set up routes for application/api

	//Create an order
	orderRouter.HandleFunc("/orders", controller.CreateOrder).Methods("POST")
	// Get a single order
	//orderRouter.HandleFunc("/order/{id}", controller.GetOrder).Methods("GET")

	//	orderRouter.HandleFunc("/order/{id}", controller.UpdateOrder).Methods("PUT")
	//	orderRouter.HandleFunc("/order/{id}", controller.GetOrder).Methods("DELETE")

	router.PathPrefix("/order").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(orderRouter),
	))
	return router

}
