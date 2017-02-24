package main

import (
	"./common"
	"./controllers"
	"os"
	//	"flag"
	"./data"
	"github.com/shopspring/decimal"
	"github.com/urfave/negroni"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	_ "strconv"
	"time"
)

// Variables for different levels of logging.
var (
	Trace   *log.Logger // Just about Anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // General Errors
)

func init() {
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
	router := InitRoutes()
	router.HandleFunc("/chat", chatHandler)
	router.HandleFunc("/status/{id}", UUIDStatus).Methods("GET")
	router.Handle("/room", globalRoom)
	n := negroni.Classic()
	n.UseHandler(router)
	server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: n,
	}
	Info.Println("API is Listening on: ", common.AppConfig.Server)
	Error.Fatal(server.ListenAndServe())
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
	formPrice := template.HTMLEscapeString(r.PostFormValue("price"))
	price, _ := decimal.NewFromString(formPrice)

	webOrder.PriceStr = price.String()

	Trace.Println("Price is :", webOrder.Price)

	//res, err := strconv.ParseFloat(price, 64)
	//if err != nil {
	println("Error parsing price string into int:", err)
	//}
	//webOrder.Price = float64(res)
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
		Trace.Println(err)
		common.DisplayAppError(w, err, err.Error(), 500)
		return
	}
	//renderTemplate(w, "admin", "base", webOrders)
	w.Header()["Location"] = []string{"/admin"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func UUIDStatus(w http.ResponseWriter, r *http.Request) {
	id := controllers.IdFromRequest(r)

	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}
	webOrder, err := repo.GetByUUID(id)
	if err != nil {
		common.DisplayAppError(w, err, err.Error(), 404)
		return
	}

	if !webOrder.Acknowledged {
		w.Write([]byte("stay"))
		return

	}
	w.Write([]byte("go"))

}
