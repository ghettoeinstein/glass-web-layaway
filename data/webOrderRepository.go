package data

import (
	"../models"
	"time"
	//  "errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//  Putting this here for now. Break this out later.
type WebOrderRepository struct {
	C *mgo.Collection
}

func (r *WebOrderRepository) NewWebOrder(order *models.WebOrder) (err error) {
	now := time.Now()
	order.Id = bson.NewObjectId()
	order.CreatedAt = now
	order.UpdatedAt = now
	err = r.C.Insert(&order)
	return
}

func (r *WebOrderRepository) GetNewOrders() ([]models.WebOrder, error) {
	var webOrders []models.WebOrder
	iter := r.C.Find(bson.M{"acknowledged": false}).Iter()

	result := models.WebOrder{}
	for iter.Next(&result) {
		webOrders = append(webOrders, result)
	}
	return webOrders, nil
}

func (r *WebOrderRepository) UpdateOrder(wo *models.WebOrder) (err error) {

	// partial update on MongoDB

	err = r.C.Update(bson.M{"_id": wo.Id},
		bson.M{"$set": bson.M{
			"first_name":   wo.FirstName,
			"last_name":    wo.LastName,
			"email":        wo.Email,
			"phone_number": wo.PhoneNumber,
			"url":          wo.URL,
			"price":        wo.Price,
			"decision":     wo.Decision,
			"acknowledged": wo.Acknowledged,
		}})

	return err
}

func (r *WebOrderRepository) GetApprovedOrders() ([]models.WebOrder, error) {
	var webOrders []models.WebOrder
	iter := r.C.Find(bson.M{"acknowledged": true, "decision": "approved"}).Iter()

	result := models.WebOrder{}
	for iter.Next(&result) {
		webOrders = append(webOrders, result)
	}
	return webOrders, nil
}

func (r *WebOrderRepository) GetDeniedOrders() ([]models.WebOrder, error) {
	var webOrders []models.WebOrder
	iter := r.C.Find(bson.M{"acknowledged": true, "decision": "denied"}).Iter()

	result := models.WebOrder{}
	for iter.Next(&result) {
		webOrders = append(webOrders, result)
	}
	return webOrders, nil
}

func (r *WebOrderRepository) GetByUUID(uuid string) (webOrder *models.WebOrder, err error) {
	err = r.C.Find(bson.M{"uuid": uuid}).One(&webOrder)
	return

}

func (r *WebOrderRepository) DeleteByUUID(uuid string) (err error) {
	err = r.C.Remove(bson.M{"uuid": uuid})
	return

}
