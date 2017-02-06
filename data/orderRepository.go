package data

import (
	"../models"
	"time"
	//  "errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//  "log"
)

type OrderRepository struct {
	C *mgo.Collection
}

func (r *OrderRepository) CreateOrderForUser(user *models.User) (order *models.Order, err error) {
	order = &models.Order{}
	err = r.NewOrder(order)

	order.User = user

	err = r.C.Update(bson.M{"_id": order.Id},
		bson.M{"$set": bson.M{
			"items":               order.Items,
			"shipped":             order.Shipped,
			"shipping_address":    order.ShippingAddress,
			"trackingNumber":      order.TrackingNumber,
			"total":               order.Total,
			"tax_rate":            order.TaxRate,
			"url":                 order.URL,
			"missed_deadline":     order.MissedDeadline,
			"email":               order.Email,
			"first_payment_due":   order.FirstPaymentDue,
			"first_payment_paid":  order.FirstPaymentPaid,
			"second_payment_due":  order.SecondPaymentDue,
			"second_payment_paid": order.SecondPaymentPaid,
			"third_payment_due":   order.ThirdPaymentDue,
			"third_payment_paid":  order.ThirdPaymentPaid,
			"monthly_payment":     order.MonthlyPayment,
			"balance":             order.Balance,
			"total_fmt":           order.TotalFmt,
			"monthly_payment_fmt": order.MonthlyPaymentFmt,
			"order_date":          order.OrderDate,
			"updated_at":          time.Now(),
		}})
	return
}

func (r *OrderRepository) NewOrder(order *models.Order) (err error) {
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	err = r.C.Insert(&order)
	return
}

func (r *OrderRepository) CreateOrderFromOrder(order *models.Order) error {
	obj_id := bson.NewObjectId()
	order.Id = obj_id

	err := r.C.Insert(&order)
	return err
}

func (r *OrderRepository) Update(order *models.Order) (err error) {

	// partial update on MongoDB

	err = r.C.Update(bson.M{"_id": order.Id},
		bson.M{"$set": bson.M{
			"missed_deadline":     order.MissedDeadline,
			"items":               order.Items,
			"shipped":             order.Shipped,
			"shipping_address":    order.ShippingAddress,
			"trackingNumber":      order.TrackingNumber,
			"total":               order.Total,
			"tax_rate":            order.TaxRate,
			"email":               order.Email,
			"uuid":                order.UUID,
			"first_payment_due":   order.FirstPaymentDue,
			"first_payment_paid":  order.FirstPaymentPaid,
			"second_payment_due":  order.SecondPaymentDue,
			"second_payment_paid": order.SecondPaymentPaid,
			"third_payment_due":   order.ThirdPaymentDue,
			"third_payment_paid":  order.ThirdPaymentPaid,
			"monthly_payment":     order.MonthlyPayment,
			"balance":             order.Balance,
			"total_fmt":           order.TotalFmt,
			"monthly_payment_fmt": order.MonthlyPaymentFmt,
			"url":        order.URL,
			"updated_at": time.Now(),
		}})

	return err
}

func (r *OrderRepository) GetNewOrders() ([]models.Order, error) {
	var orders []models.Order
	iter := r.C.Find(bson.M{}).Iter()
	result := models.Order{}
	for iter.Next(&result) {
		orders = append(orders, result)
	}
	return orders, nil
}

func (r *OrderRepository) GetForUser(user *models.User) []models.Order {
	var orders []models.Order
	iter := r.C.Find(bson.M{"email": user.Email}).Iter()
	result := models.Order{}
	for iter.Next(&result) {
		orders = append(orders, result)
	}
	return orders
}

func (r *OrderRepository) GetById(id string) (order *models.Order, err error) {
	err = r.C.FindId(bson.ObjectIdHex(id)).One(&order)
	return

}
