package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

type Product struct {
	ID       string  `json:"id"`
	UIDD     string  `json:"uuid"`
	Name     string  `json:"name"`
	Des      string  `json:"des"`
	CPU      string  `json:"cpu"`
	RAM      string  `json:"ram"`
	Disk     string  `json:"disk"`
	Price    float64 `json:"price"`
	Instock  bool    `json:"instock"`
	Setupfee float64 `json:"setupfee"`
	Discount float64 `json:"discount"`
}

type details struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type resObj struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type signUpObj struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims is the data structure for encodeing the JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type CheckoutData struct {
	ClientSecret string `json:"clientecret"`
}

type OrderData struct {
	PaymentID string `json="PaymentID"`
	ProductID string `json="ProductID"`
}

type OrderObj struct {
	OrderID string `json="order_id"`
	ProdID  string `json="prod_id"`
}

type UserDetails struct {
	Name  string `json="name"`
	Email string `json="email"`
}

var db *sql.DB
var err error

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var jwtKey []byte

func main() {

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")

	if exists {
		println(jwtSecret)
		jwtKey = []byte(jwtSecret)
	}

	db, err = sql.Open("mysql", "root:99dZ%dtw&gE@tcp(127.0.0.1:4000)/servert")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin", "Access-Control-Request-Headers"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{"http://127.0.0.1:8080"})
	creds := handlers.AllowCredentials()

	r.HandleFunc("/products/{prodID}", GetProducts).Methods("GET")
	r.HandleFunc("/products", GetProducts).Methods("GET")
	r.HandleFunc("/intrest", RegIntrest).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/signin", signin).Methods("POST", "OPTIONS")
	r.HandleFunc("/account", account).Methods("GET")
	r.HandleFunc("/acountInfo", accountInfo).Methods("GET")
	r.HandleFunc("/create-payment-intent/{prodID}", createPaymentIntent).Methods("GET")
	r.HandleFunc("/loggedIn", isLoggedIn).Methods("GET")
	r.HandleFunc("/refresh", Refresh).Methods("GET")
	r.HandleFunc("/order", Order).Methods("POST", "OPTIONS")
	r.HandleFunc("/order", GetOrders).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(header, methods, origin, creds)(r)))

}

func generateUserID() string {

	id := ksuid.New()
	userID := "user_" + id.String()

	println(userID)

	result, err := db.Query("SELECT user_id from users WHERE user_id=(?)", userID)
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

func generateOrderID() string {

	id := ksuid.New()
	orderID := "order_" + id.String()

	println(orderID)

	result, err := db.Query("SELECT order_id from orders WHERE order_id=(?)", orderID)
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var orderid string
		err := result.Scan(&orderid)
		if err != nil {
			panic(err)
		}
		if orderid != "" {
			// make new user id
			println(orderid, "already exists")
			generateOrderID()
		}
	}

	return orderID
}

func checkifemailexists(email string) bool {

	result, err := db.Query("SELECT email from users WHERE email=(?)", email)
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

// RegIntrest will add name and email to DB
func RegIntrest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	details := details{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &details)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(details.Email, details.Name)

	result, err := db.Query("INSERT INTO reg (name, email) VALUES (?,?)", details.Name, details.Email)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	log.Println(result)

	res := resObj{true, "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)
}

// GetProducts will return a list of products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id := vars["prodID"]

	if id != "" {
		var allproducts []Product

		result, err := db.Query("SELECT * from products WHERE prod_id=(?)", id)
		if err != nil {
			println(err)
		}

		defer result.Close()

		for result.Next() {
			var product Product
			err := result.Scan(&product.ID, &product.UIDD, &product.Name, &product.Des, &product.CPU, &product.RAM, &product.Disk, &product.Price, &product.Instock, &product.Setupfee, &product.Discount)
			if err != nil {
				panic(err.Error())
			}
			allproducts = append(allproducts, product)
		}

		println(len(allproducts))

		if len(allproducts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := resObj{false, "product not found"}
			json.NewEncoder(w).Encode(res)
		}

		fmt.Println(allproducts)
		json.NewEncoder(w).Encode(allproducts)
		return
	}

	// type Category product

	var allproducts []Product

	result, err := db.Query("SELECT * from products")
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var product Product
		err := result.Scan(&product.ID, &product.UIDD, &product.Name, &product.Des, &product.CPU, &product.RAM, &product.Disk, &product.Price, &product.Instock, &product.Setupfee, &product.Discount)
		if err != nil {
			panic(err.Error())
		}
		allproducts = append(allproducts, product)
	}

	fmt.Println(allproducts)
	json.NewEncoder(w).Encode(allproducts)
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	signup := signUpObj{}

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
		res := resObj{false, "Email already in use"}
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

	result, err := db.Query("INSERT INTO users (name, user_id, email, password) VALUES (?,?,?,?)", signup.Name, userID, signup.Email, hashedPassword)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")

	// issue JWT so user can login into protected routes

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
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
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	res := resObj{true, "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)

}

func signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	signin := signUpObj{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &signin)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(signin)

	// basic logging cred checking
	if signin.Email == "" || signin.Password == "" {
		log.Println("No email")
		res := resObj{false, "Problem signing in"}
		json.NewEncoder(w).Encode(res)
		return
	}

	// check password.

	result, err := db.Query("SELECT user_id,password from users WHERE email=(?)", signin.Email)
	if err != nil {
		println(err)
	}

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

			claims := &Claims{
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
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})

			res := resObj{true, "Details inserted into DB"}
			json.NewEncoder(w).Encode(res)
		}
		fmt.Println(err) // nil means it is a match
		return
	}

	res := resObj{false, "Something went wrong"}
	json.NewEncoder(w).Encode(res)

}

func account(w http.ResponseWriter, r *http.Request) {
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
	claims := &Claims{}

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

	res := resObj{true, claims.UserID}

	json.NewEncoder(w).Encode(res)
}

func createPaymentIntent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["prodID"]

	if id != "" {
		var allproducts []Product

		result, err := db.Query("SELECT * from products WHERE prod_id=(?)", id)
		if err != nil {
			println(err)
			res := resObj{false, "Couldn't find product"}
			json.NewEncoder(w).Encode(res)
			return
		}

		defer result.Close()

		for result.Next() {
			var product Product
			err := result.Scan(&product.ID, &product.UIDD, &product.Name, &product.Des, &product.CPU, &product.RAM, &product.Disk, &product.Price, &product.Instock, &product.Setupfee, &product.Discount)
			if err != nil {
				panic(err.Error())
			}
			allproducts = append(allproducts, product)
		}

		if len(allproducts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := resObj{false, "product not found"}
			json.NewEncoder(w).Encode(res)
			return
		}

		stripe.Key = "sk_test_OGXIlmLXL1Gvhpa9jqBdxutN00YB96uOjP"

		price := allproducts[0].Price
		pricePennies := int64(price * 100)
		println(price, " in £")
		println(price, " in p")

		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(pricePennies),
			Currency: stripe.String(string(stripe.CurrencyGBP)),
		}

		intent, _ := paymentintent.New(params)

		data := CheckoutData{
			ClientSecret: intent.ClientSecret,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)

	}
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) {
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
	claims := &Claims{}

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
	res := resObj{true, claims.UserID}
	json.NewEncoder(w).Encode(res)
}

// Refresh will return a new JWT when passed a vaild one
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
	claims := &Claims{}
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
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

// Order will add a new order to DB
func Order(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

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

	claims := &Claims{}

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

	order := OrderData{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &order)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	OrderID := generateOrderID()

	println(OrderID, "seletced")

	result, err := db.Query("INSERT INTO orders (order_id, user_id, payment_id, prod_id) VALUES (?,?,?,?)", OrderID, claims.UserID, order.PaymentID, order.ProductID)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")

	res := resObj{true, "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)

}

// GetOrders will return a list of orders for an account
func GetOrders(w http.ResponseWriter, r *http.Request) {
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

	claims := &Claims{}

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

	result, err := db.Query("SELECT order_id, prod_id from orders WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err)
	}

	var allorders []OrderObj

	for result.Next() {
		var order OrderObj
		err := result.Scan(&order.OrderID, &order.ProdID)
		if err != nil {
			panic(err.Error())
		}
		allorders = append(allorders, order)
	}

	println(len(allorders))

	if len(allorders) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := resObj{false, "No orders found if you think this is wrong please contact us"}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(allorders)
	json.NewEncoder(w).Encode(allorders)
}

func accountInfo(w http.ResponseWriter, r *http.Request) {
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

	claims := &Claims{}

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

	result, err := db.Query("SELECT name, email from users WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err)
	}

	var userDetailsList []UserDetails

	for result.Next() {
		var user UserDetails
		err := result.Scan(&user.Name, &user.Email)
		if err != nil {
			panic(err.Error())
		}
		userDetailsList = append(userDetailsList, user)
	}

	println(len(userDetailsList))

	if len(userDetailsList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := resObj{false, "No orders found if you think this is wrong please contact us"}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(userDetailsList)
	json.NewEncoder(w).Encode(userDetailsList)

}
