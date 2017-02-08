package controllers

import (
	"github.com/gorilla/mux"
	_ "github.com/joiggama/money"
	"net/http"
)

func IdFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	id := vars["id"]
	return id
}

// Mandrill Transactional methods
// Sends a one off email with Mandrill

func sendConf(order *models.Order, r *http.Request) error {

	user, _ := UserFromRequest(r)
	log.Println(user)
	client := m.ClientWithKey("_ZMKw0PeBC3p8jFsROTb7g")

	message := &m.Message{}
	message.AddRecipient(user.Email, "", "to")
	message.FromEmail = "receipts@getglass.co"
	message.FromName = "Glass Financial"
	message.Subject = "Order Confirmation from Glass"

	message.GlobalMergeVars = m.MapToVars(map[string]interface{}{
		"TODAYS_DATE":    time.Now().Format("01/02/06"),
		"PAYMENT_AMT":    order.MonthlyPaymentFmt,
		"ADDRESS_1":      user.Address1,
		"ADDRESS_2":      user.Address2,
		"ORDER_NO":       order.Id,
		"SALES_TAX":      order.SalesTaxFmt,
		"FIRST_PAYMENT":  order.FirstPaymentDue,
		"SECOND_PAYMENT": order.SecondPaymentDue,
		"THIRD_PAYMENT":  order.ThirdPaymentDue,
		"TOTAL":          order.TotalFmt,
	})

	_, err := client.MessagesSendTemplate(message, "new-receipt", nil)
	if err != nil {
		println(err.Error())
		return (err)
	}
	return nil
}
