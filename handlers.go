package main

import (
	"./common"
	"./controllers"
	"./data"
	"./models"
	"./random"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joiggama/money"
	"html/template"
	"log"
	"net/http"
	"time"
)

func SMSLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "sms-login", "base", nil)

}

func GETVerifySMSLogin(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	phone := r.FormValue("phone")

	otp := GenerateOTP(6, phone)
	otpStore[phone] = otp

	if _, _, err := sendMessage(phone, "Your Glass verification code is: "+otp.Passcode); err != nil {
		log.Fatal("Error sending msg: ", err)
	}
	log.Println("Sent message successfully to handset: ", phone)
	renderTemplate(w, "sms-verify", "base", phone)
}

func POSTVerifySMSLogin(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	passcode := r.FormValue("passcode")
	phone := r.FormValue("phone")
	log.Println("Starting login for:  ", phone)

	ok := verifyOTP(passcode, phone)
	if !ok {
		//Unsuccessful login

		http.Redirect(w, r, "/login", 401)
		return
	} else {

		//Successful Login
		user, err := UserByNumber(phone)
		if err != nil {

			common.DisplayAppError(w, err, "Could not find user to verify", 500)
			return
		}
		token, err := common.GenerateJWT(user.Email, "member")
		if err != nil {
			log.Println("Error creating JWT for %s %s", user.Email, err)
			http.Redirect(w, r, "/login", 500)
		}
		log.Println("Generating cookie for: ", user.Email)
		cookie := http.Cookie{Name: "Auth", Value: token, Path: "/user/", Expires: time.Now().Add(time.Hour * 24), HttpOnly: true}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/user/home", 302)
	}
	//

}

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

	//set  page to  expiration time in past, so that the page is never cached.
	// Put this into middleware  later
	ctx := r.Context()
	email := ctx.Value(common.EmailKey).(string)

	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("users")
	repo := &data.UserRepository{c}
	user, err := repo.GetByUsername(email)
	if err != nil {
		log.Println("No user found for email")
	}

	log.Println(user.FullName)

	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")

	renderTemplate(w, "profile", "base", user)

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

	//set  page to  expiration time in past, so that the page is never cached.
	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")

	renderTemplate(w, "history", "base", "")
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := controllers.NewContext()
	defer ctx.Close()

	renderTemplate(w, "orders", "base", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	//set  page to  expiration time in past, so that the page is never cached.
	w.Header().Set("Cache-Control", "no-cache,no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", " Sat, 26 Jul 1997 05:00:00 GMT")

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

	webOrder, err := repo.GetByUUID(uuid["id"])
	if err != nil {
		log.Println("Error fetching order for UUID:", err)
		common.DisplayAppError(w, err, "Error fetching order for UUID", 500)
		return
	}
	taxes := webOrder.Price * 0.0875
	serviceFee := webOrder.Price * 0.1
	termsPayload := struct {
		MonthlyPayment interface{}
		Total          interface{}
		FirstPayment   interface{}
		UUID           interface{}
	}{
		money.Format(webOrder.Price / 4),
		money.Format(webOrder.Price),
		money.Format(webOrder.Price/4 + serviceFee + taxes),
		uuid["id"],
	}
	renderTemplate(w, "terms", "base", termsPayload)
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

	//set  page to  expiration time in past, so that the page is never cached.
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

	http.Redirect(w, r, "/login", 303)

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

func userLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "Auth",
		Value:  "",
		Path:   "/user/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "Auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header()["Location"] = []string{"/login"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func emailHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	email := r.FormValue("email")

	if email == "" {
		http.Error(w, errors.New("Cannot have blank email").Error(), 500)
		return
	}

	if err := AddSubscriberToMailChimp(email); err != nil {
		log.Println("Error adding email", email, "to mailchimp: ", err)
	}

	log.Println("Email is: ", email)
	payload := struct {
		Result string `json:"result"`
	}{
		"Saved email successfully",
	}
	json.NewEncoder(w).Encode(payload)
	return
}
