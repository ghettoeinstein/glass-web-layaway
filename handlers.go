package main

import (
	"./common"
	"./controllers"
	"./data"
	"./models"
	"./random"
	"github.com/gorilla/mux"

	"html/template"
	"log"
	"net/http"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", "base", "")
}

func glassHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "glass", "base", "")

}

func postGlassHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting post")
	r.ParseMultipartForm(20 << 32)
	err := r.ParseForm()
	if err != nil {
		log.Fatal("Error parsing form")
	}
	log.Println("no errors")

	fullname := r.Form.Get("fullname")
	email := r.Form.Get("email")
	phoneNumber := r.Form.Get("phoneNumber")
	dateCreated := time.Now()
	url := r.Form.Get("url")
	url = template.HTMLEscapeString(url)
	_ = &NewUser{FullName: fullname, Email: email, PhoneNumber: phoneNumber, DateCreated: dateCreated, URL: url}

	uuid := random.GenerateUUID()

	order := &models.WebOrder{
		FullName:    fullname,
		UUID:        uuid,
		Email:       email,
		PhoneNumber: phoneNumber,
		URL:         url,
		Decision:    "undecided",
	}

	err = saveOrder(order)
	if err != nil {
		log.Println(err)
		renderTemplate(w, "glass", "base", err.Error())
		return
	}

	//go postOrderToSlack(order)
	//go textOrderToAdmins(order)

	payload := struct {
		Order    *models.WebOrder
		Redirect string
	}{
		order,
		"/decision/" + uuid,
	}
	renderTemplate(w, "countdown", "base", payload)

}

func saveOrder(order *models.WebOrder) (err error) {
	context := controllers.NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	log.Println("About to save to database")
	if err = repo.NewWebOrder(order); err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "profile", "base", "nil")
}

func congratulationsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "congratulations", "base", "")
}

func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "privacy-policy", "base", "")
}

func sorryHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "sorry", "base", "")
}

func offerHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "countdown", "base", nil)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "chat", "base", "")
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "history", "base", "")
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := controllers.NewContext()
	defer ctx.Close()

	renderTemplate(w, "orders", "base", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", "base", "")
}

func termsHandler(w http.ResponseWriter, r *http.Request) {

	uuid := mux.Vars(r)
	if uuid["id"] == "" {
		return
	}

	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	_, err := repo.GetByUUID(uuid["id"])
	if err != nil {
		log.Println("Error fetching order for UUID:", err)
		common.DisplayAppError(w, err, "Error fetching order for UUID", 500)
		return
	}
	renderTemplate(w, "terms", "base", uuid["id"])
}

func aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about-us", "base", nil)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Handler for HTTP Get - "/admin/merchants/{id}/products"
	// Returns all Tasks created by a User

	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}
	web_orders, err := repo.GetNewOrders()

	if err != nil {
		log.Println("DB Error looking up web orders :", err)
		renderTemplate(w, "orders", "base", nil)
		return
	}

	//j, err := json.Marshal(ProductsResource{Data: products})
	//if err != nil {
	//	common.DisplayAppError(
	//		w,
	//		err,
	//		"An unexpected error has occurred",
	//		500,
	//	)
	//	return
	//}
	//w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(j)

	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")
	renderTemplate(w, "orders", "base", web_orders)
}

func loggingHandler(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "Auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header()["Location"] = []string{"/login"}
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func decisionHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)

	ctx := controllers.NewContext()
	defer ctx.Close()
	c := ctx.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	webOrder, err := repo.GetByUUID(uuid["id"])
	if err != nil {
		log.Println(err)
		common.DisplayAppError(w, err, "Error getting result", 500)
		return
	}

	switch webOrder.Decision {
	case "approved":
		renderTemplate(w, "congratulations", "base", uuid["id"])

	case "denied":
		renderTemplate(w, "sorry", "base", uuid["id"])
	default:
		renderTemplate(w, "sorry", "base", uuid["id"])
	}
	return
}

func IdFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	id := vars["id"]

	return id
}
