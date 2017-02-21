package main

import (
    "fmt"
    "github.com/shopspring/decimal"
)





func test() {
    price, err := decimal.NewFromString("1000.00")
    if err != nil {
        panic(err)
    }

    quantity := decimal.NewFromFloat(1)

    fee, _ := decimal.NewFromString(".45")
    taxRate, _ := decimal.NewFromString(".08875")

    subtotal := price.Mul(quantity)

    preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1)))

    total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1)))

    fmt.Println("Subtotal:", subtotal)                      // Subtotal: 408.06
    fmt.Println("Pre-tax:", preTax)                         // Pre-tax: 422.3421
    fmt.Println("Taxes:", total.Sub(preTax))                // Taxes: 37.482861375
    fmt.Println("Total:", total)                            // Total: 459.824961375
    fmt.Println("Tax rate:", total.Sub(preTax).Div(preTax)) // Tax rate: 0.08875
}
