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

func (r *OrderRepository) CreateOrderFromWebOrder(webOrder *models.WebOrder) (order *models.Order, err error) {

	return
}

func (r *OrderRepository) CreateOrderForUser(user *models.User) (order *models.Order, err error) {
	order = &models.Order{}
	err = r.NewOrder(order)
	if err != nil {
		return nil, err
	}
	order.User = user

	err = r.C.Update(bson.M{"_id": order.Id},
		bson.M{"$set": bson.M{
			"uuid":                  order.UUID,
			"user":                  order.User,
			"shipped":               order.Shipped,
			"shipping_address":      order.ShippingAddress,
			"trackingNumber":        order.TrackingNumber,
			"total":                 order.Total,
			"tax_rate":              order.TaxRate,
			"url":                   order.URL,
			"missed_deadline":       order.MissedDeadline,
			"email":                 order.Email,
			"first_payment_due":     order.FirstPaymentDue,
			"first_payment_paid":    order.FirstPaymentPaid,
			"second_payment_due":    order.SecondPaymentDue,
			"second_payment_paid":   order.SecondPaymentPaid,
			"third_payment_due":     order.ThirdPaymentDue,
			"third_payment_paid":    order.ThirdPaymentPaid,
			"monthly_payment":       order.MonthlyPayment,
			"balance_post_creation": order.BalancePostCreation,
			"total_fmt":             order.TotalFmt,
			"monthly_payment_fmt":   order.MonthlyPaymentFmt,
			"order_date":            order.OrderDate,
			"combined_total":        order.CombinedTotal,
			"updated_at":            time.Now(),
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

func (r *OrderRepository) SaveOrder(order *models.Order) error {
	obj_id := bson.NewObjectId()
	order.Id = obj_id
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	err := r.C.Insert(&order)
	return err
}

func (r *OrderRepository) Update(order *models.Order) (err error) {

	// partial update on MongoDB

	err = r.C.Update(bson.M{"_id": order.Id},
		bson.M{"$set": bson.M{
			"missed_deadline":       order.MissedDeadline,
			"items":                 order.Items,
			"shipped":               order.Shipped,
			"shipping_address":      order.ShippingAddress,
			"trackingNumber":        order.TrackingNumber,
			"total":                 order.Total,
			"tax_rate":              order.TaxRate,
			"email":                 order.Email,
			"uuid":                  order.UUID,
			"first_payment_due":     order.FirstPaymentDue,
			"first_payment_paid":    order.FirstPaymentPaid,
			"second_payment_due":    order.SecondPaymentDue,
			"second_payment_paid":   order.SecondPaymentPaid,
			"third_payment_due":     order.ThirdPaymentDue,
			"third_payment_paid":    order.ThirdPaymentPaid,
			"monthly_payment":       order.MonthlyPayment,
			"balance_post_creation": order.BalancePostCreation,
			"balance_post_first":    order.BalancePostFirst,
			"balance_post_second":   order.BalancePostSecond,
			"total_fmt":             order.TotalFmt,
			"monthly_payment_fmt":   order.MonthlyPaymentFmt,
			"sales_tax":             order.SalesTax,
			"combined_total":        order.CombinedTotal,
			"url":                   order.URL,
			"updated_at":            time.Now(),
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

func (r *OrderRepository) GetForUser(user *models.User) []*models.Order {
	var orders []*models.Order
	iter := r.C.Find(bson.M{"email": user.Email}).Iter()
	result := &models.Order{}
	for iter.Next(&result) {
		orders = append(orders, result)
	}
	return orders
}

func (r *OrderRepository) GetByUUID(uuid string) (order *models.Order, err error) {
	err = r.C.Find(bson.M{"uuid": uuid}).One(&order)
	return

}

func (r *OrderRepository) GetById(id string) (order *models.Order, err error) {
	err = r.C.FindId(bson.ObjectIdHex(id)).One(&order)
	return

}
