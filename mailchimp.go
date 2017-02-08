package main

import (
	"errors"
	"github.com/mnbbrown/mailchimp"
	"os"
)

func AddSubscriberToMailChimp(email string) error {

	if email == "" {
		return errors.New("Cannot use blank email to subscribe.")
	}
	apiKey := os.Getenv("MAILCHIMP_KEY")

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
