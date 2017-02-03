package main

import (
	"html/template"
)

func setupTemplates() {

	templates["about-us"] = template.Must(template.ParseFiles("templates/about-us.html", "templates/base.html"))
	templates["glass"] = template.Must(template.ParseFiles("templates/glass.html", "templates/base.html"))
	templates["congratulations"] = template.Must(template.ParseFiles("templates/congratulations.html", "templates/base.html"))
	templates["countdown"] = template.Must(template.ParseFiles("templates/countdown.html", "templates/base.html"))
	templates["privacy-policy"] = template.Must(template.ParseFiles("templates/privacy-policy.html", "templates/base.html"))
	templates["sorry"] = template.Must(template.ParseFiles("templates/sorry.html", "templates/base.html"))
	templates["index"] = template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	templates["profile"] = template.Must(template.ParseFiles("templates/profile.html", "templates/base.html"))
	templates["history"] = template.Must(template.ParseFiles("templates/history.html", "templates/base.html"))
	templates["order"] = template.Must(template.ParseFiles("templates/order.html", "templates/base.html"))

	templates["home"] = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))
	templates["terms"] = template.Must(template.ParseFiles("templates/terms.html", "templates/base.html"))

	templates["chat"] = template.Must(template.ParseFiles("templates/chat.html", "templates/base.html"))
	templates["orders"] = template.Must(template.ParseFiles("templates/orders.html", "templates/base.html"))

}
