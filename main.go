package main

import (
	"./common"
	"./controllers"
	"os"
	//	"flag"

	"github.com/gorilla/mux"

	"github.com/urfave/negroni"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	go setupTemplates()
	go func() {
		globalRoom = newRoom()
		globalRoom.run()
	}()

}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Trapped panic: %s (%T) \n", err, err)
		}
	}()
	go controllers.InitTemplates()

	f, err := os.OpenFile("logs/glassLogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error  opening up logfile %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	//http.HandleFunc("/", rootHandler)

	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	//router.HandleFunc("/register", controllers.Register).Methods("POST")
	//router.ServeFiles("/assets/*filepath", http.Dir("/assets"))
	//
	//router.GET("/about-us", aboutUsHandler)
	router.HandleFunc("/about-us", aboutUsHandler).Methods("GET")
	router.HandleFunc("/sorry", sorryHandler).Methods("GET")
	//router.GET("/privacy-policy", privacyPolicyHandler)
	router.HandleFunc("/glass", glassHandler).Methods("GET")
	//router.GET("/terms", termsHandler)
	router.HandleFunc("/glass", postGlassHandler).Methods("POST")
	router.HandleFunc("/getglass", glassHandler).Methods("GET")

	router.HandleFunc("/congratulations", congratulationsHandler).Methods("GET")

	router.HandleFunc("/terms/{id}", termsHandler).Methods("GET")
	router.HandleFunc("/decision/{id}", decisionHandler).Methods("GET")
	router.HandleFunc("/offer/{id}/charge", controllers.ChargeForOffer).Methods("POST")
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/logout", userLogout).Methods("GET")
	router.HandleFunc("/users/register", controllers.Register).Methods("POST")
	router.HandleFunc("/submit-email", emailHandler).Methods("POST")

	userRouter := mux.NewRouter()

	userRouter.HandleFunc("/user/home", homeHandler).Methods("GET")
	userRouter.HandleFunc("/user/history", historyHandler).Methods("GET")
	userRouter.HandleFunc("/user/profile", profileHandler).Methods("GET")

	router.PathPrefix("/user").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(userRouter),
	))

	//router.POST("/orders", controllers.CreateOrder).Methids("POST")

	router.HandleFunc("/about", aboutUsHandler)

	//
	router.HandleFunc("/login", SMSLogin).Methods("GET")

	router.HandleFunc("/login/verify", GETVerifySMSLogin).Methods("GET")
	router.HandleFunc("/login/verify", POSTVerifySMSLogin).Methods("POST")
	router.HandleFunc("/team/login", controllers.GetLogin).Methods("GET")
	router.HandleFunc("/team/login", controllers.AdminLogin).Methods("POST")

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin/chat", chatHandler).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}", controllers.AdminDisplayOrder).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}/process", controllers.AdminProcessOrder).Methods("POST")
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

	router.StrictSlash(true)
	//router.GET("/chat", chatHandler)
	n := negroni.Classic()
	n.UseHandler(router)

	server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: n,
	}

	log.Println("API is Listening on: ", common.AppConfig.Server)
	log.Fatal(server.ListenAndServe())

}
