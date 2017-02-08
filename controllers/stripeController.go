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
	"github.com/stripe/stripe-go/invoiceitem"
	"github.com/stripe/stripe-go/sub"
	"time"

	"github.com/stripe/stripe-go/plan"
	"log"
	"net/http"
	"os"
	"strconv"
)

func NewUserFromWebOrder(o *models.WebOrder) (*models.User, error) {

	user := &models.User{
		Email:       o.Email,
		FullName:    o.FullName,
		PhoneNumber: o.PhoneNumber,
	}

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("users")
	repo := &data.UserRepository{c}

	err := repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

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
func CreateStripeCustomerWithToken(u *models.User, token string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: u.Email,
		Desc:  "Customer for " + u.Email,
	}
	log.Println("token is", "token")
	err := params.SetSource(token)
	if err != nil {
		log.Println("Could not add source to customer")
		return nil, err
	}
	stripe.Key = os.Getenv("STRIPE_KEY")

	cust, err := customer.New(params)
	if err != nil {
		log.Println(err)
		return cust, err

	}
	//Set the user's Stripe `CustomerId` Field
	u.StripeCustomer.CustomerId = cust.ID

	// send the customer id back from stripe on the send-only channel `sc`. The calling invoker of this function blocks

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
	stripe.Key = os.Getenv("STRIPE_KEY")
	log.Println("Stripe key is:", stripe.Key)
	r.ParseMultipartForm(32 << 20)
	token := r.PostFormValue("stripeToken")

	if token == "" {
		log.Fatalln("No token")
		return
	}

	id := IdFromRequest(r)

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	webOrder, err := repo.GetByUUID(id)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(webOrder.UUID)

	//Create and save a user to the database from the web order
	user, err := NewUserFromWebOrder(webOrder)
	if err != nil {
		log.Println("error creating user", err)
		return
	}

	//C reate a stripe customer from the user
	stripeCustomer, err := CreateStripeCustomerWithToken(user, token)
	if err != nil {
		log.Println("Error creating User", err)
		return
	}
	// set the user's stripe customerid to the returned customer object's
	user.StripeCustomer.CustomerId = stripeCustomer.ID

	serviceFee := int64(2000)
	// Create an invoice for the service fee of the loan for the user.
	invoiceParams := &stripe.InvoiceItemParams{
		Customer: user.StripeCustomer.CustomerId,
		Amount:   serviceFee,
		Currency: "usd",
		Desc:     "One-time service fee for plan: " + id,
	}
	_, err = invoiceitem.New(invoiceParams)
	if err != nil {
		log.Println("Error creating invoice item")
		return
	}

	p, err := plan.New(&stripe.PlanParams{
		Amount:   19999,
		Interval: "month",
		Name:     id + " Installment Plan",
		Currency: "usd",
		ID:       id,
	})

	if err != nil {
		log.Println("Error creating plan")
		return
	}
	log.Println(p)

	s, err = sub.New(&stripe.SubParams{
		Customer: user.StripeCustomer.CustomerId,
		Plan:     p.ID,
	})
	if err != nil {
		log.Println("Error creating subscription")
		return
	}

	authToken, err := common.GenerateJWT(user.Email, "Customer")
	if err != nil {
		log.Println("Error creating token")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Auth",
		Value:    authToken,
		Path:     "/user/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	w.Header()["Location"] = []string{"/user/history"}
	w.WriteHeader(http.StatusSeeOther)

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
