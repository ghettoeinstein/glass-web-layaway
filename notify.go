package main

import (
	"./models"
	"bytes"
	"fmt"
	"log"
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
	url := os.Getenv("GLASS_URL")
	msg := fmt.Sprintf("New web application submission from:" + o.FullName + " https://" + url + "/admin/orders/" + o.UUID)
	log.Println(msg)
	p := &Payload{msg}
	jsonStr := p.Print()
	log.Println(jsonStr)

	req, err := http.NewRequest("POST", TestWebHookURL, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}

}
