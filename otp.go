package main

import (
	"./common"
	"./controllers"
	"./data"
	"./models"
	"./random"
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

type OTP struct {
	PhoneNumber string    `json:"phone_number"`
	Passcode    string    `json:"otp_passcode"`
	Expiry      time.Time `json:"expiration_time"`
}

func (otp OTP) Valid() bool {

	return otp.Expiry.Before(time.Now())
}

func (otp OTP) String() string {

	return fmt.Sprintf("%s\n %s\n %s\n", otp.PhoneNumber, otp.Expiry.String(), otp.Passcode)

}

var otpStore = make(map[string]OTP)

type MyCustomClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// Make a OTP(one-time-passcode). Return a struct of type OTP
func GenerateOTP(length int, phone string) OTP {
	pin := random.RandomPin(length)

	otp := OTP{Passcode: pin, Expiry: time.Now().Add(time.Second * 180), PhoneNumber: phone}
	log.Println("Generated OTP:", otp.Passcode, otp.Expiry.String, otp.PhoneNumber)
	return otp
}

func UserByNumber(phoneNumber string) (user *models.User, err error) {
	context := controllers.NewContext()
	defer context.Close()
	c := context.DbCollection("users")

	repo := &data.UserRepository{c}
	user, err = repo.GetByPhoneNumber(phoneNumber)
	return
}

func GenerateOTPAppHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//Generate a OTP of suplied length length for given phone number

	phoneNumber := r.FormValue("phoneNumber")

	if phoneNumber == "" {
		common.DisplayAppError(w, errors.New("Phone number is blank"), "Can't lookup user", 500)
		return
	}

	user, err := UserByNumber(phoneNumber)
	if err == nil && user.Email != "" {
		log.Println("Found a user for phoneNumber", phoneNumber)
		log.Println(user)
		common.DisplayAppError(w, errors.New("Already registered"), "User exists", 200)
		return
	}

	if phoneNumber != "" {
		otp := GenerateOTP(4, phoneNumber)
		otpStore[phoneNumber] = otp
		log.Println(otpStore[phoneNumber])
		if _, _, err := sendMessage(phoneNumber, "Your Glass verification code is: "+otp.Passcode); err != nil {
			log.Fatal("Error sending msg: ", err)
		}
		json_otp, err := json.MarshalIndent(otp, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_otp)

	}
}

func GenerateOTPHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//Generate a OTP of suplied length length for given phone number

	phoneNumber := r.FormValue("phoneNumber")
	if phoneNumber != "" {
		otp := GenerateOTP(6, phoneNumber)
		otpStore[phoneNumber] = otp
		log.Println(otpStore[phoneNumber])
		if _, _, err := sendMessage(phoneNumber, "Your Glass verification code is: "+otp.Passcode); err != nil {
			log.Fatal("Error sending msg: ", err)
		}
		json_otp, err := json.MarshalIndent(otp, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_otp)

	}
}

func verifyOTP(passcode, phoneNumber string) bool {
	if otp, exists := otpStore[phoneNumber]; exists {
		if passcode == otp.Passcode {
			log.Println("Match is true")
			delete(otpStore, phoneNumber)
			return true
		}
		if passcode == "" {
			log.Println("passcode is empty")
			return false
		}

		log.Printf("Passcode is %s", passcode)
		log.Printf("Phone number is %s", phoneNumber)

		return false
	}

	log.Println("returning false from verifyOTP method")
	return false
}

//func verifyOTPHandler(w http.ResponseWriter, r *http.Request) {
//	err := r.ParseForm()
//	if err != nil {
//
//		log.Println("error parsing form")
//		common.DisplayAppError(w, err, "Error parsing form", 500)
//	}
//
//	phoneNumber := r.FormValue("phoneNumber")
//	verCode := r.FormValue("passcode")
//	log.Println(phoneNumber + " is the phone number")
//	log.Println(verCode + " is the code from form")
//	ok := verifyOTP(verCode, phoneNumber)
//
//	// User submitted correct 4 - digit pin
//	if ok {
//		log.Println("It was okay, lets make a token for regstration user")
//		token, err := common.GenerateJWT(phoneNumber, "member")
//		if err != nil {
//			log.Println("Error creating JWT for %s %s", phoneNumber, err)
//		}
//		// You have a token send via JSON
//		user_token := Token{token}
//		jsonResponse(user_token, w)
//		return
//	} else {
//		log.Println("Wrong phone number")
//		w.Header().Set("Content-Type", "application/json")
//		failAuthResp := &Response{"Your verification code is incorrect"}
//		if err := json.NewEncoder(w).Encode(failAuthResp); err != nil {
//			log.Fatal("Error: ", err)
//			http.Error(w, err.Error(), 500)
//			return
//		}
//
//	}
//}
//func OTPLoginHandler(w http.ResponseWriter, r *http.Request) {
//	err := r.ParseForm()
//	if err != nil {
//
//		log.Println("error parsing form")
//		common.DisplayAppError(w, err, "Error parsing form", 500)
//	}
//
//	phoneNumber := r.FormValue("phoneNumber")
//	verCode := r.FormValue("passcode")
//	log.Println(phoneNumber + " is the phone number")
//	log.Println(verCode + " is the code from form")
//	ok := verifyOTP(verCode, phoneNumber)
//
//	// User submitted correct 4 - digit pin
//	if ok {
//		log.Println("It was okay, lets make a token")
//		//Issue a token for their phone number
//		user, err := UserByNumber(phoneNumber)
//		if err != nil {
//			log.Println("Phone Number")
//			common.DisplayAppError(w, err, "Could not find user to verify", 500)
//			return
//		}
//
//		token, err := common.GenerateJWT(user.Email, "member")
//		if err != nil {
//			log.Println("Error creating JWT for %s %s", user.Email, err)
//		}
//		// You have a token send via JSON
//		userObject := struct {
//			Token Token        `json:"token"`
//			User  *models.User `json:"user"`
//		}{
//			Token{token},
//			user,
//		}
//
//		jsonResponse(userObject, w)
//		return
//	} else {
//		log.Println("Wrong phone number")
//		w.Header().Set("Content-Type", "application/json")
//		failAuthResp := &Response{"Your verification code is incorrect"}
//		if err := json.NewEncoder(w).Encode(failAuthResp); err != nil {
//			log.Fatal("Error: ", err)
//			http.Error(w, err.Error(), 500)
//			return
//		}
//
//	}
//}
//
