package main

import (
	"./common"
	"./controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetUserRoutes(router *mux.Router) *mux.Router {

	userRouter := mux.NewRouter()
	userRouter.HandleFunc("/user/home", homeHandler).Methods("GET")
	userRouter.HandleFunc("/user/history", historyHandler).Methods("GET")
	userRouter.HandleFunc("/user/profile", profileHandler).Methods("GET")
	userRouter.HandleFunc("/user/terms/{id}", userTermsHandler).Methods("GET")
	userRouter.HandleFunc("/user/decision/{id}", userDecisionHandler).Methods("GET")
	userRouter.HandleFunc("/user/offer/{id}/charge", controllers.ChargeCustomerForOffer).Methods("POST")
	userRouter.HandleFunc("/user/glass", userGlassHandler).Methods("GET")
	userRouter.HandleFunc("/user/glass", userPostGlassHandler).Methods("POST")
	router.PathPrefix("/user").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(userRouter),
	))
	return router
}
