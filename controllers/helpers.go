package controllers

import (
	"../models"
	"github.com/gorilla/mux"
	"github.com/joiggama/money"
	m "github.com/keighl/mandrill"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Stripe *log.Logger
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
		"PAID_TODAY":            money.Format(order.MonthlyPayment + order.ServiceFee + order.SalesTax),
		"TODAYS_DATE":           time.Now().Format("01/02/06"),
		"PHONE_NUMBER":          order.User.PhoneNumber,
		"MONTHLY_FMT":           money.Format(order.MonthlyPayment),
		"BALANCE_POST_FIRST":    money.Format(order.BalancePostFirst),
		"BALANCE_POST_SECOND":   money.Format(order.BalancePostSecond),
		"BALANCE_POST_CREATION": money.Format(order.BalancePostCreation),
		"ORDERNO":               order.UUID,
		"FULLNAME":              order.User.FirstName + " " + order.User.LastName,
		"FIRST_PAYMENT":         order.CreatedAt.Add(time.Duration(time.Hour * 24 * 30)).Format("01/02/06"),
		"SECOND_PAYMENT":        order.SecondPaymentDue,
		"THIRD_PAYMENT":         order.ThirdPaymentDue,
		"GRAND_TOTAL":           money.Format(order.Total),
		"SERVICE_FEE":           money.Format(order.ServiceFee),
		"TAXES":                 money.Format(order.SalesTax),
		"COMBINED_TOTAL":        money.Format(order.CombinedTotal),
	})
	_, err := client.MessagesSendTemplate(message, "new-receipt", nil)
	if err != nil {
		println(err.Error())
		return (err)
	}
	return nil
}
