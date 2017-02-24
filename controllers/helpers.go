package controllers

import (
	"../data"
	"../models"
	"github.com/gorilla/mux"
	_ "github.com/joiggama/money"
	m "github.com/keighl/mandrill"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Stripe *log.Logger

	Order *log.Logger
)

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
	logfile, err := os.OpenFile("logs/glassLogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error  opening up logfile %v", err)
	}

	log.SetOutput(logfile)
	Stripe = log.New(logfile,
		"[STRIPE] ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Order = log.New(logfile,
		"[ORDER] ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

func IdFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	id := vars["id"]
	return id
}

// Mandrill Transactional methods
// Sends a one off email with Mandrill

func SendConf(order *models.Order) error {

	client := m.ClientWithKey(os.Getenv("MANDRILL_KEY"))

	message := &m.Message{}
	message.AddRecipient(order.Email, "", "to")
	message.FromEmail = "receipts@getglass.co"
	message.FromName = "Glass Financial"
	message.Subject = "Order Confirmation from Glass"

	message.GlobalMergeVars = m.MapToVars(map[string]interface{}{
		//	"PAID_TODAY":            money.Format(order.MonthlyPayment + order.ServiceFee + order.SalesTax),
		"TODAYS_DATE":  order.CreatedAt.Format("01/02/06"),
		"PHONE_NUMBER": order.User.PhoneNumber,
		//	"MONTHLY_FMT":           money.Format(order.MonthlyPayment),
		//"BALANCE_POST_FIRST":    money.Format(order.BalancePostFirst),
		//"BALANCE_POST_SECOND":   money.Format(order.BalancePostSecond),
		//"BALANCE_POST_CREATION": money.Format(order.BalancePostCreation),
		"ORDERNO":        order.UUID,
		"FULLNAME":       order.User.FirstName + " " + order.User.LastName,
		"FIRST_PAYMENT":  order.CreatedAt.Add(time.Duration(time.Hour * 24 * 30)).Format("01/02/06"),
		"SECOND_PAYMENT": order.SecondPaymentDue,
		"THIRD_PAYMENT":  order.ThirdPaymentDue,
		//	"GRAND_TOTAL":           money.Format(order.Total),
		//	"SERVICE_FEE":           money.Format(order.ServiceFee),
		//	"TAXES":                 money.Format(order.SalesTax),
		//	"COMBINED_TOTAL":        money.Format(order.CombinedTotal),
	})
	_, err := client.MessagesSendTemplate(message, "new-receipt", nil)
	if err != nil {
		println(err.Error())
		return (err)
	}
	return nil
}

// Helper method for retrieving the web order from the data store.
func WebOrderForUUID(uuid string) (*models.WebOrder, error) {
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("web_orders")
	repo := &data.WebOrderRepository{c}
	webOrder, err := repo.GetByUUID(uuid)

	return webOrder, err

}

// Handle the case of the stripe errors on the terms page.
func HandleStripeError(w http.ResponseWriter, r *http.Request, id string, err error) {

	// Try to safely cast a generic error to a stripe.Error so that we can get at
	// some additional Stripe-specific information about what went wrong.
	if stripeErr, ok := err.(*stripe.Error); ok {
		// The Code field will contain a basic identifier for the failure.
		switch stripeErr.Code {
		case stripe.IncorrectNum:
			http.Redirect(w, r, "/terms/"+id+"?err=1", http.StatusSeeOther)
			return
		case stripe.InvalidNum:
			http.Redirect(w, r, "/terms/"+id+"?err=2", http.StatusSeeOther)
			return
		case stripe.InvalidExpM:
			http.Redirect(w, r, "/terms/"+id+"?err=3", http.StatusSeeOther)
			return
		case stripe.InvalidExpY:
			http.Redirect(w, r, "/terms/"+id+"?err=4", http.StatusSeeOther)
			return
		case stripe.InvalidCvc:
			http.Redirect(w, r, "/terms/"+id+"?err=5", http.StatusSeeOther)
			return
		case stripe.ExpiredCard:
			http.Redirect(w, r, "/terms/"+id+"?err=6", http.StatusSeeOther)
			return
		case stripe.IncorrectCvc:

			http.Redirect(w, r, "/terms/"+id+"?err=7", http.StatusSeeOther)
			return
		case stripe.IncorrectZip:
			http.Redirect(w, r, "/terms/"+id+"?err=8", http.StatusSeeOther)
			return
		case stripe.CardDeclined:
			http.Redirect(w, r, "/terms/"+id+"?err=9", http.StatusSeeOther)
			return
		case stripe.Missing:
			http.Redirect(w, r, "/terms/"+id+"?err=10", http.StatusSeeOther)
			return
		case stripe.ProcessingErr:
			http.Redirect(w, r, "/terms/"+id+"?err=11", http.StatusSeeOther)
			return
		}

		// The Err field can be coerced to a more specific error type with a type
		// assertion. This technique can be used to get more specialized
		// information for certain errors.
		if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
			Stripe.Printf("Card was declined with code: %v\n", cardErr.DeclineCode)
		} else {
			Stripe.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
		}
	} else {
		Stripe.Printf("Other error occurred: %v\n", err.Error())
	}
	return

}

