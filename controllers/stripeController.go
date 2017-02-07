package controllers

import (
	"../common"
	"../data"
	"../models"
	"encoding/json"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetCustomerForUser(w http.ResponseWriter, r *http.Request) {
	cust, err := customerFromRequest(r)
	if err != nil {
		log.Println("Error fetching customer for request:", err.Error)
		common.DisplayAppError(w, err, "Error fetching customer:", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(cust)
	return
}

//Pass in a pointer to a models.User, and `sc` channel to report completion of creating customer
func CreateStripeCustomer(u *models.User, sc chan<- string) (interface{}, error) {
	params := &stripe.CustomerParams{
		Desc: "Customer for " + u.Email,
	}
	stripe.Key = os.Getenv("STRIPE_KEY")

	cust, err := customer.New(params)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	//Set the user's Stripe `CustomerId` Field
	u.StripeCustomer.CustomerId = cust.ID

	// send the customer id back from stripe on the send-only channel `sc`. The calling invoker of this function blocks
	//waiting for this send.
	sc <- cust.ID
	log.Println("Created customer with id ", cust.ID)
	return cust, nil
}

// Set default payment source(card) for a customer
func SetDefaultSource(w http.ResponseWriter, r *http.Request) {
	cust, err := customerFromRequest(r)
	if err != nil {
		log.Println("Error fetching customer")
	}
	log.Println(cust)
	stripe.Key = os.Getenv("STRIPE_KEY")

	r.ParseForm()

	source := r.Form.Get("source")

	c, err := customer.Update(
		cust.ID,
		&stripe.CustomerParams{DefaultSource: source},
	)
	log.Println(c)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(402)
		return
	}
	w.WriteHeader(200)
	return

}

func AddSourceToCustomer(w http.ResponseWriter, r *http.Request) {
	cust, err := customerFromRequest(r)
	if err != nil {
		common.DisplayAppError(w, err, "Error fetching customer", 500)
		return
	}
	r.ParseForm()

	source := r.Form.Get("source")
	log.Println("Source is:", source)
	stripe.Key = os.Getenv("STRIPE_KEY")

	c, err := card.New(&stripe.CardParams{
		Customer: cust.ID,
		Token:    source,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(402)
		return
	}
	log.Println("Added card:", c)
	w.WriteHeader(200)
	return

}

func ChargeForOffer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			log.Println(key, value)
		}
	}
	http.Redirect(w, r, "/glass", 307)
}

func ChargeCustomer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cust, err := customerFromRequest(r)
	if err != nil {
		common.DisplayAppError(w, err, "Error getting customer from database", 500)
		return
	}
	source := r.Form.Get("source")
	amount := r.Form.Get("amount")
	chargeAmount, err := strconv.Atoi(amount)
	if err != nil {
		log.Println("error getting amouunt")
	}
	log.Println(" Amount is  " + amount)
	log.Println("Source for charge is:", source)
	stripe.Key = os.Getenv("STRIPE_KEY")

	chargeParams := &stripe.ChargeParams{
		Customer: cust.ID,
		Amount:   uint64(chargeAmount),
		Currency: "usd",
		Desc:     "Charge for test@getglass.co",
	}
	chargeParams.SetSource(source)
	ch, err := charge.New(chargeParams)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(402)
		return
	}

	log.Println("Charge was successful", ch)
	w.WriteHeader(200)
	return

}

func customerFromRequest(r *http.Request) (cust *stripe.Customer, err error) {
	token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &common.AppClaims{}, common.KeyFunc)
	if err != nil {
		return nil, err
	}

	username := token.Claims.(*common.AppClaims).UserName
	log.Println("Looking up user for username: ", username)
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("users")

	repo := &data.UserRepository{c}
	user, err := repo.GetByUsername(username)
	if err != nil {
		log.Println("error fetching user from db:", err.Error())
		return nil, err
	}

	// We have user at this point
	log.Println("User stripe customer id:", user.CustomerId)

	stripe.Key = os.Getenv("STRIPE_KEY")

	cust, err = customer.Get(user.CustomerId, nil)

	return
}
