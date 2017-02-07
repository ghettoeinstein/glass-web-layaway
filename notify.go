package main

import (
	"./models"
	"bytes"
	"fmt"
	"log"
	"net/http"
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

func postSlack(o *models.WebOrder) {

	msg := fmt.Sprintf("New web application submission http://localhost:9090/admin/orders/" + o.UUID)
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

func postOrderToSlack(o *models.WebOrder) {

	msg := fmt.Sprintf("New web application submission by  http://localhost:9090/admin/orders/" + o.UUID)
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
