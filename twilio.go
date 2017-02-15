package main

import (
	"./models"
	"fmt"
	"github.com/subosito/twilio"

	"os"
)

var (
	AccountSid = "ACc4e7478e390bbe34b39a38ea94c3f259"
	AuthToken  = "070ceb76a3921552d3f76fea270cd1c2"
)

const (
	from            = "+16502065606"
	passcode_length = 4
)

var admins = []string{"+13234236654", "+17185219161"}

var url = os.Getenv("GLASS_URL")
var order_url = ".https://" + url + "getglass.co/admin/orders/"

func textOrderToAdmins(o *models.WebOrder) {
	for _, admin := range admins {
		sendMessage(admin, fmt.Sprintf("New web application submission from:"+o.FirstName+"  "+o.LastName+".https://"+url+"/admin/orders/"+o.UUID+"."))
	}
}

func sendMessage(to, msg string) (*twilio.Message, bool, error) {
	env := os.Getenv("DEVELOPMENT")
	if env == "1" {
		Trace.Println(" Skipping SMS, due to being in development environment")
		return nil, false, nil
	}

	c := twilio.NewClient(AccountSid, AuthToken, nil)

	from := "+16502065606"
	params := twilio.MessageParams{
		Body: msg,
	}
	s, resp, err := c.Messages.Send(from, to, params)
	if err != nil {
		Error.Println(s, resp, err)
		return nil, false, err
	}

	return s, true, nil

}
