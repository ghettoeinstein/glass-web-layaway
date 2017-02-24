package main

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func test() {
	price, err := decimal.NewFromString("99.99")
	if err != nil {
		panic(err)
	}

	quantity := decimal.NewFromFloat(2)

	fee, _ := decimal.NewFromString(".1")
	taxRate, _ := decimal.NewFromString(".08875")

	subtotal := price.Mul(quantity)

	preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1))).Truncate(2)

	total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1)))

	fmt.Println("Subtotal:", subtotal)                               // Subtotal: 408.06
	fmt.Println("Pre-tax:", preTax)                                  // Pre-tax: 422.3421
	fmt.Println("Taxes:", total.Sub(preTax).Round(2))                // Taxes: 37.482861375
	fmt.Println("Total:", total)                                     // Total: 459.824961375
	fmt.Println("Tax rate:", total.Sub(preTax).Div(preTax).Round(5)) // Tax rate: 0.08875
}
