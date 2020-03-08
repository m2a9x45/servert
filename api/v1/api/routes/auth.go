package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../database"
	"../models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

func IsLoggedIn(w http.ResponseWriter, r *http.Request) {
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

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	res := models.ResObj{Success: true, Message: claims.UserID}
	json.NewEncoder(w).Encode(res)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	signup := models.SignUpObj{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &signup)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(signup)

	// check signup details
	// https://github.com/go-ozzo/ozzo-validation
	// check that email doesn't already exit.

	emailInDB := checkifemailexists(signup.Email)
	if emailInDB != false {
		res := models.ResObj{Success: false, Message: "Email already in use"}
		json.NewEncoder(w).Encode(res)
		return
	}

	password := []byte(signup.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	userID := generateUserID()

	println(userID, "seletced")

	result, err := database.DBCon.Query("INSERT INTO users (name, user_id, email, password) VALUES (?,?,?,?)", signup.Name, userID, signup.Email, hashedPassword)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")

	// issue JWT so user can login into protected routes

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	res := models.ResObj{Success: true, Message: "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)

}

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	signin := models.SignUpObj{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &signin)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(signin, "here")

	// basic logging cred checking
	if signin.Email == "" || signin.Password == "" {
		log.Println("No email")
		res := models.ResObj{Success: false, Message: "Problem signing in"}
		json.NewEncoder(w).Encode(res)
		return
	}

	// check password.

	result, err := database.DBCon.Query("SELECT user_id,password from users WHERE email=(?)", signin.Email)
	if err != nil {
		println(err)
	}

	log.Println("DB check complete")

	defer result.Close()

	for result.Next() {
		var hashedPassword string
		var userID string
		err := result.Scan(&userID, &hashedPassword)
		if err != nil {
			panic(err)
		}
		// Comparing the password with the hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(signin.Password))
		if err == nil {
			expirationTime := time.Now().Add(5 * time.Minute)

			claims := &models.Claims{
				UserID: userID,
				StandardClaims: jwt.StandardClaims{
					// In JWT, the expiry time is expressed as unix milliseconds
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Create the JWT string
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				// If there is an error in creating the JWT return an internal server error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  expirationTime,
				HttpOnly: true,
				Path:     "/",
			})

			res := models.ResObj{Success: true, Message: "Details inserted into DB"}
			json.NewEncoder(w).Encode(res)
		}
		fmt.Println(err) // nil means it is a match
		return
	}

	res := models.ResObj{Success: false, Message: "Something went wrong"}
	json.NewEncoder(w).Encode(res)

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 75*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	})
}

func Reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	reset := models.Reset{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &reset)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(reset)

	// check if user exits
	emailInDB := checkifemailexists(reset.Email)
	if emailInDB != true {
		res := models.ResObj{Success: false, Message: "It doesn't look like you have an account with us"}
		json.NewEncoder(w).Encode(res)
		return
	}

	// send email with one time reset password link

	// create random token
	token := ksuid.New() // maybe check if this is allrady in the DB
	// Make a function that check if the new generated token allready exits in the BD

	now := time.Now()
	secs := now.Unix() + 1800

	fmt.Println(secs)

	// add token and email along with expirey time into the DB

	result, err := database.DBCon.Query("SELECT email from reset WHERE email=(?)", reset.Email)
	if err != nil {
		println(err)
	}

	var allemails []string

	for result.Next() {
		var email string
		err := result.Scan(&email)
		if err != nil {
			panic(err.Error())
		}
		allemails = append(allemails, email)
	}

	if len(allemails) != 0 {
		// means user has allready tried to reset pass word and we shold update token and expiry
		result, err := database.DBCon.Query("UPDATE reset SET token = (?), expires = (?) WHERE email = (?)", token.String(), secs, reset.Email)
		if err != nil {
			log.Fatal("Error wilst inserting into DB", err)
		}

		defer result.Close()

		fmt.Println("Inserted Into DB")
	} else {
		result, err := database.DBCon.Query("INSERT INTO reset (email, token, expires) VALUES (?,?,?)", reset.Email, token.String(), secs)
		if err != nil {
			log.Fatal("Error wilst inserting into DB", err)
		}

		defer result.Close()

		fmt.Println("Inserted Into DB")
	}

	// Send email

	url := "http://127.0.0.1:8080/web/website/signin/reset/index.html?t=" + token.String()

	fmt.Println(url)

	res := models.ResObj{Success: false, Message: url}
	json.NewEncoder(w).Encode(res)

}

func UpdatePasswordFromReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	reset := models.ResetPassword{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &reset)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(reset.Token)

	// check DB for token
	// Get the email fro that token then update the password in the user table

	result, err := database.DBCon.Query("SELECT email, expires from reset WHERE token=(?)", reset.Token)
	if err != nil {
		println(err)
	}

	defer result.Close()

	var allresetpassworddb []models.ResetPasswordDB

	for result.Next() {
		var resetpassworddb models.ResetPasswordDB
		err := result.Scan(&resetpassworddb.Email, &resetpassworddb.Expires)
		if err != nil {
			panic(err.Error())
		}
		allresetpassworddb = append(allresetpassworddb, resetpassworddb)
	}

	fmt.Println(allresetpassworddb[0])

	now := time.Now()
	secs := now.Unix()

	password := []byte(reset.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	if secs < allresetpassworddb[0].Expires {
		fmt.Println("Token is vaild")
		// change the password to the new password

		result, err := database.DBCon.Query("UPDATE users SET password = (?) WHERE email = (?)", hashedPassword, allresetpassworddb[0].Email)
		if err != nil {
			log.Fatal("Error wilst inserting into DB", err)
		}

		defer result.Close()

		fmt.Println("Inserted Into DB")
	} else {
		fmt.Println("Token is invaild")
		res := models.ResObj{Success: false, Message: "Token has expired"}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := models.ResObj{Success: true, Message: "Password updated"}
	json.NewEncoder(w).Encode(res)
}

func generateUserID() string {

	id := ksuid.New()
	userID := "user_" + id.String()

	println(userID)

	result, err := database.DBCon.Query("SELECT user_id from users WHERE user_id=(?)", userID)
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var userid string
		err := result.Scan(&userid)
		if err != nil {
			panic(err)
		}
		if userid != "" {
			// make new user id
			println(userid, "already exists")
			generateUserID()
		}
	}

	return userID
}

func checkifemailexists(email string) bool {

	result, err := database.DBCon.Query("SELECT email from users WHERE email=(?)", email)
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
