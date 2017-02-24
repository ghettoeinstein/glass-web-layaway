//Contains helper functions for http.Handlers in the main package.

package main

import (
	"./accounting"
	"./common"
	"./controllers"
	"./data"
	"./models"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mnbbrown/mailchimp"
	"net/http"
	"os"
)

// Parse Id variable from request
func IdFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	id := vars["id"]

	return id
}

// Lookup, and serialize a user from the database.
func UserFromRequest(r *http.Request) (user *models.User, err error) {
	ctx := r.Context()

	email := ctx.Value(common.EmailKey).(string)

	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("users")
	repo := &data.UserRepository{c}
	user, err = repo.GetByUsername(email)
	return
}

func OrdersForUser(user *models.User) (orders []*models.Order, err error) {

	context := controllers.NewContext()
	defer context.Close()

	c := context.DbCollection("orders")
	repo := &data.OrderRepository{c}
	orders = repo.GetForUser(user)
	if len(orders) == 0 {
		err = errors.New("No orders found for user.")
	}
	return
}

func NewTermsPayload(ov *accounting.OrderValues, uuid, flash string) interface{} {
	termsPayload := struct {
		MonthlyPayment interface{}
		Total          interface{}
		FirstPayment   interface{}
		UUID           interface{}
		PublishableKey string
		Flash          string
	}{
		ov.MonthlyPaymentFmt,
		ov.PriceFmt,
		ov.InitialPaymentFmt,
		uuid,
		os.Getenv("STRIPE_PUB_KEY"),
		flash,
	}
	return termsPayload
}

// Helper method for retrieving the web order from the data store.
func WebOrderForUUID(uuid string) (*models.WebOrder, error) {
	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}
	webOrder, err := repo.GetByUUID(uuid)

	return webOrder, err

}

func AllWebOrders() ([]models.WebOrder, error) {
	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}
	return repo.GetAll()

}

func ParseFlash(req *http.Request) string {
	flash := req.URL.Query().Get("err")
	var flashMessage string
	switch flash {
	case "1":
		flashMessage = "Incorrect card number. Please enter the correct number, or enter a different card"
	case "2":
		flashMessage = "Invalid card"
	case "3":
		flashMessage = "Invalid Expiration Month."
	case "4":
		flashMessage = "Invalid expiration month"
	case "5":
		flashMessage = "Invalid CVC"
	case "6":
		flashMessage = "Expired Card"
	case "7":
		flashMessage = "Incorrect CVC"
	case "8":
		flashMessage = "Incorrect ZIP"
	case "9":
		flashMessage = "Card declined, please try a different card."
	case "10":
		flashMessage = "There was an error please try again."
	case "11":
		flashMessage = "Process error, please try again later."
	case "12":
		flashMessage = "There was an error proccessing your card. Please re-enter card details"
	default:
		flashMessage = ""
	}
	return flashMessage
}

func saveOrder(order *models.WebOrder) (err error) {
	context := controllers.NewContext()
	defer context.Close()

	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}

	Trace.Printf("About to save web order %s to database", order.UUID)
	if err = repo.NewWebOrder(order); err != nil {
		Error.Println(err.Error())
		return err
	}

	return nil
}

func AddSubscriberToMailChimp(email string) error {

	if email == "" {
		return errors.New("Cannot use blank email to subscribe.")
	}
	apiKey := os.Getenv("MAILCHIMP_KEY")

	if apiKey == "" {
		Error.Println("API Key not found for this environment, adding of subscriber will fail.")
	}

	client, err := mailchimp.NewClient(apiKey, nil)
	if err != nil {
		return err
	}
	_, err = client.Subscribe(email, "3f1750ea4a")
	if err != nil {
		return err
	}

	return nil

}
