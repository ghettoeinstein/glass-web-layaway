package controllers

import (
	"../models"
)

//Models for JSON resources
type (
	//For Post - /user/register
	UserResource struct {
		Data models.User `json:"data"`
	}

	OrderResource struct {
		Items []models.Item `json:"items"`

		OrderPayload OrderPayload `json:"data"`
	}

	OrderPayload struct {
		Address           string  `json:"address"`
		Email             string  `json:"email"`
		OSVersion         string  `json:"os_version"`
		MonthlyPayment    int     `json:"monthly_payment"`
		MonthlyPaymentFmt string  `json:"monthly_formatted"`
		ServiceFee        string  `json:"service_fee"`
		Amount            int     `json:"amount"`
		AppName           string  `json:"app_name"`
		IsPhone           bool    `json:"is_phone"`
		IsiPad            bool    `json:"is_ipad"`
		OrderTotal        int     `json:"total"`
		SalesTaxFmt       string  `json:"sales_tax_fmt"`
		SalesTax          int     `json:"sales_tax"`
		TotalFmt          string  `json:"total_formatted"`
		ItemCount         string  `json:"item_count"`
		TaxRate           float64 `json:"tax_rate"`
	}
	//For Post /user/login/app
	AppUserResource struct {
		Data  models.User `json:"data"`
		Token string      `json:"token"`
	}

	//For Post - /user/login
	LoginResource struct {
		Data LoginModel `json:"data"`
	}
	//Response for authorized user Post - /user/login
	AuthUserResource struct {
		Data AuthUserModel `json:"data"`
	}
	// For Post/Put - /tasks
	// For Get - /tasks/id

	MerchantResource struct {
		Data models.Merchant `json:"data"`
	}
	// For Get - /tasks
	MerchantsResource struct {
		Data []models.Merchant `json:"data"`
	}

	//Model for authentication
	LoginModel struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//Model for authorized user with access token
	AuthUserModel struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}
	OrderResponse struct {
		Count   int    `json:"count"`
		Message string `json:"message"`
	}
)
