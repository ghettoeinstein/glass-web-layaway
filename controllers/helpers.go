package controllers

import (
	"../models"
	"github.com/gorilla/mux"
	"github.com/joiggama/money"
	m "github.com/keighl/mandrill"

	"net/http"
	"time"
)

func IdFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	id := vars["id"]
	return id
}

// Mandrill Transactional methods
// Sends a one off email with Mandrill

func sendConf(order *models.Order) error {

	client := m.ClientWithKey("_ZMKw0PeBC3p8jFsROTb7g")

	message := &m.Message{}
	message.AddRecipient(order.Email, "", "to")
	message.FromEmail = "receipts@getglass.co"
	message.FromName = "Glass Financial"
	message.Subject = "Order Confirmation from Glass"

	message.GlobalMergeVars = m.MapToVars(map[string]interface{}{
		"TODAYS_DATE":           time.Now().Format("01/02/06"),
		"PHONE_NUMBER":          order.User.PhoneNumber,
		"MONTHLY_FMT":           money.Format(order.MonthlyPayment),
		"BALANCE_POST_FIRST":    money.Format(order.BalancePostFirst),
		"BALANCE_POST_SECOND":   money.Format(order.BalancePostSecond),
		"BALANCE_POST_CREATION": money.Format(order.BalancePostCreation),
		"ORDERNO":               order.UUID,
		"FULLNAME":              order.User.FullName,
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
