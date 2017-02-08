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

	templates["admin-login"] = template.Must(template.ParseFiles("templates/admin-login.html", "templates/base.html"))

	templates["order"] = template.Must(template.ParseFiles("templates/order.html", "templates/base.html"))

	templates["approved"] = template.Must(template.ParseFiles("templates/approved.html", "templates/base.html"))

	templates["denied"] = template.Must(template.ParseFiles("templates/denied.html", "templates/base.html"))

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
