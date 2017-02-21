package controllers

import (
	"../models"
	"testing"
)

func TestSendConf(t *testing.T) {

	order := &models.Order{}

	if err := SendConf(order); err != nil {
		t.Errorf("Test failed")
	}
}
