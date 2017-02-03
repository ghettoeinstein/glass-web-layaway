package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func IdFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	id := vars["id"]
	return id
}
