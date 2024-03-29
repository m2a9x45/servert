package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"./database"
	"./models"
	"./routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var jwtKey []byte

func main() {

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	// dbURL, _ := os.LookupEnv("DB_URL")
	dbURL, _ := os.LookupEnv("DB_URL_DEV")

	if exists {
		println(jwtSecret)
		println(dbURL)
		jwtKey = []byte(jwtSecret)
	}

	var err error
	database.DBCon, err = sql.Open("mysql", dbURL)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin", "Access-Control-Request-Headers"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"})
	origin := handlers.AllowedOrigins([]string{"https://servert.co.uk", "http://127.0.0.1:8080"})
	creds := handlers.AllowCredentials()

	// auth.go
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/isloggedin", routes.IsLoggedIn).Methods("GET")
	auth.HandleFunc("/signup", routes.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signin", routes.Signin).Methods("POST", "OPTIONS")
	auth.HandleFunc("/refresh", routes.Refresh).Methods("GET")
	auth.HandleFunc("/reset", routes.Reset).Methods("POST", "OPTIONS")
	auth.HandleFunc("/restpassword", routes.UpdatePasswordFromReset).Methods("PATCH", "OPTIONS")

	auth.HandleFunc("/stafflogin", routes.SigninStaff).Methods("POST", "OPTIONS")
	auth.HandleFunc("/isstaffloggedin", routes.IsLoggedInStaff).Methods("GET")

	//account.go
	account := r.PathPrefix("/account").Subrouter()
	account.Use(routes.Userauth)
	account.HandleFunc("/account", routes.Account).Methods("GET")
	account.HandleFunc("/accountinfo", routes.AccountInfo).Methods("GET")
	account.HandleFunc("/receipt/{receiptID}", routes.Getreceipt).Methods("GET")
	account.HandleFunc("/customercards", routes.Getcustomercards).Methods("GET")

	r.HandleFunc("/intrest", routes.Intrest).Methods("POST", "OPTIONS")

	//products.go
	r.HandleFunc("/products/{prodID}", routes.GetProducts).Methods("GET")
	r.HandleFunc("/products", routes.GetProducts).Methods("GET")

	//orders.go
	order := r.PathPrefix("/order").Subrouter()
	order.Use(routes.Userauth)
	order.HandleFunc("/create-payment-intent/{prodID}/{dur}/{cardID}", routes.CreatePaymentIntent).Methods("GET")
	order.HandleFunc("/create-payment-intent/{prodID}/{dur}/", routes.CreatePaymentIntent).Methods("GET")
	order.HandleFunc("/makeorder", routes.MakeOrder).Methods("POST", "OPTIONS")
	order.HandleFunc("/getorders", routes.GetOrders).Methods("GET")

	// internal
	internal := r.PathPrefix("/internal").Subrouter()
	internal.Use(routes.Staffauth)
	internal.HandleFunc("/tasks", routes.GetTasks).Methods("GET")
	internal.HandleFunc("/refresh", routes.RefreshStaff).Methods("GET")

	// servers
	servers := r.PathPrefix("/servers").Subrouter()
	servers.Use(routes.Userauth)
	servers.HandleFunc("/active", routes.GetServer).Methods("GET")
	r.HandleFunc("/", hello).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(header, methods, origin, creds)(r)))
}

func hello(w http.ResponseWriter, r *http.Request) {
	res := &models.Hello{Name: "servert-api"}
	json.NewEncoder(w).Encode(res)
}
