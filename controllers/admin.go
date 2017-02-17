package controllers

import (
	"../common"
	"../data"
	"../models"
	"html/template"
	"log"
	"strconv"
	"time"
	//"github.com/gorilla/mux"
	//	"gopkg.in/mgo.v2/bson"

	"net/http"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "admin-login", "base", nil)
}

func AdminGetNewOrders(w http.ResponseWriter, r *http.Request) {

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	var orders []models.WebOrder

	orders, err := repo.GetNewOrders()
	if err != nil {
		common.DisplayAppError(w, err, "Could not retrieve orders. Contact IT", 500)
		return
	}

	log.Println(len(orders))
	//set  page to  expiration time in past, so that the page is never cached.
	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")
	renderTemplate(w, "orders", "base", orders)

}

func AdminProcessOrder(w http.ResponseWriter, r *http.Request) {

	id := IdFromRequest(r)

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	webOrder, err := repo.GetByUUID(id)
	if err != nil {
		common.DisplayAppError(w, err, err.Error(), 500)
		return
	}

	price := template.HTMLEscapeString(r.PostFormValue("price"))
	notes := template.HTMLEscapeString(r.PostFormValue("notes"))
	log.Println(notes)
	res, err := strconv.ParseFloat(price, 64)
	if err != nil {

		println("Error parsing price string into int:", err)
	}

	webOrder.Notes = notes
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

	//renderTemplate(w, "admin", "base", webOrders)

	w.Header()["Location"] = []string{"/admin"}
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func AdminGetDeniedOrders(w http.ResponseWriter, r *http.Request) {

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	var orders []models.WebOrder

	orders, err := repo.GetDeniedOrders()
	if err != nil {
		common.DisplayAppError(w, err, "Could not retrieve orders. Contact IT", 500)
		return
	}

	//set  page to  expiration time in past, so that the page is never cached.
	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")
	renderTemplate(w, "denied", "base", orders)

}

func AdminGetApprovedOrders(w http.ResponseWriter, r *http.Request) {

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	var orders []models.WebOrder

	orders, err := repo.GetApprovedOrders()
	if err != nil {
		log.Println("Error:", err)
		common.DisplayAppError(w, err, "Could not retrieve approved orders. Contact IT", 500)
		return
	}

	//set  page to  expiration time in past, so that the page is never cached.
	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")
	renderTemplate(w, "approved", "base", orders)

}

func AdminDisplayOrder(w http.ResponseWriter, r *http.Request) {
	//_, err := customerFromRequest(r)
	//if err != nil {
	//	log.Println("Error fetching customer for request:", err.Error)
	//	common.DisplayAppError(w, err, "Error creating  order for customer", 500)
	//	return
	//}
	id := IdFromRequest(r)

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	web_order, err := repo.GetByUUID(id)
	if err != nil {
		common.DisplayAppError(w, err, "Error fetching  web order.", 500)
		return
	}

	renderTemplate(w, "order", "base", web_order)

}

func AdminGetEditOrder(w http.ResponseWriter, r *http.Request) {
	//_, err := customerFromRequest(r)
	//if err != nil {
	//	log.Println("Error fetching customer for request:", err.Error)
	//	common.DisplayAppError(w, err, "Error creating  order for customer", 500)
	//	return
	//}

	context := NewContext()
	defer context.Close()

	//Use helper method `IdFromRequest to get the product id
	webOrderId := IdFromRequest(r)
	log.Println(webOrderId)
	c := context.DbCollection("web_orders")

	//Repository for web orders.
	repo := &data.WebOrderRepository{c}
	webOrder, err := repo.GetByUUID(webOrderId)
	if err != nil {
		common.DisplayAppError(w, err, err.Error(), 500)
		return
	}

	renderTemplate(w, "order", "base", webOrder)
	return

}

func AdminDeleteOrder(w http.ResponseWriter, r *http.Request) {

	id := IdFromRequest(r)

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	err := repo.DeleteByUUID(id)
	if err != nil {
		common.DisplayAppError(w, err, "Error deleting order", 500)
		return
	}
	http.Redirect(w, r, r.Referer(), 303)
}

// Used to process front-end login
func AdminLogin(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	var username string
	var password string
	if r.Form.Get("email") != "" {
		username = r.Form.Get("email")
	}

	if r.Form.Get("password") != "" {
		password = r.Form.Get("password")
	}

	log.Println("Starting login")
	var dataResource LoginResource
	var token string

	_ = dataResource.Data
	loginUser := models.User{
		Email:    username,
		Password: password,
	}
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("users")
	repo := &data.UserRepository{C: col}
	// Authenticate the login user

	log.Println(loginUser)
	user, err := repo.Login(loginUser)
	if err != nil {
		log.Println("Error logging in: ", err)
		http.Redirect(w, r, "/team/login", 301)

		return
	}
	// Generate JWT token
	token, err = common.GenerateJWT(user.Email, "administrator")
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Eror while generating the access token",
			500,
		)
		return
	}

	cookie := http.Cookie{Name: "Auth", Path: "/", Value: token, Expires: time.Now().Add(time.Hour * 24), HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/admin", 302)

}