// Create an order for a logged out/new user.
func NewOrder(webOrder *models.WebOrder, user *models.User) *models.Order {

	order := &models.Order{
		Total:               webOrder.Price().String(),
		BalancePostCreation: webOrder.Price().Mul(decimal.NewFromFloat(0.75)).String(),
		BalancePostFirst:    webOrder.Price().Div(decimal.NewFromFloat(2)).String(),
		BalancePostSecond:   webOrder.Price().Div(decimal.NewFromFloat(4)).String(),
		User:                user,
		Email:               webOrder.Email,
		URL:                 webOrder.URL,
		UUID:                webOrder.UUID,
		CustomerId:          user.StripeCustomer.CustomerId,
		SalesTax:            webOrder.Price().Mul(decimal.NewFromFloat(0.0875)).String(),
		MonthlyPayment:      webOrder.Price().Div(decimal.NewFromFloat(4)).String(),
		MonthlyPaymentFmt:   webOrder.Price().Div(decimal.NewFromFloat(4)).String(),
		FirstPaymentDue:     time.Now().Add(time.Hour * 24 * 30).Format("01/02/06"),
		SecondPaymentDue:    time.Now().Add(time.Hour * 24 * 60).Format("01/02/06"),
		ThirdPaymentDue:     time.Now().Add(time.Hour * 24 * 90).Format("01/02/06"),
	}

	return order

}

// Create a Order for a logged in user.
func NewOrderForUser(webOrder *models.WebOrder, user *models.User, cust *stripe.Customer) *models.Order {

	order := &models.Order{
		Total:               webOrder.Price().String(),
		BalancePostCreation: webOrder.Price().Mul(decimal.NewFromFloat(0.75)).String(),
		BalancePostFirst:    webOrder.Price().Div(decimal.NewFromFloat(2)).String(),
		BalancePostSecond:   webOrder.Price().Div(decimal.NewFromFloat(4)).String(),
		User:                user,
		Email:               webOrder.Email,
		URL:                 webOrder.URL,
		UUID:                webOrder.UUID,
		CustomerId:          cust.ID,
		SalesTax:            webOrder.Price().Mul(decimal.NewFromFloat(0.0875)).String(),
		MonthlyPayment:      webOrder.Price().Div(decimal.NewFromFloat(4)).String(),
		MonthlyPaymentFmt:   webOrder.Price().Div(decimal.NewFromFloat(4)).String(),
		FirstPaymentDue:     time.Now().Add(time.Hour * 24 * 30).Format("01/02/06"),
		SecondPaymentDue:    time.Now().Add(time.Hour * 24 * 60).Format("01/02/06"),
		ThirdPaymentDue:     time.Now().Add(time.Hour * 24 * 90).Format("01/02/06"),
	}

	return order

}

// Create a user from a webb order.
func NewUserFromWebOrder(o *models.WebOrder) *models.User {

	return &models.User{
		Id:          bson.NewObjectId(),
		Email:       o.Email,
		FirstName:   o.FirstName,
		LastName:    o.LastName,
		PhoneNumber: o.PhoneNumber,
	}
}
