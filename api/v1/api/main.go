package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"./database"
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

	if exists {
		println(jwtSecret)
		jwtKey = []byte(jwtSecret)
	}

	var err error
	database.DBCon, err = sql.Open("mysql", "root:99dZ%dtw&gE@tcp(127.0.0.1:4000)/servert")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin", "Access-Control-Request-Headers"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{"http://127.0.0.1:8080"})
	creds := handlers.AllowCredentials()

	// auth.go
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/isloggedin", routes.IsLoggedIn).Methods("GET")
	r.HandleFunc("/signup", routes.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signin", routes.Signin).Methods("POST", "OPTIONS")
	r.HandleFunc("/refresh", routes.Refresh).Methods("GET")

	//account.go
	account := r.PathPrefix("/account").Subrouter()
	account.HandleFunc("/intrest", routes.Intrest).Methods("POST", "OPTIONS")
	account.HandleFunc("/account", routes.Account).Methods("GET")
	account.HandleFunc("/accountinfo", routes.AccountInfo).Methods("GET")

	//products.go
	r.HandleFunc("/products/{prodID}", routes.GetProducts).Methods("GET")
	r.HandleFunc("/products", routes.GetProducts).Methods("GET")

	//orders.go
	r.HandleFunc("/create-payment-intent/{prodID}", routes.CreatePaymentIntent).Methods("GET")
	r.HandleFunc("/makeorder", routes.MakeOrder).Methods("POST", "OPTIONS")
	r.HandleFunc("/getorders", routes.GetOrders).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(header, methods, origin, creds)(r)))
}
