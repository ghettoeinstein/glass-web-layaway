package main

import (
	"time"
)

type NewUser struct {
	FullName       string    `json:"full_name" bson:"full_name`
	Email          string    `json:"email" bson:"email"`
	PhoneNumber    string    `json:"phone_number" bson:"phone_number`
	URL            string    `json:"url" bson:"url"`
	Address        string    `json:"address" bson:"address"`
	UUID           string    `json:"uuid" bson:"uuid"`
	DateCreated    time.Time `json:"date_created" bson:"date_created"`
	Approved       bool      `json:"approved" bson:"approved"`
	TimerExpiresAt time.Time `json:"timer_expires_at" bson:"timer_expires_at"`
	State          bool
}

type Order struct {
	User        *NewUser  `json:"user" bson:"user"`
	UUID        string    `json:"uuid" bson:"uuid"`
	DateCreated time.Time `json:"date_created"`
	ExpireTime  time.Time `json:"expire_time" bson:"expire_time"`
	Expired     bool      `json:"expired" bson:"expired"`
	Redirect    string
}
