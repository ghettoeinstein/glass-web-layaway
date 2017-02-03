package main

import (
	"./common"
	"./controllers"

	//	"flag"
	"fmt"
	"github.com/gorilla/mux"

	"github.com/urfave/negroni"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"time"
)

func (o Order) String() string {
	return fmt.Sprintf(o.UUID)
}

var templates map[string]*template.Template

func init() {

	rand.Seed(time.Now().UTC().UnixNano())

	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	go setupTemplates()

}

func renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	//Ensure the template exists in the `templates` map
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Trapped panic: %s (%T) \n", err, err)
		}
	}()

	go setupLogging()

	//http.HandleFunc("/", rootHandler)

	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/", rootHandler)

	router.HandleFunc("/admin-token", genAdminHandler).Methods("GET")

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	//
	//router.GET("/about-us", aboutUsHandler)
	//router.GET("/sorry", sorryHandler)
	//router.GET("/privacy-policy", privacyPolicyHandler)
	router.HandleFunc("/glass", glassHandler).Methods("GET")
	//router.GET("/terms", termsHandler)
	router.HandleFunc("/glass", postGlassHandler).Methods("POST")
	router.HandleFunc("/getglass", glassHandler).Methods("GET")
	router.HandleFunc("/home", homeHandler).Methods("GET")
	router.HandleFunc("/congratulations", congratulationsHandler).Methods("GET")
	router.HandleFunc("/history", historyHandler).Methods("GET")
	router.HandleFunc("/profile", profileHandler).Methods("GET")
	router.HandleFunc("/terms", termsHandler).Methods("GET")
	router.HandleFunc("/", rootHandler)
	//router.POST("/orders", controllers.CreateOrder).Methids("POST")
	//router.GET("/offer/:uuid", offerHandler)

	//
	router.HandleFunc("/login", controllers.GetLogin).Methods("GET")
	go controllers.InitTemplates()
	adminRouter := mux.NewRouter()

	adminRouter.HandleFunc("/admin/orders/{id}", controllers.AdminDisplayOrder).Methods("GET")
	adminRouter.HandleFunc("/admin", adminHandler)
	adminRouter.HandleFunc("/admin/orders", ordersHandler)
	adminRouter.HandleFunc("/admin/register", adminHandler)

	router.PathPrefix("/admin").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(adminRouter),
	))

	//router.GET("/chat", chatHandler)
	n := negroni.Classic()
	n.UseHandler(router)

	server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: n,
	}

	//go globalRoom.run()

	log.Println("API is Listening on: ", common.AppConfig.Server)
	log.Fatal(server.ListenAndServe())

}
