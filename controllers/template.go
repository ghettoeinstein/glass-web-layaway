package controllers

import (
	"html/template"
	"log"
	"net/http"
)

var templates map[string]*template.Template

func InitTemplates() {

	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["admin"] = template.Must(template.ParseFiles("templates/orders.html", "templates/base.html"))

	templates["login"] = template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))

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
