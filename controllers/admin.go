package controllers

import (
	"../common"
	"../data"
	"../models"
	"encoding/json"
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
