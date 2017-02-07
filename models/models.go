package models

import (
	"fmt"
	"time"
	//	"github.com/stripe/stripe-go/plan"
	"gopkg.in/mgo.v2/bson"
)

type (
	WebOrder struct {
		Id             bson.ObjectId `bson:"_id",omitempty" json:"id`
		FullName       string        `json:"full_name" bson:"full_name`
		Email          string        `json:"email" bson:"email"`
		PhoneNumber    string        `json:"phone_number" bson:"phone_number`
		URL            string        `json:"url" bson:"url"`
		Address        string        `json:"address" bson:"address"`
		UUID           string        `json:"uuid" bson:"uuid"`
		DateCreated    time.Time     `json:"date_created" bson:"date_created"`
		Decision       string        `json:"decision" bson:"decision"`
		TimerExpiresAt time.Time     `json:"timer_expires_at" bson:"timer_expires_at"`
		CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
		UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
		Acknowledged   bool          `json:"acknowledged" bson:"acknowledged"`
		Price          int           `json:"price" bson:"price"`
	}
	User struct {
		Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
		FirstName   string        `json:"first_name" bson:"first_name"`
		LastName    string        `json:"last_name"  bson:"last_name"`
		Email       string        `json:"email"      bson:"email"`
		PhoneNumber string        `json:"phone_number" bson:"phone_number"`
		Password    string        `json:"password,omitempty"`
		Address1    string        `json:"address_1" bson:"address_1"`
		Address2    string        `json:"address_2" bson:"address_2"`
		City        string        `json:"city" bson:"city"`

		State        string `json:"state" bson:"state"`
		ZipCode      string `json:"zip_code" bson:"zip_code"`
		DOB          string `json:"dob" 	bson:date_of_birth"`
		HashPassword []byte `json:"hashpassword,omitempty "`
		StripeCustomer
		CreatedAt time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	}
	Admin struct {
		User
	}

	Category struct {
		Name string
	}

	Merchant struct {
		Id        bson.ObjectId `json:"merchant_id" bson:"_id,omitempty"`
		Website   string        `json:"website" bson:"website"`
		Category  string        `json:"category" bson:"category" `
		LogoURL   string        `json:"logo_url"    bson:"logo_url"`
		BrandName string        `json:"brand_name"  bson:"brand_name"`
		Active    bool          `json:"active"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}
	Item struct {
		ProductId   string `json:"product_id"`
		URL         string `json:"url" bson:"url"`
		Price       string `json:"price"  bson:"price"`
		Quantity    string `json:"quantity"  bson:"quantity"`
		ProductName string `json:"product_name" bson:"product_name`
	}

	Order struct {
		Id                bson.ObjectId  `json:"order_id" bson:"_id,omitempty"`
		Balance           int            `json:"balance" bson:"balance"`
		User              *User          `json:"user" bson:"user"`
		Shipped           bool           `json:"shipped" bson:"shipped"`
		ShippingAddress   string         `json:"shipping_address" bson:"shipping_address"`
		TrackingNumber    string         `json:"tracking_number" bson:"tracking_number"`
		Items             []Item         `json:"items" bson:"items"`
		Total             int            `json:"total" bson:"total"`
		TotalFmt          string         `json:"total_fmt" bson:"total_fmt"`
		TaxRate           float64        `json:"tax_rate" bson:"tax_rate"`
		OrderDate         string         `json:"order_date" bson:"order_date"`
		StripePlanID      string         `json:"stripe_plan_id" bson:"stripe_plan_id"`
		CreatedAt         time.Time      `json:"created_at" bson:"created_at"`
		UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
		SalesTax          int            `json:"sales_tax" bson:"sales_tax"`
		SalesTaxFmt       string         `json:"sales_tax_fmt" bson:"sales_tax_fmt"`
		Email             string         `json:"email" bson:"email"`
		UUID              string         `json:"uuid" bson:"uuid"`
		FirstPaymentPaid  bool           `json:"first_payment"  bson:"first_payment"`
		FirstPaymentDue   string         `json:"first_payment"  bson:"first_payment_due"`
		SecondPaymentPaid bool           `json:"second_payment" bson:"second_payment"`
		SecondPaymentDue  string         `json:"second_payment" bson:"second_payment_due"`
		ThirdPaymentPaid  bool           `json:"third_payment" bson:"third_payment"`
		ThirdPaymentDue   string         `json:"third_payment" bson:"third_payment_due"`
		MonthlyPayment    int            `json:"monthly_payment" bson:"monthly_payment"`
		MonthlyPaymentFmt string         `json:"monthly_payment_fmt" bson:"monthly_payment_fmt"`
		Payments          []Payment      `json:"payments" bson:"payments"`
		Customer          StripeCustomer `json:"customer" bson:"customer"`
		MissedDeadline    bool           `json:"missed_deadline" bson:"missed_deadline"`
		URL               string         `json:"url" bson:"url"`
	}

	Payment struct {
		ChargeID    string    `json:"-"`
		DateCharged time.Time `json:"date_charged"`
	}

	StripeCustomer struct {
		CustomerId string `json:"customer_id" bson:"customer_id"`
	}
	Role struct {
	}
)

func (u User) isAdmin() bool  { return false }
func (a Admin) isAdmin() bool { return true }

func (o Order) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s", o.User.FirstName, o.User.LastName, o.User.Email, o.Total, o.SalesTax, o.MonthlyPayment, o.FirstPaymentDue, o.SecondPaymentDue, o.ThirdPaymentDue)
}
