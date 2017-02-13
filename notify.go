package main

import (
	"./models"
	"bytes"
	"fmt"

	"net/http"
	"os"
)

const (
	ServiceFee     = 0.10
	ProdWebHookURL = "https://hooks.slack.com/services/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"

	TestWebHookURL = "https://hooks.slack.com/services/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"
)

type Payload struct {
	Msg string
}

func (p *Payload) Print() string {
	return fmt.Sprintf(`{"text":"%s"}`,
		p.Msg)
}

func postOrderToSlack(o *models.WebOrder) {

	env := os.Getenv("DEVELOPMENT")
	if env == "1" {
		Trace.Println(" Skipping Slack, due to being in development environment")
		return
	}

	url := os.Getenv("GLASS_URL")
	msg := fmt.Sprintf("New web application submission from:" + o.FirstName + " " + o.LastName + " https://" + url + "/admin/orders/" + o.UUID)
	Info.Println(msg)
	p := &Payload{msg}
	jsonStr := p.Print()

	req, err := http.NewRequest("POST", TestWebHookURL, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		Error.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		Error.Println(err)
	}

}
