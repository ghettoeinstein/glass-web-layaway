package controllers

import (
	"../common"
	"../data"

	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var dataResource UserResource
	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid User data",
			500,
		)
		return
	}
	user := &dataResource.Data
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("users")
	repo := &data.UserRepository{C: col}

	//custChan := make(chan string, 1)
	//// Insert User document
	//_, err = CreateStripeCustomer(user, custChan)
	//if err != nil {
	//
	//	common.DisplayAppError(
	//		w,
	//		err,
	//		"Eror while generating the access token",
	//		500,
	//	)
	//	return
	//
	//}
	//
	//result := <-custChan
	//user.CustomerId = result
	err = repo.CreateUser(user)
	if err != nil {
		log.Fatalln(err)
	}
	// Clean-up the hashpassword to eliminate it from response JSON
	user.HashPassword = nil

	j, err := json.Marshal(UserResource{Data: *user})
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)

}
