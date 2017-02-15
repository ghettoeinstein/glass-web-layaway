package main

import (
	"./common"
	"./controllers"
	"os"
	//	"flag"
	"./data"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/urfave/negroni"
	"html/template"

	"log"
	"math/rand"
	"net/http"
	"time"
)

// Variables for different levels of logging.
var (
	Trace   *log.Logger // Just about Anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger
)

func init() {
	if roomStore == nil {
		roomStore = make(map[string]*room)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	go controllers.InitTemplates()

	go setupTemplates()
	go func() {
		globalRoom = newRoom("global")
		globalRoom.run()
	}()

}

func controllerChannel(out <-chan string, r *room) {
	for {
		select {
		case uuid := <-out:
			Trace.Println("Attempting to send UUID:", uuid)
			r.forward <- []byte(uuid)
		default:
		}
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Trapped panic: %s (%T) \n", err, err)
		}
	}()

	f, err := os.OpenFile("logs/glassLogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error  opening up logfile %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	Trace = log.New(f,
		"[TRACE] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(f,
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(f,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(f,
		"[WARNING] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	//http.HandleFunc("/", rootHandler)

	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	//router.HandleFunc("/register", controllers.Register).Methods("POST")
	//router.ServeFiles("/assets/*filepath", http.Dir("/assets"))
	//
	//router.GET("/about-us", aboutUsHandler)
	//router.GET("/tos", termsHandler)
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
	//router.HandleFunc("/users/register", controllers.Register).Methods("POST")
	router.HandleFunc("/submit-email", emailHandler).Methods("POST")

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

	router.HandleFunc("/about", aboutUsHandler)

	//
	router.HandleFunc("/login", SMSLogin).Methods("GET")

	router.HandleFunc("/login/verify", GETVerifySMSLogin).Methods("GET")
	router.HandleFunc("/login/verify", POSTVerifySMSLogin).Methods("POST")
	router.HandleFunc("/team/login", controllers.GetLogin).Methods("GET")
	router.HandleFunc("/team/login", controllers.AdminLogin).Methods("POST")
	router.HandleFunc("/chat", chatHandler)
	router.Handle("/room", globalRoom)

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin/chat", chatHandler).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}", controllers.AdminDisplayOrder).Methods("GET")
	adminRouter.HandleFunc("/admin/orders/{id}/process", AdminProcessOrder).Methods("POST")
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

	Info.Println("API is Listening on: ", common.AppConfig.Server)
	log.Fatal(server.ListenAndServe())

}

func AdminProcessOrder(w http.ResponseWriter, r *http.Request) {

	id := controllers.IdFromRequest(r)

	context := controllers.NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	webOrder, err := repo.GetByUUID(id)
	if err != nil {
		common.DisplayAppError(w, err, err.Error(), 500)
		return
	}

	price := template.HTMLEscapeString(r.PostFormValue("price"))
	res, err := strconv.ParseFloat(price, 64)
	if err != nil {

		println("Error parsing price string into int:", err)
	}

	webOrder.Price = float64(res)

	decision := r.PostFormValue("decision")
	switch decision {
	case "approve":
		webOrder.Decision = "approved"

	case "deny":
		webOrder.Decision = "denied"
	default:
		webOrder.Decision = "denied"
	}
	webOrder.Acknowledged = true

	err = repo.UpdateOrder(webOrder)
	if err != nil {
		common.DisplayAppError(w, err, err.Error(), 500)
		return
	}
	globalRoom.forward <- []byte(id)

	//renderTemplate(w, "admin", "base", webOrders)

	w.Header()["Location"] = []string{"/admin"}
	w.WriteHeader(http.StatusTemporaryRedirect)

}
