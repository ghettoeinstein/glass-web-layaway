package common

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func NewContextWithEmail(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, EmailKey, token.Claims.(*AppClaims).UserName)
}

type Key int

// using asymmetric crypto/RSA keyss
const (
	EmailKey Key = 0
	// openssl genrsa -out app.rsa 1024
	privKeyPath = "keys/app.rsa"
	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	pubKeyPath    = "keys/app.rsa.pub"
	sessionLength = 24 * 2 * time.Hour
)

// private key for signing and public key for verification
var (
	verifyKey, signKey []byte
	privateKey         *rsa.PrivateKey
	publicKey          *rsa.PublicKey
)

// Read the key files before starting http handlers
func initKeys() {
	var err error
	signKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("error reading private Key: ", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKey)
	if err != nil {
		log.Fatal("could not convert rsa private key")
	}

	verifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKey)
	if err != nil {
		log.Fatal("could not convert rsa key")
	}

}

type AppClaims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(name, role string) (string, error) {
	// Create the Claims
	claims := AppClaims{
		name,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    "admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

//Admin tokens are only good for 2 hours
func GenerateAdminJWT(name, role string) (string, error) {
	// Create the Claims
	claims := AppClaims{
		name,
		"Administrator",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    "admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

// Function for verifying the token's validity. We only sign with the private key(RSA), so this is all we need to verify the
// token has not been tampered with
func KeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return publicKey, nil
}

// Authorize Middleware for validating JWT tokens
func Authorize(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//Authorize for Web login
	if cookie, err := r.Cookie("Auth"); err != nil {
		if r.URL.Path == "/admin/" || r.URL.Path == "/admin" {
			http.Redirect(w, r, "/team/login", 307)
			return
		} else {
			http.Redirect(w, r, "/login", 307)
			return
		}
	} else {
		token, err := jwt.ParseWithClaims(cookie.Value, &AppClaims{}, KeyFunc)
		if err != nil {
			switch err.(type) {

			case *jwt.ValidationError: // JWT validation error
				vErr := err.(*jwt.ValidationError)

				switch vErr.Errors {
				case jwt.ValidationErrorExpired: //JWT expired
					http.Redirect(w, r, "/login", 307)
					return
				default:
					DisplayAppError(w,
						err,
						"Error while parsing the Access Token!",
						500,
					)
					return
				}

			default:
				DisplayAppError(w,
					err,
					"Error while parsing Access Token!",
					500)
				return
			}
		}
		if token.Valid {
			// Set user name to HTTP context

			ctx := NewContextWithEmail(r.Context(), token)
			next(w, r.WithContext(ctx))
			return
		} else {
			DisplayAppError(
				w,
				err,
				"Invalid Access Token",
				401,
			)
		}
		next(w, r)
		return
	}

	// Get  & validate token from request (API)
	token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &AppClaims{}, KeyFunc)

	if err != nil {
		switch err.(type) {

		case *jwt.ValidationError: // JWT validation error
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired: //JWT expired
				DisplayAppError(
					w,
					err,
					"Access Token is expired, get a new Token",
					401,
				)
				return

			default:
				DisplayAppError(w,
					err,
					"Error while parsing the Access Token!",
					500,
				)
				return
			}

		default:
			DisplayAppError(w,
				err,
				"Error while parsing Access Token!",
				500)
			return
		}

	}

	if token.Valid {
		ctx := NewContextWithEmail(r.Context(), token)
		next(w, r.WithContext(ctx))
	} else {
		DisplayAppError(
			w,
			err,
			"Invalid Access Token",
			401,
		)
	}
}

// TokenFromAuthHeader is a "TokenExtractor" that takes a given request and extracts
// the JWT token from the Authorization header in API
func TokenFromAuthHeader(r *http.Request) (string, error) {
	// Look for an Authorization header
	if ah := r.Header.Get("Authorization"); ah != "" {
		// Should be a bearer token
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			return ah[7:], nil
		}
	}
	return "", errors.New("No token in the HTTP request")
}
