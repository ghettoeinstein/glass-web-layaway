package main

import (
	"./common"
	"./controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetAdminRoutes(router *mux.Router) *mux.Router {

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin/chat", chatHandler).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}", controllers.AdminDisplayOrder).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}/process", AdminProcessOrder).Methods("POST")
	adminRouter.HandleFunc("/admin/orders/{id}/confirmation", controllers.SendConfirmationForOrder).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}/delete", controllers.AdminDeleteOrder).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/decision/approved", controllers.AdminGetApprovedOrders).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/decision/denied", controllers.AdminGetDeniedOrders).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}/edit", controllers.AdminGetEditOrder).Methods("GET")
	adminRouter.HandleFunc("/admin", adminHandler)
	adminRouter.HandleFunc("/admin/orders/new", adminHandler).Methods("GET")
	adminRouter.HandleFunc("/admin/register", adminHandler)
	adminRouter.Handle("/admin/room", globalRoom)
	router.PathPrefix("/admin").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(adminRouter),
	))
	return router
}
