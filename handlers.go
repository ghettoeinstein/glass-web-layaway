package main

import (
	"./common"
	"./random"
	"encoding/json"
	"fmt"
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
	r.ParseMultipartForm(20 << 32)
	err := r.ParseForm()
	if err != nil {
		log.Fatal("Error parsing form")
	}
	log.Println("no errors")
	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			fmt.Println(key, value)
		}
	}
	_ = r.Form.Get("offerId")
	fullname := r.Form.Get("fullname")
	email := r.Form.Get("email")
	phoneNumber := r.Form.Get("phoneNumber")
	dateCreated := time.Now()
	url := r.Form.Get("url")
	url = template.HTMLEscapeString(url)
	user := &NewUser{FullName: fullname, Email: email, PhoneNumber: phoneNumber, DateCreated: dateCreated, URL: url}

	uuid, err := random.GenerateUUID()
	if err != nil {
		return
	}

	order := &Order{
		User: user,
		UUID: uuid,

		DateCreated: dateCreated,
		ExpireTime:  dateCreated.Add(time.Duration(2 * time.Minute)),
		Expired:     false,
	}

	postOrderToSlack(order)
	textOrderToAdmins(order)

	payload := struct {
		Order *Order
	}{
		order,
	}
	renderTemplate(w, "countdown", "base", payload)

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
	renderTemplate(w, "orders", "base", "")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", "base", "")
}

func termsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "terms", "base", "")
}

func viewOrderHandler(w http.ResponseWriter, r *http.Request) {

	uuid, err := random.GenerateUUID()
	if err != nil {
		return
	}

	order := &NewUser{
		UUID:        uuid,
		URL:         "http://www.sephora.com/ruby-pink-gel-coat-P414001?skuId=1890474&icid2=just%20arrived:p414001",
		FullName:    "Juan Carlos",
		PhoneNumber: "3234236651",
		Email:       "juancarlos@yahoo.com",
	}
	renderTemplate(w, "order", "base", order)
}

func aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about-us", "base", nil)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "orders", "base", nil)
}

func loggingHandler(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func genAdminHandler(w http.ResponseWriter, r *http.Request) {
	token, err := common.GenerateAdminJWT("caleb@getglass.co", "admin")
	if err != nil {
		common.DisplayAppError(w, err, "Error making admin token", 500)
	}

	cookie := http.Cookie{Name: "Auth", Value: token, Expires: time.Now().Add(time.Hour * 24 * 30), HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/admin/orders", 302)

}

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := random.GenerateUUID()
	if err != nil {
		return
	}

	UUIDPayload := struct {
		UUID string `json:"uuid"`
	}{
		uuid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UUIDPayload)
	return
}
