package controllers

import (
	"../common"
	"../data"
	"../models"
	//	"github.com/shopspring/decimal"

	_ "github.com/joiggama/money"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/invoiceitem"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	//	"os"
	//"time"
)

func GetCustomerForUser(w http.ResponseWriter, r *http.Request) {

	return
}

//Pass in a pointer to a models.User, and token string to facilitate creating customer
func CreateStripeCustomerWithToken(u *models.User, token string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: u.Email,
		Desc:  "Customer for " + u.Email,
	}
	Stripe.Println("token is", token)
	err := params.SetSource(token)
	if err != nil {
		Stripe.Println("Could not add source to customer")
		return nil, err
	}

	cust, err := customer.New(params)
	if err != nil {
		return nil, err

	}
	//Set the user's Stripe `CustomerId` Field
	u.StripeCustomer.CustomerId = cust.ID

	// send the customer id back from stripe on the send-only channel `sc`. The calling invoker of this function blocks

	Stripe.Println("Created customer with id ", cust.ID)
	return cust, nil
}

// Set default payment source(card) for a customer
//func SetDefaultSource(w http.ResponseWriter, r *http.Request) {
//	cust, err := customerFromRequest(r)
//	if err != nil {
//		log.Println("Error fetching customer")
//	}
//	log.Println(cust)
//	stripe.Key = os.Getenv("STRIPE_KEY")
//
//	r.ParseForm()
//
//	source := r.Form.Get("source")
//
//	c, err := customer.Update(
//		cust.ID,
//		&stripe.CustomerParams{DefaultSource: source},
//	)
//	log.Println(c)
//	if err != nil {
//		log.Println(err.Error())
//		w.WriteHeader(402)
//		return
//	}
//	w.WriteHeader(200)
//	return
//
//}
//
//func AddSourceToCustomer(w http.ResponseWriter, r *http.Request) {
//	cust, err := customerFromRequest(r)
//	if err != nil {
//		common.DisplayAppError(w, err, "Error fetching customer", 500)
//		return
//	}
//	r.ParseForm()
//
//	source := r.Form.Get("source")
//	log.Println("Source is:", source)
//	stripe.Key = os.Getenv("STRIPE_KEY")
//
//	c, err := card.New(&stripe.CardParams{
//		Customer: cust.ID,
//		Token:    source,
//	})
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(402)
//		return
//	}
//	log.Println("Added card:", c)
//	w.WriteHeader(200)
//	return
//
//}

//

