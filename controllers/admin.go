package controllers

import (
	"../common"
	"../data"
	"../models"
	"encoding/json"
	"log"
	"time"
	//"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"net/http"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", "base", nil)
}
func AdminDisplayOrder(w http.ResponseWriter, r *http.Request) {
	//_, err := customerFromRequest(r)
	//if err != nil {
	//	log.Println("Error fetching customer for request:", err.Error)
	//	common.DisplayAppError(w, err, "Error creating  order for customer", 500)
	//	return
	//}

	context := NewContext()
	defer context.Close()

	c := context.DbCollection("orders")
	repo := &data.OrderRepository{c}
	sampleOrder := &models.Order{Id: bson.NewObjectId()}

	if err := repo.NewOrder(sampleOrder); err != nil {

		common.DisplayAppError(w, err, "Error saving order.", 500)
		return
	}

	w.Header().Set("Content-Type", "Application/json; charset=utf-8")

	json.NewEncoder(w).Encode(sampleOrder)
	return

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
	user, err := repo.Login(loginUser)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid login credentials",
			401,
		)
		return
	}
	// Generate JWT token
	token, err = common.GenerateJWT(user.Email, "member")
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Eror while generating the access token",
			500,
		)
		return
	}
	log.Println("ijjjjj")
	cookie := http.Cookie{Name: "Auth", Value: token, Expires: time.Now().Add(time.Hour * 3), HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/admin", 302)

}
