package data

import (
	"../models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserRepository struct {
	C *mgo.Collection
}

func (r *UserRepository) CreateUser(user *models.User) error {
	obj_id := bson.NewObjectId()
	user.Id = obj_id

	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.HashPassword = hpass
	//clear the incoming text password

	user.Password = ""

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	err = r.C.Insert(&user)
	return err
}

func (r *UserRepository) Update(user *models.User) error {
	// partial update on MongoDB
	err := r.C.Update(bson.M{"_id": user.Id},
		bson.M{"$set": bson.M{
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"email":           user.Email,
			"date_of_birth":   user.DOB,
			"address_1":       user.Address1,
			"address_2":       user.Address2,
			"city":            user.City,
			"state":           user.State,
			"zip_code":        user.ZipCode,
			"phone_number":    user.PhoneNumber,
			"updated_at":      time.Now(),
			"stripe_customer": user.StripeCustomer,
		}})
	return err
}

func (r *UserRepository) Login(user models.User) (u models.User, err error) {

	err = r.C.Find(bson.M{"email": user.Email}).One(&u)
	if err != nil {
		return
	}
	// Validate password
	err = bcrypt.CompareHashAndPassword(u.HashPassword, []byte(user.Password))
	if err != nil {
		u = models.User{}
	}
	return
}

func (r *UserRepository) GetByUsername(username string) (user *models.User, err error) {
	err = r.C.Find(bson.M{"email": username}).One(&user)
	return
}

func (r *UserRepository) GetByPhoneNumber(phoneNumber string) (user *models.User, err error) {
	err = r.C.Find(bson.M{"phone_number": phoneNumber}).One(&user)
	return
}
