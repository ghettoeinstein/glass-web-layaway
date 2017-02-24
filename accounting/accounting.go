package accounting

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const (
	TaxRate = 0.085
)

func ConvertIntPrice(p float64) interface{} {
	formatted := (p * 100) / 100
	return formatted
}

// From https://gist.github.com/DavidVaini/10308388
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func ServiceFee(price, serviceFee decimal.Decimal) decimal.Decimal {

	return price.Mul(serviceFee)
}

type OrderValues struct {
	Price             int64
	PriceFmt          string
	ServiceFee        int64
	ServiceFeeFmt     string
	Taxes             int64
	TaxesFmt          string
	Total             int64
	TotalFmt          string
	InitialPayment    int64
	InitialPaymentFmt string
	MonthlyPayment    int64
	MonthlyPaymentFmt string
}

func (ov OrderValues) String() string {
	return fmt.Sprintf(" Price: %s\n Svc Fee: %s\n Taxes: %s\n Total: %s\n Initial Payment: %s\n Monthly Payment: %s\n",
		ov.PriceFmt,
		ov.ServiceFeeFmt,
		ov.TaxesFmt,
		ov.TotalFmt,
		ov.InitialPaymentFmt,
		ov.MonthlyPaymentFmt,
	)
}

func OrderValuesFromPrice(p string) *OrderValues {

	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	price, _ := decimal.NewFromString(p)
	priceFmt, _ := price.Float64()
	serviceFee, _ := decimal.NewFromString("0.10")

	taxRate := decimal.NewFromFloat(TaxRate)
	taxtotal := price.Mul(taxRate.Add(decimal.NewFromFloat(1)))
	taxtotal.Float64()

	taxes := taxtotal.Sub(price)
	taxesFmt, _ := taxes.Float64()

	svcFee := ServiceFee(price, serviceFee)
	svcFmt, _ := svcFee.Float64()

	total := svcFee.Add(price).Add(taxRate.Mul(price))

	grandTotal, _ := total.Float64()

	dueToday := price.Div(decimal.NewFromFloat(float64(4))).Add(svcFee).Add(taxtotal.Sub(price))
	dueFloat, _ := dueToday.Float64()

	var f float64
	subtotal := price.Div(decimal.NewFromFloat(float64(1)))
	f, _ = subtotal.Float64()

	//		difference := subtotal.Sub(price)

	priceInt := int64(Round(priceFmt*100, .5, 0))
	svcFeeInt := int64(Round(svcFmt*100, .5, 0))
	taxInt := int64(Round(taxesFmt*100, .5, 0))
	totalInt := int64(Round(grandTotal*100, .5, 0))
	dueTodayInt := int64(Round(dueFloat*100, .5, 0))
	monthlyPayment := int64(Round((f/4)*100, .5, 0))

	ov := &OrderValues{
		Price:             priceInt,
		PriceFmt:          ac.FormatMoney(priceFmt),
		ServiceFee:        svcFeeInt,
		ServiceFeeFmt:     ac.FormatMoney(svcFmt),
		Taxes:             taxInt,
		TaxesFmt:          ac.FormatMoney(taxesFmt),
		Total:             totalInt,
		TotalFmt:          ac.FormatMoney(grandTotal),
		InitialPayment:    dueTodayInt,
		InitialPaymentFmt: ac.FormatMoney(dueFloat),
		MonthlyPayment:    monthlyPayment,
		MonthlyPaymentFmt: ac.FormatMoney(f / 4),
	}
	return ov
}

func (ov OrderValues) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, ov)
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed%s %s  in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func loggingHandler2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("started logging handler 2")
		next.ServeHTTP(w, r)
		log.Printf("Completed logging handler 2")
	})
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	http.HandleFunc("/favicon.ico", faviconHandler)

	orderValues := OrderValuesFromPrice(os.Args[1])
	http.Handle("/", loggingHandler2(loggingHandler(orderValues)))
	http.ListenAndServe(":8000", nil)
}
