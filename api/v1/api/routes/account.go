package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../database"
	"../models"

	jwt "github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var jwtKey []byte

func checkifemailexistsreg(email string) bool {

	result, err := database.DBCon.Query("SELECT email from reg WHERE email=(?)", email)
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var emailfound string
		err := result.Scan(&emailfound)
		if err != nil {
			panic(err)
		}
		if emailfound != "" {
			// make new user id
			println(emailfound, "already exists")
			return true
		}
	}

	return false

}

func Intrest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	details := models.Details{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &details)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(details.Email, details.Name)

	data := data{Name: details.Name, Email: details.Email}

	err = data.Validate()

	if err != nil {
		fmt.Println(err)
		res := models.ResObj{Success: false, Message: "something went wrong"}
		json.NewEncoder(w).Encode(res)
		return
	}

	emailInDB := checkifemailexistsreg(details.Email)
	if emailInDB != false {
		res := models.ResObj{Success: false, Message: "Email already in use"}
		json.NewEncoder(w).Encode(res)
		return
	}

	result, err := database.DBCon.Query("INSERT INTO reg (name, email) VALUES (?,?)", details.Name, details.Email)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	log.Println(result)

	res := models.ResObj{Success: true, Message: "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)
}

func Account(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &models.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res := models.ResObj{Success: true, Message: claims.UserID}

	json.NewEncoder(w).Encode(res)
}

func AccountInfo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := database.DBCon.Query("SELECT name, email from users WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err)
	}

	var userDetailsList []models.UserDetails

	for result.Next() {
		var user models.UserDetails
		err := result.Scan(&user.Name, &user.Email)
		if err != nil {
			panic(err.Error())
		}
		userDetailsList = append(userDetailsList, user)
	}

	println(len(userDetailsList))

	if len(userDetailsList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := models.ResObj{Success: false, Message: "No orders found if you think this is wrong please contact us"}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(userDetailsList)
	json.NewEncoder(w).Encode(userDetailsList)

}

type data struct {
	Name  string
	Email string
}

func (d data) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required, validation.Length(2, 45)),
		validation.Field(&d.Email, validation.Required, is.Email),
	)
}