func ChargeNewCustomerForOffer(w http.ResponseWriter, r *http.Request) {
	log.Println("Stripe key is:", stripe.Key)
	r.ParseMultipartForm(32 << 20)
	token := r.PostFormValue("stripeToken")

	id := IdFromRequest(r)
	if token == "" {
		Stripe.Println("No token in request, exiting handler.")
		http.Redirect(w, r, "/terms/"+id+"?err=12", http.StatusSeeOther)
		return
	}

	//Get the web Order from the database
	webOrder, err := WebOrderForUUID
	if err != nil {
		Trace.Println(err)
		return
	}

	//Create and save a user to the database from the web order
	user := NewUserFromWebOrder(webOrder)
	user.Id = bson.NewObjectId()

	//Create a stripe customer from the user
	stripeCustomer, err := CreateStripeCustomerWithToken(user, token)
	if err != nil {
		HandleStripeError(w, r, id, err)
		return
	}

	// set the user's stripe customerid to the returned customer object's
	user.StripeCustomer.CustomerId = stripeCustomer.ID
	c := context.DbCollection("users")
	userRepo := &data.UserRepository{c}

	err = userRepo.CreateUser(user)
	if err != nil {
		log.Println("Error persisting user:", err)
	}

	order := NewOrder(webOrder, user)

	c = context.DbCollection("orders")
	orderRepo := &data.OrderRepository{c}

	serviceFee, _ := decimal.NewFromString(order.Total).Mul(decimal.NewFromFloat(0.10))
	order.ServiceFee = serviceFee.String()

	order.CombinedTotal, _ = decimal.NewFromString(order.Total).Add(order.SalesTax.Add(order.ServiceFee))
	// Create an invoice for the Glass Service Fee(10%) of the total cost of goods  for the user.

	invoiceItem1, err := invoiceitem.New(&stripe.InvoiceItemParams{
		Customer: user.StripeCustomer.CustomerId,
		Amount:   order.ServiceFee.IntPart(),
		Currency: "usd",
		Desc:     "One-time service fee for plan: " + id,
	})
	if err != nil {
		log.Println("Error creating invoice item")
		http.Error(w, err.Error(), 500)
		return
	}

	invoiceItem2, err := invoiceitem.New(&stripe.InvoiceItemParams{
		Customer: user.StripeCustomer.CustomerId,
		Amount:   order.SalesTax.IntPart(),
		Currency: "usd",
		Desc:     "One-time taxes(8.75%) for plan: " + id,
	})
	if err != nil {
		Stripe.Println("Error creating invoice item")
		http.Error(w, err.Error(), 500)
		return
	}

	order.InvoiceItems = append(order.InvoiceItems, invoiceItem1, invoiceItem2)

	p, err := plan.New(&stripe.PlanParams{
		Amount:   uint64(order.MonthlyPayment.IntPart()),
		Interval: "month",
		Name:     id + " Installment Plan",
		Currency: "usd",
		ID:       id,
	})
	if err != nil {
		Stripe.Println("Error creating plan")
		return
	}
	Stripe.Println(p)

	order.PlanID = p.ID

	// Create a subscription and attach the plan for the order to the sub.
	s, err := sub.New(&stripe.SubParams{
		Customer: user.StripeCustomer.CustomerId,
		Plan:     p.ID,
	})

	if err != nil {
		Stripe.Println("Error creating subscription:", s)

		w.Header()["Location"] = []string{}
		http.Redirect(w, r, "/terms/"+id+"/flash=1", http.StatusSeeOther)

		return
	}

	order.SubscriptionID = s.ID

	err = orderRepo.SaveOrder(order)
	if err != nil {
		log.Println("Error saving order:", order)
	}
	Stripe.Println("Saved order successfully: ", order.UUID)

	err = SendConf(order)
	if err != nil {
		Stripe.Println("Error sending confirmation email: ", err)
	}

	// If all goes well create a cookie for the user to be able to login. Set to expire in one day.
	authToken, err := common.GenerateJWT(user.Email, "Customer")
	if err != nil {
		Stripe.Println("Error creating token")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Auth",
		Value:    authToken,
		Path:     "/user/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	w.Header()["Location"] = []string{"/confirmation"}
	renderTemplate(w, "confirmation", "base", order)

}

func ChargeCustomerForOffer(w http.ResponseWriter, r *http.Request) {
	//	stripe.Key = os.Getenv("STRIPE_KEY")
	//	log.Println("Stripe key is:", stripe.Key)
	//	r.ParseMultipartForm(32 << 20)
	//
	//	id := IdFromRequest(r)
	//
	//	context := NewContext()
	//	defer context.Close()
	//	c := context.DbCollection("web_orders")
	//	repo := &data.WebOrderRepository{c}
	//
	//	// get the web Order from the database
	//	webOrder, err := repo.GetByUUID(id)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	ctx := r.Context()
	//	email := ctx.Value(common.EmailKey).(string)
	//
	//	c = context.DbCollection("users")
	//	userRepo := &data.UserRepository{c}
	//	user, err := userRepo.GetByUsername(email)
	//	if err != nil {
	//		log.Println("No user found for email")
	//	}
	//
	//	cust, err := customer.Get(user.StripeCustomer.CustomerId, nil)
	//	if err != nil {
	//		log.Println("Error fetching customer data")
	//	}
	//
	//	order := NewOrderForUser(webOrder, user, cust)
	//
	//	c = context.DbCollection("orders")
	//	orderRepo := &data.OrderRepository{c}
	//
	//	serviceFee := order.Total.Mul(decimal.NewFromFloat(0.10))
	//	log.Println("Service Fee is:", serviceFee)
	//	order.ServiceFee = serviceFee
	//
	//	// Create an invoice for the Glas Service Fee(10%) of the total cost of goods  for the user.
	//
	//	invoiceItem1, err := invoiceitem.New(&stripe.InvoiceItemParams{
	//		Customer: user.StripeCustomer.CustomerId,
	//		Amount:   order.ServiceFee.IntPart(),
	//		Currency: "usd",
	//		Desc:     "One-time service fee for plan: " + id,
	//	})
	//	if err != nil {
	//		log.Println("Error creating invoice item")
	//		http.Error(w, err.Error(), 500)
	//		return
	//	}
	//
	//	invoiceItem2, err := invoiceitem.New(&stripe.InvoiceItemParams{
	//		Customer: user.StripeCustomer.CustomerId,
	//		Amount:   order.SalesTax.IntPart(),
	//		Currency: "usd",
	//		Desc:     "One-time taxes(8.75%) for plan: " + id,
	//	})
	//	if err != nil {
	//		Stripe.Println("Error creating invoice item")
	//		http.Error(w, err.Error(), 500)
	//		return
	//	}
	//
	//	order.InvoiceItems = append(order.InvoiceItems, invoiceItem1, invoiceItem2)
	//
	//	p, err := plan.New(&stripe.PlanParams{
	//		Amount:   uint64(order.MonthlyPayment.IntPart()),
	//		Interval: "month",
	//		Name:     id + " Installment Plan",
	//		Currency: "usd",
	//		ID:       id,
	//	})
	//	if err != nil {
	//		log.Println("Error creating plan")
	//		return
	//	}
	//	log.Println(p)
	//
	//	order.PlanID = p.ID
	//
	//	// Create a subscription and attach the plan for the order to the sub.
	//	s, err := sub.New(&stripe.SubParams{
	//		Customer: user.StripeCustomer.CustomerId,
	//		Plan:     p.ID,
	//	})
	//	if err != nil {
	//		log.Println("Error creating subscription:", s)
	//
	//		w.Header()["Location"] = []string{"/terms/" + id}
	//		w.WriteHeader(http.StatusSeeOther)
	//
	//		return
	//	}
	//
	//	order.SubscriptionID = s.ID
	//
	//	err = orderRepo.SaveOrder(order)
	//	if err != nil {
	//		log.Println("Error saving order:", order)
	//	}
	//	log.Println("Saved order successfully: ", order.UUID)
	//
	//	err = SendConf(order)
	//	if err != nil {
	//		log.Println("Error sending confirmation email:", err)
	//	}
	//
	//	http.Redirect(w, r, "/user/history", 303)

}
